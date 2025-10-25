// ไฟล์: product_repository.dart
import '../domain/product.dart';

abstract class ProductRepository {
  Future<List<Product>> list({
    String q = '',
    // ‼️ 1. ลบ 'category' ออกจาก interface
    double min = 0,
    double max = 99999,
    Map<String, List<String>> filters = const {},
  });
}

/// Mock repo – replace with HTTP later
class MockProductRepository implements ProductRepository {
  final _data = <Product>[
    Product(
      id: 1,
      name: 'Fresh Kale Bundle',
      price: 59,
      image: 'https://picsum.photos/400/300?1',
      category: 'Veggies',
      brand: 'Kazi Farmas',
    ),
    Product(
      id: 2,
      name: 'Organic Mango Pack',
      price: 99,
      image: 'https://picsum.photos/400/300?2',
      category: 'Fruits',
      brand: 'Malee',
    ),
    Product(
      id: 3,
      name: 'Cherry Tomatoes',
      price: 49,
      image: 'https://picsum.photos/400/300?3',
      category: 'Veggies',
      brand: 'Individual Collection',
    ),
    Product(
      id: 4,
      name: 'Thai Basil Leaves',
      price: 25,
      image: 'https://picsum.photos/400/300?4',
      category: 'Herbs',
    ),
    Product(
      id: 5,
      name: 'Avocado Set',
      price: 120,
      image: 'https://picsum.photos/400/300?5',
      category: 'Fruits',
      brand: 'Ifod',
    ),
    Product(
      id: 10,
      name: 'Egg Chicken Red',
      price: 1.99,
      image: 'https://picsum.photos/400/300?10',
      category: 'Dairy',
      brand: 'Kazi Farmas',
    ),
    Product(
      id: 11,
      name: 'Egg Chicken White',
      price: 1.50,
      image: 'https://picsum.photos/400/300?11',
      category: 'Dairy',
      brand: 'Individual Collection',
    ),
    Product(
      id: 12,
      name: 'Malee Tangerine Orange Juice',
      price: 65,
      image: 'https://picsum.photos/400/300?12',
      category: 'Beverages',
      brand: 'Malee',
    ),
    Product(
      id: 13,
      name: 'Malee Peach Juice',
      price: 65,
      image: 'https://picsum.photos/400/300?13',
      category: 'Beverages',
      brand: 'Malee',
    ),
  ];

  @override
  Future<List<Product>> list({
    String q = '',
    // ‼️ 2. ลบ 'category' ออกจาก parameter
    double min = 0,
    double max = 99999,
    Map<String, List<String>> filters = const {},
  }) async {
    await Future.delayed(const Duration(milliseconds: 120));

    // 3. ดึงค่าฟิลเตอร์ที่เลือก
    final filterCategories = filters['categories'] ?? [];
    final filterBrands = filters['brands'] ?? [];

    return _data.where((p) {
      // ‼️ 4. FIX: ค้นหาจาก Search Bar (q)
      // ให้ค้นหาทั้ง p.name และ p.category
      final byQ =
          q.isEmpty ||
          p.name.toLowerCase().contains(q.toLowerCase()) ||
          p.category.toLowerCase().contains(q.toLowerCase());

      // ‼️ 5. FIX: ลบ byC (ตัวปัญหา) ทิ้ง
      // final byC = category == 'All' || p.category == category; // 👈 ลบทิ้ง

      // 6. ฟิลเตอร์จาก Price Range (ถ้ามี)
      final byP = p.price >= min && p.price <= max;

      // ‼️ 7. FIX: ฟิลเตอร์จากหน้า FilterPage (Categories Checkbox)
      // แก้ไขให้เช็ค p.category ตรงๆ (จาก 'Fruits', 'Dairy' ฯลฯ)
      final byFilterCategory =
          filterCategories.isEmpty || filterCategories.contains(p.category);

      // 8. ฟิลเตอร์จากหน้า FilterPage (Brands Checkbox)
      final byFilterBrand =
          filterBrands.isEmpty ||
          (p.brand != null && filterBrands.contains(p.brand));

      // 9. คืนค่าสินค้าที่ตรงทุกเงื่อนไข (ลบ byC ออก)
      return byQ && byP && byFilterCategory && byFilterBrand;
    }).toList();
  }
}
