#!/usr/bin/env python3
"""
Comprehensive test generator for all remaining backend modules
Generates complete test files based on module analysis
"""

import os
import re
from pathlib import Path

BASE_DIR = Path(__file__).parent.parent
MODULES_DIR = BASE_DIR / "modules"

# Module configurations with specific test requirements
MODULE_CONFIGS = {
    "recipes": {
        "dto_prefix": "Recipe",
        "entity": "Recipe",
        "requires": [],
        "fields": ["name", "instructions", "kcal"],
    },
    "reviews": {
        "dto_prefix": "Review",
        "entity": "Review",
        "requires": ["User", "Product"],
        "fields": ["rating", "comment", "user_id", "product_id"],
    },
    "notifications": {
        "dto_prefix": "Notification",
        "entity": "Notification",
        "requires": ["User"],
        "fields": ["title", "message", "type", "user_id"],
    },
    "meal_plans": {
        "dto_prefix": "MealPlan",
        "entity": "MealPlan",
        "requires": ["User"],
        "fields": ["name", "start_date", "end_date", "user_id"],
    },
    "meal_plan_items": {
        "dto_prefix": "MealPlanItem",
        "entity": "MealPlanItem",
        "requires": ["MealPlan", "Recipe"],
        "fields": ["meal_plan_id", "recipe_id", "day", "meal_type"],
    },
    "order_items": {
        "dto_prefix": "OrderItem",
        "entity": "OrderItem",
        "requires": ["Order", "Product"],
        "fields": ["order_id", "product_id", "quantity", "price"],
    },
    "recipe_items": {
        "dto_prefix": "RecipeItem",
        "entity": "RecipeItem",
        "requires": ["Recipe", "Product"],
        "fields": ["recipe_id", "product_id", "quantity", "unit"],
    },
    "bundle_items": {
        "dto_prefix": "BundleItem",
        "entity": "BundleItem",
        "requires": ["Bundle", "Product"],
        "fields": ["bundle_id", "product_id", "quantity"],
    },
}

def generate_repo_test_content(module_name, config):
    """Generate repository test content"""
    package = module_name
    dto_prefix = config["dto_prefix"]
    entity = config["entity"]
    requires = config.get("requires", [])
    
    # Generate test data setup
    setup_code = ""
    if requires:
        setup_code = "// Create required entities\n"
        for req in requires:
            req_lower = req.lower()
            if req == "User":
                setup_code += f'\t{req_lower}, err := client.{req}.Create().\n'
                setup_code += '\t\tSetID(uuid.New()).\n'
                setup_code += '\t\tSetEmail("test@example.com").\n'
                setup_code += '\t\tSetName("Test User").\n'
                setup_code += '\t\tSetPassword("password").\n'
                setup_code += '\t\tSave(ctx)\n'
                setup_code += '\trequire.NoError(t, err)\n\n'
            elif req == "Product":
                setup_code += f'\t// Create vendor first\n'
                setup_code += '\tvendor, err := client.Vendor.Create().\n'
                setup_code += '\t\tSetID(uuid.New()).\n'
                setup_code += '\t\tSetName("Test Vendor").\n'
                setup_code += '\t\tSetContact("vendor@example.com").\n'
                setup_code += '\t\tSave(ctx)\n'
                setup_code += '\trequire.NoError(t, err)\n\n'
                setup_code += f'\t{req_lower}, err := client.{req}.Create().\n'
                setup_code += '\t\tSetID(uuid.New()).\n'
                setup_code += '\t\tSetName("Test Product").\n'
                setup_code += '\t\tSetSku("TEST-001").\n'
                setup_code += '\t\tSetPrice(99.99).\n'
                setup_code += '\t\tSetIsActive(true).\n'
                setup_code += f'\t\tSetVendor(vendor).\n'
                setup_code += '\t\tSave(ctx)\n'
                setup_code += '\trequire.NoError(t, err)\n\n'
            elif req == "Order":
                setup_code += f'\t// Create user first\n'
                setup_code += '\tuser, err := client.User.Create().\n'
                setup_code += '\t\tSetID(uuid.New()).\n'
                setup_code += '\t\tSetEmail("test@example.com").\n'
                setup_code += '\t\tSetName("Test User").\n'
                setup_code += '\t\tSetPassword("password").\n'
                setup_code += '\t\tSave(ctx)\n'
                setup_code += '\trequire.NoError(t, err)\n\n'
                setup_code += f'\t{req_lower}, err := client.{req}.Create().\n'
                setup_code += '\t\tSetID(uuid.New()).\n'
                setup_code += '\t\tSetOrderNo("ORD-001").\n'
                setup_code += '\t\tSetStatus("pending").\n'
                setup_code += '\t\tSetSubtotal(100.0).\n'
                setup_code += '\t\tSetShippingFee(10.0).\n'
                setup_code += '\t\tSetDiscount(0.0).\n'
                setup_code += '\t\tSetTotal(110.0).\n'
                setup_code += '\t\tAddUser(user).\n'
                setup_code += '\t\tSave(ctx)\n'
                setup_code += '\trequire.NoError(t, err)\n\n'
            elif req == "Bundle":
                setup_code += f'\t{req_lower}, err := client.{req}.Create().\n'
                setup_code += '\t\tSetID(uuid.New()).\n'
                setup_code += '\t\tSetName("Test Bundle").\n'
                setup_code += '\t\tSetPrice(99.99).\n'
                setup_code += '\t\tSetIsActive(true).\n'
                setup_code += '\t\tSave(ctx)\n'
                setup_code += '\trequire.NoError(t, err)\n\n'
            elif req == "MealPlan":
                setup_code += f'\t// Create user first\n'
                setup_code += '\tuser, err := client.User.Create().\n'
                setup_code += '\t\tSetID(uuid.New()).\n'
                setup_code += '\t\tSetEmail("test@example.com").\n'
                setup_code += '\t\tSetName("Test User").\n'
                setup_code += '\t\tSetPassword("password").\n'
                setup_code += '\t\tSave(ctx)\n'
                setup_code += '\trequire.NoError(t, err)\n\n'
                setup_code += f'\t{req_lower}, err := client.{req}.Create().\n'
                setup_code += '\t\tSetID(uuid.New()).\n'
                setup_code += '\t\tSetName("Test Meal Plan").\n'
                setup_code += '\t\tAddUser(user).\n'
                setup_code += '\t\tSave(ctx)\n'
                setup_code += '\trequire.NoError(t, err)\n\n'
            elif req == "Recipe":
                setup_code += f'\t{req_lower}, err := client.{req}.Create().\n'
                setup_code += '\t\tSetID(uuid.New()).\n'
                setup_code += '\t\tSetName("Test Recipe").\n'
                setup_code += '\t\tSetKcal(500).\n'
                setup_code += '\t\tSave(ctx)\n'
                setup_code += '\trequire.NoError(t, err)\n\n'
    
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

{setup_code}
	// Create test {package}
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

{setup_code}
	// TODO: Implement test
}}

func TestEntRepo_Create(t *testing.T) {{
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

{setup_code}
	// TODO: Implement test
}}

func TestEntRepo_Update(t *testing.T) {{
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

{setup_code}
	// TODO: Implement test
}}

func TestEntRepo_Delete(t *testing.T) {{
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

{setup_code}
	// TODO: Implement test
}}
'''

def generate_service_test_content(module_name, config):
    """Generate service test content"""
    package = module_name
    dto_prefix = config["dto_prefix"]
    
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

def generate_controller_test_content(module_name, config):
    """Generate controller test content"""
    package = module_name
    dto_prefix = config["dto_prefix"]
    resource_name = package.replace("_", "-")
    resource_plural = resource_name + "s" if not resource_name.endswith("s") else resource_name
    controller_method = "List" + "".join(word.capitalize() for word in resource_plural.split("-"))
    
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
			app.Get("/{resource_name}", controller.List{resource_plural.replace("-", "").title().replace(" ", "")})

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
    print("🔧 Generating comprehensive tests for all remaining backend modules")
    print("=" * 70)
    
    for module_name, config in MODULE_CONFIGS.items():
        module_path = MODULES_DIR / module_name
        if not module_path.exists():
            print(f"⚠️  Module {module_name} not found, skipping...")
            continue
        
        print(f"\n📝 Generating tests for {module_name}...")
        
        # Generate repo_test.go
        repo_test_file = module_path / "repo_test.go"
        if not repo_test_file.exists():
            repo_content = generate_repo_test_content(module_name, config)
            repo_test_file.write_text(repo_content)
            print(f"✅ Created {repo_test_file}")
        else:
            print(f"⚠️  {repo_test_file} already exists, skipping...")
        
        # Generate service_test.go
        service_test_file = module_path / "service_test.go"
        if not service_test_file.exists():
            service_content = generate_service_test_content(module_name, config)
            service_test_file.write_text(service_content)
            print(f"✅ Created {service_test_file}")
        else:
            print(f"⚠️  {service_test_file} already exists, skipping...")
        
        # Generate controller_test.go
        controller_test_file = module_path / "controller_test.go"
        if not controller_test_file.exists():
            controller_content = generate_controller_test_content(module_name, config)
            controller_test_file.write_text(controller_content)
            print(f"✅ Created {controller_test_file}")
        else:
            print(f"⚠️  {controller_test_file} already exists, skipping...")
    
    print("\n🎉 Test generation completed!")
    print("⚠️  Remember to implement the TODO items in the generated test files")
    print("💡 Run 'go test ./modules/...' to verify tests")

if __name__ == "__main__":
    main()

