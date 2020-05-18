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
variable "email_post_office_bucket" {
  type        = string
  default     = ""
  description = "S3 bucket SES uses to save emails."
}

variable "email_post_office_prefix" {
  type        = string
  default     = "post-office"
  description = "S3 prefix SES uses to save emails."
}

variable "email_ses_notification_topic" {
  type = string

}

####################################################################################
# Locals
locals {
  app_name    = "gopher-mail"
  full_domain = "${var.sub_domain}.${var.base_domain}"

  s3_origin_id          = "gopher-mail-web"
  api_gateway_origin_id = "gopher-mail-api"

  tags = {
    App = "gopher-mail"
  }
}

####################################################################################
# Outputs

output "lambda_archive_bucket" {
  value = aws_s3_bucket.lambda_archive.bucket
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
  count   = 2 # 1 for the base domain, and 1 for the subject_alternative_names
  zone_id = data.aws_route53_zone.base_domain_zone.zone_id

  name    = element(aws_acm_certificate.shared_certificate.domain_validation_options, count.index).resource_record_name
  type    = element(aws_acm_certificate.shared_certificate.domain_validation_options, count.index).resource_record_type
  records = [element(aws_acm_certificate.shared_certificate.domain_validation_options, count.index).resource_record_value]

  ttl = 60
}

resource "aws_acm_certificate_validation" "shared_certificat_valid" {
  provider = aws.east

  certificate_arn         = aws_acm_certificate.shared_certificate.arn
  validation_record_fqdns = [aws_route53_record.certificate_validation_record[0].fqdn]
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
  bucket_prefix = local.full_domain
  acl           = "private"

  tags = local.tags
}

resource "aws_s3_bucket_policy" "gopher_mail_web" {
  bucket = aws_s3_bucket.gopher_mail_web.id
  policy = data.aws_iam_policy_document.gopher_mail_web.json
}

resource "aws_cloudfront_distribution" "gopher_mail" {
  origin {
    domain_name = local.full_domain
    origin_id   = local.s3_origin_id

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.gopher_mail_web.cloudfront_access_identity_path
    }
  }

  enabled         = true
  is_ipv6_enabled = true
  comment         = "Gopher Mail cloudfront distribution that brings all parts under one domain."

  default_root_object = "index.html"

  aliases = [local.full_domain]

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = local.s3_origin_id

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

# SES
####################################################################################

# TODO: Add SES configuration to support creating gopher-mail from scratch.



# Lambda related resources
####################################################################################

data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    sid    = "AllowLambdaAssumeRole"
    effect = "Allow"

    actions = [
      "sts::AssumeRole",
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
      "${var.email_post_office_bucket}/*"
    ]
  }

  statement {
    sid    = "S3List"
    effect = "Allow"

    actions = [
      "s3:ListBucket",
    ]
    resources = [
      var.email_post_office_bucket
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

  lifecycle {
    ignore_changes = [
      source
    ]
  }
}

resource "aws_cloudwatch_log_group" "postmaster" {
  name_prefix       = "/aws/lambda/${local.app_name}-postmaster"
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

  handler = "bin/postmaster"
  runtime = "go1.x"

  memory_size = 512

  environment {
    variables = {
      DOMAIN             = var.base_domain
      POST_OFFICE_BUCKET = var.email_post_office_bucket
      POST_OFFICE_PREFIX = var.email_post_office_prefix
    }
  }

  depends_on = [
    aws_cloudwatch_log_group.postmaster,
    aws_iam_role_policy.postmaster_log_permissions,
    aws_s3_bucket.lambda_archive
  ]
}

resource "aws_lambda_permission" "sns_email_trigger" {
  statement_id  = "SNSTriggerPermission"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.postmaster.function_name
  principal     = "sns.amazonaws.com"
  source_arn    = var.email_ses_notification_topic
}

resource "aws_sns_topic_subscription" "postmaster_email_notification" {
  topic_arn = var.email_ses_notification_topic
  protocol  = "lambda"
  endpoint  = aws_lambda_function.postmaster.arn
}

# API
####################################################################################
