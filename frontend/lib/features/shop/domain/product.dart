class Product {
  final String id;
  final String name;
  final double price;
  final String image;
  final String category;
  final String description;
  final String unitLabel;
  final String vendorId;
  final String vendorName;
  final String categoryId;
  final String categoryName;
  final int stockQuantity;
  final bool isInStock;

  Product({
    required this.id,
    required this.name,
    required this.price,
    required this.image,
    required this.category,
    required this.description,
    required this.unitLabel,
    required this.vendorId,
    required this.vendorName,
    required this.categoryId,
    required this.categoryName,
    required this.stockQuantity,
    required this.isInStock,
  });

  // Factory constructor for backward compatibility with existing mock data
  factory Product.fromLegacy({
    required int id,
    required String name,
    required double price,
    required String image,
    required String category,
  }) {
    return Product(
      id: id.toString(),
      name: name,
      price: price,
      image: image,
      category: category,
      description: '',
      unitLabel: 'kg',
      vendorId: '',
      vendorName: '',
      categoryId: '',
      categoryName: category,
      stockQuantity: 100,
      isInStock: true,
    );
  }

  // Convert from API DTO
  factory Product.fromShopDTO(dynamic shopProductDTO) {
    return Product(
      id: shopProductDTO.id,
      name: shopProductDTO.name,
      price: shopProductDTO.price,
      image: shopProductDTO.imageUrl,
      category: shopProductDTO.categoryName,
      description: shopProductDTO.description,
      unitLabel: shopProductDTO.unitLabel,
      vendorId: shopProductDTO.vendorId,
      vendorName: shopProductDTO.vendorName,
      categoryId: shopProductDTO.categoryId,
      categoryName: shopProductDTO.categoryName,
      stockQuantity: shopProductDTO.stockQuantity,
      isInStock: shopProductDTO.isInStock,
    );
  }

  Product copyWith({
    String? id,
    String? name,
    double? price,
    String? image,
    String? category,
    String? description,
    String? unitLabel,
    String? vendorId,
    String? vendorName,
    String? categoryId,
    String? categoryName,
    int? stockQuantity,
    bool? isInStock,
  }) {
    return Product(
      id: id ?? this.id,
      name: name ?? this.name,
      price: price ?? this.price,
      image: image ?? this.image,
      category: category ?? this.category,
      description: description ?? this.description,
      unitLabel: unitLabel ?? this.unitLabel,
      vendorId: vendorId ?? this.vendorId,
      vendorName: vendorName ?? this.vendorName,
      categoryId: categoryId ?? this.categoryId,
      categoryName: categoryName ?? this.categoryName,
      stockQuantity: stockQuantity ?? this.stockQuantity,
      isInStock: isInStock ?? this.isInStock,
    );
  }
}
