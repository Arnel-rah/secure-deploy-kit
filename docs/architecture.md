# Architecture

## Design goals

1. **Least exposure by default** — nothing is publicly reachable except
   the load balancer's HTTP(S) listeners.
2. **Modularity** — network, security, compute, and database are
   independent Terraform modules, each testable and reusable on its own.
3. **Zero cost to learn** — the whole stack runs on LocalStack; nothing
   requires a paid AWS account to build, test, or demonstrate.
4. **Security as a layer, not an afterthought** — WAF and OIDC are
   first-class parts of the stack, not bolted on later.

## Layers

### Network (`terraform/modules/network`)
A VPC with public subnets (hosting the load balancer and NAT Gateway)
and private subnets (hosting the application containers and database).
Only the public subnets have a route to the internet gateway; private
subnets reach the internet outbound only, via the NAT Gateway, and
accept no unsolicited inbound traffic.

**Why split public/private instead of one flat subnet?**
It's the standard defense-in-depth pattern: even if the application or
database has a misconfiguration, it isn't directly reachable from the
internet — the load balancer is the only public entry point, and it
only forwards to the app tier, never the database tier.

### Security (`terraform/modules/security`)
Three security groups (ALB, app, database) chained so each tier only
accepts traffic from the tier in front of it — never directly from the
internet except at the ALB. A WAF Web ACL adds request-level filtering
in front of the ALB. See [`THREAT_MODEL.md`](../security/THREAT_MODEL.md)
for the reasoning behind each rule.

### Compute (`terraform/modules/compute`)
ECS Fargate was chosen over EC2 or Lambda for this kit because:
- No server patching/maintenance (unlike EC2)
- No cold-start or execution-time constraints that complicate a
  long-running Spring Boot API (unlike Lambda)
- Straightforward mapping to a container image built by any CI pipeline

The roadmap includes an EKS-based alternative for teams that need
Kubernetes-specific tooling.

### Database (`terraform/modules/database`)
RDS Postgres, private-subnet only, encrypted at rest, with an
auto-generated password by default (via the `random` provider) so
there's no secret to accidentally commit. `multi_az` and
`deletion_protection` are off by default in `dev` (cost/simplicity) and
on by default in `prod`.

### Auth (`auth/oidc`)
Authentication is delegated to an OIDC provider rather than
implemented as custom JWT logic in the application. See
[`auth/oidc/README.md`](../auth/oidc/README.md) for the reasoning and
setup.

## Why Terraform over Ansible/Pulumi/CDK

Terraform was chosen for:
- Wide LocalStack support via `tflocal`, which is central to the
  "zero-cost-to-learn" goal
- Declarative state tracking, useful for a project revisited across a
  full semester
- The most widely adopted IaC tool in the current job market, so the
  skills transfer directly

## Known limitations

- LocalStack's service emulation isn't 100% identical to real AWS in
  every edge case (particularly for newer or less common services) —
  always sanity-check any final deployment against real AWS before
  treating it as production-ready.
- This kit optimizes for learning and portfolio use, not for
  high-availability production traffic. `docs/cost-estimate.md` and
  the `prod` environment sketch what scaling up would require.
