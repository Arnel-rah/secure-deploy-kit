# Cost estimate

## What this project actually costs, as built

**$0.** Every resource in this kit runs against LocalStack, which
emulates AWS locally — there is no cloud bill for building, testing,
running the CI pipeline, or demonstrating the project, as long as it
stays on LocalStack. This section exists to show a realistic
understanding of what running it for real would cost, without ever
having had to pay for it.

## Estimated monthly cost on real AWS (dev-scale, us-east-1)

Prices below are approximate, based on published AWS on-demand pricing
at time of writing — always check the [AWS Pricing Calculator](https://calculator.aws/)
for current numbers before committing to a real deployment.

| Resource | Configuration | Approx. monthly cost |
|---|---|---|
| NAT Gateway | 1 gateway, low traffic | ~$32 + data processing |
| Application Load Balancer | 1 ALB, low traffic | ~$16 + LCU usage |
| ECS Fargate | 1 task, 0.25 vCPU / 512 MB, always on | ~$9 |
| RDS Postgres | db.t3.micro, single-AZ, 20 GB | Free tier eligible for 12 months, then ~$13 |
| WAF | 1 Web ACL + 3 rules, low traffic | ~$6 + request volume |
| Data transfer | Depends heavily on traffic | Variable |
| **Total (after free tier)** | | **~$75–90/month** |

**During the AWS free tier (first 12 months on a new account):** RDS
and some ECS/Fargate usage can be substantially or fully offset,
bringing the realistic cost closer to **$50–60/month**, dominated by
the NAT Gateway and ALB, which are not free-tier eligible.

## Where the cost actually comes from

The NAT Gateway and ALB are the two line items that don't have a
meaningful free tier and that dominate the bill even at near-zero
traffic — this is a common surprise for people budgeting their first
real AWS deployment. Two honest mitigations, with trade-offs:

- **Skip the NAT Gateway** by giving the app tier a public IP directly
  (loses the "private subnet, no direct exposure" property from the
  threat model — not recommended without adjusting the security
  posture accordingly).
- **Skip the ALB** by exposing the ECS service through a Fargate
  public IP + security group only (loses WAF integration, load
  balancing, and TLS termination conveniences).

Neither shortcut is used in this kit's default configuration — the
point of `dev`/`prod` being separate environments is that you can
make different cost/security trade-offs deliberately per environment,
rather than compromising the design to save money everywhere.

## Practical takeaway

This is exactly why the kit defaults to LocalStack: it lets you build,
break, and rebuild this same architecture indefinitely for free, and
only incur real cost once — deliberately — when you actually need a
publicly reachable, real-AWS deployment (e.g. for a live demo).
