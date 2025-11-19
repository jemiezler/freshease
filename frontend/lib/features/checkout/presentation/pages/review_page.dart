import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/state/cart_controller.dart';
import '../../../../core/state/checkout_controller.dart';

class ReviewPage extends StatelessWidget {
  const ReviewPage({super.key});

  @override
  Widget build(BuildContext context) {
    final cart = CartScope.of(context);
    final co = CheckoutScope.of(context);

    return AnimatedBuilder(
      animation: Listenable.merge([cart, co]),
      builder: (_, _) {
        final isPlan = co.isPlanCheckout;
        final ready = co.isReadyForReview && (isPlan || cart.count > 0);
        return Scaffold(
          appBar: AppBar(title: const Text('Review & Confirm')),
          body: LayoutBuilder(
            builder: (context, c) {
              final isWide = c.maxWidth >= 1000;

              final addressCard = _SectionCard(
                title: 'Shipping Address',
                action: TextButton(
                  onPressed: () => context.go('/cart/checkout/address'),
                  child: const Text('Edit'),
                ),
                child: co.shippingAddress == null
                    ? const Text('No address set')
                    : _AddressView(addr: co.shippingAddress!),
              );

              final paymentCard = _SectionCard(
                title: 'Payment',
                action: TextButton(
                  onPressed: () => context.go('/cart/checkout/payment'),
                  child: const Text('Edit'),
                ),
                child: _PaymentView(
                  method: co.paymentMethod,
                  last4: co.paymentLast4,
                ),
              );

              // Use selected plan if present, else cart items
              final itemsCard = _SectionCard(
                title: co.isPlanCheckout ? 'Subscription Plan' : 'Items',
                child: isPlan
                    ? _PlanSummaryBlock(
                        title: co.selectedPlan!.title,
                        subtitle: co.selectedPlan!.subtitle ?? '',
                        price: co.selectedPlan!.price,
                      )
                    : co.isPlanCheckout
                    ? Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            co.selectedPlan!.title,
                            style: const TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.w700,
                            ),
                          ),
                          if (co.selectedPlan!.subtitle != null) ...[
                            const SizedBox(height: 6),
                            Text(
                              co.selectedPlan!.subtitle!,
                              style: const TextStyle(color: Colors.grey),
                            ),
                          ],
                          const SizedBox(height: 12),
                          Row(
                            children: [
                              const Text(
                                'Price:',
                                style: TextStyle(
                                  fontWeight: FontWeight.w600,
                                  fontSize: 15,
                                ),
                              ),
                              const Spacer(),
                              Text(
                                '฿${co.selectedPlan!.price.toStringAsFixed(2)}',
                                style: const TextStyle(
                                  fontWeight: FontWeight.w700,
                                  color: Colors.green,
                                ),
                              ),
                            ],
                          ),
                        ],
                      )
                    : Column(
                        children: [
                          for (final line in cart.lines)
                            Padding(
                              padding: const EdgeInsets.symmetric(vertical: 8),
                              child: Row(
                                children: [
                                  SizedBox(
                                    width: 48,
                                    height: 48,
                                    child: ClipRRect(
                                      borderRadius: BorderRadius.circular(8),
                                      child: Image.network(
                                        line.product.image,
                                        fit: BoxFit.cover,
                                      ),
                                    ),
                                  ),
                                  const SizedBox(width: 12),
                                  Expanded(
                                    child: Text(
                                      '${line.product.name} × ${line.qty}',
                                      maxLines: 2,
                                      overflow: TextOverflow.ellipsis,
                                    ),
                                  ),
                                  Text(
                                    '฿${line.lineTotal.toStringAsFixed(2)}',
                                    style: const TextStyle(
                                      fontWeight: FontWeight.w700,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          const Divider(height: 24),
                          _Line('Subtotal', cart.subtotal),
                          _Line('Shipping', cart.shipping),
                          if (cart.promoDiscount > 0)
                            _Line(
                              'Discount',
                              -cart.promoDiscount,
                              accent: true,
                            ),
                          _Line('VAT (7%)', cart.vat),
                          const Divider(height: 24),
                          _Line('Total', cart.total, bold: true, big: true),
                        ],
                      ),
              );

              final placeOrder = SizedBox(
                width: double.infinity,
                child: FilledButton.icon(
                  icon: const Icon(Icons.check_circle_outline),
                  label: const Text('Place Order'),
                  onPressed: !ready
                      ? null
                      : () async {
                          final result = await co.placeOrder();
                          // clear both controllers
                          cart.clear();
                          co.clear();
                          ScaffoldMessenger.of(context).showSnackBar(
                            SnackBar(
                              content: Text(
                                co.isPlanCheckout
                                    ? 'Subscription activated!'
                                    : 'Order ${result.orderId} placed successfully!',
                              ),
                            ),
                          );
                          context.go(
                            '/checkout/confirmation?id=${result.orderId}',
                          );
                        },
                  style: FilledButton.styleFrom(
                    minimumSize: const Size.fromHeight(48),
                  ),
                ),
              );

              if (isWide) {
                return Center(
                  child: ConstrainedBox(
                    constraints: const BoxConstraints(maxWidth: 1100),
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Row(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Expanded(
                            flex: 7,
                            child: Column(
                              children: [addressCard, paymentCard, itemsCard],
                            ),
                          ),
                          const SizedBox(width: 16),
                          Expanded(
                            flex: 4,
                            child: Card(
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(16),
                              ),
                              child: Padding(
                                padding: const EdgeInsets.all(16),
                                child: Column(
                                  mainAxisSize: MainAxisSize.min,
                                  children: [
                                    const Text(
                                      'Almost there!',
                                      style: TextStyle(
                                        fontSize: 18,
                                        fontWeight: FontWeight.w700,
                                      ),
                                    ),
                                    const SizedBox(height: 12),
                                    placeOrder,
                                  ],
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                );
              }

              // narrow layout
              return ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  addressCard,
                  paymentCard,
                  itemsCard,
                  const SizedBox(height: 12),
                  placeOrder,
                ],
              );
            },
          ),
        );
      },
    );
  }
}

class _AddressView extends StatelessWidget {
  final Address addr;
  const _AddressView({required this.addr});

  @override
  Widget build(BuildContext context) {
    return Text(
      '${addr.fullName}\n${addr.line1}${addr.line2 != null ? '\n${addr.line2}' : ''}\n'
      '${addr.subDistrict}, ${addr.district}, ${addr.province} ${addr.postalCode}\n'
      '☎ ${addr.phone}',
      style: const TextStyle(height: 1.35),
    );
  }
}

class _PaymentView extends StatelessWidget {
  final PaymentMethod? method;
  final String? last4;
  const _PaymentView({required this.method, required this.last4});

  @override
  Widget build(BuildContext context) {
    if (method == PaymentMethod.cod) {
      return const Text('Cash on Delivery');
    } else if (method == PaymentMethod.card) {
      return Text('Card •••• ${last4 ?? '????'} (demo)');
    }
    return const Text('No payment selected');
  }
}

class _SectionCard extends StatelessWidget {
  final String title;
  final Widget child;
  final Widget? action;
  const _SectionCard({required this.title, required this.child, this.action});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Text(
                  title,
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w700,
                  ),
                ),
                const Spacer(),
                if (action != null) action!,
              ],
            ),
            const SizedBox(height: 12),
            child,
          ],
        ),
      ),
    );
  }
}

class _Line extends StatelessWidget {
  final String label;
  final double amount;
  final bool bold;
  final bool big;
  final bool accent;
  const _Line(
    this.label,
    this.amount, {
    this.bold = false,
    this.big = false,
    this.accent = false,
  });

  @override
  Widget build(BuildContext context) {
    final style = TextStyle(
      fontWeight: bold ? FontWeight.w800 : FontWeight.w600,
      fontSize: big ? 18 : 14,
      color: accent ? Colors.green : null,
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

class _PlanSummaryBlock extends StatelessWidget {
  final String title;
  final String subtitle;
  final double price;
  const _PlanSummaryBlock({
    required this.title,
    required this.subtitle,
    required this.price,
  });

  @override
  Widget build(BuildContext context) {
    // Plan pricing model: no shipping, no promo, VAT 7%
    final vat = price * 0.07;
    final total = price + vat;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          title,
          style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w700),
        ),
        if (subtitle.isNotEmpty) ...[
          const SizedBox(height: 6),
          Text(subtitle, style: const TextStyle(color: Colors.grey)),
        ],
        const SizedBox(height: 12),
        const Divider(height: 24),
        _Line('Subtotal', price),
        _Line('VAT (7%)', vat),
        const Divider(height: 24),
        _Line('Total', total, bold: true, big: true),
      ],
    );
  }
}
