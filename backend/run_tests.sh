#!/bin/bash

# Test runner script for Freshease backend
# This script runs all unit tests, integration tests, and generates coverage reports

set -e

echo "🧪 Starting Freshease Backend Test Suite"
echo "========================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
COVERAGE_DIR="coverage"
COVERAGE_FILE="$COVERAGE_DIR/coverage.out"
COVERAGE_HTML="$COVERAGE_DIR/coverage.html"

# Create coverage directory
mkdir -p $COVERAGE_DIR

echo -e "${BLUE}📁 Created coverage directory: $COVERAGE_DIR${NC}"

# Function to run tests with coverage
run_tests() {
    local test_type=$1
    local test_path=$2
    local test_name=$3
    
    echo -e "\n${YELLOW}🔍 Running $test_name...${NC}"
    
    if [ "$test_type" = "unit" ]; then
        go test -v -race -coverprofile="$COVERAGE_DIR/${test_name}_coverage.out" -covermode=atomic $test_path
    elif [ "$test_type" = "integration" ]; then
        go test -v -race -tags=integration -coverprofile="$COVERAGE_DIR/${test_name}_coverage.out" -covermode=atomic $test_path
    fi
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ $test_name passed${NC}"
    else
        echo -e "${RED}❌ $test_name failed${NC}"
        exit 1
    fi
}

# Function to merge coverage files
merge_coverage() {
    echo -e "\n${BLUE}📊 Merging coverage reports...${NC}"
    
    echo "mode: atomic" > $COVERAGE_FILE
    
    for file in $COVERAGE_DIR/*_coverage.out; do
        if [ -f "$file" ]; then
            tail -n +2 "$file" >> $COVERAGE_FILE
        fi
    done
    
    echo -e "${GREEN}✅ Coverage reports merged${NC}"
}

# Function to generate HTML coverage report
generate_html_coverage() {
    echo -e "\n${BLUE}📈 Generating HTML coverage report...${NC}"
    
    go tool cover -html=$COVERAGE_FILE -o $COVERAGE_HTML
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ HTML coverage report generated: $COVERAGE_HTML${NC}"
    else
        echo -e "${RED}❌ Failed to generate HTML coverage report${NC}"
        exit 1
    fi
}

# Function to show coverage summary
show_coverage_summary() {
    echo -e "\n${BLUE}📊 Coverage Summary:${NC}"
    go tool cover -func=$COVERAGE_FILE
}

# Function to run all tests
run_all_tests() {
    echo -e "\n${YELLOW}🚀 Running all tests...${NC}"
    
    # Run unit tests
    run_tests "unit" "./modules/users/..." "Users Unit Tests"
    run_tests "unit" "./modules/products/..." "Products Unit Tests"
    run_tests "unit" "./internal/common/middleware/..." "Middleware Unit Tests"
    
    # Run integration tests
    run_tests "integration" "./internal/common/testutils/..." "Integration Tests"
    
    # Run any other test packages
    run_tests "unit" "./internal/common/..." "Common Package Tests"
}

# Function to run specific test suite
run_specific_tests() {
    local suite=$1
    
    case $suite in
        "users")
            run_tests "unit" "./modules/users/..." "Users Tests"
            ;;
        "products")
            run_tests "unit" "./modules/products/..." "Products Tests"
            ;;
        "middleware")
            run_tests "unit" "./internal/common/middleware/..." "Middleware Tests"
            ;;
        "integration")
            run_tests "integration" "./internal/common/testutils/..." "Integration Tests"
            ;;
        *)
            echo -e "${RED}❌ Unknown test suite: $suite${NC}"
            echo "Available suites: users, products, middleware, integration"
            exit 1
            ;;
    esac
}

# Function to run tests with race detection
run_race_tests() {
    echo -e "\n${YELLOW}🏃 Running tests with race detection...${NC}"
    
    go test -v -race ./...
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Race detection tests passed${NC}"
    else
        echo -e "${RED}❌ Race detection tests failed${NC}"
        exit 1
    fi
}

# Function to run benchmarks
run_benchmarks() {
    echo -e "\n${YELLOW}⚡ Running benchmarks...${NC}"
    
    go test -v -bench=. -benchmem ./...
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Benchmarks completed${NC}"
    else
        echo -e "${RED}❌ Benchmarks failed${NC}"
        exit 1
    fi
}

# Function to clean up test artifacts
cleanup() {
    echo -e "\n${BLUE}🧹 Cleaning up test artifacts...${NC}"
    
    # Remove test coverage files (keep HTML report)
    find $COVERAGE_DIR -name "*_coverage.out" -delete
    
    echo -e "${GREEN}✅ Cleanup completed${NC}"
}

# Function to show help
show_help() {
    echo "Freshease Backend Test Runner"
    echo "Usage: $0 [OPTIONS] [SUITE]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -a, --all           Run all tests (default)"
    echo "  -r, --race          Run tests with race detection"
    echo "  -b, --bench         Run benchmarks"
    echo "  -c, --coverage      Generate coverage report"
    echo "  -s, --suite SUITE   Run specific test suite"
    echo "  --clean             Clean up test artifacts"
    echo ""
    echo "Test Suites:"
    echo "  users               Users module tests"
    echo "  products            Products module tests"
    echo "  middleware          Middleware tests"
    echo "  integration         Integration tests"
    echo ""
    echo "Examples:"
    echo "  $0                  # Run all tests"
    echo "  $0 -s users         # Run only users tests"
    echo "  $0 -r               # Run tests with race detection"
    echo "  $0 -c               # Generate coverage report"
    echo "  $0 -b               # Run benchmarks"
}

# Main execution
main() {
    local run_all=true
    local run_race=false
    local run_bench=false
    local generate_coverage=false
    local suite=""
    local clean=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -a|--all)
                run_all=true
                shift
                ;;
            -r|--race)
                run_race=true
                run_all=false
                shift
                ;;
            -b|--bench)
                run_bench=true
                run_all=false
                shift
                ;;
            -c|--coverage)
                generate_coverage=true
                run_all=false
                shift
                ;;
            -s|--suite)
                suite="$2"
                run_all=false
                shift 2
                ;;
            --clean)
                clean=true
                run_all=false
                shift
                ;;
            *)
                echo -e "${RED}❌ Unknown option: $1${NC}"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Execute based on options
    if [ "$clean" = true ]; then
        cleanup
        exit 0
    fi
    
    if [ "$run_race" = true ]; then
        run_race_tests
    elif [ "$run_bench" = true ]; then
        run_benchmarks
    elif [ "$generate_coverage" = true ]; then
        run_all_tests
        merge_coverage
        generate_html_coverage
        show_coverage_summary
    elif [ -n "$suite" ]; then
        run_specific_tests "$suite"
    elif [ "$run_all" = true ]; then
        run_all_tests
        merge_coverage
        generate_html_coverage
        show_coverage_summary
    fi
    
    echo -e "\n${GREEN}🎉 Test suite completed successfully!${NC}"
    echo -e "${BLUE}📊 Coverage report available at: $COVERAGE_HTML${NC}"
}

# Run main function with all arguments
main "$@"
