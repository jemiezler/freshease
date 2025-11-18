import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/features/account/data/models/user_profile_dto.dart';
import 'package:frontend/features/auth/data/models/user_dto.dart';
import 'package:frontend/features/cart/data/models/cart_dtos.dart';

void main() {
  group('UserProfileDto', () {
    test('fromJson creates DTO with all fields', () {
      final json = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Test User',
        'phone': '+1234567890',
        'bio': 'Test bio',
        'avatar': 'avatar.jpg',
        'cover': 'cover.jpg',
        'date_of_birth': '1990-01-01T00:00:00Z',
        'sex': 'male',
        'goal': 'weight_loss',
        'height_cm': 175.0,
        'weight_kg': 70.0,
        'status': 'active',
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
      };

      final dto = UserProfileDto.fromJson(json);

      expect(dto.id, 'user-1');
      expect(dto.email, 'test@example.com');
      expect(dto.name, 'Test User');
      expect(dto.phone, '+1234567890');
      expect(dto.bio, 'Test bio');
      expect(dto.avatar, 'avatar.jpg');
      expect(dto.cover, 'cover.jpg');
      expect(dto.dateOfBirth, isNotNull);
      expect(dto.sex, 'male');
      expect(dto.goal, 'weight_loss');
      expect(dto.heightCm, 175.0);
      expect(dto.weightKg, 70.0);
      expect(dto.status, 'active');
    });

    test('fromJson handles missing optional fields', () {
      final json = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Test User',
        'status': 'active',
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
      };

      final dto = UserProfileDto.fromJson(json);

      expect(dto.id, 'user-1');
      expect(dto.phone, isNull);
      expect(dto.bio, isNull);
      expect(dto.avatar, isNull);
      expect(dto.dateOfBirth, isNull);
      expect(dto.status, 'active');
    });

    test('fromJson defaults status to active when missing', () {
      final json = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Test User',
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
      };

      final dto = UserProfileDto.fromJson(json);

      expect(dto.status, 'active');
    });

    test('toJson returns correct format', () {
      final dto = UserProfileDto(
        id: 'user-1',
        email: 'test@example.com',
        name: 'Test User',
        phone: '+1234567890',
        status: 'active',
        createdAt: DateTime.parse('2024-01-01T00:00:00Z'),
        updatedAt: DateTime.parse('2024-01-01T00:00:00Z'),
      );

      final json = dto.toJson();

      expect(json['id'], 'user-1');
      expect(json['email'], 'test@example.com');
      expect(json['name'], 'Test User');
      expect(json['phone'], '+1234567890');
      expect(json['status'], 'active');
    });

    test('toEntity converts to UserProfile', () {
      final dto = UserProfileDto(
        id: 'user-1',
        email: 'test@example.com',
        name: 'Test User',
        status: 'active',
        createdAt: DateTime.parse('2024-01-01T00:00:00Z'),
        updatedAt: DateTime.parse('2024-01-01T00:00:00Z'),
      );

      final entity = dto.toEntity();

      expect(entity.id, 'user-1');
      expect(entity.email, 'test@example.com');
      expect(entity.name, 'Test User');
    });
  });

  group('UserDto', () {
    test('fromJson creates DTO with all fields', () {
      final json = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Test User',
        'avatar': 'avatar.jpg',
      };

      final dto = UserDto.fromJson(json);

      expect(dto.id, 'user-1');
      expect(dto.email, 'test@example.com');
      expect(dto.name, 'Test User');
      expect(dto.avatar, 'avatar.jpg');
    });

    test('fromJson handles missing optional fields', () {
      final json = {
        'id': 'user-1',
        'email': 'test@example.com',
      };

      final dto = UserDto.fromJson(json);

      expect(dto.id, 'user-1');
      expect(dto.email, 'test@example.com');
      expect(dto.name, isNull);
      expect(dto.avatar, isNull);
    });

    test('toEntity converts to User', () {
      final dto = UserDto(
        id: 'user-1',
        email: 'test@example.com',
        name: 'Test User',
        avatar: 'avatar.jpg',
      );

      final entity = dto.toEntity();

      expect(entity.id, 'user-1');
      expect(entity.email, 'test@example.com');
      expect(entity.name, 'Test User');
      expect(entity.avatar, 'avatar.jpg');
    });
  });

  group('CartItemDTO', () {
    test('fromJson creates DTO correctly', () {
      final json = {
        'id': 'item-1',
        'product_id': 'product-1',
        'product_name': 'Test Product',
        'product_image': 'image.jpg',
        'product_price': 99.99,
        'quantity': 2,
        'line_total': 199.98,
      };

      final dto = CartItemDTO.fromJson(json);

      expect(dto.id, 'item-1');
      expect(dto.productId, 'product-1');
      expect(dto.productName, 'Test Product');
      expect(dto.productPrice, 99.99);
      expect(dto.quantity, 2);
      expect(dto.lineTotal, 199.98);
    });

    test('toJson returns correct format', () {
      final dto = CartItemDTO(
        id: 'item-1',
        productId: 'product-1',
        productName: 'Test Product',
        productImage: 'image.jpg',
        productPrice: 99.99,
        quantity: 2,
        lineTotal: 199.98,
      );

      final json = dto.toJson();

      expect(json['id'], 'item-1');
      expect(json['product_id'], 'product-1');
      expect(json['product_name'], 'Test Product');
      expect(json['product_price'], 99.99);
      expect(json['quantity'], 2);
      expect(json['line_total'], 199.98);
    });
  });

  group('CartDTO', () {
    test('fromJson creates DTO correctly', () {
      final json = {
        'id': 'cart-1',
        'status': 'active',
        'subtotal': 100.0,
        'shipping': 10.0,
        'tax': 7.0,
        'total': 117.0,
        'items': [
          {
            'id': 'item-1',
            'product_id': 'product-1',
            'product_name': 'Test Product',
            'product_image': 'image.jpg',
            'product_price': 50.0,
            'quantity': 2,
            'line_total': 100.0,
          }
        ],
        'promo_code': 'SAVE10',
        'promo_discount': 10.0,
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
      };

      final dto = CartDTO.fromJson(json);

      expect(dto.id, 'cart-1');
      expect(dto.status, 'active');
      expect(dto.subtotal, 100.0);
      expect(dto.shipping, 10.0);
      expect(dto.tax, 7.0);
      expect(dto.total, 117.0);
      expect(dto.items.length, 1);
      expect(dto.promoCode, 'SAVE10');
      expect(dto.promoDiscount, 10.0);
    });

    test('fromJson handles missing promo code', () {
      final json = {
        'id': 'cart-1',
        'status': 'active',
        'subtotal': 100.0,
        'shipping': 10.0,
        'tax': 7.0,
        'total': 117.0,
        'items': [],
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
      };

      final dto = CartDTO.fromJson(json);

      expect(dto.promoCode, isNull);
      expect(dto.promoDiscount, 0.0);
    });

    test('toJson returns correct format', () {
      final dto = CartDTO(
        id: 'cart-1',
        status: 'active',
        subtotal: 100.0,
        shipping: 10.0,
        tax: 7.0,
        total: 117.0,
        items: [],
        promoCode: 'SAVE10',
        promoDiscount: 10.0,
        createdAt: DateTime.parse('2024-01-01T00:00:00Z'),
        updatedAt: DateTime.parse('2024-01-01T00:00:00Z'),
      );

      final json = dto.toJson();

      expect(json['id'], 'cart-1');
      expect(json['status'], 'active');
      expect(json['subtotal'], 100.0);
      expect(json['promo_code'], 'SAVE10');
      expect(json['promo_discount'], 10.0);
    });
  });

  group('AddToCartRequest', () {
    test('toJson returns correct format', () {
      final request = AddToCartRequest(
        productId: 'product-1',
        quantity: 2,
      );

      final json = request.toJson();

      expect(json['product_id'], 'product-1');
      expect(json['quantity'], 2);
    });
  });

  group('UpdateCartItemRequest', () {
    test('toJson returns correct format', () {
      final request = UpdateCartItemRequest(
        cartItemId: 'item-1',
        quantity: 3,
      );

      final json = request.toJson();

      expect(json['cart_item_id'], 'item-1');
      expect(json['quantity'], 3);
    });
  });

  group('ApplyPromoRequest', () {
    test('toJson returns correct format', () {
      final request = ApplyPromoRequest(promoCode: 'SAVE10');

      final json = request.toJson();

      expect(json['promo_code'], 'SAVE10');
    });
  });
}

