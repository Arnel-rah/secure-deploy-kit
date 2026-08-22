#!/usr/bin/env bash
set -euo pipefail

TARGET_URL="${1:-}"

if [[ -z "$TARGET_URL" ]]; then
  echo "Usage: $0 <target-url>"
  echo "Example: $0 http://localhost:8080"
  exit 1
fi

RESULTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/results"
mkdir -p "$RESULTS_DIR"
TIMESTAMP="$(date +%Y%m%d-%H%M%S)"

echo "Running OWASP ZAP baseline scan against: $TARGET_URL"
echo "Results will be saved to: $RESULTS_DIR"

docker run --rm \
  --network host \
  -v "$RESULTS_DIR:/zap/wrk/:rw" \
  ghcr.io/zaproxy/zaproxy:stable \
  zap-baseline.py \
  -t "$TARGET_URL" \
  -r "report-${TIMESTAMP}.html" \
  -J "report-${TIMESTAMP}.json" \
  -I

echo ""
echo "Scan complete. Review:"
echo "  - $RESULTS_DIR/report-${TIMESTAMP}.html"
echo "  - $RESULTS_DIR/report-${TIMESTAMP}.json"
echo ""
echo "Note: the baseline scan is passive/light. For deeper testing"
echo "(active scan), see the ZAP docs — active scanning can be"
echo "disruptive and should only ever be run against your own"
echo "environment, never a third party's."
