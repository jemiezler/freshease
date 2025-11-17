#!/bin/bash

# Comprehensive test generator for all backend modules
# This script generates complete test files for modules missing tests

set -e

echo "🔧 Generating comprehensive tests for all backend modules"
echo "========================================================"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

MODULES=(
    "deliveries"
    "payments"
    "recipes"
    "reviews"
    "notifications"
    "meal_plans"
    "meal_plan_items"
    "order_items"
    "recipe_items"
    "bundle_items"
    "uploads"
)

# Function to generate repo_test.go
generate_repo_test() {
    local module=$1
    local package=$(basename $module)
    local file="modules/${module}/repo_test.go"
    
    if [ -f "$file" ]; then
        echo -e "${YELLOW}⚠️  $file already exists, skipping...${NC}"
        return
    fi
    
    cat > "$file" << 'EOF'
package PACKAGE_NAME

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
	// TODO: Add specific test data based on module requirements

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
	// TODO: Implement based on module requirements

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
    sed -i "s/PACKAGE_NAME/$package/g" "$file"
    echo -e "${GREEN}✅ Created $file${NC}"
}

# Function to generate service_test.go
generate_service_test() {
    local module=$1
    local package=$(basename $module)
    local file="modules/${module}/service_test.go"
    
    if [ -f "$file" ]; then
        echo -e "${YELLOW}⚠️  $file already exists, skipping...${NC}"
        return
    fi
    
    cat > "$file" << 'EOF'
package PACKAGE_NAME

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

func (m *MockRepository) List(ctx context.Context) ([]*GetDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, u *CreateDTO) (*GetDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, u *UpdateDTO) (*GetDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name: "success - returns list",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", context.Background()).Return([]*GetDTO{}, nil)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", context.Background()).Return([]*GetDTO(nil), errors.New("database error"))
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := NewService(mockRepo)
			result, err := svc.List(context.Background())

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
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
    # Replace placeholders based on module
    local dto_name=$(echo "$package" | sed 's/_\([a-z]\)/\U\1/g' | sed 's/^\([a-z]\)/\U\1/')
    sed -i "s/PACKAGE_NAME/$package/g" "$file"
    sed -i "s/GetDTO/Get${dto_name}DTO/g" "$file"
    sed -i "s/CreateDTO/Create${dto_name}DTO/g" "$file"
    sed -i "s/UpdateDTO/Update${dto_name}DTO/g" "$file"
    echo -e "${GREEN}✅ Created $file${NC}"
}

# Function to generate controller_test.go
generate_controller_test() {
    local module=$1
    local package=$(basename $module)
    local file="modules/${module}/controller_test.go"
    
    if [ -f "$file" ]; then
        echo -e "${YELLOW}⚠️  $file already exists, skipping...${NC}"
        return
    fi
    
    cat > "$file" << 'EOF'
package PACKAGE_NAME

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

func (m *MockService) List(ctx context.Context) ([]*GetDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateDTO) (*GetDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateDTO) (*GetDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_List(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - returns list",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", context.Background()).Return([]*GetDTO{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "Retrieved Successfully"},
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", context.Background()).Return([]*GetDTO(nil), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"message": "database error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/RESOURCE", controller.ListRESOURCE)

			req := httptest.NewRequest(http.MethodGet, "/RESOURCE", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
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
    local dto_name=$(echo "$package" | sed 's/_\([a-z]\)/\U\1/g' | sed 's/^\([a-z]\)/\U\1/')
    local resource_name=$(echo "$package" | sed 's/_/-/g')
    sed -i "s/PACKAGE_NAME/$package/g" "$file"
    sed -i "s/GetDTO/Get${dto_name}DTO/g" "$file"
    sed -i "s/CreateDTO/Create${dto_name}DTO/g" "$file"
    sed -i "s/UpdateDTO/Update${dto_name}DTO/g" "$file"
    sed -i "s/RESOURCE/$resource_name/g" "$file"
    echo -e "${GREEN}✅ Created $file${NC}"
}

# Main execution
main() {
    for module in "${MODULES[@]}"; do
        if [ ! -d "modules/${module}" ]; then
            echo -e "${YELLOW}⚠️  Module ${module} not found, skipping...${NC}"
            continue
        fi
        
        echo -e "\n${BLUE}📝 Generating tests for ${module}...${NC}"
        
        generate_repo_test "$module"
        generate_service_test "$module"
        generate_controller_test "$module"
    done
    
    echo -e "\n${GREEN}🎉 Test generation completed!${NC}"
    echo -e "${YELLOW}⚠️  Remember to implement the TODO items in the generated test files${NC}"
    echo -e "${BLUE}💡 Run 'go test ./modules/...' to verify tests${NC}"
}

main

