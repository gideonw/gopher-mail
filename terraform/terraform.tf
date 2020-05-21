terraform {
  required_version = ">= 0.12"
}

####################################################################################
# Variables
variable "aws_region" {
  type = string
}

variable "s3_bucket_prefix" {
  description = "Prefix used when creating the S3 buckets to make them globally unique."
  type        = string
}

variable "base_domain" {
  description = "Base domain must be an already existing Route53 Hosted zone."
  type        = string
  default     = ""
}

variable "sub_domain" {
  type    = string
  default = ""
}

# Email
variable "email_mailbox_bucket" {
  type        = string
  default     = ""
  description = "S3 bucket SES uses to save emails."
}

variable "email_post_office_prefix" {
  type        = string
  default     = "post-office"
  description = "S3 prefix SES uses to save emails."
}

variable "email_mailbox_prefix" {
  type        = string
  default     = "mailbox"
  description = "S3 prefix postmaster uses to sort emails."
}

####################################################################################
# Locals
locals {
  app_name    = "gopher-mail"
  full_domain = "${var.sub_domain}.${var.base_domain}"

  s3_origin_id          = "gopher-mail-web"
  api_gateway_origin_id = "gopher-mail-api"

  dash_domain  = replace(var.base_domain, ".", "-")
  email_bucket = var.email_mailbox_bucket != "" ? var.email_mailbox_bucket : "${local.dash_domain}-email"

  mailman_routes = [
    "GET /api/emails",
    "GET /api/email/{emailID}",
  ]

  tags = {
    App = "gopher-mail"
  }
}

####################################################################################
# Outputs

output "lambda_archive_bucket" {
  value = aws_s3_bucket.lambda_archive.bucket
}

output "apigateway_endpoint" {
  value = aws_apigatewayv2_api.gopher_mail.api_endpoint
}

####################################################################################
# Provider
provider "aws" {
  version = "~> 2.0"
  region  = var.aws_region
}

provider "aws" {
  version = "~> 2.0"
  alias   = "east"
  region  = "us-east-1"
}

provider "random" {}

data "aws_caller_identity" "account" {}

data "aws_region" "selected" {
  name = var.aws_region
}

####################################################################################
# Resources
####################################################################################

# Gopher mail web resources
####################################################################################
data "aws_route53_zone" "base_domain_zone" {
  name         = var.base_domain
  private_zone = false
}

resource "aws_acm_certificate" "shared_certificate" {
  provider = aws.east

  validation_method = "DNS"
  domain_name       = var.base_domain
  subject_alternative_names = [
    "*.${var.base_domain}"
  ]

  tags = local.tags

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "certificate_validation_record" {
  zone_id = data.aws_route53_zone.base_domain_zone.zone_id

  name    = aws_acm_certificate.shared_certificate.domain_validation_options[0].resource_record_name
  type    = aws_acm_certificate.shared_certificate.domain_validation_options[0].resource_record_type
  records = [aws_acm_certificate.shared_certificate.domain_validation_options[0].resource_record_value]

  ttl = 60
}

resource "aws_acm_certificate_validation" "shared_certificat_valid" {
  provider = aws.east

  certificate_arn         = aws_acm_certificate.shared_certificate.arn
  validation_record_fqdns = [aws_route53_record.certificate_validation_record.fqdn]
}

resource "aws_cloudfront_origin_access_identity" "gopher_mail_web" {
  comment = "Gopher Mail static website s3 access identity."
}

data "aws_iam_policy_document" "gopher_mail_web" {
  statement {
    sid       = "CloudFrontAccess"
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.gopher_mail_web.arn}/*"]

    principals {
      type        = "AWS"
      identifiers = [aws_cloudfront_origin_access_identity.gopher_mail_web.iam_arn]
    }
  }
}

# Gopher Mail static website bucket
resource "aws_s3_bucket" "gopher_mail_web" {
  bucket_prefix = "${local.full_domain}-"
  acl           = "private"

  tags = local.tags
}

resource "aws_s3_bucket_policy" "gopher_mail_web" {
  bucket = aws_s3_bucket.gopher_mail_web.id
  policy = data.aws_iam_policy_document.gopher_mail_web.json
}

resource "random_string" "cf_verify_header" {
  length  = 16
  special = false
}

resource "random_string" "cf_verify_value" {
  length           = 16
  special          = true
  override_special = ")(-+*%$#"
}

resource "aws_cloudfront_distribution" "gopher_mail" {
  enabled = true

  is_ipv6_enabled     = true
  comment             = "Gopher Mail cloudfront distribution that brings all parts under one domain."
  aliases             = [local.full_domain]
  default_root_object = "index.html"

  origin {
    domain_name = aws_s3_bucket.gopher_mail_web.bucket_regional_domain_name
    origin_id   = local.s3_origin_id

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.gopher_mail_web.cloudfront_access_identity_path
    }
  }

  origin {
    domain_name = replace(aws_apigatewayv2_api.gopher_mail.api_endpoint, "https://", "")
    origin_id   = "apigateway-${aws_apigatewayv2_api.gopher_mail.id}"

    custom_header {
      name  = random_string.cf_verify_header.result
      value = random_string.cf_verify_value.result
    }

    custom_origin_config {
      http_port  = "80"
      https_port = "443"

      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  default_cache_behavior {
    target_origin_id = local.s3_origin_id
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
  }

  ordered_cache_behavior {
    target_origin_id = "apigateway-${aws_apigatewayv2_api.gopher_mail.id}"
    path_pattern     = "/api/*"
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
  }

  price_class = "PriceClass_200"
  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = aws_acm_certificate.shared_certificate.arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2018"
  }

  tags = local.tags

  depends_on = [aws_route53_record.mail_domain_verification]
}

resource "aws_route53_record" "gopher_mail" {
  zone_id = data.aws_route53_zone.base_domain_zone.zone_id
  name    = local.full_domain
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.gopher_mail.domain_name
    zone_id                = aws_cloudfront_distribution.gopher_mail.hosted_zone_id
    evaluate_target_health = false
  }
}

# Email S3 and SES
####################################################################################

data "aws_iam_policy_document" "mailbox" {
  statement {
    sid    = "AllowSESPutObject"
    effect = "Allow"

    actions = [
      "s3:PutObject"
    ]

    principals {
      type        = "Service"
      identifiers = ["ses.amazonaws.com"]
    }

    resources = [
      "${aws_s3_bucket.mailbox.arn}/*"
    ]

    condition {
      test     = "StringEquals"
      variable = "aws:Referer"

      values = [
        data.aws_caller_identity.account.account_id
      ]
    }
  }
}

resource "aws_s3_bucket" "mailbox" {
  bucket_prefix = "${local.email_bucket}-"
  acl           = "private"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  tags = local.tags
}

# explicitly set these in-case they are changed, on terraform apply it will be reapplied
resource "aws_s3_bucket_public_access_block" "force_private" {
  bucket = aws_s3_bucket.mailbox.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_policy" "ses_perms" {
  bucket = aws_s3_bucket.mailbox.id
  policy = data.aws_iam_policy_document.mailbox.json
}

# SES
resource "aws_ses_domain_identity" "mail_domain" {
  domain = var.base_domain

  depends_on = [aws_lambda_function.postmaster] // Postmaster needs to be available to get sns notifications
}

resource "aws_route53_record" "mail_domain_inbound_mx" {
  zone_id = data.aws_route53_zone.base_domain_zone.zone_id
  name    = var.base_domain
  type    = "MX"
  ttl     = "1800"
  records = ["10 inbound-smtp.${data.aws_region.selected.name}.amazonaws.com"] # use data region to ensure that the region is correct
}

resource "aws_route53_record" "mail_domain_verification" {
  zone_id = data.aws_route53_zone.base_domain_zone.zone_id
  name    = "_amazonses.${var.base_domain}"
  type    = "TXT"
  ttl     = "1800"
  records = [aws_ses_domain_identity.mail_domain.verification_token]
}

resource "aws_ses_domain_identity_verification" "mail_domain_identity_verified" {
  domain = aws_ses_domain_identity.mail_domain.domain

  depends_on = [aws_route53_record.mail_domain_verification]
}

resource "aws_ses_domain_dkim" "mail_domain" {
  domain = aws_ses_domain_identity.mail_domain.domain
}

resource "aws_route53_record" "mail_domain_dkim" {
  count = 3

  zone_id = data.aws_route53_zone.base_domain_zone.zone_id
  name    = "${element(aws_ses_domain_dkim.mail_domain.dkim_tokens, count.index)}._domainkey.${var.base_domain}"
  type    = "CNAME"
  ttl     = "1800"
  records = ["${element(aws_ses_domain_dkim.mail_domain.dkim_tokens, count.index)}.dkim.amazonses.com"]
}

resource "aws_ses_domain_mail_from" "mail_domain" {
  domain           = aws_ses_domain_identity.mail_domain.domain
  mail_from_domain = "bounce.${aws_ses_domain_identity.mail_domain.domain}"
}

resource "aws_route53_record" "mail_domain_from_mx" {
  zone_id = data.aws_route53_zone.base_domain_zone.zone_id
  name    = aws_ses_domain_mail_from.mail_domain.mail_from_domain
  type    = "MX"
  ttl     = "1800"
  records = ["10 feedback-smtp.${data.aws_region.selected.name}.amazonses.com"] # use data region to ensure that the region is correct
}

resource "aws_route53_record" "mail_domain_from_txt" {
  zone_id = data.aws_route53_zone.base_domain_zone.zone_id
  name    = aws_ses_domain_mail_from.mail_domain.mail_from_domain
  type    = "TXT"
  ttl     = "1800"
  records = ["v=spf1 include:amazonses.com ~all"]
}

# resource "aws_ses_receipt_rule_set" "mail_domain" {
#   rule_set_name = "${local.dash_domain}-rules"
# }

resource "aws_ses_receipt_rule" "mail_domain" {
  enabled       = true
  name          = "${local.dash_domain}-save-to-s3"
  rule_set_name = "default-rule-set"

  recipients = [
    var.base_domain
  ]
  scan_enabled = true

  s3_action {
    position = 1

    bucket_name       = aws_s3_bucket.mailbox.id
    object_key_prefix = var.email_post_office_prefix
    topic_arn         = aws_sns_topic.new_email.arn
  }

  depends_on = [
    aws_s3_bucket_policy.ses_perms
  ]
}

resource "aws_sns_topic" "new_email" {
  name_prefix  = "${local.dash_domain}-new-email-"
  display_name = "${local.dash_domain}-new-email"

  tags = local.tags
}

# Lambda related resources
####################################################################################

data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    sid    = "AllowLambdaAssumeRole"
    effect = "Allow"

    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "lambda_logging" {
  statement {
    sid    = "AllowLambdaLogging"
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    # TODO: Change the first star to a prefix that can be used by all lambdas
    resources = [
      "arn:aws:logs:*:*:*"
    ]
  }
}

# Lambda zip archive bucket
resource "aws_s3_bucket" "lambda_archive" {
  bucket_prefix = "${local.app_name}-lambda-archive-"
  acl           = "private"

  versioning {
    enabled = true
  }

  tags = local.tags
}

# Postmaster Lambda
data "aws_iam_policy_document" "postmaster" {
  statement {
    sid    = "S3ReadWrite"
    effect = "Allow"

    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject",
    ]
    resources = [
      "${aws_s3_bucket.mailbox.arn}/*"
    ]
  }

  statement {
    sid    = "S3List"
    effect = "Allow"

    actions = [
      "s3:ListBucket",
    ]
    resources = [
      aws_s3_bucket.mailbox.arn
    ]
  }
}

data "aws_s3_bucket_object" "postmaster" {
  bucket = aws_s3_bucket.lambda_archive.id
  key    = "postmaster/postmaster.zip"

  depends_on = [
    aws_s3_bucket_object.postmaster
  ]
}

resource "aws_s3_bucket_object" "postmaster" {
  bucket = aws_s3_bucket.lambda_archive.id
  key    = "postmaster/postmaster.zip"

  source = "${path.root}/../bin/postmaster.zip"

  # lifecycle {
  #   ignore_changes = [
  #     source
  #   ]
  # }
}

resource "aws_cloudwatch_log_group" "postmaster" {
  name              = "/aws/lambda/${aws_lambda_function.postmaster.function_name}"
  retention_in_days = 14
}

resource "aws_iam_role" "postmaster" {
  name = "${local.app_name}-postmaster"

  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

resource "aws_iam_role_policy" "postmaster_service_permissions" {
  name = "postmaster-service-permissions"
  role = aws_iam_role.postmaster.id

  policy = data.aws_iam_policy_document.postmaster.json
}

resource "aws_iam_role_policy" "postmaster_log_permissions" {
  name = "postmaster-log-permissions"
  role = aws_iam_role.postmaster.id

  policy = data.aws_iam_policy_document.lambda_logging.json
}

resource "aws_lambda_function" "postmaster" {
  function_name = "${local.app_name}-postmaster"

  s3_bucket         = aws_s3_bucket.lambda_archive.id
  s3_key            = data.aws_s3_bucket_object.postmaster.key
  s3_object_version = data.aws_s3_bucket_object.postmaster.version_id

  role = aws_iam_role.postmaster.arn

  handler = "postmaster"
  runtime = "go1.x"

  memory_size = 512

  environment {
    variables = {
      DOMAIN             = var.base_domain
      MAILBOX_BUCKET     = aws_s3_bucket.mailbox.id
      POST_OFFICE_PREFIX = var.email_post_office_prefix
      MAILBOX_PREFIX     = var.email_mailbox_prefix
    }
  }

  depends_on = [
    # aws_cloudwatch_log_group.postmaster,
    aws_iam_role_policy.postmaster_log_permissions,
    aws_s3_bucket.lambda_archive,
    aws_s3_bucket_object.postmaster
  ]

  tags = local.tags
}

resource "aws_lambda_permission" "sns_email_trigger" {
  statement_id  = "SNSTriggerPermission"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.postmaster.function_name
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.new_email.arn
}

resource "aws_sns_topic_subscription" "postmaster_email_notification" {
  topic_arn = aws_sns_topic.new_email.arn
  protocol  = "lambda"
  endpoint  = aws_lambda_function.postmaster.arn
}


# Mailman Lambda
data "aws_iam_policy_document" "mailman" {
  statement {
    sid    = "S3ReadWrite"
    effect = "Allow"

    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject",
    ]
    resources = [
      "${aws_s3_bucket.mailbox.arn}/*"
    ]
  }

  statement {
    sid    = "S3List"
    effect = "Allow"

    actions = [
      "s3:ListBucket",
    ]
    resources = [
      aws_s3_bucket.mailbox.arn
    ]
  }
}

data "aws_s3_bucket_object" "mailman" {
  bucket = aws_s3_bucket.lambda_archive.id
  key    = "mailman/mailman.zip"

  depends_on = [
    aws_s3_bucket_object.mailman
  ]
}

resource "aws_s3_bucket_object" "mailman" {
  bucket = aws_s3_bucket.lambda_archive.id
  key    = "mailman/mailman.zip"

  source = "${path.root}/../bin/mailman.zip"

  # lifecycle {
  #   ignore_changes = [
  #     source
  #   ]
  # }
}

resource "aws_cloudwatch_log_group" "mailman" {
  name              = "/aws/lambda/${aws_lambda_function.mailman.function_name}"
  retention_in_days = 14
}

resource "aws_iam_role" "mailman" {
  name = "${local.app_name}-mailman"

  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

resource "aws_iam_role_policy" "mailman_service_permissions" {
  name = "mailman-service-permissions"
  role = aws_iam_role.mailman.id

  policy = data.aws_iam_policy_document.mailman.json
}

resource "aws_iam_role_policy" "mailman_log_permissions" {
  name = "mailman-log-permissions"
  role = aws_iam_role.mailman.id

  policy = data.aws_iam_policy_document.lambda_logging.json
}

resource "aws_lambda_function" "mailman" {
  function_name = "${local.app_name}-mailman"

  s3_bucket         = aws_s3_bucket.lambda_archive.id
  s3_key            = data.aws_s3_bucket_object.mailman.key
  s3_object_version = data.aws_s3_bucket_object.mailman.version_id

  role = aws_iam_role.mailman.arn

  handler = "mailman"
  runtime = "go1.x"

  memory_size = 512

  environment {
    variables = {
      DOMAIN           = var.base_domain
      MAILBOX_BUCKET   = aws_s3_bucket.mailbox.id
      MAILBOX_PREFIX   = var.email_mailbox_prefix
      CF_VERIFY_HEADER = random_string.cf_verify_header.result
      CF_VERIFY_VALUE  = random_string.cf_verify_value.result
    }
  }

  depends_on = [
    # aws_cloudwatch_log_group.mailman,
    aws_iam_role_policy.mailman_log_permissions,
    aws_s3_bucket.lambda_archive,
    aws_s3_bucket_object.mailman
  ]

  tags = local.tags
}

resource "aws_lambda_permission" "mailman_apigateway_invoke" {
  statement_id  = "APIGatewayPermission"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.mailman.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.gopher_mail.execution_arn}/*"
}

# API
####################################################################################
resource "aws_apigatewayv2_api" "gopher_mail" {
  name          = "gopher-mail"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id = aws_apigatewayv2_api.gopher_mail.id
  name   = "$default"

  auto_deploy = true

  # https://github.com/terraform-providers/terraform-provider-aws/issues/12893
  lifecycle {
    ignore_changes = [deployment_id, default_route_settings]
  }
}

resource "aws_apigatewayv2_integration" "mailman_lambda" {
  api_id = aws_apigatewayv2_api.gopher_mail.id

  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.mailman.invoke_arn

  lifecycle {
    ignore_changes = [passthrough_behavior]
  }
}

resource "aws_apigatewayv2_route" "mailman_routes" {
  count = length(local.mailman_routes)

  api_id = aws_apigatewayv2_api.gopher_mail.id

  route_key = element(local.mailman_routes, count.index)
  target    = "integrations/${aws_apigatewayv2_integration.mailman_lambda.id}"
}
