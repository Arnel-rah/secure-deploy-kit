variable "project_name" {
  type    = string
  default = "secure-deploy-kit"
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

# Prod defaults to real AWS. Flip to true only if you deliberately want
# to dry-run the prod config against LocalStack first.
variable "use_localstack" {
  type    = bool
  default = false
}

variable "aws_access_key" {
  description = "Leave default and use a real AWS profile/credentials chain in prod"
  type        = string
  default     = null
}

variable "aws_secret_key" {
  type      = string
  default   = null
  sensitive = true
}

variable "container_image" {
  type = string
}

variable "app_port" {
  type    = number
  default = 8080
}

variable "enable_waf" {
  type    = bool
  default = true
}
