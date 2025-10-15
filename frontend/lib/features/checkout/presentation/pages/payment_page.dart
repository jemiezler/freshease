import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/state/checkout_controller.dart';

class PaymentPage extends StatefulWidget {
  const PaymentPage({super.key});

  @override
  State<PaymentPage> createState() => _PaymentPageState();
}

class _PaymentPageState extends State<PaymentPage> {
  PaymentMethod? _method;
  final _last4 = TextEditingController();

  @override
  void dispose() {
    _last4.dispose();
    super.dispose();
  }

  void _submit() {
    if (_method == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please select a payment method')),
      );
      return;
    }
    if (_method == PaymentMethod.card) {
      final v = _last4.text.trim();
      if (!RegExp(r'^\d{4}$').hasMatch(v)) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Enter last 4 card digits')),
        );
        return;
      }
      CheckoutScope.of(context).setPayment(_method!, last4: v);
    } else {
      CheckoutScope.of(context).setPayment(_method!);
    }
    context.go('/cart/checkout/review');
  }

  @override
  Widget build(BuildContext context) {
    final saved = CheckoutScope.of(context).paymentMethod;
    _method ??= saved;

    final isCard = _method == PaymentMethod.card;

    return Scaffold(
      appBar: AppBar(title: const Text('Payment Method')),
      body: Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 900),
          child: ListView(
            padding: const EdgeInsets.all(16),
            children: [
              Card(
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: LayoutBuilder(
                    builder: (context, c) {
                      final isWide = c.maxWidth >= 700;
                      final column = Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Choose how you’d like to pay',
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.w700,
                            ),
                          ),
                          const SizedBox(height: 12),
                          RadioListTile<PaymentMethod>(
                            value: PaymentMethod.cod,
                            groupValue: _method,
                            onChanged: (m) => setState(() => _method = m),
                            title: const Text('Cash on Delivery'),
                            subtitle: const Text(
                              'Pay in cash when your order arrives',
                            ),
                          ),
                          const SizedBox(height: 8),
                          RadioListTile<PaymentMethod>(
                            value: PaymentMethod.card,
                            groupValue: _method,
                            onChanged: (m) => setState(() => _method = m),
                            title: const Text('Credit / Debit Card'),
                            subtitle: const Text(
                              'Demo mode – card not charged',
                            ),
                          ),
                          if (isCard) ...[
                            const SizedBox(height: 12),
                            const Text(
                              'Card (demo)',
                              style: TextStyle(fontWeight: FontWeight.w600),
                            ),
                            const SizedBox(height: 8),
                            Row(
                              children: [
                                Expanded(
                                  child: TextField(
                                    controller: _last4,
                                    keyboardType: TextInputType.number,
                                    decoration: const InputDecoration(
                                      labelText: 'Last 4 digits',
                                      hintText: '1234',
                                    ),
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ],
                      );
                      if (isWide) {
                        return Row(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Expanded(child: column),
                            const SizedBox(width: 16),
                            Expanded(child: _PaymentTips()),
                          ],
                        );
                      }
                      return column;
                    },
                  ),
                ),
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton(
                      onPressed: () => context.pop(),
                      child: const Text('Back'),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: FilledButton.icon(
                      onPressed: _submit,
                      icon: const Icon(Icons.chevron_right),
                      label: const Text('Continue'),
                      style: FilledButton.styleFrom(
                        minimumSize: const Size.fromHeight(48),
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _PaymentTips extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Card(
      color: Theme.of(context).colorScheme.surfaceContainerHighest,
      child: const Padding(
        padding: EdgeInsets.all(16),
        child: Text(
          'Tip: In demo mode, no real payments are processed. '
          'Selecting “Card” will only record the last 4 digits for display on the review page.',
          style: TextStyle(height: 1.3),
        ),
      ),
    );
  }
}
