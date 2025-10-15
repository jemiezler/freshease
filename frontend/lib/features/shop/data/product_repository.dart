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
    Product(
      id: 1,
      name: 'Fresh Kale Bundle',
      price: 59,
      image: 'https://picsum.photos/400/300?1',
      category: 'Veggies',
    ),
    Product(
      id: 2,
      name: 'Organic Mango Pack',
      price: 99,
      image: 'https://picsum.photos/400/300?2',
      category: 'Fruits',
    ),
    Product(
      id: 3,
      name: 'Cherry Tomatoes',
      price: 49,
      image: 'https://picsum.photos/400/300?3',
      category: 'Veggies',
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
