#!/bin/bash

# Test Runner Script for GoFiber Skeleton
# This script provides a comprehensive testing workflow

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to run tests with proper error handling
run_tests() {
    local test_name="$1"
    local test_command="$2"

    print_status "Running $test_name..."

    if eval "$test_command"; then
        print_success "$test_name passed"
        return 0
    else
        print_error "$test_name failed"
        return 1
    fi
}

# Main execution
main() {
    local test_type="${1:-all}"

    print_status "Starting GoFiber Skeleton Test Runner"
    print_status "Test type: $test_type"
    echo ""

    # Check prerequisites
    print_status "Checking prerequisites..."

    if ! command_exists go; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi

    if ! command_exists golangci-lint; then
        print_warning "golangci-lint not found, installing..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    fi

    print_success "Prerequisites check completed"
    echo ""

    # Ensure dependencies are up to date
    print_status "Updating dependencies..."
    go mod download
    go mod tidy

    # Generate mocks
    print_status "Generating mocks..."
    make mocks

    # Run tests based on type
    case "$test_type" in
        "unit")
            print_status "Running unit tests only..."
            run_tests "User Domain Unit Tests" "make test-unit-user"
            run_tests "Post Domain Unit Tests" "make test-unit-post"
            ;;

        "integration")
            print_status "Running integration tests with Docker..."
            print_warning "Note: Integration tests will use Docker containers for database"

            # Check if Docker is running
            if ! docker info >/dev/null 2>&1; then
                print_error "Docker is not running. Please start Docker and try again."
                exit 1
            fi

            # Run integration tests with Docker
            run_tests "Database Integration Tests (Docker)" "make test-integration-docker"
            ;;

        "lint")
            print_status "Running lint checks only..."
            run_tests "Code Formatting" "make fmt-check"
            run_tests "Vet" "make vet"
            run_tests "Linter" "make lint"
            ;;

        "coverage")
            print_status "Running tests with coverage..."
            run_tests "Test Coverage" "make test-coverage"
            print_status "Coverage report generated: coverage.html"
            ;;

        "docker-integration")
            print_status "Running Docker-based integration tests..."

            if ! docker info >/dev/null 2>&1; then
                print_error "Docker is not running. Please start Docker and try again."
                exit 1
            fi

            run_tests "Docker Integration Tests" "make test-integration-docker"
            ;;

        "all"|*)
            print_status "Running complete test suite..."

            # Code quality checks
            run_tests "Code Formatting" "make fmt-check"
            run_tests "Vet" "make vet"

            # Unit tests
            run_tests "User Domain Unit Tests" "make test-unit-user"
            run_tests "Post Domain Unit Tests" "make test-unit-post"

            # Coverage
            run_tests "Test Coverage" "make test-coverage"

            # Integration tests (Docker-based)
            if command_exists docker && docker info >/dev/null 2>&1; then
                print_status "Running Docker integration tests..."
                make test-integration-docker
            else
                print_warning "Docker not available, skipping integration tests"
                print_status "Run with './scripts/test-runner.sh integration' to run integration tests manually"
            fi
            ;;
    esac

    echo ""
    print_status "Test suite completed!"

    if [ "$test_type" = "coverage" ] || [ "$test_type" = "all" ]; then
        print_status "Coverage reports generated:"
        print_status "  - HTML: coverage.html"
        print_status "  - Text: coverage.txt (run 'make test-report')"
    fi
}

# Show usage
usage() {
    echo "Usage: $0 [test_type]"
    echo ""
    echo "Test types:"
    echo "  unit              - Run unit tests only"
    echo "  integration       - Run integration tests with Docker"
    echo "  docker-integration - Run Docker-based integration tests only"
    echo "  lint             - Run lint checks only"
    echo "  coverage         - Run tests with coverage report"
    echo "  all              - Run complete test suite (default)"
    echo ""
    echo "Examples:"
    echo "  $0           # Run all tests"
    echo "  $0 unit     # Run unit tests only"
    echo "  $0 coverage # Run tests with coverage"
    echo ""
    echo "Environment variables:"
    echo "  INTEGRATION_TESTS=1  - Enable integration tests"
}

# Parse command line arguments
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    usage
    exit 0
fi

# Run main function
main "$@"