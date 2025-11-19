import 'package:flutter/material.dart';
import 'package:frontend/core/state/checkout_controller.dart';
import 'package:go_router/go_router.dart';
import '../../data/models/bundle_dtos.dart';
import '../../data/repositories/plans_repository.dart';

class PlanDetailPage extends StatefulWidget {
  final BundleDTO? bundle;
  final String? bundleId;
  const PlanDetailPage({super.key, this.bundle, this.bundleId});

  @override
  State<PlanDetailPage> createState() => _PlanDetailPageState();
}

class _PlanDetailPageState extends State<PlanDetailPage> {
  final PlansRepository _repository = PlansRepository();
  BundleDTO? _bundle;
  List<BundleItemDTO> _bundleItems = [];
  bool _isLoading = true;
  String? _error;
  Color? _color;

  @override
  void initState() {
    super.initState();
    _bundle = widget.bundle;
    if (_bundle != null) {
      _loadBundleDetails();
    } else if (widget.bundleId != null) {
      _loadBundleById(widget.bundleId!);
    } else {
      setState(() {
        _isLoading = false;
        _error = 'Bundle not found';
      });
    }
  }

  Future<void> _loadBundleById(String bundleId) async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final bundle = await _repository.getBundle(bundleId);
      setState(() {
        _bundle = bundle;
      });
      await _loadBundleDetails();
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  Future<void> _loadBundleDetails() async {
    if (_bundle == null) return;

    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final bundleWithItems = await _repository.getBundleWithItems(_bundle!.id);
      setState(() {
        _bundle = bundleWithItems.bundle;
        _bundleItems = bundleWithItems.items;
        _isLoading = false;
        // Generate a color based on bundle ID hash
        _color = _getColorForBundle(_bundle!.id);
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  Color _getColorForBundle(String id) {
    final colors = [Colors.green, Colors.orange, Colors.teal, Colors.blue, Colors.purple];
    final hash = id.hashCode;
    return colors[hash.abs() % colors.length];
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return Scaffold(
        appBar: AppBar(title: const Text('Plan Details')),
        body: const Center(child: CircularProgressIndicator()),
      );
    }

    if (_error != null || _bundle == null) {
      return Scaffold(
        appBar: AppBar(title: const Text('Plan Details')),
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                _error ?? 'Bundle not found',
                style: TextStyle(color: Colors.red[700]),
              ),
              if (_error != null) ...[
                const SizedBox(height: 16),
                ElevatedButton(
                  onPressed: _loadBundleDetails,
                  child: const Text('Retry'),
                ),
              ],
            ],
          ),
        ),
      );
    }

    final color = _color ?? Colors.green;

    return Scaffold(
      appBar: AppBar(title: Text(_bundle!.name)),
      body: Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 900),
          child: ListView(
            padding: const EdgeInsets.all(24),
            children: [
              Text(
                _bundle!.name,
                style: const TextStyle(
                  fontSize: 26,
                  fontWeight: FontWeight.w800,
                ),
              ),
              if (_bundle!.description != null && _bundle!.description!.isNotEmpty) ...[
                const SizedBox(height: 8),
                Text(
                  _bundle!.description!,
                  style: TextStyle(fontSize: 16, color: Colors.grey[600]),
                ),
              ],
              const SizedBox(height: 16),
              Container(
                decoration: BoxDecoration(
                  color: color.withValues(alpha: .1),
                  borderRadius: BorderRadius.circular(12),
                ),
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Bundle Items:',
                      style: TextStyle(
                        fontWeight: FontWeight.w700,
                        fontSize: 18,
                      ),
                    ),
                    const SizedBox(height: 12),
                    if (_bundleItems.isEmpty)
                      const Padding(
                        padding: EdgeInsets.all(8.0),
                        child: Text(
                          'No items in this bundle',
                          style: TextStyle(color: Colors.grey),
                        ),
                      )
                    else
                      ..._bundleItems.map((item) => Padding(
                            padding: const EdgeInsets.only(bottom: 8),
                            child: Row(
                              children: [
                                const Icon(Icons.check, color: Colors.green),
                                const SizedBox(width: 8),
                                Expanded(
                                  child: Text(
                                    'Product ID: ${item.productId} (Qty: ${item.qty})',
                                    style: const TextStyle(fontSize: 14),
                                  ),
                                ),
                              ],
                            ),
                          )),
                  ],
                ),
              ),
              const SizedBox(height: 24),
              Text(
                'Price: ฿${_bundle!.price.toStringAsFixed(0)}',
                style: const TextStyle(
                  fontSize: 22,
                  fontWeight: FontWeight.w700,
                  color: Colors.green,
                ),
              ),
              const SizedBox(height: 24),
              FilledButton.icon(
                onPressed: () {
                  final co = CheckoutScope.of(context);
                  co.setPlanCheckout(
                    PlanOrder(
                      id: _bundle!.id,
                      title: _bundle!.name,
                      price: _bundle!.price,
                      subtitle: _bundle!.description,
                    ),
                  );

                  // Navigate directly into checkout flow (starting at address)
                  context.go('/cart/checkout/address');
                },
                icon: const Icon(Icons.shopping_bag_outlined),
                label: const Text('Subscribe Now'),
                style: FilledButton.styleFrom(
                  minimumSize: const Size.fromHeight(48),
                  backgroundColor: color.withValues(alpha: .9),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
