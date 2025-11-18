import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/features/shop/data/models/shop_dtos.dart';

void main() {
  group('ShopProductDTO', () {
    test('fromJson creates DTO with all fields', () {
      final json = {
        'id': 'product-1',
        'name': 'Test Product',
        'price': 99.99,
        'description': 'Test description',
        'image_url': 'image.jpg',
        'unit_label': 'kg',
        'is_active': 'active',
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
        'vendor_id': 'vendor-1',
        'vendor_name': 'Test Vendor',
        'category_id': 'cat-1',
        'category_name': 'Test Category',
        'stock_quantity': 100,
        'is_in_stock': true,
      };

      final dto = ShopProductDTO.fromJson(json);

      expect(dto.id, 'product-1');
      expect(dto.name, 'Test Product');
      expect(dto.price, 99.99);
      expect(dto.description, 'Test description');
      expect(dto.imageUrl, 'image.jpg');
      expect(dto.unitLabel, 'kg');
      expect(dto.isActive, 'active');
      expect(dto.vendorId, 'vendor-1');
      expect(dto.vendorName, 'Test Vendor');
      expect(dto.categoryId, 'cat-1');
      expect(dto.categoryName, 'Test Category');
      expect(dto.stockQuantity, 100);
      expect(dto.isInStock, true);
    });

    test('toJson returns correct format', () {
      final dto = ShopProductDTO(
        id: 'product-1',
        name: 'Test Product',
        price: 99.99,
        description: 'Test description',
        imageUrl: 'image.jpg',
        unitLabel: 'kg',
        isActive: 'active',
        createdAt: DateTime.parse('2024-01-01T00:00:00Z'),
        updatedAt: DateTime.parse('2024-01-01T00:00:00Z'),
        vendorId: 'vendor-1',
        vendorName: 'Test Vendor',
        categoryId: 'cat-1',
        categoryName: 'Test Category',
        stockQuantity: 100,
        isInStock: true,
      );

      final json = dto.toJson();

      expect(json['id'], 'product-1');
      expect(json['name'], 'Test Product');
      expect(json['price'], 99.99);
      expect(json['image_url'], 'image.jpg');
      expect(json['unit_label'], 'kg');
      expect(json['vendor_id'], 'vendor-1');
      expect(json['category_name'], 'Test Category');
      expect(json['stock_quantity'], 100);
      expect(json['is_in_stock'], true);
    });
  });

  group('ShopCategoryDTO', () {
    test('fromJson creates DTO correctly', () {
      final json = {
        'id': 'cat-1',
        'name': 'Test Category',
        'description': 'Test description',
      };

      final dto = ShopCategoryDTO.fromJson(json);

      expect(dto.id, 'cat-1');
      expect(dto.name, 'Test Category');
      expect(dto.description, 'Test description');
    });

    test('toJson returns correct format', () {
      final dto = ShopCategoryDTO(
        id: 'cat-1',
        name: 'Test Category',
        description: 'Test description',
      );

      final json = dto.toJson();

      expect(json['id'], 'cat-1');
      expect(json['name'], 'Test Category');
      expect(json['description'], 'Test description');
    });
  });

  group('ShopVendorDTO', () {
    test('fromJson creates DTO correctly', () {
      final json = {
        'id': 'vendor-1',
        'name': 'Test Vendor',
        'email': 'vendor@example.com',
        'phone': '+1234567890',
        'address': '123 Street',
        'city': 'City',
        'state': 'State',
        'country': 'Country',
        'postal_code': '12345',
        'website': 'https://example.com',
        'logo_url': 'logo.jpg',
        'description': 'Description',
        'is_active': 'active',
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
      };

      final dto = ShopVendorDTO.fromJson(json);

      expect(dto.id, 'vendor-1');
      expect(dto.name, 'Test Vendor');
      expect(dto.email, 'vendor@example.com');
      expect(dto.phone, '+1234567890');
      expect(dto.address, '123 Street');
      expect(dto.city, 'City');
      expect(dto.state, 'State');
      expect(dto.country, 'Country');
      expect(dto.postalCode, '12345');
      expect(dto.website, 'https://example.com');
      expect(dto.logoUrl, 'logo.jpg');
      expect(dto.description, 'Description');
      expect(dto.isActive, 'active');
    });

    test('toJson returns correct format', () {
      final dto = ShopVendorDTO(
        id: 'vendor-1',
        name: 'Test Vendor',
        email: 'vendor@example.com',
        phone: '+1234567890',
        address: '123 Street',
        city: 'City',
        state: 'State',
        country: 'Country',
        postalCode: '12345',
        website: 'https://example.com',
        logoUrl: 'logo.jpg',
        description: 'Description',
        isActive: 'active',
        createdAt: DateTime.parse('2024-01-01T00:00:00Z'),
        updatedAt: DateTime.parse('2024-01-01T00:00:00Z'),
      );

      final json = dto.toJson();

      expect(json['id'], 'vendor-1');
      expect(json['name'], 'Test Vendor');
      expect(json['email'], 'vendor@example.com');
      expect(json['postal_code'], '12345');
      expect(json['logo_url'], 'logo.jpg');
    });
  });

  group('ShopSearchFilters', () {
    test('toQueryParams includes only non-null fields', () {
      final filters = ShopSearchFilters(
        searchTerm: 'test',
        minPrice: 10.0,
        maxPrice: 100.0,
        limit: 20,
        offset: 0,
      );

      final params = filters.toQueryParams();

      expect(params['search'], 'test');
      expect(params['min_price'], 10.0);
      expect(params['max_price'], 100.0);
      expect(params.containsKey('limit'), false); // limit=20 is default, not included
      expect(params.containsKey('offset'), false); // offset=0 is default, not included
    });

    test('toQueryParams excludes null fields and default values', () {
      final filters = ShopSearchFilters(
        limit: 20,
        offset: 0,
      );

      final params = filters.toQueryParams();

      expect(params.containsKey('search'), false);
      expect(params.containsKey('min_price'), false);
      expect(params.containsKey('max_price'), false);
      expect(params.containsKey('limit'), false); // default value not included
      expect(params.containsKey('offset'), false); // default value not included
    });

    test('toQueryParams includes non-default limit and offset', () {
      final filters = ShopSearchFilters(
        limit: 50,
        offset: 10,
      );

      final params = filters.toQueryParams();

      expect(params['limit'], 50);
      expect(params['offset'], 10);
    });

    test('toQueryParams includes all filter fields when provided', () {
      final filters = ShopSearchFilters(
        categoryId: 'cat-1',
        vendorId: 'vendor-1',
        minPrice: 10.0,
        maxPrice: 100.0,
        searchTerm: 'test',
        inStock: true,
        limit: 30,
        offset: 5,
      );

      final params = filters.toQueryParams();

      expect(params['category_id'], 'cat-1');
      expect(params['vendor_id'], 'vendor-1');
      expect(params['min_price'], 10.0);
      expect(params['max_price'], 100.0);
      expect(params['search'], 'test');
      expect(params['in_stock'], true);
      expect(params['limit'], 30);
      expect(params['offset'], 5);
    });
  });

  group('ShopSearchResponse', () {
    test('fromJson creates response correctly', () {
      final json = {
        'products': [
          {
            'id': 'product-1',
            'name': 'Test Product',
            'price': 99.99,
            'description': 'Description',
            'image_url': 'image.jpg',
            'unit_label': 'kg',
            'is_active': 'active',
            'created_at': '2024-01-01T00:00:00Z',
            'updated_at': '2024-01-01T00:00:00Z',
            'vendor_id': 'vendor-1',
            'vendor_name': 'Vendor',
            'category_id': 'cat-1',
            'category_name': 'Category',
            'stock_quantity': 100,
            'is_in_stock': true,
          }
        ],
        'total': 1,
        'limit': 20,
        'offset': 0,
        'has_more': false,
      };

      final response = ShopSearchResponse.fromJson(json);

      expect(response.products.length, 1);
      expect(response.total, 1);
      expect(response.limit, 20);
      expect(response.offset, 0);
      expect(response.hasMore, false);
    });
  });
}

