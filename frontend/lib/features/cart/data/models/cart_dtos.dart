class CartItemDTO {
  final String id;
  final String productId;
  final String productName;
  final String productImage;
  final double productPrice;
  final int quantity;
  final double lineTotal;

  CartItemDTO({
    required this.id,
    required this.productId,
    required this.productName,
    required this.productImage,
    required this.productPrice,
    required this.quantity,
    required this.lineTotal,
  });

  factory CartItemDTO.fromJson(Map<String, dynamic> json) {
    // Handle product_image which can be null or a string
    final productImageValue = json['product_image'];
    final productImage = productImageValue is String 
        ? productImageValue 
        : productImageValue?.toString() ?? '';
    
    // Handle UUID fields - ensure they're strings
    final idValue = json['id'];
    final idString = idValue is String 
        ? idValue 
        : idValue?.toString() ?? '';
    
    final productIdValue = json['product_id'];
    final productIdString = productIdValue is String 
        ? productIdValue 
        : productIdValue?.toString() ?? '';
    
    return CartItemDTO(
      id: idString,
      productId: productIdString,
      productName: json['product_name'] as String? ?? '',
      productImage: productImage,
      productPrice: (json['product_price'] as num?)?.toDouble() ?? 0.0,
      quantity: json['quantity'] as int? ?? 0,
      lineTotal: (json['line_total'] as num?)?.toDouble() ?? 0.0,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'product_id': productId,
      'product_name': productName,
      'product_image': productImage,
      'product_price': productPrice,
      'quantity': quantity,
      'line_total': lineTotal,
    };
  }
}

class CartDTO {
  final String id;
  final String status;
  final double subtotal;
  final double shipping;
  final double tax;
  final double total;
  final List<CartItemDTO> items;
  final String? promoCode;
  final double promoDiscount;
  final DateTime createdAt;
  final DateTime updatedAt;

  CartDTO({
    required this.id,
    required this.status,
    required this.subtotal,
    required this.shipping,
    required this.tax,
    required this.total,
    required this.items,
    this.promoCode,
    required this.promoDiscount,
    required this.createdAt,
    required this.updatedAt,
  });

  factory CartDTO.fromJson(Map<String, dynamic> json) {
    // Handle UUID id field - ensure it's a string
    final idValue = json['id'];
    final idString = idValue is String 
        ? idValue 
        : idValue?.toString() ?? '';
    
    return CartDTO(
      id: idString,
      status: json['status'] as String? ?? 'pending',
      subtotal: (json['subtotal'] as num?)?.toDouble() ?? 0.0,
      shipping: (json['shipping'] as num?)?.toDouble() ?? 0.0,
      tax: (json['tax'] as num?)?.toDouble() ?? 0.0,
      total: (json['total'] as num?)?.toDouble() ?? 0.0,
      items: (json['items'] as List?)
          ?.map((item) => CartItemDTO.fromJson(item as Map<String, dynamic>))
          .toList() ?? [],
      promoCode: json['promo_code'] as String?,
      promoDiscount: (json['promo_discount'] as num?)?.toDouble() ?? 0.0,
      createdAt: json['created_at'] != null 
          ? DateTime.parse(json['created_at'] as String)
          : DateTime.now(),
      updatedAt: json['updated_at'] != null
          ? DateTime.parse(json['updated_at'] as String)
          : DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'status': status,
      'subtotal': subtotal,
      'shipping': shipping,
      'tax': tax,
      'total': total,
      'items': items.map((item) => item.toJson()).toList(),
      'promo_code': promoCode,
      'promo_discount': promoDiscount,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}

class AddToCartRequest {
  final String productId;
  final int quantity;

  AddToCartRequest({required this.productId, required this.quantity});

  Map<String, dynamic> toJson() {
    // Ensure product_id is a clean string (no whitespace, proper UUID format)
    final cleanProductId = productId.trim();
    if (cleanProductId.isEmpty) {
      throw Exception('Product ID cannot be empty');
    }
    return {
      'product_id': cleanProductId,
      'quantity': quantity,
    };
  }
}

class UpdateCartItemRequest {
  final String cartItemId;
  final int quantity;

  UpdateCartItemRequest({required this.cartItemId, required this.quantity});

  Map<String, dynamic> toJson() {
    return {'cart_item_id': cartItemId, 'quantity': quantity};
  }
}

class ApplyPromoRequest {
  final String promoCode;

  ApplyPromoRequest({required this.promoCode});

  Map<String, dynamic> toJson() {
    return {'promo_code': promoCode};
  }
}
