#!/usr/bin/env bash


set -euo pipefail

ENVIRONMENT="${1:-dev}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_DIR="$ROOT_DIR/terraform/environments/$ENVIRONMENT"

if [[ ! -d "$ENV_DIR" ]]; then
  echo "Unknown environment: $ENVIRONMENT (expected 'dev' or 'prod')"
  exit 1
fi

echo "Destroying infrastructure in environment: $ENVIRONMENT"
cd "$ENV_DIR"

if [[ "$ENVIRONMENT" == "prod" ]]; then
  read -r -p "This will destroy PROD infrastructure. Type 'destroy-prod' to confirm: " CONFIRM
  if [[ "$CONFIRM" != "destroy-prod" ]]; then
    echo "Aborted."
    exit 1
  fi
  terraform destroy
else
  command -v tflocal &> /dev/null && tflocal destroy -auto-approve || terraform destroy -auto-approve
fi

if [[ "$ENVIRONMENT" == "dev" ]] && docker ps --format '{{.Names}}' | grep -q '^localstack$'; then
  read -r -p "Stop the LocalStack container too? [y/N] " STOP_LS
  if [[ "$STOP_LS" =~ ^[Yy]$ ]]; then
    docker stop localstack && docker rm localstack
    echo "LocalStack stopped."
  fi
fi

echo "Done."
