import 'package:flutter/material.dart';
import 'package:frontend/core/widgets/global_appbar.dart';
import '../../../../core/state/cart_controller.dart';
import '../../domain/product.dart';

class ProductDetailPage extends StatelessWidget {
  final Product product;
  const ProductDetailPage({super.key, required this.product});

  @override
  Widget build(BuildContext context) {
    final cart = CartScope.of(context);

    return Scaffold(
      appBar: GlobalAppBar(title: product.name),
      body: LayoutBuilder(
        builder: (context, c) {
          final isWide = c.maxWidth > 700;
          final image = ClipRRect(
            borderRadius: BorderRadius.circular(16),
            child: Image.network(
              product.image,
              fit: BoxFit.cover,
              height: isWide ? 400 : 250,
              width: double.infinity,
            ),
          );

          final details = Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                product.name,
                style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                '฿${product.price.toStringAsFixed(2)}',
                style: const TextStyle(
                  fontSize: 22,
                  fontWeight: FontWeight.w700,
                  color: Colors.green,
                ),
              ),
              const SizedBox(height: 16),
              Text(
                'Category: ${product.category}',
                style: const TextStyle(color: Colors.grey),
              ),
              const SizedBox(height: 16),
              const Text(
                'Description',
                style: TextStyle(fontWeight: FontWeight.w600, fontSize: 16),
              ),
              const SizedBox(height: 8),
              Text(
                'Fresh and locally sourced ${product.name}. Packed daily for maximum freshness and nutrition.',
                style: const TextStyle(fontSize: 15, height: 1.4),
              ),
              const SizedBox(height: 24),
              FilledButton.icon(
                onPressed: () {
                  cart.add(product);
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('${product.name} added to cart')),
                  );
                },
                icon: const Icon(Icons.add_shopping_cart),
                label: const Text('Add to Cart'),
                style: FilledButton.styleFrom(
                  minimumSize: const Size(double.infinity, 48),
                ),
              ),
            ],
          );

          if (isWide) {
            return Padding(
              padding: const EdgeInsets.all(24),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Expanded(flex: 3, child: image),
                  const SizedBox(width: 24),
                  Expanded(flex: 4, child: details),
                ],
              ),
            );
          } else {
            return ListView(
              padding: const EdgeInsets.all(16),
              children: [image, const SizedBox(height: 20), details],
            );
          }
        },
      ),
    );
  }
}
