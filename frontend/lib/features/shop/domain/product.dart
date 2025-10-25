class Product {
  final int id;
  final String name;
  final double price;
  final String image;
  final String category;

  // 1. เพิ่มฟิลด์ description: ใช้ String? เพื่อให้เป็น optional (อนุญาตให้เป็น null ได้)
  final String? description;

  Product({
    required this.id,
    required this.name,
    required this.price,
    required this.image,
    required this.category,

    // 2. รับค่า description ใน Constructor
    this.description,
  });
}
