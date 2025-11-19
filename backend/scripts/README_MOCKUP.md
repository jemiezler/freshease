# Mockup Data Script

This script creates mockup data by posting to the server via HTTP API calls.

## Usage

### Prerequisites
- The backend server must be running
- Default server URL: `http://localhost:8080/api`
- Can be overridden with `API_BASE_URL` environment variable

### Run the script

```bash
# Using default URL (http://localhost:8080/api)
go run scripts/mockup_data.go

# Using custom server URL
API_BASE_URL=http://your-server:8080/api go run scripts/mockup_data.go
```

### What it creates

The script creates the following mockup data:

1. **Vendors** (5 vendors)
   - FreshFarm Co.
   - Organic Valley
   - Local Market
   - Green Grocers
   - Farm Fresh Direct

2. **Categories** (8 categories)
   - Fruits, Vegetables, Dairy, Meat, Grains, Beverages, Snacks, Organic

3. **Products** (20+ products)
   - Fruits: Apples, Bananas, Strawberries, Blueberries, Oranges
   - Vegetables: Carrots, Broccoli, Spinach, Tomatoes, Bell Peppers
   - Dairy: Milk, Yogurt, Cheese, Butter, Eggs
   - Meat: Chicken, Beef, Salmon, Pork
   - Each product includes inventory quantities and reorder levels

4. **Users** (3 users)
   - John Doe (weight_loss goal)
   - Jane Smith (maintenance goal)
   - Bob Wilson (weight_gain goal)

5. **Recipes** (5 recipes)
   - Healthy Breakfast Bowl
   - Grilled Chicken Salad
   - Salmon with Quinoa
   - Vegetable Stir Fry
   - Pasta Primavera

6. **Bundles** (4 bundles)
   - Breakfast Bundle
   - Salad Lover's Pack
   - Protein Power Pack
   - Fresh Fruit Basket

7. **Carts** (1 active cart for first user)

8. **Orders** (2 sample orders)

## Notes

- The script performs a health check before creating data
- If any entity creation fails, it logs a warning and continues
- All IDs are generated using UUIDs
- The script maintains a list of created IDs for reference
- Some entities may require authentication - adjust the script if needed

## Error Handling

The script will:
- Check server health before starting
- Log warnings for individual failures but continue processing
- Exit with error if critical steps fail (vendors, categories, etc.)

## Customization

To customize the mockup data:
1. Edit the data arrays in the respective functions
2. Modify the `create*` functions to add more entities
3. Adjust quantities, prices, and other attributes as needed


