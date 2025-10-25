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
      image:
          'https://i.pinimg.com/1200x/cb/b7/19/cbb719ec3451b7e94f64ebc591c3fa2a.jpg',
      category: 'Veggies',
      brand: 'Kazi Farmas',
    ),
    Product(
      id: 2,
      name: 'Organic Mango Pack',
      price: 99,
      image:
          'https://i.pinimg.com/736x/d6/d0/85/d6d085ab76d5b69e8f3d4e8baa07150d.jpg',
      category: 'Fruits',
      brand: 'Malee',
    ),
    Product(
      id: 3,
      name: 'Cherry Tomatoes',
      price: 49,
      image:
          'https://i.pinimg.com/736x/99/2c/22/992c22f91ce2bc1070d05744680aca40.jpg',
      category: 'Veggies',
      brand: 'Individual Collection',
    ),
    Product(
      id: 4,
      name: 'Thai Basil Leaves',
      price: 25,
      image:
          'https://i.pinimg.com/736x/15/7b/31/157b31079d0a41445cdd88774423a9e6.jpg',
      category: 'Herbs',
    ),
    Product(
      id: 5,
      name: 'Avocado Set',
      price: 120,
      image:
          'https://i.pinimg.com/736x/55/76/7c/55767cd8be626988eb85d81d58a02010.jpg',
      category: 'Fruits',
      brand: 'Ifod',
    ),
    Product(
      id: 10,
      name: 'Egg Chicken Red',
      price: 1.99,
      image:
          'https://i.pinimg.com/736x/46/b8/7d/46b87dc1b5990413635bd9822e447240.jpg',
      category: 'Dairy',
      brand: 'Kazi Farmas',
    ),
    Product(
      id: 11,
      name: 'Egg Chicken White',
      price: 1.50,
      image:
          'https://i.pinimg.com/736x/78/ef/a4/78efa4deccca25ea4be54556c2e7c6fa.jpg',
      category: 'Dairy',
      brand: 'Individual Collection',
    ),
    Product(
      id: 12,
      name: 'Malee Tangerine Orange Juice',
      price: 65,
      image:
          'https://www.shopping-d.com/cdn/shop/products/MaleeTangerineOrangeJuicewithOrangePulpSize1L_8853333001815_17000.png?v=1594787116',
      category: 'Beverages',
      brand: 'Malee',
    ),
    Product(
      id: 13,
      name: 'Malee Peach Juice',
      price: 65,
      image:
          'https://media-stark.gourmetmarketthailand.com/products/cover/8853333013597-1.webp',
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
