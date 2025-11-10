import 'dart:async';
import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/state/cart_controller.dart';
import 'package:frontend/core/widgets/banner_carousel.dart';
import 'package:frontend/core/widgets/global_appbar.dart';
import 'package:frontend/core/widgets/search_pill.dart';
import 'package:frontend/core/widgets/design_system/soft_chip.dart';
import 'package:frontend/core/widgets/design_system/soft_button.dart';
import 'package:frontend/core/widgets/design_system/soft_card.dart';
import 'package:frontend/core/theme/design_tokens.dart';
import 'package:frontend/core/health/health_controller.dart';
import 'package:frontend/features/shop/data/product_repository.dart';
import 'package:frontend/features/shop/domain/product.dart';
import 'package:frontend/features/shop/widgets/product_card.dart';
import 'package:go_router/go_router.dart';
import 'package:get_it/get_it.dart';

class ShopPage extends StatefulWidget {
  const ShopPage({super.key});
  @override
  State<ShopPage> createState() => _ShopPageState();
}

class _ShopPageState extends State<ShopPage> {
  final _repo = GetIt.instance<ProductRepository>();
  final _healthController = GetIt.instance<HealthController>();
  final _search = TextEditingController();
  String _category = 'All';
  RangeValues _range = const RangeValues(0, 150);
  Timer? _debounce;
  List<Product> _items = [];
  bool _isLoading = false;
  final List<BannerItem> _banners = const [
    BannerItem(imageUrl: 'https://picsum.photos/1200/400?1', route: '/promo/1'),
    BannerItem(imageUrl: 'https://picsum.photos/1200/400?2', route: '/promo/2'),
    BannerItem(
      imageUrl: 'https://picsum.photos/1200/400?3',
      // or custom onTap if you don’t want a route:
      // onTap: () => debugPrint('clicked banner 3'),
    ),
  ];
  List<String> get _chips => const [
    'All',
    'Prepared Food',
    'Veggies',
    'Fruits',
  ];

  Future<void> _load() async {
    if (_isLoading) return;

    setState(() => _isLoading = true);

    try {
      final list = await _repo.list(
        q: _search.text,
        category: _category,
        min: _range.start,
        max: _range.end,
      );
      setState(() => _items = list);
    } catch (e) {
      setState(() => _items = []);
    } finally {
      setState(() => _isLoading = false);
    }
  }

  @override
  void initState() {
    super.initState();
    _load();

    // Trigger meal plan generation when user visits shop page
    _healthController.triggerAutoGeneration();

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
            child: SearchPill(
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
                // --- Banner carousel ---
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(16, 0, 16, 8),
                    child: BannerCarousel(
                      items: _banners,
                      borderRadius: DesignTokens.radiusMedium,
                      autoPlay: true,
                      autoPlayInterval: const Duration(seconds: 4),
                      onPageChanged: (_) {},
                    ),
                  ),
                ),

                SliverToBoxAdapter(
                  child: SizedBox(
                    height: 50,
                    child: ListView.separated(
                      padding: const EdgeInsets.symmetric(
                        horizontal: DesignTokens.paddingMedium,
                        vertical: DesignTokens.paddingSmall,
                      ),
                      scrollDirection: Axis.horizontal,
                      itemCount: _chips.length,
                      separatorBuilder: (_, _) => const SizedBox(width: 8),
                      itemBuilder: (_, i) {
                        final label = _chips[i];
                        final selected = _category == label;
                        return SoftChip(
                          label: label,
                          isSelected: selected,
                          onTap: () => setState(() => _category = label),
                        );
                      },
                    ),
                  ),
                ),

                if (_isLoading)
                  const SliverFillRemaining(
                    hasScrollBody: false,
                    child: Center(child: CircularProgressIndicator()),
                  )
                else
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

                if (!_isLoading && _items.isEmpty)
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
      backgroundColor: Colors.transparent,
      builder: (_) => StatefulBuilder(
        builder: (context, setSheet) {
          return SoftCard(
            margin: EdgeInsets.zero,
            borderRadius: DesignTokens.radiusLarge,
            child: Padding(
              padding: const EdgeInsets.all(DesignTokens.paddingLarge),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Text(
                    'Filters',
                    style: TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.w700,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  const SizedBox(height: DesignTokens.paddingLarge),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text(
                        'Price range',
                        style: TextStyle(
                          fontSize: 16,
                          color: AppColors.textPrimary,
                        ),
                      ),
                      Text(
                        '฿${temp.start.toInt()} – ฿${temp.end.toInt()}',
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: AppColors.primary,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: DesignTokens.paddingMedium),
                  RangeSlider(
                    values: temp,
                    min: 0,
                    max: 200,
                    divisions: 20,
                    activeColor: AppColors.primary,
                    labels: RangeLabels(
                      '฿${temp.start.toInt()}',
                      '฿${temp.end.toInt()}',
                    ),
                    onChanged: (v) => setSheet(() => temp = v),
                  ),
                  const SizedBox(height: DesignTokens.paddingLarge),
                  Row(
                    children: [
                      Expanded(
                        child: SoftButton(
                          label: 'Reset',
                          isPrimary: false,
                          onPressed: () {
                            setState(() => _range = const RangeValues(0, 150));
                            Navigator.pop(context);
                          },
                        ),
                      ),
                      const SizedBox(width: DesignTokens.paddingMedium),
                      Expanded(
                        child: SoftButton(
                          label: 'Apply',
                          isPrimary: true,
                          onPressed: () {
                            setState(() => _range = temp);
                            Navigator.pop(context);
                          },
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
