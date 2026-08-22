# This file documents the WAF rule reasoning. The actual Terraform
# resources live in terraform/modules/security/main.tf (aws_wafv2_web_acl) —
# kept there so the security module stays the single source of truth
# and doesn't get applied twice.
#
# Rule summary (see THREAT_MODEL.md → R3, R5 for the risks these address):
#
# 1. AWSManagedRulesCommonRuleSet
#    Baseline protection against a broad set of common exploits:
#    cross-site scripting, local file inclusion, oversized request bodies.
#
# 2. AWSManagedRulesSQLiRuleSet
#    Targeted SQL injection detection, layered on top of parameterized
#    queries in the application (this is defense-in-depth, not a
#    replacement for safe query building).
#
# 3. rate-limit-per-ip (custom rate-based rule)
#    Blocks any single IP exceeding `waf_rate_limit` requests per 5
#    minutes — mitigates basic denial-of-service and credential-
#    stuffing patterns.
#
# To extend: add further managed rule groups (e.g.
# AWSManagedRulesKnownBadInputsRuleSet, AWSManagedRulesLinuxRuleSet if
# the backend runs on Linux-specific paths) directly in
# terraform/modules/security/main.tf, following the same `rule` block
# pattern, with an incrementing `priority`.
