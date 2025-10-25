import 'package:flutter/material.dart';
import 'package:frontend/features/shop/domain/product.dart';

class CartLine {
  final Product product;
  int qty;
  CartLine({required this.product, this.qty = 1});
  double get lineTotal => product.price * qty;
}

class CartController extends ChangeNotifier {
  final List<CartLine> _lines = [];
  String? _promoCode;

  // --- getters
  List<CartLine> get lines => List.unmodifiable(_lines);
  int get itemKinds => _lines.length;
  int get count => _lines.fold(0, (s, l) => s + l.qty);
  double get subtotal => _lines.fold(0.0, (s, l) => s + l.lineTotal);

  // Simple business rules (tweak as needed)
  double get shipping =>
      subtotal == 0 ? 0 : (subtotal >= 200 ? 0 : 20); // free shipping over ฿200
  double get promoDiscount {
    if (_promoCode == null) return 0;
    if (_promoCode!.toUpperCase() == 'FRESH10') return subtotal * 0.10;
    if (_promoCode!.toUpperCase() == 'FREESHIP') return shipping;
    return 0;
  }

  double get vat => (subtotal - promoDiscount) * 0.07; // 7% VAT after discount
  double get total => (subtotal - promoDiscount) + shipping + vat;

  String? get promoCode => _promoCode;

  // --- mutators
  // 🎯 แก้ไข: เปลี่ยนชื่อพารามิเตอร์จาก 'qty' เป็น 'quantity'
  void add(Product p, {int quantity = 1}) {
    final i = _lines.indexWhere((l) => l.product.id == p.id);
    if (i >= 0) {
      // 🎯 ใช้ 'quantity' แทน 'qty'
      _lines[i].qty += quantity;
    } else {
      // 🎯 ใช้ 'quantity' แทน 'qty' ใน Constructor
      _lines.add(CartLine(product: p, qty: quantity));
    }
    notifyListeners();
  }

  // 🎯 แก้ไข: เปลี่ยนชื่อพารามิเตอร์จาก 'qty' เป็น 'quantity'
  void decrement(Product p, {int quantity = 1}) {
    final i = _lines.indexWhere((l) => l.product.id == p.id);
    if (i < 0) return;
    // 🎯 ใช้ 'quantity' แทน 'qty'
    _lines[i].qty -= quantity;
    if (_lines[i].qty <= 0) _lines.removeAt(i);
    notifyListeners();
  }

  void setQty(Product p, int qty) {
    if (qty <= 0) return remove(p);
    final i = _lines.indexWhere((l) => l.product.id == p.id);
    if (i >= 0) {
      _lines[i].qty = qty;
      notifyListeners();
    }
  }

  void remove(Product p) {
    _lines.removeWhere((l) => l.product.id == p.id);
    notifyListeners();
  }

  void clear() {
    _lines.clear();
    _promoCode = null;
    notifyListeners();
  }

  bool applyPromo(String code) {
    final c = code.trim().toUpperCase();
    // Accept only known codes (sample)
    if (c == 'FRESH10' || c == 'FREESHIP') {
      _promoCode = c;
      notifyListeners();
      return true;
    }
    return false;
  }

  void removePromo() {
    _promoCode = null;
    notifyListeners();
  }
}

/// Inherited notifier to access CartController without external deps.
class CartScope extends InheritedNotifier<CartController> {
  final CartController controller;
  const CartScope({super.key, required this.controller, required super.child})
    : super(notifier: controller);

  static CartController of(BuildContext context) {
    final scope = context.dependOnInheritedWidgetOfExactType<CartScope>();
    assert(scope != null, 'CartScope not found in context');
    return scope!.controller;
  }
}
