# Threat model

This is an EBIOS-inspired risk analysis of the infrastructure this kit
deploys. EBIOS structures the analysis around: what needs protecting,
who might attack it, how, and what mitigates each risk. This isn't a
full formal EBIOS Risk Manager study — it's a lightweight, practical
adaptation sized for a small API deployment.

## 1. Assets (what needs protecting)

| Asset | Why it matters |
|---|---|
| API application data (in RDS) | User/business data — confidentiality and integrity are critical |
| Authentication tokens / sessions | Compromise = full account takeover |
| Infrastructure credentials (IAM roles, DB password) | Compromise = full environment takeover |
| Service availability | Downtime has direct cost/reputation impact |
| Source code and Terraform state | Contains architecture and potentially sensitive config |

## 2. Threat sources (who might attack, and how)

| Source | Motivation | Likely method |
|---|---|---|
| Opportunistic scanners / bots | Automated, no specific target | Port scanning, credential stuffing, known CVEs |
| Script kiddies | Low-skill, high-volume attempts | SQL injection, XSS on unpatched endpoints |
| Targeted attacker | Specific interest in the data or service | Reconnaissance + tailored exploit chain |
| Insider (accidental) | No malicious intent | Misconfigured security group, leaked credential in a commit |

For a student/portfolio-scale project, the realistic threat profile is
dominated by the first two rows — automated and low-skill attackers —
rather than sophisticated targeted actors. The design below is scoped
accordingly, while still following patterns that scale to more serious
threat models.

## 3. Risk scenarios and mitigations

### R1 — Direct database exposure
**Scenario**: the RDS instance is reachable from the internet due to a
misconfigured security group or `publicly_accessible = true`.
**Impact**: high (full data exposure).
**Mitigation**: the `database` module hardcodes `publicly_accessible = false`
and only opens the DB security group to the app security group
(`terraform/modules/security/main.tf`). No public IP is ever assigned.

### R2 — Unauthenticated API access
**Scenario**: an endpoint is deployed without authentication, exposing
data or state-changing operations.
**Impact**: high.
**Mitigation**: OIDC resource-server validation at the application layer
(see `auth/oidc/`), enforced before any business logic runs. Health
check endpoints are the only explicitly public route.

### R3 — Common web exploits (SQLi, XSS, path traversal)
**Scenario**: an attacker probes the API for classic OWASP Top 10
vulnerabilities.
**Impact**: medium to high depending on the flaw.
**Mitigation**: WAF managed rule groups (`AWSManagedRulesCommonRuleSet`,
`AWSManagedRulesSQLiRuleSet`) filter traffic before it reaches the
application. This is a defense-in-depth layer, not a substitute for
secure coding — see `security/scans/` for the scan scripts used to
validate the application layer itself.

### R4 — Credential leakage via source control
**Scenario**: a database password, API key, or `.tfvars` file with
secrets gets committed to the repository.
**Impact**: high.
**Mitigation**: `.gitignore` excludes `*.tfvars` and `.env*` by default;
the database module auto-generates a random password via the `random`
provider when none is supplied, so there's nothing to accidentally
commit in the common case.

### R5 — Denial of service / resource exhaustion
**Scenario**: a flood of requests degrades or takes down the service.
**Impact**: medium (availability).
**Mitigation**: WAF rate-based rule blocking any single IP exceeding
the configured threshold (`waf_rate_limit`, default 2000 req / 5 min).

### R6 — Lateral movement after a single-container compromise
**Scenario**: the application container is compromised (e.g. via a
dependency vulnerability); the attacker tries to pivot further.
**Impact**: high if unmitigated.
**Mitigation**: the ECS task role (`aws_iam_role.app_task`) is scoped
to only what the application needs — no wildcard IAM permissions.
Private subnet placement with no public IP limits direct inbound
reachability of the container itself.

### R7 — Terraform state exposure
**Scenario**: Terraform state (which can contain sensitive values,
e.g. the generated DB password) leaks.
**Impact**: high.
**Mitigation**: local state is gitignored for solo/dev use; the prod
backend documents migration to an encrypted S3 backend with restricted
access when the project graduates beyond solo use.

## 4. What's explicitly out of scope (for now)

- Formal penetration testing by a third party
- Multi-region / disaster-recovery design
- Compliance-grade logging and audit trail retention
- Protection against a well-resourced, targeted nation-state-level actor

These are noted honestly rather than glossed over — a threat model
that pretends to cover everything is less useful than one that's
explicit about its boundaries.

## 5. Validation

The mitigations above are checked, not just asserted, via:
- `security/scans/run-owasp-scan.sh` — automated OWASP ZAP baseline scan
- Checkov IaC scanning in CI (`.github/workflows/ci.yml`)
- Manual review of security group rules after every `terraform plan`
