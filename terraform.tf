terraform {
  required_version = ">= 0.12"

  backend "s3" {
    bucket = "terraform-state-bucket"
    key    = "gopher-mail.tfstate"
  }
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

# Lambda related resources
####################################################################################

# Lambda zip archive bucket
resource "aws_s3_bucket" "lambda_archive" {
  bucket_prefix = "${local.app_name}-lambda-archive-"
  acl           = "private"

  versioning {
    enabled = true
  }

  tags = local.tags
}

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

# Gopher mail lambda resources
####################################################################################
