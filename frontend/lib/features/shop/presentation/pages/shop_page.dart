// ไฟล์: shoppage.dart

import 'package:flutter/material.dart';
import 'package:frontend/core/state/cart_controller.dart';
import 'package:frontend/features/shop/widgets/product_card.dart';
import 'package:go_router/go_router.dart';

import '../../../shop/data/product_repository.dart';
import '../../../shop/domain/product.dart';
import '../../widgets/search_pill_widget.dart'; // 🎯 Widget ช่องค้นหาที่ใช้ร่วมกัน

// ⚠️ กำหนดสีธีมใหม่ตามภาพ UI
const Color _primaryColor = Color(0xFF53B175); // <--- 🎨 แก้ไขสีเขียวตาม UI
const Color _selectedChipColor = Color(0xFFE8F5E9); // <--- 🎨 แก้ไขสีชิป
const Color _scaffoldBgColor = Color(
  0xFFFAFAFA,
); // <--- 🎨 สีพื้นหลัง (เกือบขาว)
const Color AppColors_primary = _primaryColor;

class ShopPage extends StatefulWidget {
  const ShopPage({super.key});
  @override
  State<ShopPage> createState() => _ShopPageState();
}

class _ShopPageState extends State<ShopPage> {
  final _repo = MockProductRepository();
  String _category = 'All';
  List<Product> _items = [];

  // ‼️ เพิ่ม State สำหรับ Banner
  int _bannerCurrentPage = 0;
  final PageController _bannerController = PageController();
  final List<String> _bannerImages = [
    'assets/images/fresh_veg_banner.png',
    'assets/images/fresh_veg_banner.png',
    'assets/images/fresh_veg_banner.png',
  ];

  // 🎯 ข้อมูล Exclusive Offers
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
      'image': 'assets/images/apple.png',
    },
  ];

  // 🎯 ข้อมูล Best Selling (เพิ่มตามรูป) // <--- 🛍️ เพิ่มข้อมูลส่วน Best Selling
  final List<Map<String, dynamic>> _bestSellingData = const [
    {
      'id': 8,
      'name': 'Red Pepper',
      'category': 'Veggies',
      'description': '1kg, Priceg',
      'price': 4.99, // ⚠️ ราคาในรูปไม่ชัด
      'image': 'assets/images/pepper.png', // ⚠️ ต้องมีรูปนี้ใน assets
    },
    {
      'id': 6, // ใช้ซ้ำได้เพื่อทดสอบ
      'name': 'Organic Bananas',
      'category': 'Fruits',
      'description': '7pcs, Priceg',
      'price': 4.99,
      'image': 'assets/images/banana.png',
    },
  ];

  Product _mapToProduct(Map<String, dynamic> data) {
    return Product(
      id: data['id'] as int,
      name: data['name'] as String,
      category: data['category'] as String,
      price: data['price'] as double,
      image: data['image'] as String? ?? '',
    );
  }

  List<String> get _chips => const ['All', 'Fruits', 'Veggies', 'Herbs'];

  // 🎯 (ฟังก์ชัน _load() นี้ถูกต้องแล้ว)
  Future<void> _load() async {
    final Map<String, List<String>> filters = {
      'categories': _category == 'All' ? [] : [_category],
      'brands': [],
    };
    final list = await _repo.list(q: '', filters: filters);
    if (mounted) {
      setState(() => _items = list);
    }
  }

  @override
  void initState() {
    super.initState();
    _load();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: _scaffoldBgColor, // <--- 🎨 แก้ไขสีพื้นหลัง
      body: RefreshIndicator(
        onRefresh: _load,
        child: LayoutBuilder(
          builder: (context, constraints) {
            final width = constraints.maxWidth;
            int crossAxisCount = (width > 600) ? 3 : 2;
            double aspectRatio = (width > 600) ? 0.75 : 0.72;

            return CustomScrollView(
              slivers: [
                // --- ‼️ ลบ SliverAppBar และแทนที่ด้วย Title + Search ---
                SliverToBoxAdapter(
                  child: SafeArea(
                    bottom: false,
                    child: Padding(
                      padding: const EdgeInsets.only(
                        top: 16.0,
                      ), // <--- 📏 ระยะห่าง
                      child: Center(
                        child: Text(
                          'FreshEase',
                          style: Theme.of(context).textTheme.headlineSmall
                              ?.copyWith(fontWeight: FontWeight.w700),
                        ),
                      ),
                    ),
                  ),
                ),
                // --- ช่องค้นหา (Search Store) ---
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(
                      16,
                      16,
                      16,
                      16,
                    ), // <--- 📏 ระยะห่าง
                    child: SearchPill(
                      readOnly: true,
                      showFilter: false,
                      onTap: () {
                        context.go('/explore');
                      },
                    ),
                  ),
                ),
                // --- จบส่วน Title/Search ใหม่ ---

                // ‼️ --- 1. แสดง Banner/Offers เฉพาะหน้า 'All' ---
                if (_category == 'All') ...[
                  // --- แบนเนอร์ (Fresh Vegetables) ---
                  SliverToBoxAdapter(child: _buildBanner()),

                  // --- Exclusive Offer Section ---
                  SliverToBoxAdapter(
                    child: Padding(
                      padding: const EdgeInsets.fromLTRB(
                        16,
                        0,
                        16,
                        0,
                      ), // <--- 📏 ระยะห่าง
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
                            onTap: () {
                              context.go('/explore');
                            },
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
                      height: 270, // <--- 📏 ความสูง (ปรับได้ตาม UI)
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
                            width: 150, // <--- 📏 ความกว้างการ์ด (ปรับได้)
                            child: ProductCard(
                              product: p,
                              productDetail: data['description'],
                              onTap: () =>
                                  context.go('/shop/product/${p.id}', extra: p),
                              onAdd: () {
                                CartScope.of(context).add(p);
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

                  // --- ‼️ Best Selling Section (แก้ไขตามรูป) ---
                  SliverToBoxAdapter(
                    child: Padding(
                      padding: const EdgeInsets.fromLTRB(
                        16,
                        0,
                        16,
                        0,
                      ), // <--- 📏 ระยะห่าง
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
                            onTap: () {
                              context.go('/explore');
                            },
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
                  // --- ‼️ เพิ่ม List แนวนอนสำหรับ Best Selling ---
                  SliverToBoxAdapter(
                    child: SizedBox(
                      height: 270, // <--- 📏 ความสูง (ปรับได้)
                      child: ListView.separated(
                        scrollDirection: Axis.horizontal,
                        padding: const EdgeInsets.all(16),
                        itemCount: _bestSellingData.length,
                        separatorBuilder: (context, index) =>
                            const SizedBox(width: 12),
                        itemBuilder: (context, index) {
                          final data = _bestSellingData[index];
                          final p = _mapToProduct(data);
                          return SizedBox(
                            width: 150,
                            child: ProductCard(
                              product: p,
                              productDetail: data['description'],
                              onTap: () =>
                                  context.go('/shop/product/${p.id}', extra: p),
                              onAdd: () {
                                CartScope.of(context).add(p);
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
                ], // ‼️ --- สิ้นสุดการซ่อน (if) ---
                // --- Chip Filters (แสดงตลอด) ---
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
                          backgroundColor: Colors.white, // <--- 🎨
                          selectedColor: _selectedChipColor,
                          shape: StadiumBorder(
                            side: BorderSide(
                              color: selected
                                  ? _primaryColor
                                  : Colors.grey.shade300,
                            ),
                          ),
                          labelStyle: TextStyle(
                            color: selected
                                ? _primaryColor // <--- 🎨
                                : Colors.grey.shade700,
                            fontWeight: selected
                                ? FontWeight.w700
                                : FontWeight.w500,
                          ),
                        );
                      },
                    ),
                  ),
                ),

                // --- Product Grid (แสดงตลอด และอัปเดตตาม _items) ---
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
                          CartScope.of(context).add(p);
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
                      mainAxisSpacing: 12, // <--- 📏 เพิ่มช่องไฟ
                      crossAxisSpacing: 12, // <--- 📏 เพิ่มช่องไฟ
                      childAspectRatio: aspectRatio,
                    ),
                  ),
                ),
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

  // ‼️ --- Widget สำหรับสร้างแบนเนอร์ ---
  Widget _buildBanner() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(
        16,
        0,
        16,
        16,
      ), // <--- 📏 แก้ไขระยะห่าง
      child: AspectRatio(
        aspectRatio: 16 / 7,
        child: Stack(
          children: [
            PageView.builder(
              controller: _bannerController,
              itemCount: _bannerImages.length,
              onPageChanged: (index) {
                setState(() => _bannerCurrentPage = index);
              },
              itemBuilder: (context, index) {
                return Container(
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(12),
                    image: DecorationImage(
                      image: AssetImage(_bannerImages[index]),
                      fit: BoxFit.cover,
                      onError: (e, stack) =>
                          debugPrint('Banner image failed to load'),
                    ),
                  ),
                );
              },
            ),
            // ตัวบอกตำแหน่ง (Dots)
            Positioned(
              bottom: 10,
              left: 0,
              right: 0,
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: List.generate(_bannerImages.length, (index) {
                  return Container(
                    width: 8,
                    height: 8,
                    margin: const EdgeInsets.symmetric(horizontal: 4),
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      color: _bannerCurrentPage == index
                          ? _primaryColor // <--- 🎨
                          : Colors.grey.shade400, // <--- 🎨
                    ),
                  );
                }),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
