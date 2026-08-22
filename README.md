# secure-deploy-kit

Infrastructure-as-code kit for deploying secure REST APIs on AWS —
built for developers and students in resource-constrained environments
who want production-grade practices without enterprise-scale budgets.

Runs entirely against a local AWS emulator ([LocalStack](https://www.localstack.cloud/))
by default — no AWS account, no credit card, no cloud bill required to
learn or demonstrate the whole stack. A guide for switching to a real
AWS account is included for when you're ready to go live.

## Why this exists

Getting a backend project onto real cloud infrastructure — with proper
networking, CI/CD, and security — is a steep learning curve, especially
without access to expensive courses or enterprise tooling. This kit
packages the patterns I learned building and securing several API
projects (Spring Boot, Go) into a reusable, documented baseline, and
proves that production-grade infrastructure practices don't require a
cloud budget to learn or demonstrate.

## What's inside

- **Terraform modules**: VPC, subnets, security groups, NAT Gateway
- **CI/CD pipeline** (GitHub Actions): build, test, deploy
- **Security layer**: OIDC authentication, WAF rules, a documented
  threat model (EBIOS-inspired)
- **Vulnerability testing scripts** (OWASP-aligned checks)
- **Local-first workflow**: everything runs on LocalStack by default;
  a documented path to real AWS is provided separately

## Who it's for

Students and junior developers who have a working API and want to
learn how to deploy it securely on AWS — without having to piece
together a dozen tutorials, and without paying for cloud infrastructure
just to learn how it works.

## Architecture

```
                         ┌─────────────────────┐
                         │   GitHub Actions     │
                         │   (CI/CD pipeline)   │
                         └──────────┬───────────┘
                                    │ deploy
                                    ▼
   Internet ──▶ [ WAF ] ──▶ [ ALB / API Gateway ] ──▶ [ ECS Fargate ]
                                                            │
                                                    ┌───────┴───────┐
                                                    │  Private VPC   │
                                                    │  subnet        │
                                                    └───────┬───────┘
                                                            ▼
                                                     [ RDS Postgres ]

   Auth: OIDC provider validates every request before it reaches
   the application layer.
```

See [`docs/architecture.md`](docs/architecture.md) for the detailed
reasoning behind each choice, and [`security/THREAT_MODEL.md`](security/THREAT_MODEL.md)
for the risk analysis.

## Quick start (LocalStack — free, no AWS account)

```bash
# 1. Clone and enter the repo
git clone https://github.com/Arnel-rah/secure-deploy-kit.git
cd secure-deploy-kit

# 2. Start LocalStack (requires Docker + a LocalStack Student/Ultimate license
#    via the GitHub Student Developer Pack, or the free Community tier for a
#    reduced service set)
scripts/bootstrap.sh

# 3. Deploy the infrastructure
cd terraform/environments/dev
tflocal init
tflocal apply

# 4. Tear everything down when you're done (no cost either way, but good habit)
scripts/destroy.sh
```

Full walkthrough: [`docs/getting-started.md`](docs/getting-started.md).

## Switching to real AWS

The same Terraform code works against real AWS — only the backend and
provider endpoints change. See the "Going to real AWS" section in
[`docs/getting-started.md`](docs/getting-started.md), and
[`docs/cost-estimate.md`](docs/cost-estimate.md) for what it would
actually cost.

## Security

- [`security/THREAT_MODEL.md`](security/THREAT_MODEL.md) — EBIOS-inspired
  risk analysis
- [`security/waf-rules/`](security/waf-rules) — baseline WAF rules
- [`security/scans/`](security/scans) — OWASP-aligned scan scripts and results
- [`auth/oidc/`](auth/oidc) — how authentication is wired in

## Roadmap

- [ ] Add Kubernetes (EKS) deployment target as an alternative to ECS
- [ ] Add automated cost-estimate diffing on PRs (Infracost)
- [ ] Publish a companion write-up on the threat model and lessons learned
- [ ] Package a one-command "starter project" generator

## License

MIT — see [`LICENSE`](LICENSE).
# secure-deploy-kit

Infrastructure-as-code kit for deploying secure REST APIs on AWS —
built for developers and students in resource-constrained environments
who want production-grade practices without enterprise-scale budgets.

Runs entirely against a local AWS emulator ([LocalStack](https://www.localstack.cloud/))
by default — no AWS account, no credit card, no cloud bill required to
learn or demonstrate the whole stack. A guide for switching to a real
AWS account is included for when you're ready to go live.

## Why this exists

Getting a backend project onto real cloud infrastructure — with proper
networking, CI/CD, and security — is a steep learning curve, especially
without access to expensive courses or enterprise tooling. This kit
packages the patterns I learned building and securing several API
projects (Spring Boot, Go) into a reusable, documented baseline, and
proves that production-grade infrastructure practices don't require a
cloud budget to learn or demonstrate.

## What's inside

- **Terraform modules**: VPC, subnets, security groups, NAT Gateway
- **CI/CD pipeline** (GitHub Actions): build, test, deploy
- **Security layer**: OIDC authentication, WAF rules, a documented
  threat model (EBIOS-inspired)
- **Vulnerability testing scripts** (OWASP-aligned checks)
- **Local-first workflow**: everything runs on LocalStack by default;
  a documented path to real AWS is provided separately

## Who it's for

Students and junior developers who have a working API and want to
learn how to deploy it securely on AWS — without having to piece
together a dozen tutorials, and without paying for cloud infrastructure
just to learn how it works.

## Architecture

```
                         ┌─────────────────────┐
                         │   GitHub Actions     │
                         │   (CI/CD pipeline)   │
                         └──────────┬───────────┘
                                    │ deploy
                                    ▼
   Internet ──▶ [ WAF ] ──▶ [ ALB / API Gateway ] ──▶ [ ECS Fargate ]
                                                            │
                                                    ┌───────┴───────┐
                                                    │  Private VPC   │
                                                    │  subnet        │
                                                    └───────┬───────┘
                                                            ▼
                                                     [ RDS Postgres ]

   Auth: OIDC provider validates every request before it reaches
   the application layer.
```

See [`docs/architecture.md`](docs/architecture.md) for the detailed
reasoning behind each choice, and [`security/THREAT_MODEL.md`](security/THREAT_MODEL.md)
for the risk analysis.

## Quick start (LocalStack — free, no AWS account)

```bash
# 1. Clone and enter the repo
git clone https://github.com/Arnel-rah/secure-deploy-kit.git
cd secure-deploy-kit

# 2. Start LocalStack (requires Docker + a LocalStack Student/Ultimate license
#    via the GitHub Student Developer Pack, or the free Community tier for a
#    reduced service set)
scripts/bootstrap.sh

# 3. Deploy the infrastructure
cd terraform/environments/dev
tflocal init
tflocal apply

# 4. Tear everything down when you're done (no cost either way, but good habit)
scripts/destroy.sh
```

Full walkthrough: [`docs/getting-started.md`](docs/getting-started.md).

## Switching to real AWS

The same Terraform code works against real AWS — only the backend and
provider endpoints change. See the "Going to real AWS" section in
[`docs/getting-started.md`](docs/getting-started.md), and
[`docs/cost-estimate.md`](docs/cost-estimate.md) for what it would
actually cost.

## Security

- [`security/THREAT_MODEL.md`](security/THREAT_MODEL.md) — EBIOS-inspired
  risk analysis
- [`security/waf-rules/`](security/waf-rules) — baseline WAF rules
- [`security/scans/`](security/scans) — OWASP-aligned scan scripts and results
- [`auth/oidc/`](auth/oidc) — how authentication is wired in

## Roadmap

- [ ] Add Kubernetes (EKS) deployment target as an alternative to ECS
- [ ] Add automated cost-estimate diffing on PRs (Infracost)
- [ ] Publish a companion write-up on the threat model and lessons learned
- [ ] Package a one-command "starter project" generator

## License

MIT — see [`LICENSE`](LICENSE).
