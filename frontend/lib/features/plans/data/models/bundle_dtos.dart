class BundleDTO {
  final String id;
  final String name;
  final String? description;
  final double price;
  final bool isActive;

  BundleDTO({
    required this.id,
    required this.name,
    this.description,
    required this.price,
    required this.isActive,
  });

  factory BundleDTO.fromJson(Map<String, dynamic> json) {
    return BundleDTO(
      id: json['id'] as String,
      name: json['name'] as String,
      description: json['description'] as String?,
      price: (json['price'] as num).toDouble(),
      isActive: json['is_active'] as bool,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'description': description,
      'price': price,
      'is_active': isActive,
    };
  }
}

class BundleItemDTO {
  final String id;
  final int qty;
  final String bundleId;
  final String productId;

  BundleItemDTO({
    required this.id,
    required this.qty,
    required this.bundleId,
    required this.productId,
  });

  factory BundleItemDTO.fromJson(Map<String, dynamic> json) {
    return BundleItemDTO(
      id: json['id'] as String,
      qty: json['qty'] as int,
      bundleId: json['bundle_id'] as String,
      productId: json['product_id'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'qty': qty,
      'bundle_id': bundleId,
      'product_id': productId,
    };
  }
}

class BundleWithItemsDTO {
  final BundleDTO bundle;
  final List<BundleItemDTO> items;

  BundleWithItemsDTO({
    required this.bundle,
    required this.items,
  });
}

