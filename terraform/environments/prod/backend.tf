# Real remote state for prod: uncomment and fill in once you have a real
# AWS account + an S3 bucket and DynamoDB table dedicated to Terraform state.
# See docs/getting-started.md → "Going to real AWS".
#
# terraform {
#   backend "s3" {
#     bucket         = "your-terraform-state-bucket"
#     key            = "secure-deploy-kit/prod/terraform.tfstate"
#     region         = "us-east-1"
#     dynamodb_table = "your-terraform-locks"
#     encrypt        = true
#   }
# }

terraform {
  backend "local" {
    path = "terraform.tfstate"
  }
}
