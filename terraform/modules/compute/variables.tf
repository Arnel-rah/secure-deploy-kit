variable "project_name" {
  type = string
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "vpc_id" {
  type = string
}

variable "public_subnet_ids" {
  type = list(string)
}

variable "private_subnet_ids" {
  type = list(string)
}

variable "alb_security_group_id" {
  type = string
}

variable "app_security_group_id" {
  type = string
}

variable "waf_web_acl_arn" {
  description = "ARN of the WAF Web ACL to attach to the ALB, or null to skip"
  type        = string
  default     = null
}

variable "attach_waf" {
  description = "Whether to attach a WAF Web ACL to the ALB. Kept separate from waf_web_acl_arn because that ARN is only known after apply — using it directly in a count expression makes the plan non-deterministic."
  type        = bool
  default     = false
}

variable "container_name" {
  type    = string
  default = "app"
}

variable "container_image" {
  description = "Container image to deploy, e.g. ghcr.io/you/your-api:latest"
  type        = string
}

variable "container_env" {
  description = "Environment variables passed to the container"
  type        = map(string)
  default     = {}
}

variable "app_port" {
  type    = number
  default = 8080
}

variable "health_check_path" {
  type    = string
  default = "/actuator/health"
}

variable "task_cpu" {
  description = "Fargate task CPU units (256 = 0.25 vCPU)"
  type        = string
  default     = "256"
}

variable "task_memory" {
  description = "Fargate task memory in MB"
  type        = string
  default     = "512"
}

variable "desired_count" {
  type    = number
  default = 1
}

variable "task_execution_role_arn" {
  description = "IAM role ECS uses to pull images and write logs"
  type        = string
}

variable "task_role_arn" {
  description = "IAM role the running container assumes (from the security module)"
  type        = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
