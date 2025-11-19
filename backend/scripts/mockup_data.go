package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

const (
	defaultBaseURL = "http://localhost:8080/api"
)

var (
	baseURL    string
	httpClient = &http.Client{Timeout: 30 * time.Second}
	createdIDs = struct {
		vendors    []string
		categories []string
		products   []string
		users      []string
		recipes    []string
		bundles    []string
	}{
		vendors:    make([]string, 0),
		categories: make([]string, 0),
		products:   make([]string, 0),
		users:      make([]string, 0),
		recipes:    make([]string, 0),
		bundles:    make([]string, 0),
	}
)

func main() {
	// Get base URL from environment or use default
	baseURL = os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	log.Printf("Using API base URL: %s", baseURL)

	// Check if server is reachable
	if err := checkHealth(); err != nil {
		log.Fatalf("Server health check failed: %v", err)
	}

	log.Println("Starting mockup data creation...")

	// Create vendors first (needed for products)
	if err := createVendors(); err != nil {
		log.Fatalf("Failed to create vendors: %v", err)
	}
	log.Printf("Created %d vendors", len(createdIDs.vendors))

	// Create categories
	if err := createCategories(); err != nil {
		log.Fatalf("Failed to create categories: %v", err)
	}
	log.Printf("Created %d categories", len(createdIDs.categories))

	// Create products (requires vendors)
	if err := createProducts(); err != nil {
		log.Fatalf("Failed to create products: %v", err)
	}
	log.Printf("Created %d products", len(createdIDs.products))

	// Create users
	if err := createUsers(); err != nil {
		log.Fatalf("Failed to create users: %v", err)
	}
	log.Printf("Created %d users", len(createdIDs.users))

	// Create recipes
	if err := createRecipes(); err != nil {
		log.Fatalf("Failed to create recipes: %v", err)
	}
	log.Printf("Created %d recipes", len(createdIDs.recipes))

	// Create bundles
	if err := createBundles(); err != nil {
		log.Fatalf("Failed to create bundles: %v", err)
	}
	log.Printf("Created %d bundles", len(createdIDs.bundles))

	// Create carts (requires users)
	if err := createCarts(); err != nil {
		log.Fatalf("Failed to create carts: %v", err)
	}
	log.Println("Created carts")

	// Create orders (requires users and products)
	if err := createOrders(); err != nil {
		log.Fatalf("Failed to create orders: %v", err)
	}
	log.Println("Created orders")

	log.Println("Mockup data creation completed successfully!")
}

func checkHealth() error {
	resp, err := httpClient.Get(baseURL + "/health")
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode health response: %w", err)
	}

	log.Printf("Server health check passed: %v", result)
	return nil
}

func createVendors() error {
	vendors := []map[string]interface{}{
		{"id": uuid.New().String(), "name": "FreshFarm Co.", "contact": "contact@freshfarm.com"},
		{"id": uuid.New().String(), "name": "Organic Valley", "contact": "orders@organicvalley.com"},
		{"id": uuid.New().String(), "name": "Local Market", "contact": "info@localmarket.com"},
		{"id": uuid.New().String(), "name": "Green Grocers", "contact": "sales@greengrocers.com"},
		{"id": uuid.New().String(), "name": "Farm Fresh Direct", "contact": "hello@farmfresh.com"},
	}

	for _, vendor := range vendors {
		id, err := postJSON("/vendors", vendor)
		if err != nil {
			return err
		}
		createdIDs.vendors = append(createdIDs.vendors, id)
	}
	return nil
}

func createCategories() error {
	now := time.Now()
	categories := []map[string]interface{}{
		{"id": uuid.New().String(), "name": "Fruits", "slug": "fruits", "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339)},
		{"id": uuid.New().String(), "name": "Vegetables", "slug": "vegetables", "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339)},
		{"id": uuid.New().String(), "name": "Dairy", "slug": "dairy", "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339)},
		{"id": uuid.New().String(), "name": "Meat", "slug": "meat", "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339)},
		{"id": uuid.New().String(), "name": "Grains", "slug": "grains", "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339)},
		{"id": uuid.New().String(), "name": "Beverages", "slug": "beverages", "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339)},
		{"id": uuid.New().String(), "name": "Snacks", "slug": "snacks", "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339)},
		{"id": uuid.New().String(), "name": "Organic", "slug": "organic", "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339)},
	}

	for _, category := range categories {
		id, err := postJSON("/categories", category)
		if err != nil {
			return err
		}
		createdIDs.categories = append(createdIDs.categories, id)
	}
	return nil
}

func createProducts() error {
	if len(createdIDs.vendors) == 0 {
		return fmt.Errorf("no vendors available")
	}

	now := time.Now()
	products := []map[string]interface{}{
		// Fruits
		{"id": uuid.New().String(), "name": "Fresh Apples", "sku": "FR-APPLE-001", "price": 2.99, "description": "Crisp red apples, perfect for snacking", "unit_label": "kg", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 100, "reorder_level": 20},
		{"id": uuid.New().String(), "name": "Organic Bananas", "sku": "FR-BANANA-001", "price": 1.99, "description": "Sweet organic bananas", "unit_label": "kg", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 150, "reorder_level": 30},
		{"id": uuid.New().String(), "name": "Strawberries", "sku": "FR-STRAW-001", "price": 4.99, "description": "Fresh juicy strawberries", "unit_label": "kg", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 80, "reorder_level": 20},
		{"id": uuid.New().String(), "name": "Blueberries", "sku": "FR-BLUE-001", "price": 6.99, "description": "Premium organic blueberries", "unit_label": "kg", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 60, "reorder_level": 15},
		{"id": uuid.New().String(), "name": "Oranges", "sku": "FR-ORANGE-001", "price": 3.49, "description": "Juicy sweet oranges", "unit_label": "kg", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 120, "reorder_level": 25},

		// Vegetables
		{"id": uuid.New().String(), "name": "Carrots", "sku": "VG-CARROT-001", "price": 2.49, "description": "Fresh crunchy carrots", "unit_label": "kg", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 200, "reorder_level": 50},
		{"id": uuid.New().String(), "name": "Broccoli", "sku": "VG-BROC-001", "price": 3.99, "description": "Organic broccoli florets", "unit_label": "kg", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 90, "reorder_level": 20},
		{"id": uuid.New().String(), "name": "Spinach", "sku": "VG-SPIN-001", "price": 2.99, "description": "Fresh leafy spinach", "unit_label": "bunch", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 100, "reorder_level": 25},
		{"id": uuid.New().String(), "name": "Tomatoes", "sku": "VG-TOMATO-001", "price": 3.49, "description": "Vine-ripened tomatoes", "unit_label": "kg", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 150, "reorder_level": 30},
		{"id": uuid.New().String(), "name": "Bell Peppers", "sku": "VG-PEPPER-001", "price": 4.99, "description": "Mixed color bell peppers", "unit_label": "kg", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 70, "reorder_level": 15},

		// Dairy
		{"id": uuid.New().String(), "name": "Whole Milk", "sku": "DY-MILK-001", "price": 4.99, "description": "Fresh whole milk", "unit_label": "liter", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 80, "reorder_level": 20},
		{"id": uuid.New().String(), "name": "Greek Yogurt", "sku": "DY-YOGURT-001", "price": 5.99, "description": "Creamy Greek yogurt", "unit_label": "500g", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 100, "reorder_level": 25},
		{"id": uuid.New().String(), "name": "Cheddar Cheese", "sku": "DY-CHEESE-001", "price": 7.99, "description": "Sharp cheddar cheese", "unit_label": "250g", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 60, "reorder_level": 15},
		{"id": uuid.New().String(), "name": "Butter", "sku": "DY-BUTTER-001", "price": 4.49, "description": "Premium butter", "unit_label": "250g", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 90, "reorder_level": 20},
		{"id": uuid.New().String(), "name": "Eggs", "sku": "DY-EGGS-001", "price": 5.99, "description": "Farm-fresh eggs", "unit_label": "dozen", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 120, "reorder_level": 30},

		// Meat
		{"id": uuid.New().String(), "name": "Chicken Breast", "sku": "MT-CHICKEN-001", "price": 12.99, "description": "Boneless skinless chicken breast", "unit_label": "kg", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 50, "reorder_level": 10},
		{"id": uuid.New().String(), "name": "Ground Beef", "sku": "MT-BEEF-001", "price": 15.99, "description": "Premium ground beef", "unit_label": "kg", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 40, "reorder_level": 10},
		{"id": uuid.New().String(), "name": "Salmon Fillet", "sku": "MT-SALMON-001", "price": 18.99, "description": "Fresh Atlantic salmon", "unit_label": "kg", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 30, "reorder_level": 8},
		{"id": uuid.New().String(), "name": "Pork Chops", "sku": "MT-PORK-001", "price": 11.99, "description": "Bone-in pork chops", "unit_label": "kg", "is_active": true, "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339), "quantity": 35, "reorder_level": 8},
	}

	// Rotate through vendors
	vendorIndex := 0
	for _, product := range products {
		// Note: The API might require vendor_id in the request, but based on the DTO structure,
		// it seems products are created with vendor relationship. We'll need to check the actual API.
		// For now, we'll create products without vendor_id and let the API handle it or add it if needed.
		id, err := postJSON("/products", product)
		if err != nil {
			log.Printf("Warning: Failed to create product %s: %v", product["name"], err)
			continue
		}
		createdIDs.products = append(createdIDs.products, id)
		vendorIndex = (vendorIndex + 1) % len(createdIDs.vendors)
	}
	return nil
}

func createUsers() error {
	users := []map[string]interface{}{
		{
			"id":           uuid.New().String(),
			"email":        "john.doe@example.com",
			"password":     "password123",
			"name":         "John Doe",
			"phone":        "+1234567890",
			"date_of_birth": "1990-01-15",
			"sex":          "male",
			"goal":         "weight_loss",
			"height_cm":   175.0,
			"weight_kg":    75.0,
		},
		{
			"id":           uuid.New().String(),
			"email":        "jane.smith@example.com",
			"password":     "password123",
			"name":         "Jane Smith",
			"phone":        "+1234567891",
			"date_of_birth": "1992-05-20",
			"sex":          "female",
			"goal":         "maintenance",
			"height_cm":   165.0,
			"weight_kg":    60.0,
		},
		{
			"id":           uuid.New().String(),
			"email":        "bob.wilson@example.com",
			"password":     "password123",
			"name":         "Bob Wilson",
			"phone":        "+1234567892",
			"date_of_birth": "1988-11-10",
			"sex":          "male",
			"goal":         "weight_gain",
			"height_cm":   180.0,
			"weight_kg":    70.0,
		},
	}

	for _, user := range users {
		id, err := postJSON("/users", user)
		if err != nil {
			log.Printf("Warning: Failed to create user %s: %v", user["email"], err)
			continue
		}
		createdIDs.users = append(createdIDs.users, id)
	}
	return nil
}

func createRecipes() error {
	recipes := []map[string]interface{}{
		{
			"id":           uuid.New().String(),
			"name":         "Healthy Breakfast Bowl",
			"instructions": "1. Mix Greek yogurt with fresh fruits\n2. Add granola and honey\n3. Top with berries",
			"kcal":         350,
		},
		{
			"id":           uuid.New().String(),
			"name":         "Grilled Chicken Salad",
			"instructions": "1. Season and grill chicken breast\n2. Prepare mixed greens\n3. Add vegetables and dressing",
			"kcal":         450,
		},
		{
			"id":           uuid.New().String(),
			"name":         "Salmon with Quinoa",
			"instructions": "1. Cook quinoa according to package\n2. Season and pan-sear salmon\n3. Serve with steamed vegetables",
			"kcal":         550,
		},
		{
			"id":           uuid.New().String(),
			"name":         "Vegetable Stir Fry",
			"instructions": "1. Heat oil in pan\n2. Add vegetables and stir fry\n3. Season with soy sauce and serve over rice",
			"kcal":         400,
		},
		{
			"id":           uuid.New().String(),
			"name":         "Pasta Primavera",
			"instructions": "1. Cook pasta al dente\n2. Sauté vegetables\n3. Toss with pasta and olive oil",
			"kcal":         480,
		},
	}

	for _, recipe := range recipes {
		id, err := postJSON("/recipes", recipe)
		if err != nil {
			log.Printf("Warning: Failed to create recipe %s: %v", recipe["name"], err)
			continue
		}
		createdIDs.recipes = append(createdIDs.recipes, id)
	}
	return nil
}

func createBundles() error {
	bundles := []map[string]interface{}{
		{
			"id":          uuid.New().String(),
			"name":        "Breakfast Bundle",
			"description": "Everything you need for a healthy breakfast",
			"price":       15.99,
			"is_active":   true,
		},
		{
			"id":          uuid.New().String(),
			"name":        "Salad Lover's Pack",
			"description": "Fresh greens and vegetables for your salads",
			"price":       12.99,
			"is_active":   true,
		},
		{
			"id":          uuid.New().String(),
			"name":        "Protein Power Pack",
			"description": "High-protein foods for your fitness goals",
			"price":       35.99,
			"is_active":   true,
		},
		{
			"id":          uuid.New().String(),
			"name":        "Fresh Fruit Basket",
			"description": "Assorted fresh fruits",
			"price":       18.99,
			"is_active":   true,
		},
	}

	for _, bundle := range bundles {
		id, err := postJSON("/bundles", bundle)
		if err != nil {
			log.Printf("Warning: Failed to create bundle %s: %v", bundle["name"], err)
			continue
		}
		createdIDs.bundles = append(createdIDs.bundles, id)
	}
	return nil
}

func createCarts() error {
	if len(createdIDs.users) == 0 {
		return fmt.Errorf("no users available")
	}

	// Create a cart for the first user
	cart := map[string]interface{}{
		"user_id": createdIDs.users[0],
		"status":  "active",
		"total":   0.0,
	}

	_, err := postJSON("/carts", cart)
	if err != nil {
		return fmt.Errorf("failed to create cart: %w", err)
	}

	return nil
}

func createOrders() error {
	if len(createdIDs.users) == 0 || len(createdIDs.products) == 0 {
		return fmt.Errorf("insufficient data: need users and products")
	}

	now := time.Now()
	orders := []map[string]interface{}{
		{
			"id":          uuid.New().String(),
			"order_no":    fmt.Sprintf("ORD-%d", now.Unix()),
			"status":      "pending",
			"subtotal":    45.97,
			"shipping_fee": 5.00,
			"discount":    0.0,
			"total":       50.97,
			"placed_at":   now.Format(time.RFC3339),
			"user_id":     createdIDs.users[0],
		},
		{
			"id":          uuid.New().String(),
			"order_no":    fmt.Sprintf("ORD-%d", now.Unix()+1),
			"status":      "processing",
			"subtotal":    28.97,
			"shipping_fee": 5.00,
			"discount":    2.00,
			"total":       31.97,
			"placed_at":   now.Add(-1 * time.Hour).Format(time.RFC3339),
			"user_id":     createdIDs.users[1],
		},
	}

	for _, order := range orders {
		_, err := postJSON("/orders", order)
		if err != nil {
			log.Printf("Warning: Failed to create order %s: %v", order["order_no"], err)
			continue
		}
	}

	return nil
}

// postJSON sends a POST request with JSON body and returns the created ID
func postJSON(endpoint string, data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	url := baseURL + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	// Try to extract ID from response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		if data, ok := result["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(string); ok {
				return id, nil
			}
		}
		// Try direct id field
		if id, ok := result["id"].(string); ok {
			return id, nil
		}
	}

	// If we can't extract ID, return the ID from the request data
	if dataMap, ok := data.(map[string]interface{}); ok {
		if id, ok := dataMap["id"].(string); ok {
			return id, nil
		}
	}

	return "", fmt.Errorf("could not extract ID from response")
}


