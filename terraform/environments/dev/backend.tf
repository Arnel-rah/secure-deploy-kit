# Local state is fine for solo learning/dev use on LocalStack.
# For a shared/prod setup, switch to an S3 backend (real AWS) or
# LocalStack's emulated S3 backend — see docs/getting-started.md.
terraform {
  backend "local" {
    path = "terraform.tfstate"
  }
}
