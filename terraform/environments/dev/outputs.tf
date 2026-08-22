output "alb_dns_name" {
  description = "URL to hit the deployed API through"
  value       = module.compute.alb_dns_name
}

output "db_endpoint" {
  value = module.database.db_endpoint
}

output "vpc_id" {
  value = module.network.vpc_id
}
