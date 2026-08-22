terraform {
  required_version = ">= 1.6"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
}

# NOTE: when using `tflocal` (LocalStack's Terraform wrapper), the provider
# below is transparently patched to point at LocalStack's endpoints — you
# do not need to edit this file to switch between local and real AWS.
# `tflocal` is a thin wrapper around `terraform` that injects the
# LocalStack endpoint config at plan/apply time.
#
# Running against real AWS: use plain `terraform` instead of `tflocal`,
# and make sure your AWS credentials are configured (aws configure).
provider "aws" {
  region                      = var.aws_region
  access_key                  = var.aws_access_key
  secret_key                  = var.aws_secret_key
  skip_credentials_validation = var.use_localstack
  skip_metadata_api_check     = var.use_localstack
  skip_requesting_account_id  = var.use_localstack

  default_tags {
    tags = local.common_tags
  }
}
