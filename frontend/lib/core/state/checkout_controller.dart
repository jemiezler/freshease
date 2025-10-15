import 'package:flutter/material.dart';

class Address {
  final String fullName;
  final String phone;
  final String line1;
  final String? line2;
  final String subDistrict;
  final String district;
  final String province;
  final String postalCode;

  const Address({
    required this.fullName,
    required this.phone,
    required this.line1,
    this.line2,
    required this.subDistrict,
    required this.district,
    required this.province,
    required this.postalCode,
  });
}

enum PaymentMethod { cod, card }

class OrderResult {
  final String orderId;
  final DateTime createdAt;
  const OrderResult(this.orderId, this.createdAt);
}

class PlanOrder {
  final int id;
  final String title;
  final double price;
  final String subtitle;
  const PlanOrder({
    required this.id,
    required this.title,
    required this.price,
    required this.subtitle,
  });
}

class CheckoutController extends ChangeNotifier {
  Address? _shippingAddress;
  PaymentMethod? _paymentMethod;
  String? _paymentLast4; // for card (placeholder)

  Address? get shippingAddress => _shippingAddress;
  PaymentMethod? get paymentMethod => _paymentMethod;
  String? get paymentLast4 => _paymentLast4;

  PlanOrder? _selectedPlan;
  PlanOrder? get selectedPlan => _selectedPlan;

  bool get isPlanCheckout => _selectedPlan != null;

  void setPlanCheckout(PlanOrder p) {
    _selectedPlan = p;
    notifyListeners();
  }

  void clearPlan() {
    _selectedPlan = null;
    notifyListeners();
  }

  void clear() {
    _shippingAddress = null;
    _paymentMethod = null;
    _paymentLast4 = null;
    _selectedPlan = null;
    notifyListeners();
  }

  void setShipping(Address a) {
    _shippingAddress = a;
    notifyListeners();
  }

  void setPayment(PaymentMethod m, {String? last4}) {
    _paymentMethod = m;
    _paymentLast4 = last4;
    notifyListeners();
  }

  bool get isReadyForReview =>
      _shippingAddress != null && _paymentMethod != null;

  /// Simulate placing an order – replace with a real API call later.
  Future<OrderResult> placeOrder() async {
    await Future.delayed(const Duration(milliseconds: 350));
    final ts = DateTime.now();
    final id = 'FE-${ts.millisecondsSinceEpoch.toString().substring(7)}';
    return OrderResult(id, ts);
  }
}

/// Inherited scope stays the same
class CheckoutScope extends InheritedNotifier<CheckoutController> {
  final CheckoutController controller;
  const CheckoutScope({
    super.key,
    required this.controller,
    required super.child,
  }) : super(notifier: controller);

  static CheckoutController of(BuildContext context) {
    final scope = context.dependOnInheritedWidgetOfExactType<CheckoutScope>();
    assert(scope != null, 'CheckoutScope not found in context');
    return scope!.controller;
  }
}
