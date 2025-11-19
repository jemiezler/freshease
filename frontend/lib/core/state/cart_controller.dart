import 'package:flutter/material.dart';
import 'package:frontend/features/shop/domain/product.dart';
import 'package:frontend/features/cart/data/cart_repository.dart';
import 'package:frontend/features/cart/data/models/cart_dtos.dart';

class CartLine {
  final Product product;
  int qty;
  CartLine({required this.product, this.qty = 1});
  double get lineTotal => product.price * qty;
}

class CartController extends ChangeNotifier {
  final CartRepository _repository;
  CartDTO? _cart;
  bool _isLoading = false;
  String? _error;

  CartController(this._repository) {
    _loadCart();
  }

  // --- getters
  CartDTO? get cart => _cart;
  bool get isLoading => _isLoading;
  String? get error => _error;

  List<CartLine> get lines {
    if (_cart == null) return [];
    return _cart!.items
        .map(
          (item) => CartLine(
            product: Product(
              id: item.productId,
              name: item.productName,
              price: item.productPrice,
              image: item.productImage,
              category: 'Unknown',
              description: '',
              unitLabel: 'kg',
              vendorId: '',
              vendorName: '',
              categoryId: '',
              categoryName: 'Unknown',
              stockQuantity: 100,
              isInStock: true,
            ),
            qty: item.quantity,
          ),
        )
        .toList();
  }

  int get itemKinds => _cart?.items.length ?? 0;
  int get count {
    if (_cart == null) return 0;
    return _cart!.items.fold(0, (sum, item) => sum + item.quantity);
  }

  double get subtotal => _cart?.subtotal ?? 0.0;
  double get shipping => _cart?.shipping ?? 0.0;
  double get promoDiscount => _cart?.promoDiscount ?? 0.0;
  double get vat => _cart?.tax ?? 0.0;
  double get total => _cart?.total ?? 0.0;
  String? get promoCode => _cart?.promoCode;

  Future<void> _loadCart() async {
    _setLoading(true);
    try {
      _cart = await _repository.getCart();
      _error = null;
    } catch (e) {
      _error = e.toString();
    } finally {
      _setLoading(false);
    }
  }

  void _setLoading(bool loading) {
    _isLoading = loading;
    notifyListeners();
  }

  // --- mutators
  Future<void> add(Product p, {int qty = 1}) async {
    _setLoading(true);
    try {
      _cart = await _repository.addToCart(p, quantity: qty);
      _error = null;
      notifyListeners(); // Explicitly notify after cart update
    } catch (e) {
      _error = e.toString();
      notifyListeners(); // Notify even on error so UI can show error state
    } finally {
      _setLoading(false);
    }
  }

  Future<void> decrement(Product p, {int qty = 1}) async {
    final item = _cart?.items.firstWhere(
      (item) => item.productId == p.id,
      orElse: () => throw Exception('Item not found'),
    );

    if (item != null) {
      final newQuantity = item.quantity - qty;
      if (newQuantity <= 0) {
        await remove(p);
      } else {
        await setQty(p, newQuantity);
      }
    }
  }

  Future<void> setQty(Product p, int qty) async {
    if (qty <= 0) return remove(p);

    final item = _cart?.items.firstWhere(
      (item) => item.productId == p.id,
      orElse: () => throw Exception('Item not found'),
    );

    if (item != null) {
      _setLoading(true);
      try {
        _cart = await _repository.updateCartItem(item.id, qty);
        _error = null;
      } catch (e) {
        _error = e.toString();
      } finally {
        _setLoading(false);
      }
    }
  }

  Future<void> remove(Product p) async {
    final item = _cart?.items.firstWhere(
      (item) => item.productId == p.id,
      orElse: () => throw Exception('Item not found'),
    );

    if (item != null) {
      _setLoading(true);
      try {
        _cart = await _repository.removeCartItem(item.id);
        _error = null;
      } catch (e) {
        _error = e.toString();
      } finally {
        _setLoading(false);
      }
    }
  }

  Future<void> clear() async {
    _setLoading(true);
    try {
      _cart = await _repository.clearCart();
      _error = null;
    } catch (e) {
      _error = e.toString();
    } finally {
      _setLoading(false);
    }
  }

  Future<bool> applyPromo(String code) async {
    _setLoading(true);
    try {
      _cart = await _repository.applyPromoCode(code);
      _error = null;
      return true;
    } catch (e) {
      _error = e.toString();
      return false;
    } finally {
      _setLoading(false);
    }
  }

  Future<void> removePromo() async {
    _setLoading(true);
    try {
      _cart = await _repository.removePromoCode();
      _error = null;
    } catch (e) {
      _error = e.toString();
    } finally {
      _setLoading(false);
    }
  }

  Future<void> refresh() async {
    await _loadCart();
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
