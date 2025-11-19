import 'models/cart_dtos.dart';
import 'sources/cart_api.dart';
import '../../shop/domain/product.dart';

abstract class CartRepository {
  Future<CartDTO> getCart();
  Future<CartDTO> addToCart(Product product, {int quantity = 1});
  Future<CartDTO> updateCartItem(String cartItemId, int quantity);
  Future<CartDTO> removeCartItem(String cartItemId);
  Future<CartDTO> applyPromoCode(String promoCode);
  Future<CartDTO> removePromoCode();
  Future<CartDTO> clearCart();
}

/// Real repository implementation using the backend API
class ApiCartRepository implements CartRepository {
  final CartApiService _apiService;

  ApiCartRepository(this._apiService);

  @override
  Future<CartDTO> getCart() async {
    try {
      return await _apiService.getCart();
    } catch (e) {
      // Return empty cart on error
      return _createEmptyCart();
    }
  }

  @override
  Future<CartDTO> addToCart(Product product, {int quantity = 1}) async {
    try {
      final request = AddToCartRequest(
        productId: product.id,
        quantity: quantity,
      );
      return await _apiService.addToCart(request);
    } catch (e) {
      // Log error for debugging
      print('Error adding to cart: $e');
      // Rethrow so controller can handle the error and show it to user
      rethrow;
    }
  }

  @override
  Future<CartDTO> updateCartItem(String cartItemId, int quantity) async {
    try {
      final request = UpdateCartItemRequest(
        cartItemId: cartItemId,
        quantity: quantity,
      );
      return await _apiService.updateCartItem(request);
    } catch (e) {
      return await getCart();
    }
  }

  @override
  Future<CartDTO> removeCartItem(String cartItemId) async {
    try {
      return await _apiService.removeCartItem(cartItemId);
    } catch (e) {
      return await getCart();
    }
  }

  @override
  Future<CartDTO> applyPromoCode(String promoCode) async {
    try {
      final request = ApplyPromoRequest(promoCode: promoCode);
      return await _apiService.applyPromoCode(request);
    } catch (e) {
      return await getCart();
    }
  }

  @override
  Future<CartDTO> removePromoCode() async {
    try {
      return await _apiService.removePromoCode();
    } catch (e) {
      return await getCart();
    }
  }

  @override
  Future<CartDTO> clearCart() async {
    try {
      return await _apiService.clearCart();
    } catch (e) {
      return _createEmptyCart();
    }
  }

  CartDTO _createEmptyCart() {
    return CartDTO(
      id: '',
      status: 'pending',
      subtotal: 0.0,
      shipping: 0.0,
      tax: 0.0,
      total: 0.0,
      items: [],
      promoDiscount: 0.0,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    );
  }
}

/// Mock repository for development/testing
class MockCartRepository implements CartRepository {
  CartDTO _cart = CartDTO(
    id: 'mock-cart-1',
    status: 'pending',
    subtotal: 0.0,
    shipping: 0.0,
    tax: 0.0,
    total: 0.0,
    items: [],
    promoDiscount: 0.0,
    createdAt: DateTime.now(),
    updatedAt: DateTime.now(),
  );

  @override
  Future<CartDTO> getCart() async {
    await Future.delayed(const Duration(milliseconds: 100));
    return _cart;
  }

  @override
  Future<CartDTO> addToCart(Product product, {int quantity = 1}) async {
    await Future.delayed(const Duration(milliseconds: 200));

    // Check if item already exists
    final existingIndex = _cart.items.indexWhere(
      (item) => item.productId == product.id,
    );

    if (existingIndex >= 0) {
      // Update existing item
      final existingItem = _cart.items[existingIndex];
      final updatedItem = CartItemDTO(
        id: existingItem.id,
        productId: existingItem.productId,
        productName: existingItem.productName,
        productImage: existingItem.productImage,
        productPrice: existingItem.productPrice,
        quantity: existingItem.quantity + quantity,
        lineTotal:
            existingItem.productPrice * (existingItem.quantity + quantity),
      );

      final updatedItems = List<CartItemDTO>.from(_cart.items);
      updatedItems[existingIndex] = updatedItem;

      _cart = _recalculateCart(updatedItems);
    } else {
      // Add new item
      final newItem = CartItemDTO(
        id: 'mock-item-${DateTime.now().millisecondsSinceEpoch}',
        productId: product.id,
        productName: product.name,
        productImage: product.image,
        productPrice: product.price,
        quantity: quantity,
        lineTotal: product.price * quantity,
      );

      final updatedItems = [..._cart.items, newItem];
      _cart = _recalculateCart(updatedItems);
    }

    return _cart;
  }

  @override
  Future<CartDTO> updateCartItem(String cartItemId, int quantity) async {
    await Future.delayed(const Duration(milliseconds: 150));

    if (quantity <= 0) {
      return removeCartItem(cartItemId);
    }

    final itemIndex = _cart.items.indexWhere((item) => item.id == cartItemId);
    if (itemIndex < 0) return _cart;

    final item = _cart.items[itemIndex];
    final updatedItem = CartItemDTO(
      id: item.id,
      productId: item.productId,
      productName: item.productName,
      productImage: item.productImage,
      productPrice: item.productPrice,
      quantity: quantity,
      lineTotal: item.productPrice * quantity,
    );

    final updatedItems = List<CartItemDTO>.from(_cart.items);
    updatedItems[itemIndex] = updatedItem;

    _cart = _recalculateCart(updatedItems);
    return _cart;
  }

  @override
  Future<CartDTO> removeCartItem(String cartItemId) async {
    await Future.delayed(const Duration(milliseconds: 150));

    final updatedItems = _cart.items
        .where((item) => item.id != cartItemId)
        .toList();
    _cart = _recalculateCart(updatedItems);
    return _cart;
  }

  @override
  Future<CartDTO> applyPromoCode(String promoCode) async {
    await Future.delayed(const Duration(milliseconds: 200));

    double discount = 0.0;
    if (promoCode.toUpperCase() == 'FRESH10') {
      discount = _cart.subtotal * 0.10;
    } else if (promoCode.toUpperCase() == 'FREESHIP') {
      discount = _cart.shipping;
    }

    _cart = CartDTO(
      id: _cart.id,
      status: _cart.status,
      subtotal: _cart.subtotal,
      shipping: _cart.shipping,
      tax: _cart.tax,
      total: _cart.total,
      items: _cart.items,
      promoCode: promoCode,
      promoDiscount: discount,
      createdAt: _cart.createdAt,
      updatedAt: DateTime.now(),
    );

    return _cart;
  }

  @override
  Future<CartDTO> removePromoCode() async {
    await Future.delayed(const Duration(milliseconds: 150));

    _cart = CartDTO(
      id: _cart.id,
      status: _cart.status,
      subtotal: _cart.subtotal,
      shipping: _cart.shipping,
      tax: _cart.tax,
      total: _cart.total,
      items: _cart.items,
      promoCode: null,
      promoDiscount: 0.0,
      createdAt: _cart.createdAt,
      updatedAt: DateTime.now(),
    );

    return _cart;
  }

  @override
  Future<CartDTO> clearCart() async {
    await Future.delayed(const Duration(milliseconds: 200));

    _cart = CartDTO(
      id: _cart.id,
      status: 'pending',
      subtotal: 0.0,
      shipping: 0.0,
      tax: 0.0,
      total: 0.0,
      items: [],
      promoCode: null,
      promoDiscount: 0.0,
      createdAt: _cart.createdAt,
      updatedAt: DateTime.now(),
    );

    return _cart;
  }

  CartDTO _recalculateCart(List<CartItemDTO> items) {
    final subtotal = items.fold(0.0, (sum, item) => sum + item.lineTotal);
    final shipping = subtotal >= 200 ? 0.0 : 20.0;
    final tax = (subtotal - _cart.promoDiscount) * 0.07;
    final total = (subtotal - _cart.promoDiscount) + shipping + tax;

    return CartDTO(
      id: _cart.id,
      status: _cart.status,
      subtotal: subtotal,
      shipping: shipping,
      tax: tax,
      total: total,
      items: items,
      promoCode: _cart.promoCode,
      promoDiscount: _cart.promoDiscount,
      createdAt: _cart.createdAt,
      updatedAt: DateTime.now(),
    );
  }
}
