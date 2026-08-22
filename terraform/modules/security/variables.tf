variable "project_name" {
  type = string
}

variable "vpc_id" {
  description = "VPC ID from the network module"
  type        = string
}

variable "app_port" {
  description = "Port the application listens on"
  type        = number
  default     = 8080
}

variable "db_port" {
  description = "Database port"
  type        = number
  default     = 5432
}

variable "enable_waf" {
  description = "Whether to create the WAF Web ACL. Check LocalStack service coverage for wafv2 under your plan before enabling."
  type        = bool
  default     = true
}

variable "waf_rate_limit" {
  description = "Max requests per 5 minutes per IP before the WAF blocks it"
  type        = number
  default     = 2000
}

variable "tags" {
  type    = map(string)
  default = {}
}
