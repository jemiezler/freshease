import '../domain/product.dart';
import 'models/shop_dtos.dart';
import 'sources/shop_api.dart';

abstract class ProductRepository {
  Future<List<Product>> list({
    String q = '',
    String? categoryId,
    String category = 'All', // Deprecated, use categoryId instead
    double min = 0,
    double max = 99999,
    int limit = 20,
    int offset = 0,
  });

  Future<Product?> getProduct(String productId);
  Future<List<ShopCategoryDTO>> getCategories();
  Future<List<ShopVendorDTO>> getVendors();
}

/// Real repository implementation using the backend API
class ApiProductRepository implements ProductRepository {
  final ShopApiService _apiService;

  ApiProductRepository(this._apiService);

  @override
  Future<List<Product>> list({
    String q = '',
    String? categoryId,
    String category = 'All',
    double min = 0,
    double max = 99999,
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      final filters = ShopSearchFilters(
        categoryId: categoryId,
        searchTerm: q.isNotEmpty ? q : null,
        minPrice: min > 0 ? min : null,
        maxPrice: max < 99999 ? max : null,
        limit: limit,
        offset: offset,
      );

      final response = await _apiService.searchProducts(filters);

      return response.products.map((dto) => Product.fromShopDTO(dto)).toList();
    } catch (e) {
      // Fallback to empty list on error
      return [];
    }
  }

  @override
  Future<Product?> getProduct(String productId) async {
    try {
      final dto = await _apiService.getProduct(productId);
      return Product.fromShopDTO(dto);
    } catch (e) {
      return null;
    }
  }

  @override
  Future<List<ShopCategoryDTO>> getCategories() async {
    try {
      return await _apiService.getCategories();
    } catch (e) {
      return [];
    }
  }

  @override
  Future<List<ShopVendorDTO>> getVendors() async {
    try {
      return await _apiService.getVendors();
    } catch (e) {
      return [];
    }
  }
}

/// Mock repo – fallback for development/testing
class MockProductRepository implements ProductRepository {
  final _data = <Product>[
    Product.fromLegacy(
      id: 1,
      name: 'Organic Mixed Salad Bowl',
      price: 89,
      image:
          'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSKi4F589upoZSmqvnpdO2RsgAhWFS45WFe4A&s',
      category: 'Salad',
    ),
    Product.fromLegacy(
      id: 2,
      name: 'Organic Romaine Lettuce',
      price: 45,
      image:
          'https://farmlinkhawaii.com/cdn/shop/files/organic_20romaine-800x533.jpg?v=1724305179',
      category: 'Veggies',
    ),
    Product.fromLegacy(
      id: 3,
      name: 'Fresh Organic Strawberry Mix',
      price: 120,
      image:
          'https://img.freepik.com/premium-photo/fresh-organic-summer-mix-strawberry-sweet-cherry-round-stone-plate-dark-textured-surface_114309-2179.jpg',
      category: 'Fruits',
    ),
    Product.fromLegacy(
      id: 4,
      name: 'Avocado & Spinach Raw Set',
      price: 135,
      image:
          'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRL_PFJxbihQiBW0O_kxPVEf94uDtN4-gBJBA&s',
      category: 'Raw Material',
    ),
    Product.fromLegacy(
      id: 5,
      name: 'Tropical Fruit Fusion Pack',
      price: 149,
      image:
          'https://www.packd.co.uk/cdn/shop/files/PACKD-organic-tropical-fruit_a2f8ead5-e78b-4d0c-9e72-bd57871b2b05_800x.png?v=1751978174',
      category: 'Fruits',
    ),
  ];

  @override
  Future<List<Product>> list({
    String q = '',
    String? categoryId,
    String category = 'All',
    double min = 0,
    double max = 99999,
    int limit = 20,
    int offset = 0,
  }) async {
    await Future.delayed(const Duration(milliseconds: 120));
    return _data.where((p) {
      final byQ = q.isEmpty || p.name.toLowerCase().contains(q.toLowerCase());
      final byC = category == 'All' || p.category == category;
      final byP = p.price >= min && p.price <= max;
      return byQ && byC && byP;
    }).toList();
  }

  @override
  Future<Product?> getProduct(String productId) async {
    await Future.delayed(const Duration(milliseconds: 100));
    try {
      return _data.firstWhere((p) => p.id == productId);
    } catch (e) {
      return null;
    }
  }

  @override
  Future<List<ShopCategoryDTO>> getCategories() async {
    await Future.delayed(const Duration(milliseconds: 100));
    return [
      ShopCategoryDTO(id: '1', name: 'All', description: 'All products'),
      ShopCategoryDTO(
        id: '2',
        name: 'Prepared Food',
        description: 'Ready to eat meals',
      ),
      ShopCategoryDTO(
        id: '3',
        name: 'Veggies',
        description: 'Fresh vegetables',
      ),
      ShopCategoryDTO(id: '4', name: 'Fruits', description: 'Fresh fruits'),
      ShopCategoryDTO(
        id: '5',
        name: 'Raw Material',
        description: 'Raw ingredients',
      ),
    ];
  }

  @override
  Future<List<ShopVendorDTO>> getVendors() async {
    await Future.delayed(const Duration(milliseconds: 100));
    return [
      ShopVendorDTO(
        id: '1',
        name: 'Fresh Farm Co.',
        email: 'contact@freshfarm.com',
        phone: '+1234567890',
        address: '123 Farm Road',
        city: 'Farm City',
        state: 'FC',
        country: 'USA',
        postalCode: '12345',
        website: 'https://freshfarm.com',
        logoUrl: 'https://example.com/logo1.png',
        description: 'Fresh organic produce',
        isActive: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      ),
    ];
  }
}
