# Getting started

## Prerequisites

- Docker (with WSL2 integration if you're on Windows)
- Terraform >= 1.6
- Python + pip (for `tflocal`)
- A LocalStack account — free Community tier works for a reduced
  service set; the **Student plan** (via the GitHub Student Developer
  Pack) unlocks the full Ultimate-tier service set at no cost

## 1. Get LocalStack running

```bash
scripts/bootstrap.sh
```

If you have a Student plan license, export your token first so the
script picks it up:

```bash
export LOCALSTACK_AUTH_TOKEN="your-token-here"
scripts/bootstrap.sh
```

## 2. Configure your deployment

```bash
cd terraform/environments/dev
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` — at minimum, set `container_image` to point
at your own API's container image (e.g. published to GitHub Container
Registry via your CI pipeline).

## 3. Deploy

```bash
tflocal init
tflocal apply
```

`tflocal` is a thin wrapper that runs plain Terraform but automatically
redirects AWS API calls to your local LocalStack instance — you don't
need to hand-edit endpoints.

Once applied, the ALB DNS name is printed as an output:

```bash
terraform output alb_dns_name
```

(LocalStack's emulated ALB DNS name resolves within the LocalStack
network — see LocalStack's docs on accessing services for how to hit
it from your host machine.)

## 4. Tear it down

```bash
scripts/destroy.sh dev
```

## Checking what LocalStack supports

Not every AWS service/feature is emulated identically across LocalStack
tiers. Before relying on a specific resource (e.g. `aws_wafv2_web_acl`),
check the [LocalStack service coverage page](https://docs.localstack.cloud/references/coverage/)
for your plan. If a resource isn't supported, set the corresponding
`enable_*` variable to `false` (e.g. `enable_waf = false`) and note it
in your own deployment notes — the Terraform code degrades gracefully
when these flags are off.

## Going to real AWS

When you're ready to deploy against a real AWS account:

1. Set up AWS credentials (`aws configure`, or an IAM role in CI).
2. In `terraform/environments/prod/variables.tf`, `use_localstack`
   already defaults to `false`.
3. Uncomment the S3 backend block in
   `terraform/environments/prod/backend.tf` and fill in a real bucket
   + DynamoDB table for state locking.
4. Run plain `terraform` (not `tflocal`) from `terraform/environments/prod`.
5. Read [`docs/cost-estimate.md`](cost-estimate.md) first — know what
   you're about to be billed for.

## Running the security scans

Once the API is deployed and reachable:

```bash
security/scans/run-owasp-scan.sh http://<your-alb-dns-name>
```

See [`security/THREAT_MODEL.md`](../security/THREAT_MODEL.md) for the
reasoning behind what's being tested.
