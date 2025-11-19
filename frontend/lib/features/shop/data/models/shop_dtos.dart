class ShopProductDTO {
  final String id;
  final String name;
  final double price;
  final String description;
  final String imageUrl;
  final String unitLabel;
  final String isActive;
  final DateTime createdAt;
  final DateTime updatedAt;
  final String vendorId;
  final String vendorName;
  final String categoryId;
  final String categoryName;
  final int stockQuantity;
  final bool isInStock;

  ShopProductDTO({
    required this.id,
    required this.name,
    required this.price,
    required this.description,
    required this.imageUrl,
    required this.unitLabel,
    required this.isActive,
    required this.createdAt,
    required this.updatedAt,
    required this.vendorId,
    required this.vendorName,
    required this.categoryId,
    required this.categoryName,
    required this.stockQuantity,
    required this.isInStock,
  });

  factory ShopProductDTO.fromJson(Map<String, dynamic> json) {
    return ShopProductDTO(
      id: json['id'] as String,
      name: json['name'] as String,
      price: (json['price'] as num).toDouble(),
      description: json['description'] as String,
      imageUrl: json['image_url'] as String,
      unitLabel: json['unit_label'] as String,
      isActive: json['is_active'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
      vendorId: json['vendor_id'] as String,
      vendorName: json['vendor_name'] as String,
      categoryId: json['category_id'] as String,
      categoryName: json['category_name'] as String,
      stockQuantity: json['stock_quantity'] as int,
      isInStock: json['is_in_stock'] as bool,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'price': price,
      'description': description,
      'image_url': imageUrl,
      'unit_label': unitLabel,
      'is_active': isActive,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
      'vendor_id': vendorId,
      'vendor_name': vendorName,
      'category_id': categoryId,
      'category_name': categoryName,
      'stock_quantity': stockQuantity,
      'is_in_stock': isInStock,
    };
  }
}

class ShopCategoryDTO {
  final String id;
  final String name;
  final String description;

  ShopCategoryDTO({
    required this.id,
    required this.name,
    required this.description,
  });

  factory ShopCategoryDTO.fromJson(Map<String, dynamic> json) {
    return ShopCategoryDTO(
      id: json['id'] as String,
      name: json['name'] as String,
      description: json['description'] as String? ?? json['slug'] as String? ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {'id': id, 'name': name, 'description': description};
  }
}

class ShopVendorDTO {
  final String id;
  final String name;
  final String email;
  final String phone;
  final String address;
  final String city;
  final String state;
  final String country;
  final String postalCode;
  final String website;
  final String logoUrl;
  final String description;
  final String isActive;
  final DateTime createdAt;
  final DateTime updatedAt;

  ShopVendorDTO({
    required this.id,
    required this.name,
    required this.email,
    required this.phone,
    required this.address,
    required this.city,
    required this.state,
    required this.country,
    required this.postalCode,
    required this.website,
    required this.logoUrl,
    required this.description,
    required this.isActive,
    required this.createdAt,
    required this.updatedAt,
  });

  factory ShopVendorDTO.fromJson(Map<String, dynamic> json) {
    return ShopVendorDTO(
      id: json['id'] as String,
      name: json['name'] as String,
      email: json['email'] as String,
      phone: json['phone'] as String,
      address: json['address'] as String,
      city: json['city'] as String,
      state: json['state'] as String,
      country: json['country'] as String,
      postalCode: json['postal_code'] as String,
      website: json['website'] as String,
      logoUrl: json['logo_url'] as String,
      description: json['description'] as String,
      isActive: json['is_active'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'email': email,
      'phone': phone,
      'address': address,
      'city': city,
      'state': state,
      'country': country,
      'postal_code': postalCode,
      'website': website,
      'logo_url': logoUrl,
      'description': description,
      'is_active': isActive,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}

class ShopSearchResponse {
  final List<ShopProductDTO> products;
  final int total;
  final int limit;
  final int offset;
  final bool hasMore;

  ShopSearchResponse({
    required this.products,
    required this.total,
    required this.limit,
    required this.offset,
    required this.hasMore,
  });

  factory ShopSearchResponse.fromJson(Map<String, dynamic> json) {
    return ShopSearchResponse(
      products: (json['products'] as List)
          .map((item) => ShopProductDTO.fromJson(item as Map<String, dynamic>))
          .toList(),
      total: json['total'] as int,
      limit: json['limit'] as int,
      offset: json['offset'] as int,
      hasMore: json['has_more'] as bool,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'products': products.map((item) => item.toJson()).toList(),
      'total': total,
      'limit': limit,
      'offset': offset,
      'has_more': hasMore,
    };
  }
}

class ShopSearchFilters {
  final String? categoryId;
  final String? vendorId;
  final double? minPrice;
  final double? maxPrice;
  final String? searchTerm;
  final bool? inStock;
  final int limit;
  final int offset;

  ShopSearchFilters({
    this.categoryId,
    this.vendorId,
    this.minPrice,
    this.maxPrice,
    this.searchTerm,
    this.inStock,
    this.limit = 20,
    this.offset = 0,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};

    if (categoryId != null) params['category_id'] = categoryId;
    if (vendorId != null) params['vendor_id'] = vendorId;
    if (minPrice != null) params['min_price'] = minPrice;
    if (maxPrice != null) params['max_price'] = maxPrice;
    if (searchTerm != null && searchTerm!.isNotEmpty)
      // ignore: curly_braces_in_flow_control_structures
      params['search'] = searchTerm;
    if (inStock != null) params['in_stock'] = inStock;
    if (limit != 20) params['limit'] = limit;
    if (offset != 0) params['offset'] = offset;

    return params;
  }
}
