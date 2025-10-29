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
    return CartItemDTO(
      id: json['id'] as String,
      productId: json['product_id'] as String,
      productName: json['product_name'] as String,
      productImage: json['product_image'] as String,
      productPrice: (json['product_price'] as num).toDouble(),
      quantity: json['quantity'] as int,
      lineTotal: (json['line_total'] as num).toDouble(),
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
    return CartDTO(
      id: json['id'] as String,
      status: json['status'] as String,
      subtotal: (json['subtotal'] as num).toDouble(),
      shipping: (json['shipping'] as num).toDouble(),
      tax: (json['tax'] as num).toDouble(),
      total: (json['total'] as num).toDouble(),
      items: (json['items'] as List)
          .map((item) => CartItemDTO.fromJson(item as Map<String, dynamic>))
          .toList(),
      promoCode: json['promo_code'] as String?,
      promoDiscount: (json['promo_discount'] as num?)?.toDouble() ?? 0.0,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
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
    return {'product_id': productId, 'quantity': quantity};
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
