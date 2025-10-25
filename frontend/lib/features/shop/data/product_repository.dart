import '../domain/product.dart';

abstract class ProductRepository {
  Future<List<Product>> list({
    String q = '',
    String category = 'All',
    double min = 0,
    double max = 99999,
  });
}

/// Mock repo – replace with HTTP later
class MockProductRepository implements ProductRepository {
  final _data = <Product>[
    // --- สินค้าที่มีอยู่เดิม ---
    Product(
      id: 1,
      name: 'Fresh Kale Bundle',
      price: 59,
      image: 'https://picsum.photos/400/300?1',
      category: 'Veggies',
      description: 'A bundle of fresh, locally sourced Kale.',
    ),
    Product(
      id: 2,
      name: 'Organic Mango Pack',
      price: 99,
      image: 'https://picsum.photos/400/300?2',
      category: 'Fruits',
      description: 'A pack of sweet and ripe organic mangoes.',
    ),
    Product(
      id: 3,
      name: 'Cherry Tomatoes',
      price: 49,
      image: 'https://picsum.photos/400/300?3',
      category: 'Veggies',
      description: 'Perfect for salads and snacking.',
    ),
    Product(
      id: 4,
      name: 'Thai Basil Leaves',
      price: 25,
      image: 'https://picsum.photos/400/300?4',
      category: 'Herbs',
      description: 'Fragrant and fresh Thai Basil leaves for cooking.',
    ),
    Product(
      id: 5,
      name: 'Avocado Set',
      price: 120,
      image: 'https://picsum.photos/400/300?5',
      category: 'Fruits',
      description: 'A set of 3 creamy Hass avocados.',
    ),

    // --- เพิ่มสินค้า Exclusive Offer ตามรูปภาพ UI ---
    Product(
      id: 6,
      name: 'Organic Bananas',
      price: 4.99,
      image: 'https://picsum.photos/400/300?6',
      category: 'Fruits',
      description: '7pcs, Priceg', // <--- เพิ่ม description
    ),
    Product(
      id: 7,
      name: 'Red Apple',
      price: 4.99,
      image: 'https://picsum.photos/400/300?7',
      category: 'Fruits',
      description: '1kg, Priceg', // <--- เพิ่ม description
    ),

    // --- เพิ่มสินค้า Best Selling ตามรูปภาพ UI ---
    Product(
      id: 8,
      name: 'Red Bell Pepper',
      price: 3.50,
      image: 'https://picsum.photos/400/300?8',
      category: 'Veggies',
      description: '3pcs, Priceg', // <--- เพิ่ม description
    ),
  ];

  @override
  Future<List<Product>> list({
    String q = '',
    String category = 'All',
    double min = 0,
    double max = 99999,
  }) async {
    await Future.delayed(const Duration(milliseconds: 120));
    return _data.where((p) {
      final byQ = q.isEmpty || p.name.toLowerCase().contains(q.toLowerCase());
      final byC = category == 'All' || p.category == category;
      final byP = p.price >= min && p.price <= max;
      return byQ && byC && byP;
    }).toList();
  }
}
