locals {
  common_tags = {
    Project     = var.project_name
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

module "network" {
  source = "../../modules/network"

  project_name       = var.project_name
  enable_nat_gateway = true
  tags               = local.common_tags
}

module "security" {
  source = "../../modules/security"

  project_name = var.project_name
  vpc_id       = module.network.vpc_id
  app_port     = var.app_port
  enable_waf   = var.enable_waf
  tags         = local.common_tags
}

# IAM role ECS uses to pull the image and write logs.
# Kept minimal on purpose (no wildcard permissions).
resource "aws_iam_role" "ecs_execution" {
  name = "${var.project_name}-ecs-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
    }]
  })

  tags = local.common_tags
}

resource "aws_iam_role_policy_attachment" "ecs_execution" {
  role       = aws_iam_role.ecs_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

module "database" {
  source = "../../modules/database"

  project_name               = var.project_name
  private_subnet_ids         = module.network.private_subnet_ids
  database_security_group_id = module.security.database_security_group_id
  tags                       = local.common_tags
}

module "compute" {
  source = "../../modules/compute"

  project_name            = var.project_name
  aws_region              = var.aws_region
  vpc_id                  = module.network.vpc_id
  public_subnet_ids       = module.network.public_subnet_ids
  private_subnet_ids      = module.network.private_subnet_ids
  alb_security_group_id   = module.security.alb_security_group_id
  app_security_group_id   = module.security.app_security_group_id
  waf_web_acl_arn         = module.security.waf_web_acl_arn
  attach_waf              = var.enable_waf
  container_image         = var.container_image
  app_port                = var.app_port
  task_execution_role_arn = aws_iam_role.ecs_execution.arn
  task_role_arn           = module.security.app_task_role_arn

  container_env = {
    SPRING_PROFILES_ACTIVE = "prod"
    DB_HOST                = module.database.db_endpoint
    DB_NAME                = module.database.db_name
    OIDC_ISSUER_URL        = var.oidc_issuer_url
  }

  tags = local.common_tags
}
