import 'dart:async';

import 'package:flutter/material.dart';
// import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/state/cart_controller.dart';
import 'package:frontend/features/shop/widgets/product_card.dart';
import 'package:go_router/go_router.dart';

import '../../../shop/data/product_repository.dart';
import '../../../shop/domain/product.dart';

// ⚠️ กำหนดสีธีมใหม่ตามภาพ UI (Pale Green/Olive Green)
const Color _primaryColor = Color(0xFF90B56D);
const Color _selectedChipColor = Color(0xFFDCE8CB);

// ใช้ค่าจาก AppColors เดิมหากมี
const Color AppColors_primary = _primaryColor;

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

  // แก้ไข: เปลี่ยนไปใช้ List<Map> เพื่อให้สอดคล้องกับคลาส Product ปัจจุบัน
  final List<Map<String, dynamic>> _exclusiveOffersData = const [
    {
      'id': 6,
      'name': 'Organic Bananas',
      'category': 'Fruits',
      'description': '7pcs, Priceg',
      'price': 4.99,
      'image': 'assets/images/banana.png',
    },
    {
      'id': 7,
      'name': 'Red Apple',
      'category': 'Fruits',
      'description': '1kg, Priceg',
      'price': 4.99,
      // ✅ แก้ไขคีย์จาก 'image:' เป็น 'image' เพื่อแก้ปัญหา Null Error
      'image': 'assets/images/apple.png',
    },
  ];

  final List<Map<String, dynamic>> _bestSellingData = const [
    {
      'id': 8,
      'name': 'Red Bell Pepper',
      'category': 'Veggies',
      'description': '3pcs, Priceg',
      'price': 3.50,
      'image': 'assets/images/pepper.png',
    },
  ];

  Product _mapToProduct(Map<String, dynamic> data) {
    return Product(
      id: data['id'] as int,
      name: data['name'] as String,
      category: data['category'] as String,
      price: data['price'] as double,
      // ✅ ใช้ ?? '' เพื่อป้องกัน Null ในกรณีที่ข้อมูลผิดพลาด
      image: data['image'] as String? ?? '',
    );
  }
  // --------------------------------------------------------------------

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
    return Scaffold(
      backgroundColor: Colors.white, // พื้นหลังของ Scaffold เป็นสีขาวสะอาด
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
                // --- SliverAppBar สำหรับส่วนหัว (Title + Search Bar) ---
                SliverAppBar(
                  backgroundColor: Colors.white, // พื้นหลังเป็นสีขาว
                  automaticallyImplyLeading: false,
                  pinned: true,
                  elevation: 0,
                  expandedHeight: 0,
                  toolbarHeight: 50, // กำหนดความสูงของ toolbar

                  title: SafeArea(
                    bottom: false,
                    child: Center(
                      child: Text(
                        'FreshEase',
                        style: Theme.of(context).textTheme.headlineSmall
                            ?.copyWith(fontWeight: FontWeight.w700),
                      ),
                    ),
                  ),
                  centerTitle: true,

                  // bottom สำหรับ Search Bar
                  bottom: PreferredSize(
                    preferredSize: const Size.fromHeight(64),
                    child: Padding(
                      padding: const EdgeInsets.fromLTRB(16, 0, 16, 12),
                      child: _SearchPill(
                        controller: _search,
                        onFilterTap: () {},
                        showFilter: false,
                      ),
                    ),
                  ),
                ),

                // --- แบนเนอร์ (Fresh Vegetables) ---
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(16, 8, 16, 16),
                    child: Container(
                      height: 120,
                      decoration: BoxDecoration(
                        borderRadius: BorderRadius.circular(12),
                        color: _primaryColor.withOpacity(0.5),
                      ),
                      child: const Center(
                        // ส่วนแสดงข้อความ (ถ้าคุณยังไม่ได้ใส่รูป Banner)
                        // Text('Fresh Vegetables Banner Placeholder', style: TextStyle(color: Colors.white)),
                      ),
                    ),
                  ),
                ),

                // --- Exclusive Offer Section ---
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text(
                          'Exclusive Offer',
                          style: TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        GestureDetector(
                          onTap: () => {/* Handle See all tap */},
                          child: const Text(
                            'See all',
                            style: TextStyle(
                              color: _primaryColor,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                SliverToBoxAdapter(
                  child: SizedBox(
                    // ✅ แก้ไข: เพิ่มความสูงเป็น 260 เพื่อแก้ปัญหา Overflow
                    height: 260,
                    child: ListView.separated(
                      scrollDirection: Axis.horizontal,
                      padding: const EdgeInsets.all(16),
                      itemCount: _exclusiveOffersData.length,
                      separatorBuilder: (context, index) =>
                          const SizedBox(width: 12),
                      itemBuilder: (context, index) {
                        final data = _exclusiveOffersData[index];
                        final p = _mapToProduct(data);

                        return SizedBox(
                          width: 150,
                          child: ProductCard(
                            product: p,
                            productDetail:
                                data['description'], // ส่งรายละเอียดเสริม
                            onTap: () =>
                                context.go('/shop/product/${p.id}', extra: p),
                            onAdd: () {
                              final cart = CartScope.of(context);
                              cart.add(p);
                              ScaffoldMessenger.of(context).showSnackBar(
                                SnackBar(
                                  content: Text('${p.name} added to cart'),
                                  duration: const Duration(milliseconds: 900),
                                ),
                              );
                            },
                          ),
                        );
                      },
                    ),
                  ),
                ),

                // --- Best Selling Section Title ---
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(16, 0, 16, 8),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text(
                          'Best Selling',
                          style: TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        GestureDetector(
                          onTap: () => {/* Handle See all tap */},
                          child: const Text(
                            'See all',
                            style: TextStyle(
                              color: _primaryColor,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),

                // --- Chip Filters ---
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
                          onSelected: (s) {
                            setState(() => _category = label);
                            _load();
                          },
                          // ⚠️ ปรับสไตล์ ChoiceChip ให้ตรงกับรูปภาพ
                          selectedColor:
                              _selectedChipColor, // สีพื้นหลังเขียวอ่อน
                          shape: StadiumBorder(
                            side: BorderSide(
                              color: selected
                                  ? _primaryColor // ขอบสีเขียวอ่อน
                                  : Colors.grey.shade300,
                            ),
                          ),
                          labelStyle: TextStyle(
                            color: selected
                                ? Colors.black
                                : Colors
                                      .grey
                                      .shade700, // ข้อความเป็นสีดำเมื่อเลือก
                            fontWeight: selected
                                ? FontWeight.w700
                                : FontWeight.w500,
                          ),
                        );
                      },
                    ),
                  ),
                ),

                // --- Product Grid ---
                SliverPadding(
                  padding: const EdgeInsets.all(16),
                  sliver: SliverGrid(
                    delegate: SliverChildBuilderDelegate((context, i) {
                      final p = _items[i];
                      final cart = CartScope.of(context);
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

  // ... (โค้ด _openFilterSheet และ _SearchPill)

  Future<void> _openFilterSheet() async {
    // ... (โค้ดสำหรับ Filter Sheet)
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
  final bool showFilter;
  const _SearchPill({
    required this.controller,
    required this.onFilterTap,
    this.showFilter = true,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 48,
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: Colors.grey.shade200),
        boxShadow: [
          BoxShadow(
            blurRadius: 10,
            spreadRadius: -5,
            color: Colors.black.withOpacity(0.05),
            offset: const Offset(0, 5),
          ),
        ],
      ),
      padding: const EdgeInsets.symmetric(horizontal: 14),
      child: Row(
        children: [
          Icon(Icons.search, size: 22, color: Colors.grey.shade600),
          const SizedBox(width: 8),
          Expanded(
            child: TextField(
              controller: controller,
              decoration: InputDecoration(
                border: InputBorder.none,
                hintText: 'Search Store',
                hintStyle: TextStyle(color: Colors.grey.shade500),
              ),
              textInputAction: TextInputAction.search,
            ),
          ),
          if (showFilter)
            IconButton(
              icon: Icon(Icons.tune_rounded, color: Colors.grey.shade600),
              onPressed: onFilterTap,
            ),
        ],
      ),
    );
  }
}
