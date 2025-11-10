# Database Seed Script

This script populates the database with mock data including products, vendors, categories, recipes, meal plans, bundles, and related entities.

## What Gets Seeded

The seed script creates:

- **5 Vendors**: FreshFarm Co., Organic Valley, Local Market, Green Grocers, Farm Fresh Direct
- **8 Categories**: Fruits, Vegetables, Dairy, Meat, Grains, Beverages, Snacks, Organic
- **28 Products**: A variety of fresh produce, dairy, meat, grains, and beverages
- **Product Categories**: Links products to appropriate categories
- **Inventories**: Stock levels for all products
- **8 Recipes**: Healthy meals including breakfast bowls, salads, protein dishes, etc.
- **Recipe Items**: Links products to recipes with quantities
- **1 User**: Test user (john.doe@example.com, password: password123)
- **2 Meal Plans**: Weekly meal plans for the next 2 weeks
- **Meal Plan Items**: Daily meals (breakfast, lunch, dinner) for each week
- **5 Bundles**: Product bundles (Breakfast Bundle, Salad Lover's Pack, Protein Power Pack, etc.)
- **Bundle Items**: Links products to bundles with quantities

## Usage

### Prerequisites

1. Make sure your database is running and accessible
2. Set the `DATABASE_URL` environment variable or ensure your `.env` file has the correct database connection string
3. The database schema should be migrated (run migrations if needed)

### Running the Seed Script

#### Option 1: Build and Run

```bash
cd backend
go build -o bin/seed ./cmd/seed
./bin/seed
```

#### Option 2: Run Directly with Go

```bash
cd backend
go run ./cmd/seed/main.go
```

#### Option 3: Using Make (if Makefile exists)

```bash
make seed
```

## Environment Variables

The script uses the same configuration as the main application:

- `DATABASE_URL`: PostgreSQL connection string (default: `postgres://postgres:user1234@localhost:5432/trail-teller_db?sslmode=disable`)
- Other environment variables from `.env` file

## Notes

- The script will create new data each time it runs. If you run it multiple times, you may encounter duplicate key errors.
- The test user credentials are:
  - Email: `john.doe@example.com`
  - Password: `password123`
- All products are linked to vendors and have inventory entries
- Recipes are linked to products via recipe items
- Meal plans are created for a test user with meals scheduled for breakfast, lunch, and dinner
- Bundles contain multiple products with specified quantities

## Troubleshooting

### Database Connection Issues

If you encounter connection errors:
1. Verify your database is running
2. Check the `DATABASE_URL` environment variable
3. Ensure the database exists and is accessible

### Duplicate Key Errors

If you see duplicate key errors:
1. Clear the existing data from the database
2. Or modify the seed script to check for existing data before creating

### Foreign Key Constraints

If you encounter foreign key constraint errors:
1. Ensure the database schema is properly migrated
2. Check that all required relationships are properly set up

## Customization

You can modify the seed script to:
- Add more products, vendors, or categories
- Change the mock data values
- Add additional recipes or meal plans
- Customize bundle contents
- Add more users with different profiles

Edit `cmd/seed/main.go` to customize the seeded data.
