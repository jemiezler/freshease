package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"freshease/backend/internal/common/config"
	"freshease/backend/internal/common/db"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"freshease/backend/ent"
)

func main() {
	cfg := config.Load()

	// Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, closeDB, err := db.NewEntClientPGX(ctx, cfg.DatabaseURL, cfg.Ent.Debug)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := closeDB(context.Background()); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	log.Println("Starting database seeding...")

	// Seed vendors
	vendors, err := seedVendors(ctx, client)
	if err != nil {
		log.Fatalf("Failed to seed vendors: %v", err)
	}
	log.Printf("Created %d vendors", len(vendors))

	// Seed categories
	categories, err := seedCategories(ctx, client)
	if err != nil {
		log.Fatalf("Failed to seed categories: %v", err)
	}
	log.Printf("Created %d categories", len(categories))

	// Seed products
	products, err := seedProducts(ctx, client, vendors)
	if err != nil {
		log.Fatalf("Failed to seed products: %v", err)
	}
	log.Printf("Created %d products", len(products))

	// Seed product categories
	err = seedProductCategories(ctx, client, products, categories)
	if err != nil {
		log.Fatalf("Failed to seed product categories: %v", err)
	}
	log.Println("Created product categories")

	// Seed inventories
	err = seedInventories(ctx, client, products, vendors)
	if err != nil {
		log.Fatalf("Failed to seed inventories: %v", err)
	}
	log.Println("Created inventories")

	// Seed recipes
	recipes, err := seedRecipes(ctx, client)
	if err != nil {
		log.Fatalf("Failed to seed recipes: %v", err)
	}
	log.Printf("Created %d recipes", len(recipes))

	// Seed recipe items
	err = seedRecipeItems(ctx, client, recipes, products)
	if err != nil {
		log.Fatalf("Failed to seed recipe items: %v", err)
	}
	log.Println("Created recipe items")

	// Seed user
	user, err := seedUser(ctx, client)
	if err != nil {
		log.Fatalf("Failed to seed user: %v", err)
	}
	log.Println("Created user")

	// Seed meal plans
	mealPlans, err := seedMealPlans(ctx, client, user)
	if err != nil {
		log.Fatalf("Failed to seed meal plans: %v", err)
	}
	log.Printf("Created %d meal plans", len(mealPlans))

	// Seed meal plan items
	err = seedMealPlanItems(ctx, client, mealPlans, recipes)
	if err != nil {
		log.Fatalf("Failed to seed meal plan items: %v", err)
	}
	log.Println("Created meal plan items")

	// Seed bundles
	bundles, err := seedBundles(ctx, client)
	if err != nil {
		log.Fatalf("Failed to seed bundles: %v", err)
	}
	log.Printf("Created %d bundles", len(bundles))

	// Seed bundle items
	err = seedBundleItems(ctx, client, bundles, products)
	if err != nil {
		log.Fatalf("Failed to seed bundle items: %v", err)
	}
	log.Println("Created bundle items")

	log.Println("Database seeding completed successfully!")
}

func seedVendors(ctx context.Context, client *ent.Client) ([]*ent.Vendor, error) {
	vendorData := []struct {
		name    string
		contact string
	}{
		{"FreshFarm Co.", "contact@freshfarm.com"},
		{"Organic Valley", "orders@organicvalley.com"},
		{"Local Market", "info@localmarket.com"},
		{"Green Grocers", "sales@greengrocers.com"},
		{"Farm Fresh Direct", "hello@farmfresh.com"},
	}

	vendors := make([]*ent.Vendor, 0, len(vendorData))
	for _, v := range vendorData {
		vendor, err := client.Vendor.Create().
			SetID(uuid.New()).
			SetName(v.name).
			SetContact(v.contact).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create vendor %s: %w", v.name, err)
		}
		vendors = append(vendors, vendor)
	}
	return vendors, nil
}

func seedCategories(ctx context.Context, client *ent.Client) ([]*ent.Category, error) {
	categoryData := []struct {
		name string
		slug string
	}{
		{"Fruits", "fruits"},
		{"Vegetables", "vegetables"},
		{"Dairy", "dairy"},
		{"Meat", "meat"},
		{"Grains", "grains"},
		{"Beverages", "beverages"},
		{"Snacks", "snacks"},
		{"Organic", "organic"},
	}

	categories := make([]*ent.Category, 0, len(categoryData))
	for _, c := range categoryData {
		category, err := client.Category.Create().
			SetID(uuid.New()).
			SetName(c.name).
			SetSlug(c.slug).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create category %s: %w", c.name, err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func seedProducts(ctx context.Context, client *ent.Client, vendors []*ent.Vendor) ([]*ent.Product, error) {
	now := time.Now()
	productData := []struct {
		name        string
		sku         string
		price       float64
		description string
		unitLabel   string
		isActive    bool
		vendorIndex int
	}{
		// Fruits
		{"Fresh Apples", "FR-APPLE-001", 2.99, "Crisp red apples, perfect for snacking", "kg", true, 0},
		{"Organic Bananas", "FR-BANANA-001", 1.99, "Sweet organic bananas", "kg", true, 1},
		{"Strawberries", "FR-STRAW-001", 4.99, "Fresh juicy strawberries", "kg", true, 0},
		{"Blueberries", "FR-BLUE-001", 6.99, "Premium organic blueberries", "kg", true, 1},
		{"Oranges", "FR-ORANGE-001", 3.49, "Juicy sweet oranges", "kg", true, 2},
		{"Grapes", "FR-GRAPE-001", 5.99, "Seedless green grapes", "kg", true, 0},

		// Vegetables
		{"Carrots", "VG-CARROT-001", 2.49, "Fresh crunchy carrots", "kg", true, 2},
		{"Broccoli", "VG-BROC-001", 3.99, "Organic broccoli florets", "kg", true, 1},
		{"Spinach", "VG-SPIN-001", 2.99, "Fresh leafy spinach", "bunch", true, 3},
		{"Tomatoes", "VG-TOMATO-001", 3.49, "Vine-ripened tomatoes", "kg", true, 0},
		{"Bell Peppers", "VG-PEPPER-001", 4.99, "Mixed color bell peppers", "kg", true, 2},
		{"Lettuce", "VG-LETTUCE-001", 2.99, "Crisp romaine lettuce", "head", true, 3},

		// Dairy
		{"Whole Milk", "DY-MILK-001", 4.99, "Fresh whole milk", "liter", true, 4},
		{"Greek Yogurt", "DY-YOGURT-001", 5.99, "Creamy Greek yogurt", "500g", true, 4},
		{"Cheddar Cheese", "DY-CHEESE-001", 7.99, "Sharp cheddar cheese", "250g", true, 4},
		{"Butter", "DY-BUTTER-001", 4.49, "Premium butter", "250g", true, 4},
		{"Eggs", "DY-EGGS-001", 5.99, "Farm-fresh eggs", "dozen", true, 0},

		// Meat
		{"Chicken Breast", "MT-CHICKEN-001", 12.99, "Boneless skinless chicken breast", "kg", true, 0},
		{"Ground Beef", "MT-BEEF-001", 15.99, "Premium ground beef", "kg", true, 0},
		{"Salmon Fillet", "MT-SALMON-001", 18.99, "Fresh Atlantic salmon", "kg", true, 0},
		{"Pork Chops", "MT-PORK-001", 11.99, "Bone-in pork chops", "kg", true, 0},

		// Grains
		{"Brown Rice", "GR-RICE-001", 3.99, "Organic brown rice", "kg", true, 1},
		{"Quinoa", "GR-QUINOA-001", 8.99, "Premium quinoa", "kg", true, 1},
		{"Whole Wheat Bread", "GR-BREAD-001", 4.99, "Fresh whole wheat bread", "loaf", true, 2},
		{"Pasta", "GR-PASTA-001", 2.99, "Organic pasta", "500g", true, 1},

		// Beverages
		{"Orange Juice", "BV-OJ-001", 4.99, "Fresh squeezed orange juice", "liter", true, 2},
		{"Green Tea", "BV-TEA-001", 5.99, "Organic green tea", "box", true, 1},
		{"Coffee Beans", "BV-COFFEE-001", 12.99, "Premium coffee beans", "500g", true, 2},
	}

	products := make([]*ent.Product, 0, len(productData))
	for _, p := range productData {
		vendor := vendors[p.vendorIndex%len(vendors)]
		product, err := client.Product.Create().
			SetID(uuid.New()).
			SetName(p.name).
			SetSku(p.sku).
			SetPrice(p.price).
			SetDescription(p.description).
			SetUnitLabel(p.unitLabel).
			SetIsActive(p.isActive).
			SetCreatedAt(now).
			SetUpdatedAt(now).
			SetVendor(vendor).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create product %s: %w", p.name, err)
		}
		products = append(products, product)
	}
	return products, nil
}

func seedProductCategories(ctx context.Context, client *ent.Client, products []*ent.Product, categories []*ent.Category) error {
	// Create a map for easy category lookup
	categoryMap := make(map[string]*ent.Category)
	for _, cat := range categories {
		categoryMap[cat.Slug] = cat
	}

	// Map products to categories based on their names/SKUs
	categoryAssignments := []struct {
		productIndices []int
		categorySlug   string
	}{
		{[]int{0, 1, 2, 3, 4, 5}, "fruits"},                    // First 6 products are fruits
		{[]int{6, 7, 8, 9, 10, 11}, "vegetables"},              // Next 6 are vegetables
		{[]int{12, 13, 14, 15, 16}, "dairy"},                   // Next 5 are dairy
		{[]int{17, 18, 19, 20}, "meat"},                        // Next 4 are meat
		{[]int{21, 22, 23, 24}, "grains"},                      // Next 4 are grains
		{[]int{25, 26, 27}, "beverages"},                       // Next 3 are beverages
		{[]int{1, 3, 7, 13, 21, 22, 23, 26}, "organic"},       // Some organic products
	}

	for _, assignment := range categoryAssignments {
		category, exists := categoryMap[assignment.categorySlug]
		if !exists {
			continue
		}

		for _, idx := range assignment.productIndices {
			if idx >= len(products) {
				continue
			}
			_, err := client.Product_category.Create().
				SetID(uuid.New()).
				SetProduct(products[idx]).
				SetCategory(category).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to create product category: %w", err)
			}
		}
	}
	return nil
}

func seedInventories(ctx context.Context, client *ent.Client, products []*ent.Product, vendors []*ent.Vendor) error {
	now := time.Now()
	for _, product := range products {
		// Get the vendor from the product
		vendor, err := product.QueryVendor().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed to get vendor for product %s: %w", product.Name, err)
		}

		// Create inventory with random quantities
		quantity := 50 + (len(product.Name) % 200) // Vary quantity based on product name
		reorderLevel := quantity / 4

		_, err = client.Inventory.Create().
			SetID(uuid.New()).
			SetQuantity(quantity).
			SetReorderLevel(reorderLevel).
			SetUpdatedAt(now).
			SetProduct(product).
			SetVendor(vendor).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create inventory for product %s: %w", product.Name, err)
		}
	}
	return nil
}

func seedRecipes(ctx context.Context, client *ent.Client) ([]*ent.Recipe, error) {
	recipeData := []struct {
		name         string
		instructions string
		kcal         int
	}{
		{
			"Healthy Breakfast Bowl",
			"1. Mix Greek yogurt with fresh fruits\n2. Add granola and honey\n3. Top with berries",
			350,
		},
		{
			"Grilled Chicken Salad",
			"1. Season and grill chicken breast\n2. Prepare mixed greens\n3. Add vegetables and dressing",
			450,
		},
		{
			"Salmon with Quinoa",
			"1. Cook quinoa according to package\n2. Season and pan-sear salmon\n3. Serve with steamed vegetables",
			550,
		},
		{
			"Vegetable Stir Fry",
			"1. Heat oil in pan\n2. Add vegetables and stir fry\n3. Season with soy sauce and serve over rice",
			400,
		},
		{
			"Pasta Primavera",
			"1. Cook pasta al dente\n2. Sauté vegetables\n3. Toss with pasta and olive oil",
			480,
		},
		{
			"Fruit Smoothie Bowl",
			"1. Blend frozen fruits\n2. Pour into bowl\n3. Top with fresh fruits and granola",
			320,
		},
		{
			"Beef and Rice Bowl",
			"1. Cook rice\n2. Season and cook ground beef\n3. Serve together with vegetables",
			600,
		},
		{
			"Caesar Salad",
			"1. Prepare romaine lettuce\n2. Add croutons and parmesan\n3. Toss with Caesar dressing",
			380,
		},
	}

	recipes := make([]*ent.Recipe, 0, len(recipeData))
	for _, r := range recipeData {
		recipe, err := client.Recipe.Create().
			SetID(uuid.New()).
			SetName(r.name).
			SetInstructions(r.instructions).
			SetKcal(r.kcal).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create recipe %s: %w", r.name, err)
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func seedRecipeItems(ctx context.Context, client *ent.Client, recipes []*ent.Recipe, products []*ent.Product) error {
	// Map recipes to products
	recipeItems := []struct {
		recipeIndex int
		productIndices []int
		amounts     []float64
		units       []string
	}{
		{0, []int{13, 1, 2}, []float64{200, 100, 50}, []string{"g", "g", "g"}},           // Breakfast Bowl: yogurt, banana, strawberries
		{1, []int{17, 11, 9}, []float64{200, 100, 150}, []string{"g", "g", "g"}},        // Chicken Salad: chicken, lettuce, tomatoes
		{2, []int{19, 22, 7}, []float64{250, 150, 100}, []string{"g", "g", "g"}},        // Salmon Quinoa: salmon, quinoa, broccoli
		{3, []int{6, 7, 9, 21}, []float64{200, 150, 200, 100}, []string{"g", "g", "g", "g"}}, // Stir Fry: carrots, broccoli, tomatoes, rice
		{4, []int{24, 9, 7}, []float64{200, 150, 100}, []string{"g", "g", "g"}},         // Pasta: pasta, tomatoes, broccoli
		{5, []int{1, 2, 3}, []float64{150, 100, 80}, []string{"g", "g", "g"}},           // Smoothie: banana, strawberries, blueberries
		{6, []int{18, 21, 9}, []float64{200, 150, 100}, []string{"g", "g", "g"}},        // Beef Bowl: beef, rice, tomatoes
		{7, []int{11, 23}, []float64{150, 2}, []string{"g", "slices"}},                  // Caesar: lettuce, bread (for croutons)
	}

	for _, ri := range recipeItems {
		if ri.recipeIndex >= len(recipes) {
			continue
		}
		recipe := recipes[ri.recipeIndex]

		for i, productIdx := range ri.productIndices {
			if productIdx >= len(products) || i >= len(ri.amounts) || i >= len(ri.units) {
				continue
			}

			_, err := client.Recipe_item.Create().
				SetID(uuid.New()).
				SetAmount(ri.amounts[i]).
				SetUnit(ri.units[i]).
				SetRecipe(recipe).
				SetProduct(products[productIdx]).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to create recipe item: %w", err)
			}
		}
	}
	return nil
}

func seedUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	password := string(hashedPassword)
	now := time.Now()
	goal := "weight_loss"
	sex := "male"
	height := 175.0
	weight := 75.0

	user, err := client.User.Create().
		SetID(uuid.New()).
		SetName("John Doe").
		SetEmail("john.doe@example.com").
		SetPassword(password).
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetNillableGoal(&goal).
		SetNillableSex(&sex).
		SetNillableHeightCm(&height).
		SetNillableWeightKg(&weight).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func seedMealPlans(ctx context.Context, client *ent.Client, user *ent.User) ([]*ent.Meal_plan, error) {
	now := time.Now()
	// Create meal plans for the next 2 weeks
	mealPlans := make([]*ent.Meal_plan, 0, 2)

	for i := 0; i < 2; i++ {
		weekStart := now.AddDate(0, 0, i*7)
		// Set to Monday of that week
		// Weekday: Sunday=0, Monday=1, ..., Saturday=6
		// To get Monday: subtract (weekday - 1), but handle Sunday specially
		weekday := int(weekStart.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday becomes 7
		}
		weekStart = weekStart.AddDate(0, 0, -weekday+1)
		goal := "weight_loss"

		mealPlan, err := client.Meal_plan.Create().
			SetID(uuid.New()).
			SetWeekStart(weekStart).
			SetNillableGoal(&goal).
			SetUser(user).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create meal plan: %w", err)
		}
		mealPlans = append(mealPlans, mealPlan)
	}
	return mealPlans, nil
}

func seedMealPlanItems(ctx context.Context, client *ent.Client, mealPlans []*ent.Meal_plan, recipes []*ent.Recipe) error {
	slots := []string{"breakfast", "lunch", "dinner"}
	daysPerWeek := 7

	for _, mealPlan := range mealPlans {
		weekStart := mealPlan.WeekStart
		for day := 0; day < daysPerWeek; day++ {
			currentDay := weekStart.AddDate(0, 0, day)
			for slotIdx, slot := range slots {
				recipeIdx := (day + slotIdx) % len(recipes)
				if recipeIdx >= len(recipes) {
					continue
				}

				_, err := client.Meal_plan_item.Create().
					SetID(uuid.New()).
					SetDay(currentDay).
					SetSlot(slot).
					SetMealPlan(mealPlan).
					SetRecipe(recipes[recipeIdx]).
					Save(ctx)
				if err != nil {
					return fmt.Errorf("failed to create meal plan item: %w", err)
				}
			}
		}
	}
	return nil
}

func seedBundles(ctx context.Context, client *ent.Client) ([]*ent.Bundle, error) {
	bundleData := []struct {
		name        string
		description string
		price       float64
		isActive    bool
	}{
		{
			"Breakfast Bundle",
			"Everything you need for a healthy breakfast",
			15.99,
			true,
		},
		{
			"Salad Lover's Pack",
			"Fresh greens and vegetables for your salads",
			12.99,
			true,
		},
		{
			"Protein Power Pack",
			"High-protein foods for your fitness goals",
			35.99,
			true,
		},
		{
			"Fresh Fruit Basket",
			"Assorted fresh fruits",
			18.99,
			true,
		},
		{
			"Vegetable Variety Pack",
			"Mixed vegetables for cooking",
			14.99,
			true,
		},
	}

	bundles := make([]*ent.Bundle, 0, len(bundleData))
	for _, b := range bundleData {
		bundle, err := client.Bundle.Create().
			SetID(uuid.New()).
			SetName(b.name).
			SetNillableDescription(&b.description).
			SetPrice(b.price).
			SetIsActive(b.isActive).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create bundle %s: %w", b.name, err)
		}
		bundles = append(bundles, bundle)
	}
	return bundles, nil
}

func seedBundleItems(ctx context.Context, client *ent.Client, bundles []*ent.Bundle, products []*ent.Product) error {
	// Map bundles to products
	bundleItems := []struct {
		bundleIndex  int
		productIndices []int
		quantities   []int
	}{
		{0, []int{13, 1, 16, 23}, []int{1, 2, 1, 1}},           // Breakfast: yogurt, banana, eggs, bread
		{1, []int{11, 8, 9, 10}, []int{2, 2, 1, 2}},            // Salad: lettuce, spinach, tomatoes, peppers
		{2, []int{17, 19, 18}, []int{1, 1, 1}},                 // Protein: chicken, salmon, beef
		{3, []int{0, 1, 2, 3, 4}, []int{1, 2, 1, 1, 1}},       // Fruits: apples, bananas, strawberries, blueberries, oranges
		{4, []int{6, 7, 9, 10, 11}, []int{1, 1, 1, 1, 1}},     // Vegetables: carrots, broccoli, tomatoes, peppers, lettuce
	}

	for _, bi := range bundleItems {
		if bi.bundleIndex >= len(bundles) {
			continue
		}
		bundle := bundles[bi.bundleIndex]

		for i, productIdx := range bi.productIndices {
			if productIdx >= len(products) || i >= len(bi.quantities) {
				continue
			}

			_, err := client.Bundle_item.Create().
				SetID(uuid.New()).
				SetQty(bi.quantities[i]).
				SetBundle(bundle).
				SetProduct(products[productIdx]).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to create bundle item: %w", err)
			}
		}
	}
	return nil
}
