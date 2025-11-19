// lib/features/cart/presentation/pages/cart_page.dart
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/state/cart_controller.dart';

class CartPage extends StatelessWidget {
  const CartPage({super.key});

  @override
  Widget build(BuildContext context) {
    final cart = CartScope.of(context);
    return AnimatedBuilder(
      animation: cart,
      builder: (_, _) {
        return Scaffold(
          backgroundColor: Theme.of(context).scaffoldBackgroundColor,
          appBar: AppBar(
            backgroundColor: Theme.of(context).colorScheme.primary,
            title: const Text(
              'Your Cart',
              style: TextStyle(color: Colors.white),
            ),
            actions: [
              if (cart.lines.isNotEmpty)
                IconButton(
                  tooltip: 'Clear cart',
                  icon: cart.isLoading
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Icon(Icons.delete_sweep_outlined),
                  onPressed: cart.isLoading ? null : cart.clear,
                ),
            ],
          ),
          body: cart.isLoading && cart.lines.isEmpty
              ? const Center(child: CircularProgressIndicator())
              : cart.error != null
              ? _ErrorState(error: cart.error!, onRetry: cart.refresh)
              : LayoutBuilder(
                  builder: (context, c) {
                    final isWide = c.maxWidth >= 900;
                    if (cart.lines.isEmpty) {
                      return _EmptyState(onShop: () => context.go('/'));
                    }

                    final list = _CartList();
                    final summary = _SummaryCard();

                    if (isWide) {
                      return Padding(
                        padding: const EdgeInsets.all(16),
                        child: Row(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            // List
                            Expanded(
                              flex: 7,
                              child: Card(
                                clipBehavior: Clip.antiAlias,
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(16),
                                ),
                                child: list,
                              ),
                            ),
                            const SizedBox(width: 16),
                            // Summary
                            Expanded(flex: 4, child: summary),
                          ],
                        ),
                      );
                    } else {
                      // Mobile: list + sticky summary
                      return Column(
                        children: [
                          Expanded(child: list),
                          const Divider(height: 1),
                          SafeArea(child: summary),
                        ],
                      );
                    }
                  },
                ),
        );
      },
    );
  }
}

/// List of cart lines with swipe-to-delete and quantity stepper
class _CartList extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final cart = CartScope.of(context);
    return ListView.separated(
      padding: const EdgeInsets.all(16),
      itemCount: cart.lines.length,
      separatorBuilder: (_, _) => const Divider(height: 24),
      itemBuilder: (_, i) {
        final line = cart.lines[i];
        return Dismissible(
          key: ValueKey('cart-${line.product.id}'),
          direction: DismissDirection.endToStart,
          background: Container(
            alignment: Alignment.centerRight,
            padding: const EdgeInsets.symmetric(horizontal: 24),
            decoration: BoxDecoration(
              color: Colors.red.withValues(alpha: 0.12),
              borderRadius: BorderRadius.circular(12),
            ),
            child: const Icon(Icons.delete_outline, color: Colors.red),
          ),
          onDismissed: (_) => cart.remove(line.product),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // image
              ClipRRect(
                borderRadius: BorderRadius.circular(12),
                child: Image.network(
                  line.product.image,
                  width: 84,
                  height: 84,
                  fit: BoxFit.cover,
                ),
              ),
              const SizedBox(width: 12),
              // info
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      line.product.name,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                        fontWeight: FontWeight.w700,
                        fontSize: 16,
                      ),
                    ),
                    const SizedBox(height: 6),
                    Text(
                      '฿${line.product.price.toStringAsFixed(2)}',
                      style: const TextStyle(
                        color: Colors.green,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(height: 12),
                    Row(
                      children: [
                        _QtyButton(
                          icon: Icons.remove,
                          onTap: cart.isLoading
                              ? () {}
                              : () => cart.decrement(line.product),
                        ),
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 12),
                          child: Text(
                            '${line.qty}',
                            style: const TextStyle(fontWeight: FontWeight.w700),
                          ),
                        ),
                        _QtyButton(
                          icon: Icons.add,
                          onTap: cart.isLoading
                              ? () {}
                              : () async {
                                  await cart.add(line.product);
                                  if (cart.error != null && context.mounted) {
                                    ScaffoldMessenger.of(context).showSnackBar(
                                      SnackBar(
                                        content: Text('Failed to update: ${cart.error}'),
                                        backgroundColor: Colors.red,
                                        duration: const Duration(seconds: 2),
                                      ),
                                    );
                                  }
                                },
                        ),
                        const Spacer(),
                        Text(
                          '฿${line.lineTotal.toStringAsFixed(2)}',
                          style: const TextStyle(
                            fontWeight: FontWeight.w700,
                            fontSize: 16,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
              // remove icon (alternative to swipe)
              IconButton(
                tooltip: 'Remove',
                icon: const Icon(Icons.close),
                onPressed: cart.isLoading
                    ? null
                    : () => cart.remove(line.product),
              ),
            ],
          ),
        );
      },
    );
  }
}

class _QtyButton extends StatelessWidget {
  final IconData icon;
  final VoidCallback onTap;
  const _QtyButton({required this.icon, required this.onTap});
  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: 34,
      height: 34,
      child: OutlinedButton(
        style: OutlinedButton.styleFrom(padding: EdgeInsets.zero),
        onPressed: onTap,
        child: Icon(icon, size: 18),
      ),
    );
  }
}

/// Summary card with promo code + totals + checkout button
class _SummaryCard extends StatefulWidget {
  @override
  State<_SummaryCard> createState() => _SummaryCardState();
}

class _SummaryCardState extends State<_SummaryCard> {
  final _promoCtrl = TextEditingController();
  String? _error;

  @override
  void dispose() {
    _promoCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final cart = CartScope.of(context);
    return Card(
      margin: const EdgeInsets.all(16),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            // promo
            Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _promoCtrl,
                    decoration: InputDecoration(
                      labelText: 'Promo code',
                      hintText: 'FRESH10 or FREESHIP',
                      errorText: _error,
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                if (cart.promoCode == null)
                  FilledButton(
                    onPressed: cart.isLoading
                        ? null
                        : () async {
                            setState(() => _error = null);
                            final ok = await cart.applyPromo(_promoCtrl.text);
                            if (!ok) setState(() => _error = 'Invalid code');
                          },
                    child: const Text('Apply'),
                  )
                else
                  OutlinedButton(
                    onPressed: cart.isLoading ? null : cart.removePromo,
                    child: Text('Remove (${cart.promoCode})'),
                  ),
              ],
            ),
            const SizedBox(height: 16),

            _Row('Subtotal', cart.subtotal),
            _Row('Shipping', cart.shipping),
            if (cart.promoDiscount > 0)
              _Row('Discount', -cart.promoDiscount, isAccent: true),
            _Row('VAT (7%)', cart.vat),
            const Divider(height: 24),
            _Row('Total', cart.total, isBold: true, big: true),

            const SizedBox(height: 16),
            SizedBox(
              width: double.infinity,
              child: FilledButton.icon(
                onPressed: cart.count == 0
                    ? null
                    : () => context.go('/cart/checkout/address'),
                icon: const Icon(Icons.lock_outline),
                label: const Text('Checkout'),
                style: FilledButton.styleFrom(
                  minimumSize: const Size.fromHeight(48),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _Row extends StatelessWidget {
  final String label;
  final double amount;
  final bool isBold;
  final bool big;
  final bool isAccent;
  const _Row(
    this.label,
    this.amount, {
    this.isBold = false,
    this.big = false,
    this.isAccent = false,
  });

  @override
  Widget build(BuildContext context) {
    final style = TextStyle(
      fontWeight: isBold ? FontWeight.w800 : FontWeight.w600,
      fontSize: big ? 18 : 14,
      color: isAccent ? Colors.green : null,
    );
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        children: [
          Text(label, style: style),
          const Spacer(),
          Text('฿${amount.toStringAsFixed(2)}', style: style),
        ],
      ),
    );
  }
}

class _EmptyState extends StatelessWidget {
  final VoidCallback onShop;
  const _EmptyState({required this.onShop});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(
              Icons.shopping_cart_outlined,
              size: 80,
              color: Colors.grey,
            ),
            const SizedBox(height: 12),
            const Text(
              'Your cart is empty',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.w700),
            ),
            const SizedBox(height: 8),
            const Text('Find fresh picks in the shop and add them here.'),
            const SizedBox(height: 16),
            FilledButton.icon(
              onPressed: onShop,
              icon: const Icon(Icons.storefront),
              label: const Text('Go to Shop'),
            ),
          ],
        ),
      ),
    );
  }
}

class _ErrorState extends StatelessWidget {
  final String error;
  final VoidCallback onRetry;
  const _ErrorState({required this.error, required this.onRetry});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline, size: 80, color: Colors.red),
            const SizedBox(height: 12),
            const Text(
              'Something went wrong',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.w700),
            ),
            const SizedBox(height: 8),
            Text(
              error,
              textAlign: TextAlign.center,
              style: const TextStyle(color: Colors.grey),
            ),
            const SizedBox(height: 16),
            FilledButton.icon(
              onPressed: onRetry,
              icon: const Icon(Icons.refresh),
              label: const Text('Try Again'),
            ),
          ],
        ),
      ),
    );
  }
}
