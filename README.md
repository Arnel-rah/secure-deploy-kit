# secure-deploy-kit

Infrastructure-as-code kit for deploying secure REST APIs on AWS —
built for developers and students in resource-constrained environments
who want production-grade practices without enterprise-scale budgets.

## Why this exists

Getting a backend project onto real cloud infrastructure — with proper
networking, CI/CD, and security — is a steep learning curve, especially
without access to expensive courses or enterprise tooling. This kit
packages the patterns I learned building and securing several API
projects (Spring Boot/Poja, Go) into a reusable, documented baseline.

## What's inside

- Terraform modules: VPC, subnets, security groups, NAT Gateway
- CI/CD pipeline (GitHub Actions): build, test, deploy
- Security layer: OIDC authentication, AWS WAF, documented threat model (EBIOS-inspired)
- Vulnerability testing scripts (OWASP-aligned checks)
- Setup guide for AWS free-tier deployment

## Who it's for

Students and junior developers who have a working API and want to
learn how to deploy it securely on AWS — without having to piece
together a dozen tutorials.

## Getting started
...

## Threat model
...

## Roadmap
...
