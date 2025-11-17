#!/bin/bash

# Master test runner for comprehensive test coverage
# This script runs all tests and generates coverage reports

set -e

echo "🧪 Running Comprehensive Test Suite"
echo "===================================="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

COVERAGE_DIR="coverage"
COVERAGE_FILE="$COVERAGE_DIR/coverage.out"
COVERAGE_HTML="$COVERAGE_DIR/coverage.html"

mkdir -p $COVERAGE_DIR

# Function to run tests for a module
run_module_tests() {
    local module=$1
    local module_path="./modules/$module"
    
    if [ ! -d "$module_path" ]; then
        echo -e "${YELLOW}⚠️  Module $module not found, skipping...${NC}"
        return 1
    fi
    
    echo -e "\n${BLUE}📦 Testing module: $module${NC}"
    
    # Run tests with coverage
    go test -v -race -coverprofile="$COVERAGE_DIR/${module}_coverage.out" -covermode=atomic $module_path/... 2>&1 | tee "$COVERAGE_DIR/${module}_test.log"
    
    local exit_code=${PIPESTATUS[0]}
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✅ $module tests passed${NC}"
        return 0
    else
        echo -e "${RED}❌ $module tests failed${NC}"
        return $exit_code
    fi
}

# Function to run integration tests
run_integration_tests() {
    echo -e "\n${BLUE}🔗 Running integration tests...${NC}"
    
    go test -v -race -tags=integration -coverprofile="$COVERAGE_DIR/integration_coverage.out" -covermode=atomic ./internal/common/testutils/... 2>&1 | tee "$COVERAGE_DIR/integration_test.log"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Integration tests passed${NC}"
        return 0
    else
        echo -e "${RED}❌ Integration tests failed${NC}"
        return 1
    fi
}

# Function to merge coverage reports
merge_coverage() {
    echo -e "\n${BLUE}📊 Merging coverage reports...${NC}"
    
    echo "mode: atomic" > $COVERAGE_FILE
    
    for file in $COVERAGE_DIR/*_coverage.out; do
        if [ -f "$file" ]; then
            tail -n +2 "$file" >> $COVERAGE_FILE 2>/dev/null || true
        fi
    done
    
    echo -e "${GREEN}✅ Coverage reports merged${NC}"
}

# Function to generate HTML coverage report
generate_html_coverage() {
    echo -e "\n${BLUE}📈 Generating HTML coverage report...${NC}"
    
    if [ -f "$COVERAGE_FILE" ]; then
        go tool cover -html=$COVERAGE_FILE -o $COVERAGE_HTML
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✅ HTML coverage report generated: $COVERAGE_HTML${NC}"
        else
            echo -e "${RED}❌ Failed to generate HTML coverage report${NC}"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠️  No coverage file found, skipping HTML generation${NC}"
    fi
}

# Function to show coverage summary
show_coverage_summary() {
    echo -e "\n${BLUE}📊 Coverage Summary:${NC}"
    if [ -f "$COVERAGE_FILE" ]; then
        go tool cover -func=$COVERAGE_FILE | tail -1
        echo -e "\n${BLUE}Detailed coverage:${NC}"
        go tool cover -func=$COVERAGE_FILE
    else
        echo -e "${YELLOW}⚠️  No coverage file found${NC}"
    fi
}

# Main execution
main() {
    local failed_modules=()
    local total_modules=0
    local passed_modules=0
    
    # List of all modules to test
    MODULES=(
        "addresses"
        "bundles"
        "bundle_items"
        "cart_items"
        "carts"
        "categories"
        "deliveries"
        "genai"
        "inventories"
        "meal_plan_items"
        "meal_plans"
        "notifications"
        "order_items"
        "orders"
        "payments"
        "permissions"
        "product_categories"
        "products"
        "recipe_items"
        "recipes"
        "reviews"
        "roles"
        "shop"
        "uploads"
        "users"
        "vendors"
    )
    
    # Run tests for each module
    for module in "${MODULES[@]}"; do
        total_modules=$((total_modules + 1))
        if run_module_tests "$module"; then
            passed_modules=$((passed_modules + 1))
        else
            failed_modules+=("$module")
        fi
    done
    
    # Run integration tests
    run_integration_tests
    
    # Merge coverage reports
    merge_coverage
    
    # Generate HTML coverage report
    generate_html_coverage
    
    # Show coverage summary
    show_coverage_summary
    
    # Print summary
    echo -e "\n${BLUE}====================================${NC}"
    echo -e "${BLUE}Test Summary${NC}"
    echo -e "${BLUE}====================================${NC}"
    echo -e "Total modules: $total_modules"
    echo -e "${GREEN}Passed: $passed_modules${NC}"
    echo -e "${RED}Failed: ${#failed_modules[@]}${NC}"
    
    if [ ${#failed_modules[@]} -gt 0 ]; then
        echo -e "\n${RED}Failed modules:${NC}"
        for module in "${failed_modules[@]}"; do
            echo -e "  - $module"
        done
        echo -e "\n${YELLOW}⚠️  Some tests failed. Check logs in $COVERAGE_DIR/${NC}"
        exit 1
    else
        echo -e "\n${GREEN}🎉 All tests passed!${NC}"
        echo -e "${BLUE}📊 Coverage report: $COVERAGE_HTML${NC}"
        exit 0
    fi
}

main "$@"

