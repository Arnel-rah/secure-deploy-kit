variable "project_name" {
  type = string
}

variable "private_subnet_ids" {
  type = list(string)
}

variable "database_security_group_id" {
  type = string
}

variable "engine_version" {
  type    = string
  default = "16.4"
}

variable "instance_class" {
  description = "Use a free-tier eligible class (db.t3.micro / db.t4g.micro) on real AWS"
  type        = string
  default     = "db.t3.micro"
}

variable "allocated_storage" {
  type    = number
  default = 20
}

variable "db_name" {
  type    = string
  default = "appdb"
}

variable "master_username" {
  type    = string
  default = "app_admin"
}

variable "master_password" {
  description = "Leave null to auto-generate a random password (recommended)"
  type        = string
  default     = null
  sensitive   = true
}

variable "multi_az" {
  description = "Keep false on free tier / LocalStack; set true for real production"
  type        = bool
  default     = false
}

variable "skip_final_snapshot" {
  type    = bool
  default = true
}

variable "deletion_protection" {
  type    = bool
  default = false
}

variable "backup_retention_period" {
  type    = number
  default = 1
}

variable "tags" {
  type    = map(string)
  default = {}
}
