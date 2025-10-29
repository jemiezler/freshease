// test/mocks/mock_repositories.dart
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:frontend/features/account/domain/repositories/user_repository.dart';
import 'package:frontend/features/account/domain/entities/user_profile.dart';
import 'package:frontend/features/shop/data/product_repository.dart';
import 'package:frontend/features/shop/domain/product.dart';
import 'package:frontend/core/health/health_repository.dart';
import 'package:frontend/core/genai/genai_service.dart';
import 'package:frontend/core/genai/models.dart';

// Generate mocks
@GenerateMocks([
  UserRepository,
  ProductRepository,
  HealthRepository,
  GenAiService,
])
void main() {}

// Mock data for testing
class MockData {
  static UserProfile get mockUserProfile => UserProfile(
    id: 'test-user-id',
    email: 'test@example.com',
    name: 'Test User',
    phone: '+1234567890',
    bio: 'Test bio',
    avatar: 'https://example.com/avatar.jpg',
    cover: 'https://example.com/cover.jpg',
    dateOfBirth: DateTime(1990, 1, 1),
    sex: 'male',
    goal: 'weight_loss',
    heightCm: 175.0,
    weightKg: 70.0,
    status: 'active',
    createdAt: DateTime(2024, 1, 1),
    updatedAt: DateTime(2024, 1, 2),
  );

  static Product get mockProduct => Product(
    id: 'test-product-id',
    name: 'Test Product',
    description: 'Test product description',
    price: 29.99,
    image: 'https://example.com/product.jpg',
    category: 'Test Category',
    unitLabel: 'kg',
    vendorId: 'vendor-1',
    vendorName: 'Test Vendor',
    categoryId: 'cat-1',
    categoryName: 'Test Category',
    stockQuantity: 100,
    isInStock: true,
  );

  static List<Product> get mockProducts => [
    mockProduct,
    Product(
      id: 'test-product-2',
      name: 'Test Product 2',
      description: 'Test product 2 description',
      price: 39.99,
      image: 'https://example.com/product2.jpg',
      category: 'Test Category',
      unitLabel: 'kg',
      vendorId: 'vendor-1',
      vendorName: 'Test Vendor',
      categoryId: 'cat-1',
      categoryName: 'Test Category',
      stockQuantity: 50,
      isInStock: true,
    ),
  ];

  static GenAiResponse get mockGenAiResponse => GenAiResponse(
    stepsToday: 8000,
    activeKcal24h: 200.0,
    plan: [
      MealPlan(
        day: 'Monday',
        meals: {
          'Breakfast': 'Healthy breakfast',
          'Lunch': 'Healthy lunch',
          'Dinner': 'Healthy dinner',
        },
        calories: {'Breakfast': 400, 'Lunch': 600, 'Dinner': 500},
        totalCalories: 1500,
      ),
    ],
  );
}
