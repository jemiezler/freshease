#!/bin/bash

# Script to generate test files for backend modules
# This script creates test templates for modules that are missing tests

set -e

echo "🔧 Generating test files for backend modules"
echo "============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Modules that need tests (excluding those that already have tests)
MODULES=(
    "bundle_items"
    "deliveries"
    "meal_plan_items"
    "meal_plans"
    "notifications"
    "order_items"
    "orders"
    "payments"
    "recipe_items"
    "recipes"
    "reviews"
    "uploads"
    "auth/password"
)

# Template for repo_test.go
generate_repo_test() {
    local module_name=$1
    local package_name=$2
    local entity_name=$3
    
    cat > "modules/${module_name}/repo_test.go" << EOF
package ${package_name}

import (
	"context"
	"testing"

	"freshease/backend/ent/enttest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func TestEntRepo_List(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test data
	// TODO: Add test data creation based on entity schema

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestEntRepo_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test data
	// TODO: Add test data creation

	// Test FindByID
	// TODO: Implement test
}

func TestEntRepo_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// TODO: Implement test
}

func TestEntRepo_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// TODO: Implement test
}

func TestEntRepo_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// TODO: Implement test
}
EOF
    echo -e "${GREEN}✅ Created repo_test.go for ${module_name}${NC}"
}

# Template for service_test.go
generate_service_test() {
    local module_name=$1
    local package_name=$2
    
    cat > "modules/${module_name}/service_test.go" << EOF
package ${package_name}

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

// TODO: Implement mock methods based on Repository interface

func TestService_List(t *testing.T) {
	// TODO: Implement test
}

func TestService_Get(t *testing.T) {
	// TODO: Implement test
}

func TestService_Create(t *testing.T) {
	// TODO: Implement test
}

func TestService_Update(t *testing.T) {
	// TODO: Implement test
}

func TestService_Delete(t *testing.T) {
	// TODO: Implement test
}
EOF
    echo -e "${GREEN}✅ Created service_test.go for ${module_name}${NC}"
}

# Template for controller_test.go
generate_controller_test() {
    local module_name=$1
    local package_name=$2
    
    cat > "modules/${module_name}/controller_test.go" << EOF
package ${package_name}

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

// TODO: Implement mock methods based on Service interface

func TestController_List(t *testing.T) {
	// TODO: Implement test
}

func TestController_Get(t *testing.T) {
	// TODO: Implement test
}

func TestController_Create(t *testing.T) {
	// TODO: Implement test
}

func TestController_Update(t *testing.T) {
	// TODO: Implement test
}

func TestController_Delete(t *testing.T) {
	// TODO: Implement test
}
EOF
    echo -e "${GREEN}✅ Created controller_test.go for ${module_name}${NC}"
}

# Main execution
main() {
    for module in "${MODULES[@]}"; do
        # Convert module path to package name
        package_name=$(basename "$module")
        
        # Skip if module directory doesn't exist
        if [ ! -d "modules/${module}" ]; then
            echo -e "${YELLOW}⚠️  Module ${module} not found, skipping...${NC}"
            continue
        fi
        
        # Skip if tests already exist
        if [ -f "modules/${module}/repo_test.go" ]; then
            echo -e "${YELLOW}⚠️  Tests already exist for ${module}, skipping...${NC}"
            continue
        fi
        
        echo -e "\n${BLUE}📝 Generating tests for ${module}...${NC}"
        
        # Generate test files
        generate_repo_test "$module" "$package_name" "$package_name"
        generate_service_test "$module" "$package_name"
        generate_controller_test "$module" "$package_name"
    done
    
    echo -e "\n${GREEN}🎉 Test generation completed!${NC}"
    echo -e "${YELLOW}⚠️  Remember to implement the TODO items in the generated test files${NC}"
}

# Run main function
main

