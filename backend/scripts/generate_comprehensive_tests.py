#!/usr/bin/env python3
"""
Comprehensive test generator for backend modules
This script analyzes module structure and generates complete test files
"""

import os
import re
import json
from pathlib import Path
from typing import Dict, List, Optional

# Module configurations
MODULES = {
    "deliveries": {
        "dto_prefix": "Delivery",
        "entity": "Delivery",
        "requires": ["Order"],
    },
    "payments": {
        "dto_prefix": "Payment",
        "entity": "Payment",
        "requires": ["Order"],
    },
    "recipes": {
        "dto_prefix": "Recipe",
        "entity": "Recipe",
        "requires": [],
    },
    "reviews": {
        "dto_prefix": "Review",
        "entity": "Review",
        "requires": ["Product", "User"],
    },
    "notifications": {
        "dto_prefix": "Notification",
        "entity": "Notification",
        "requires": ["User"],
    },
    "meal_plans": {
        "dto_prefix": "MealPlan",
        "entity": "MealPlan",
        "requires": ["User"],
    },
    "meal_plan_items": {
        "dto_prefix": "MealPlanItem",
        "entity": "MealPlanItem",
        "requires": ["MealPlan", "Recipe"],
    },
    "order_items": {
        "dto_prefix": "OrderItem",
        "entity": "OrderItem",
        "requires": ["Order", "Product"],
    },
    "recipe_items": {
        "dto_prefix": "RecipeItem",
        "entity": "RecipeItem",
        "requires": ["Recipe", "Product"],
    },
    "bundle_items": {
        "dto_prefix": "BundleItem",
        "entity": "BundleItem",
        "requires": ["Bundle", "Product"],
    },
    "uploads": {
        "dto_prefix": "",
        "entity": "",
        "requires": [],
        "special": True,
    },
}

def generate_repo_test(module: str, config: Dict) -> str:
    """Generate repository test file"""
    package = module.split("/")[-1]
    dto_prefix = config.get("dto_prefix", package.title().replace("_", ""))
    entity = config.get("entity", dto_prefix)
    
    return f'''package {package}

import (
	"context"
	"testing"

	"freshease/backend/ent/enttest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func TestEntRepo_List(t *testing.T) {{
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test data
	// TODO: Implement based on module requirements

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.NotNil(t, result)
}}

func TestEntRepo_FindByID(t *testing.T) {{
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// TODO: Implement test
}}

func TestEntRepo_Create(t *testing.T) {{
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// TODO: Implement test
}}

func TestEntRepo_Update(t *testing.T) {{
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// TODO: Implement test
}}

func TestEntRepo_Delete(t *testing.T) {{
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// TODO: Implement test
}}
'''

def generate_service_test(module: str, config: Dict) -> str:
    """Generate service test file"""
    package = module.split("/")[-1]
    dto_prefix = config.get("dto_prefix", package.title().replace("_", ""))
    
    return f'''package {package}

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {{
	mock.Mock
}}

func (m *MockRepository) List(ctx context.Context) ([]*Get{dto_prefix}DTO, error) {{
	args := m.Called(ctx)
	return args.Get(0).([]*Get{dto_prefix}DTO), args.Error(1)
}}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*Get{dto_prefix}DTO, error) {{
	args := m.Called(ctx, id)
	if args.Get(0) == nil {{
		return nil, args.Error(1)
	}}
	return args.Get(0).(*Get{dto_prefix}DTO), args.Error(1)
}}

func (m *MockRepository) Create(ctx context.Context, u *Create{dto_prefix}DTO) (*Get{dto_prefix}DTO, error) {{
	args := m.Called(ctx, u)
	if args.Get(0) == nil {{
		return nil, args.Error(1)
	}}
	return args.Get(0).(*Get{dto_prefix}DTO), args.Error(1)
}}

func (m *MockRepository) Update(ctx context.Context, u *Update{dto_prefix}DTO) (*Get{dto_prefix}DTO, error) {{
	args := m.Called(ctx, u)
	if args.Get(0) == nil {{
		return nil, args.Error(1)
	}}
	return args.Get(0).(*Get{dto_prefix}DTO), args.Error(1)
}}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {{
	args := m.Called(ctx, id)
	return args.Error(0)
}}

func TestService_List(t *testing.T) {{
	tests := []struct {{
		name          string
		mockSetup     func(*MockRepository)
		expectedCount int
		expectedError bool
	}}{{
		{{
			name: "success - returns list",
			mockSetup: func(mockRepo *MockRepository) {{
				mockRepo.On("List", context.Background()).Return([]*Get{dto_prefix}DTO{{}}, nil)
			}},
			expectedCount: 0,
			expectedError: false,
		}},
		{{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {{
				mockRepo.On("List", context.Background()).Return([]*Get{dto_prefix}DTO(nil), errors.New("database error"))
			}},
			expectedCount: 0,
			expectedError: true,
		}},
	}}

	for _, tt := range tests {{
		t.Run(tt.name, func(t *testing.T) {{
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := NewService(mockRepo)
			result, err := svc.List(context.Background())

			if tt.expectedError {{
				assert.Error(t, err)
				assert.Nil(t, result)
			}} else {{
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}}

			mockRepo.AssertExpectations(t)
		}})
	}}
}}

func TestService_Get(t *testing.T) {{
	// TODO: Implement test
}}

func TestService_Create(t *testing.T) {{
	// TODO: Implement test
}}

func TestService_Update(t *testing.T) {{
	// TODO: Implement test
}}

func TestService_Delete(t *testing.T) {{
	// TODO: Implement test
}}
'''

def generate_controller_test(module: str, config: Dict) -> str:
    """Generate controller test file"""
    package = module.split("/")[-1]
    dto_prefix = config.get("dto_prefix", package.title().replace("_", ""))
    resource_name = package.replace("_", "-")
    resource_name_plural = resource_name + "s" if not resource_name.endswith("s") else resource_name
    
    return f'''package {package}

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
type MockService struct {{
	mock.Mock
}}

func (m *MockService) List(ctx context.Context) ([]*Get{dto_prefix}DTO, error) {{
	args := m.Called(ctx)
	return args.Get(0).([]*Get{dto_prefix}DTO), args.Error(1)
}}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*Get{dto_prefix}DTO, error) {{
	args := m.Called(ctx, id)
	if args.Get(0) == nil {{
		return nil, args.Error(1)
	}}
	return args.Get(0).(*Get{dto_prefix}DTO), args.Error(1)
}}

func (m *MockService) Create(ctx context.Context, dto Create{dto_prefix}DTO) (*Get{dto_prefix}DTO, error) {{
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {{
		return nil, args.Error(1)
	}}
	return args.Get(0).(*Get{dto_prefix}DTO), args.Error(1)
}}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto Update{dto_prefix}DTO) (*Get{dto_prefix}DTO, error) {{
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {{
		return nil, args.Error(1)
	}}
	return args.Get(0).(*Get{dto_prefix}DTO), args.Error(1)
}}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {{
	args := m.Called(ctx, id)
	return args.Error(0)
}}

func TestController_List(t *testing.T) {{
	tests := []struct {{
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedBody   map[string]interface{{}}
	}}{{
		{{
			name: "success - returns list",
			mockSetup: func(mockSvc *MockService) {{
				mockSvc.On("List", context.Background()).Return([]*Get{dto_prefix}DTO{{}}, nil)
			}},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{{"message": "Retrieved Successfully"}},
		}},
		{{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {{
				mockSvc.On("List", context.Background()).Return([]*Get{dto_prefix}DTO(nil), errors.New("database error"))
			}},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{{"message": "database error"}},
		}},
	}}

	for _, tt := range tests {{
		t.Run(tt.name, func(t *testing.T) {{
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/{resource_name}", controller.List{resource_name_plural.title().replace("-", "")})

			req := httptest.NewRequest(http.MethodGet, "/{resource_name}", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		}})
	}}
}}

func TestController_Get(t *testing.T) {{
	// TODO: Implement test
}}

func TestController_Create(t *testing.T) {{
	// TODO: Implement test
}}

func TestController_Update(t *testing.T) {{
	// TODO: Implement test
}}

func TestController_Delete(t *testing.T) {{
	// TODO: Implement test
}}
'''

def main():
    """Main function to generate all test files"""
    base_dir = Path(__file__).parent.parent
    modules_dir = base_dir / "modules"
    
    print("🔧 Generating comprehensive tests for all backend modules")
    print("=" * 60)
    
    for module, config in MODULES.items():
        module_path = modules_dir / module
        if not module_path.exists():
            print(f"⚠️  Module {module} not found, skipping...")
            continue
        
        print(f"\n📝 Generating tests for {module}...")
        
        # Generate repo_test.go
        repo_test_file = module_path / "repo_test.go"
        if not repo_test_file.exists():
            repo_test_content = generate_repo_test(module, config)
            repo_test_file.write_text(repo_test_content)
            print(f"✅ Created {repo_test_file}")
        else:
            print(f"⚠️  {repo_test_file} already exists, skipping...")
        
        # Generate service_test.go
        service_test_file = module_path / "service_test.go"
        if not service_test_file.exists():
            service_test_content = generate_service_test(module, config)
            service_test_file.write_text(service_test_content)
            print(f"✅ Created {service_test_file}")
        else:
            print(f"⚠️  {service_test_file} already exists, skipping...")
        
        # Generate controller_test.go
        controller_test_file = module_path / "controller_test.go"
        if not controller_test_file.exists():
            controller_test_content = generate_controller_test(module, config)
            controller_test_file.write_text(controller_test_content)
            print(f"✅ Created {controller_test_file}")
        else:
            print(f"⚠️  {controller_test_file} already exists, skipping...")
    
    print("\n🎉 Test generation completed!")
    print("⚠️  Remember to implement the TODO items in the generated test files")
    print("💡 Run 'go test ./modules/...' to verify tests")

if __name__ == "__main__":
    main()

