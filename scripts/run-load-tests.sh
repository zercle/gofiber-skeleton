#!/bin/bash
# Load testing script for Go Fiber application

set -e

echo "=== Go Fiber Load Testing ==="
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if k6 is installed
if ! command -v k6 &> /dev/null; then
    echo -e "${YELLOW}k6 not found. Please install from https://k6.io/docs/getting-started/installation/${NC}"
    exit 1
fi

# Configuration
BASE_URL="${BASE_URL:-http://localhost:3000}"
TEST_TYPE="${1:-auth}"

echo "Target: $BASE_URL"
echo "Test Type: $TEST_TYPE"
echo ""

case $TEST_TYPE in
  auth)
    echo "Running authentication load test..."
    k6 run --out json=load-test-results.json tests/load/auth_load_test.js
    ;;

  spike)
    echo "Running spike test..."
    k6 run --out json=spike-test-results.json tests/load/spike_test.js
    ;;

  soak)
    echo "Running soak test (this will take ~4 hours)..."
    k6 run --out json=soak-test-results.json tests/load/soak_test.js
    ;;

  all)
    echo "Running all tests..."
    echo ""
    echo "1. Authentication load test"
    k6 run --out json=auth-test-results.json tests/load/auth_load_test.js
    echo ""
    echo "2. Spike test"
    k6 run --out json=spike-test-results.json tests/load/spike_test.js
    ;;

  *)
    echo "Unknown test type: $TEST_TYPE"
    echo "Usage: $0 [auth|spike|soak|all]"
    exit 1
    ;;
esac

echo ""
echo -e "${GREEN}Load test complete!${NC}"
echo "Results saved to *-test-results.json"
