#!/bin/bash
# Security scanning script for Go Fiber application

set -e

echo "=== Go Fiber Security Scan ==="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if gosec is installed
if ! command -v gosec &> /dev/null; then
    echo -e "${YELLOW}gosec not found. Installing...${NC}"
    go install github.com/securego/gosec/v2/cmd/gosec@latest
fi

# Check if trivy is installed
if ! command -v trivy &> /dev/null; then
    echo -e "${YELLOW}trivy not found. Please install it from https://aquasecurity.github.io/trivy/${NC}"
fi

echo "1. Running gosec security scanner..."
echo "-----------------------------------"
gosec -fmt=json -out=gosec-report.json ./... || true
gosec -fmt=text ./... || true
echo ""

echo "2. Checking for dependency vulnerabilities..."
echo "----------------------------------------------"
if command -v trivy &> /dev/null; then
    trivy fs --scanners vuln --severity HIGH,CRITICAL . || true
else
    echo -e "${YELLOW}Trivy not installed. Skipping dependency scan.${NC}"
fi
echo ""

echo "3. Checking Go module vulnerabilities..."
echo "-----------------------------------------"
if command -v govulncheck &> /dev/null; then
    govulncheck ./... || true
else
    echo -e "${YELLOW}govulncheck not installed. Installing...${NC}"
    go install golang.org/x/vuln/cmd/govulncheck@latest
    govulncheck ./... || true
fi
echo ""

echo "4. Running go vet..."
echo "--------------------"
go vet ./... || true
echo ""

echo "5. Checking for common issues..."
echo "---------------------------------"
# Check for hardcoded credentials
echo "Checking for hardcoded credentials..."
grep -r -n -i "password.*=.*\"" --include="*.go" --exclude-dir={vendor,node_modules,.git} . || echo "✓ No hardcoded passwords found"

# Check for SQL injection vulnerabilities
echo "Checking for potential SQL injection..."
grep -r -n "Query.*fmt.Sprintf" --include="*.go" --exclude-dir={vendor,node_modules,.git} . || echo "✓ No obvious SQL injection patterns found"

# Check for command injection
echo "Checking for potential command injection..."
grep -r -n "exec.Command" --include="*.go" --exclude-dir={vendor,node_modules,.git} . || echo "✓ No command execution found"

echo ""
echo "=== Security Scan Complete ==="
echo ""
echo -e "${GREEN}Reports generated:${NC}"
echo "  - gosec-report.json"
echo "  - Check console output above for details"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Review gosec-report.json for detailed findings"
echo "  2. Fix HIGH and CRITICAL severity issues"
echo "  3. Consider fixing MEDIUM severity issues"
echo "  4. Run 'make lint' for additional code quality checks"
