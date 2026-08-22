variable "project_name" {
  type    = string
  default = "secure-deploy-kit"
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "use_localstack" {
  description = "Set to false when running against real AWS"
  type        = bool
  default     = true
}

variable "aws_access_key" {
  description = "Dummy value is fine for LocalStack; use a real key/profile for real AWS"
  type        = string
  default     = "test"
}

variable "aws_secret_key" {
  description = "Dummy value is fine for LocalStack; use a real key/profile for real AWS"
  type        = string
  default     = "test"
  sensitive   = true
}

variable "container_image" {
  description = "Image to deploy, e.g. ghcr.io/arnel-rah/your-api:latest"
  type        = string
}

variable "app_port" {
  type    = number
  default = 8080
}

variable "enable_waf" {
  description = "Check LocalStack service coverage (wafv2) under your plan before enabling"
  type        = bool
  default     = true
}
