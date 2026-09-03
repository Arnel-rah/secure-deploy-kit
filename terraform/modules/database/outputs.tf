output "db_endpoint" {
  value = aws_db_instance.this.endpoint
}

output "db_name" {
  value = aws_db_instance.this.db_name
}

output "db_master_password" {
  value     = coalesce(var.master_password, try(random_password.master[0].result, null))
  sensitive = true
}

output "db_username" {
  value = aws_db_instance.this.username
}

output "db_password_secret_arn" {
  value = aws_secretsmanager_secret.db_password.arn
}
