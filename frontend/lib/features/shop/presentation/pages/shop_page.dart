import 'dart:async';
import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/state/cart_controller.dart';
import 'package:frontend/core/widgets/global_appbar.dart';
import 'package:frontend/features/shop/widgets/product_card.dart';
import '../../../shop/data/product_repository.dart';
import '../../../shop/domain/product.dart';
import 'package:go_router/go_router.dart';

class ShopPage extends StatefulWidget {
  const ShopPage({super.key});
  @override
  State<ShopPage> createState() => _ShopPageState();
}

class _ShopPageState extends State<ShopPage> {
  final _repo = MockProductRepository();
  final _search = TextEditingController();
  String _category = 'All';
  RangeValues _range = const RangeValues(0, 150);
  Timer? _debounce;
  List<Product> _items = [];

  List<String> get _chips => const ['All', 'Fruits', 'Veggies', 'Herbs'];

  Future<void> _load() async {
    final list = await _repo.list(
      q: _search.text,
      category: _category,
      min: _range.start,
      max: _range.end,
    );
    setState(() => _items = list);
  }

  @override
  void initState() {
    super.initState();
    _load();
    _search.addListener(() {
      _debounce?.cancel();
      _debounce = Timer(const Duration(milliseconds: 180), _load);
    });
  }

  @override
  void dispose() {
    _debounce?.cancel();
    _search.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final cart = CartScope.of(context);

    return Scaffold(
      backgroundColor: Theme.of(context).scaffoldBackgroundColor,
      appBar: GlobalAppBar(
        title: 'FreshEase Market',
        actions: [
          IconButton(
            icon: const Icon(Icons.shopping_cart_outlined),
            onPressed: () => context.go('/cart'),
          ),
        ],
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(64),
          child: Padding(
            padding: const EdgeInsets.fromLTRB(16, 0, 16, 12),
            child: _SearchPill(
              controller: _search,
              onFilterTap: () => _openFilterSheet().then((_) => _load()),
            ),
          ),
        ),
      ),

      body: RefreshIndicator(
        onRefresh: _load,
        child: LayoutBuilder(
          builder: (context, constraints) {
            final width = constraints.maxWidth;
            int crossAxisCount = 2;
            double aspectRatio = 0.72;
            if (width > 1200) {
              crossAxisCount = 5;
              aspectRatio = 0.9;
            } else if (width > 900) {
              crossAxisCount = 4;
              aspectRatio = 0.8;
            } else if (width > 600) {
              crossAxisCount = 3;
              aspectRatio = 0.75;
            }

            return CustomScrollView(
              slivers: [
                SliverToBoxAdapter(
                  child: SizedBox(
                    height: 50,
                    child: ListView.separated(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 16,
                        vertical: 8,
                      ),
                      scrollDirection: Axis.horizontal,
                      itemCount: _chips.length,
                      separatorBuilder: (_, _) => const SizedBox(width: 8),
                      itemBuilder: (_, i) {
                        final label = _chips[i];
                        final selected = _category == label;
                        return ChoiceChip(
                          label: Text(label),
                          selected: selected,
                          onSelected: (_) => setState(() => _category = label),
                          shape: StadiumBorder(
                            side: BorderSide(
                              color: selected
                                  ? Theme.of(context).colorScheme.primary
                                  : Colors.grey.shade300,
                            ),
                          ),
                          labelStyle: TextStyle(
                            fontWeight: selected
                                ? FontWeight.w700
                                : FontWeight.w500,
                          ),
                        );
                      },
                    ),
                  ),
                ),

                SliverPadding(
                  padding: const EdgeInsets.all(16),
                  sliver: SliverGrid(
                    delegate: SliverChildBuilderDelegate((context, i) {
                      final p = _items[i];
                      return ProductCard(
                        product: p,
                        onTap: () =>
                            context.go('/shop/product/${p.id}', extra: p),
                        onAdd: () {
                          cart.add(p);
                          ScaffoldMessenger.of(context).showSnackBar(
                            SnackBar(
                              content: Text('${p.name} added to cart'),
                              duration: const Duration(milliseconds: 900),
                            ),
                          );
                        },
                      );
                    }, childCount: _items.length),
                    gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
                      crossAxisCount: crossAxisCount,
                      mainAxisSpacing: 0,
                      crossAxisSpacing: 0,
                      childAspectRatio: aspectRatio,
                    ),
                  ),
                ),

                // empty state
                if (_items.isEmpty)
                  const SliverFillRemaining(
                    hasScrollBody: false,
                    child: Center(child: Text('No results')),
                  ),
              ],
            );
          },
        ),
      ),
    );
  }

  Future<void> _openFilterSheet() async {
    var temp = _range;
    await showModalBottomSheet(
      context: context,
      showDragHandle: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (_) => StatefulBuilder(
        builder: (context, setSheet) {
          return Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Text(
                  'Filters',
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.w700),
                ),
                const SizedBox(height: 16),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('Price range'),
                    Text('฿${temp.start.toInt()} – ฿${temp.end.toInt()}'),
                  ],
                ),
                RangeSlider(
                  values: temp,
                  min: 0,
                  max: 200,
                  divisions: 20,
                  labels: RangeLabels(
                    '฿${temp.start.toInt()}',
                    '฿${temp.end.toInt()}',
                  ),
                  onChanged: (v) => setSheet(() => temp = v),
                ),
                const SizedBox(height: 8),
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton(
                        onPressed: () {
                          setState(() => _range = const RangeValues(0, 150));
                          Navigator.pop(context);
                        },
                        child: const Text('Reset'),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: ElevatedButton(
                        onPressed: () {
                          setState(() => _range = temp);
                          Navigator.pop(context);
                        },
                        child: const Text('Apply'),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 8),
              ],
            ),
          );
        },
      ),
    );
  }
}

class _SearchPill extends StatelessWidget {
  final TextEditingController controller;
  final VoidCallback onFilterTap;
  const _SearchPill({required this.controller, required this.onFilterTap});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Theme.of(
          context,
        ).colorScheme.surfaceContainerHighest.withValues(alpha: 0.95),
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            blurRadius: 12,
            spreadRadius: -2,
            color: Colors.black.withValues(alpha: 0.06),
            offset: const Offset(0, 6),
          ),
        ],
      ),
      padding: const EdgeInsets.symmetric(horizontal: 14),
      child: Row(
        children: [
          const Icon(Icons.search, size: 22),
          const SizedBox(width: 8),
          Expanded(
            child: TextField(
              controller: controller,
              decoration: const InputDecoration(
                border: InputBorder.none,
                hintText: 'Search…',
              ),
              textInputAction: TextInputAction.search,
            ),
          ),
          IconButton(
            icon: const Icon(Icons.tune_rounded),
            onPressed: onFilterTap,
          ),
        ],
      ),
    );
  }
}
