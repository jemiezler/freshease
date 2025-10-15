import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/state/checkout_controller.dart';

class AddressPage extends StatefulWidget {
  const AddressPage({super.key});

  @override
  State<AddressPage> createState() => _AddressPageState();
}

class _AddressPageState extends State<AddressPage> {
  final _formKey = GlobalKey<FormState>();

  final _fullName = TextEditingController();
  final _phone = TextEditingController();
  final _line1 = TextEditingController();
  final _line2 = TextEditingController();
  final _subDistrict = TextEditingController();
  final _district = TextEditingController();
  final _province = TextEditingController();
  final _postal = TextEditingController();

  @override
  void dispose() {
    _fullName.dispose();
    _phone.dispose();
    _line1.dispose();
    _line2.dispose();
    _subDistrict.dispose();
    _district.dispose();
    _province.dispose();
    _postal.dispose();
    super.dispose();
  }

  String? _req(String? v, {String label = 'This field'}) =>
      (v == null || v.trim().isEmpty) ? '$label is required' : null;

  String? _phoneRule(String? v) {
    if (v == null || v.trim().isEmpty) return 'Phone is required';
    final s = v.replaceAll(RegExp(r'[\s-]'), '');
    if (!RegExp(r'^\+?\d{9,15}$').hasMatch(s)) return 'Invalid phone number';
    return null;
  }

  String? _postalRule(String? v) {
    if (v == null || v.trim().isEmpty) return 'Postal code is required';
    if (!RegExp(r'^\d{5}$').hasMatch(v)) return 'Use 5 digits (e.g., 10230)';
    return null;
  }

  void _submit() {
    if (!_formKey.currentState!.validate()) return;

    final checkout = CheckoutScope.of(context);
    checkout.setShipping(
      Address(
        fullName: _fullName.text.trim(),
        phone: _phone.text.trim(),
        line1: _line1.text.trim(),
        line2: _line2.text.trim().isEmpty ? null : _line2.text.trim(),
        subDistrict: _subDistrict.text.trim(),
        district: _district.text.trim(),
        province: _province.text.trim(),
        postalCode: _postal.text.trim(),
      ),
    );
    context.go('/cart/checkout/payment');
    // ScaffoldMessenger.of(
    //   context,
    // ).showSnackBar(const SnackBar(content: Text('Shipping address saved')));

    // Next step (we'll add PaymentPage later)
    // context.go('/checkout/payment');
    // context.pop(); // or go back to cart for now
  }

  @override
  Widget build(BuildContext context) {
    final saved = CheckoutScope.of(context).shippingAddress;

    // Prefill if existing
    if (saved != null && _fullName.text.isEmpty) {
      _fullName.text = saved.fullName;
      _phone.text = saved.phone;
      _line1.text = saved.line1;
      _line2.text = saved.line2 ?? '';
      _subDistrict.text = saved.subDistrict;
      _district.text = saved.district;
      _province.text = saved.province;
      _postal.text = saved.postalCode;
    }

    return Scaffold(
      appBar: AppBar(title: const Text('Shipping Address')),
      body: LayoutBuilder(
        builder: (context, c) {
          final isWide = c.maxWidth >= 900;
          return Center(
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
                      child: Form(
                        key: _formKey,
                        child: isWide ? _WideForm(this) : _NarrowForm(this),
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
          );
        },
      ),
    );
  }
}

class _NarrowForm extends StatelessWidget {
  final _AddressPageState s;
  const _NarrowForm(this.s);

  InputDecoration _dec(String label, {String? hint}) =>
      InputDecoration(labelText: label, hintText: hint);

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        TextFormField(
          controller: s._fullName,
          decoration: _dec('Full name'),
          validator: s._req,
        ),
        const SizedBox(height: 12),
        TextFormField(
          controller: s._phone,
          decoration: _dec('Phone', hint: '+66… or 08…'),
          validator: s._phoneRule,
          keyboardType: TextInputType.phone,
        ),
        const SizedBox(height: 12),
        TextFormField(
          controller: s._line1,
          decoration: _dec('Address line 1'),
          validator: s._req,
        ),
        const SizedBox(height: 12),
        TextFormField(
          controller: s._line2,
          decoration: _dec('Address line 2 (optional)'),
        ),
        const SizedBox(height: 12),
        TextFormField(
          controller: s._subDistrict,
          decoration: _dec('Sub-district (ตำบล)'),
          validator: s._req,
        ),
        const SizedBox(height: 12),
        TextFormField(
          controller: s._district,
          decoration: _dec('District (อำเภอ)'),
          validator: s._req,
        ),
        const SizedBox(height: 12),
        TextFormField(
          controller: s._province,
          decoration: _dec('Province (จังหวัด)'),
          validator: s._req,
        ),
        const SizedBox(height: 12),
        TextFormField(
          controller: s._postal,
          decoration: _dec('Postal code'),
          validator: s._postalRule,
          keyboardType: TextInputType.number,
        ),
      ],
    );
  }
}

class _WideForm extends StatelessWidget {
  final _AddressPageState s;
  const _WideForm(this.s);

  InputDecoration _dec(String label, {String? hint}) =>
      InputDecoration(labelText: label, hintText: hint);

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        // row: name + phone
        Row(
          children: [
            Expanded(
              child: TextFormField(
                controller: s._fullName,
                decoration: _dec('Full name'),
                validator: s._req,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: TextFormField(
                controller: s._phone,
                decoration: _dec('Phone', hint: '+66… or 08…'),
                validator: s._phoneRule,
                keyboardType: TextInputType.phone,
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        // row: line1 + line2
        Row(
          children: [
            Expanded(
              child: TextFormField(
                controller: s._line1,
                decoration: _dec('Address line 1'),
                validator: s._req,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: TextFormField(
                controller: s._line2,
                decoration: _dec('Address line 2 (optional)'),
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        // row: subdistrict + district
        Row(
          children: [
            Expanded(
              child: TextFormField(
                controller: s._subDistrict,
                decoration: _dec('Sub-district (ตำบล)'),
                validator: s._req,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: TextFormField(
                controller: s._district,
                decoration: _dec('District (อำเภอ)'),
                validator: s._req,
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        // row: province + postal
        Row(
          children: [
            Expanded(
              child: TextFormField(
                controller: s._province,
                decoration: _dec('Province (จังหวัด)'),
                validator: s._req,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: TextFormField(
                controller: s._postal,
                decoration: _dec('Postal code'),
                validator: s._postalRule,
                keyboardType: TextInputType.number,
              ),
            ),
          ],
        ),
      ],
    );
  }
}
