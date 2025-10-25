// ไฟล์: product.dart

class Product {
  final int id;
  final String name;
  final double price;
  final String image;
  final String category;
  final String? description;
  final String? brand; // 👈 1. เพิ่มฟิลด์ brand

  Product({
    required this.id,
    required this.name,
    required this.price,
    required this.image,
    required this.category,
    this.description,
    this.brand, // 👈 2. เพิ่มใน constructor
  });
}
