#!/usr/bin/env bash


set -euo pipefail

echo "== secure-deploy-kit bootstrap =="

if ! command -v docker &> /dev/null; then
  echo "Docker is required but not found. Install Docker Desktop (with WSL2 integration on Windows) first."
  exit 1
fi

if ! command -v tflocal &> /dev/null; then
  echo "tflocal not found — installing via pip..."
  if [[ -n "${VIRTUAL_ENV:-}" ]]; then
    echo "Detected active virtualenv ($VIRTUAL_ENV) — installing without --user."
    pip install terraform-local
  else
    pip install --user terraform-local
  fi
fi

if docker ps --format '{{.Names}}' | grep -q '^localstack$'; then
  echo "LocalStack is already running."
else
  echo "Starting LocalStack..."
  echo "If you have a LocalStack Student/Ultimate license (via the GitHub"
  echo "Student Developer Pack), set LOCALSTACK_AUTH_TOKEN in your shell"
  echo "before running this script to unlock the full service set."
  docker run -d --name localstack \
    -p 4566:4566 \
    -e LOCALSTACK_AUTH_TOKEN="${LOCALSTACK_AUTH_TOKEN:-}" \
    -v /var/run/docker.sock:/var/run/docker.sock \
    localstack/localstack
fi

echo "Waiting for LocalStack to be ready..."
for i in {1..30}; do
  if curl -sf http://localhost:4566/_localstack/health > /dev/null; then
    echo "LocalStack is up."
    break
  fi
  sleep 2
done

echo ""
echo "Ready. Next steps:"
echo "  cd terraform/environments/dev"
echo "  cp terraform.tfvars.example terraform.tfvars   # edit as needed"
echo "  tflocal init"
echo "  tflocal apply"
