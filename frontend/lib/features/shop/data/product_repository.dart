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
      name: 'Organic Mixed Salad Bowl',
      price: 89,
      image:
          'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSKi4F589upoZSmqvnpdO2RsgAhWFS45WFe4A&s',
      category: 'Salad',
    ),
    Product(
      id: 2,
      name: 'Organic Romaine Lettuce',
      price: 45,
      image:
          'https://farmlinkhawaii.com/cdn/shop/files/organic_20romaine-800x533.jpg?v=1724305179',
      category: 'Veggies',
    ),
    Product(
      id: 3,
      name: 'Fresh Organic Strawberry Mix',
      price: 120,
      image:
          'https://img.freepik.com/premium-photo/fresh-organic-summer-mix-strawberry-sweet-cherry-round-stone-plate-dark-textured-surface_114309-2179.jpg',
      category: 'Fruits',
    ),
    Product(
      id: 4,
      name: 'Avocado & Spinach Raw Set',
      price: 135,
      image:
          'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRL_PFJxbihQiBW0O_kxPVEf94uDtN4-gBJBA&s',
      category: 'Raw Material',
    ),
    Product(
      id: 5,
      name: 'Tropical Fruit Fusion Pack',
      price: 149,
      image:
          'https://www.packd.co.uk/cdn/shop/files/PACKD-organic-tropical-fruit_a2f8ead5-e78b-4d0c-9e72-bd57871b2b05_800x.png?v=1751978174',
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
