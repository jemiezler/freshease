#!/bin/bash

# Test runner script for Freshease Frontend
# This script runs all unit tests, widget tests, and integration tests

set -e

echo "🧪 Starting Freshease Frontend Test Suite"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Generate mocks
echo -e "\n${BLUE}🔧 Generating mocks...${NC}"
flutter pub run build_runner build --delete-conflicting-outputs

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Mock generation failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Mocks generated successfully${NC}"

# Function to run tests
run_tests() {
    local test_type=$1
    local test_path=$2
    local test_name=$3
    
    echo -e "\n${YELLOW}🔍 Running $test_name...${NC}"
    
    flutter test $test_path
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ $test_name passed${NC}"
    else
        echo -e "${RED}❌ $test_name failed${NC}"
        return 1
    fi
}

# Function to run all tests
run_all_tests() {
    echo -e "\n${YELLOW}🚀 Running all tests...${NC}"
    
    # Run unit tests
    run_tests "unit" "test/unit/" "Unit Tests"
    
    # Run widget tests
    run_tests "widget" "test/widgets/" "Widget Tests"
    
    # Run integration tests (commented out as they require a device/emulator)
    # run_tests "integration" "integration_test/" "Integration Tests"
}

# Function to run tests with coverage
run_tests_with_coverage() {
    echo -e "\n${YELLOW}📊 Running tests with coverage...${NC}"
    
    flutter test --coverage
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Coverage report generated${NC}"
        echo -e "${BLUE}📊 Coverage file: coverage/lcov.info${NC}"
    else
        echo -e "${RED}❌ Coverage generation failed${NC}"
        exit 1
    fi
}

# Function to show help
show_help() {
    echo "Freshease Frontend Test Runner"
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -a, --all           Run all tests (default)"
    echo "  -u, --unit          Run only unit tests"
    echo "  -w, --widget        Run only widget tests"
    echo "  -i, --integration   Run only integration tests"
    echo "  -c, --coverage      Run tests with coverage"
    echo "  -m, --mocks         Generate mocks only"
    echo ""
    echo "Examples:"
    echo "  $0                  # Run all tests"
    echo "  $0 -u               # Run only unit tests"
    echo "  $0 -c               # Run tests with coverage"
    echo "  $0 -m               # Generate mocks only"
}

# Main execution
main() {
    local run_all=true
    local run_unit=false
    local run_widget=false
    local run_integration=false
    local generate_coverage=false
    local generate_mocks=false
    
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
            -u|--unit)
                run_unit=true
                run_all=false
                shift
                ;;
            -w|--widget)
                run_widget=true
                run_all=false
                shift
                ;;
            -i|--integration)
                run_integration=true
                run_all=false
                shift
                ;;
            -c|--coverage)
                generate_coverage=true
                run_all=false
                shift
                ;;
            -m|--mocks)
                generate_mocks=true
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
    if [ "$generate_mocks" = true ]; then
        flutter pub run build_runner build --delete-conflicting-outputs
        exit 0
    fi
    
    if [ "$generate_coverage" = true ]; then
        run_tests_with_coverage
    elif [ "$run_unit" = true ]; then
        run_tests "unit" "test/unit/" "Unit Tests"
    elif [ "$run_widget" = true ]; then
        run_tests "widget" "test/widgets/" "Widget Tests"
    elif [ "$run_integration" = true ]; then
        run_tests "integration" "integration_test/" "Integration Tests"
    elif [ "$run_all" = true ]; then
        run_all_tests
    fi
    
    echo -e "\n${GREEN}🎉 Test suite completed successfully!${NC}"
}

# Run main function with all arguments
main "$@"

