// ไฟล์: explore_page.dart

// ignore_for_file: constant_identifier_names

import 'dart:async';

import 'package:flutter/material.dart';
import 'package:frontend/core/state/cart_controller.dart';
import 'package:frontend/features/shop/data/product_repository.dart';
import 'package:frontend/features/shop/domain/product.dart';
import 'package:go_router/go_router.dart';

import '../../widgets/filter_bottom_sheet.dart';
import '../../widgets/product_card.dart';
import '../../widgets/search_pill_widget.dart';

const Color AppColors_primary = Color(0xFF90B56D);

class ExplorePage extends StatefulWidget {
  // 1. เพิ่มตัวแปรสำหรับรับค่า category จาก "See all"
  final String? category;

  const ExplorePage({super.key, this.category});

  @override
  State<ExplorePage> createState() => _ExplorePageState();
}

class _ExplorePageState extends State<ExplorePage> {
  final _repo = MockProductRepository();
  final _search = TextEditingController();
  Timer? _debounce;
  List<Product> _items = [];

  // สถานะสำหรับ Filter
  Set<String> _selectedCategories = {};
  String? _selectedBrand;

  @override
  void initState() {
    super.initState();

    // 2. ตรวจสอบว่ามี "category" ถูกส่งมาจาก "See all" หรือไม่
    if (widget.category != null && widget.category!.isNotEmpty) {
      _selectedCategories = {widget.category!};
      // โหลดสินค้าตาม category ทันที
      _load();
    }

    // 3. Listener สำหรับช่องค้นหา
    _search.addListener(() {
      _debounce?.cancel();
      _debounce = Timer(const Duration(milliseconds: 300), () {
        setState(() {
          // ถ้าผู้ใช้เริ่มพิมพ์ ให้ล้างการเลือกหมวดหมู่/แบรนด์
          if (_search.text.isNotEmpty) {
            _selectedCategories.clear();
            _selectedBrand = null;
          }
          _load();
        });
      });
      // อัปเดต UI ทันทีที่พิมพ์เพื่อซ่อน Grid
      setState(() {});
    });
  }

  @override
  void dispose() {
    _debounce?.cancel();
    _search.dispose();
    super.dispose();
  }

  // 4. ฟังก์ชัน _load() ที่ถูกต้อง (ใช้ q และ category)
  Future<void> _load() async {
    final String query = _search.text;
    final String category = _selectedCategories.isNotEmpty
        ? _selectedCategories.first
        : 'All';

    final list = await _repo.list(q: query, category: category);

    if (mounted) {
      setState(() => _items = list);
    }
  }

  // 5. ฟังก์ชันเปิด Filter (เหมือนเดิม)
  void _openFilterSheet() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) {
        return FilterBottomSheet(
          initialCategories: _selectedCategories.toList(),
          initialBrand: _selectedBrand,
          onApplyFilter: (newCategories, newBrand) {
            setState(() {
              _selectedCategories = newCategories.toSet();
              _selectedBrand = newBrand;
              _search.clear(); // ล้างช่องค้นหาเมื่อใช้ Filter
            });
            _load(); // โหลดใหม่
            Navigator.pop(context);
          },
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    // 6. ตรวจสอบสถานะการแสดงผล
    final bool isSearching = _search.text.isNotEmpty;
    final bool isFiltering =
        !isSearching &&
        (_selectedCategories.isNotEmpty || _selectedBrand != null);

    return Scaffold(
      backgroundColor: Colors.white,
      body: CustomScrollView(
        slivers: [
          SliverAppBar(
            backgroundColor: Colors.white,
            elevation: 0,
            pinned: true,

            // 🎯 7. ‼️ นี่คือปุ่มย้อนกลับ ‼️
            leading: IconButton(
              icon: const Icon(Icons.arrow_back_ios_new, color: Colors.black),
              onPressed: () => context.pop(), // สั่ง GoRouter ให้ย้อนกลับ
            ),

            centerTitle: true,
            title: Text(
              'Find Products',
              style: Theme.of(
                context,
              ).textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.w700),
            ),
            bottom: PreferredSize(
              preferredSize: const Size.fromHeight(64),
              child: Padding(
                padding: const EdgeInsets.fromLTRB(16, 0, 16, 12),
                child: SearchPill(
                  controller: _search,
                  readOnly: false, // 👈 พิมพ์ได้
                  showFilter: true, // 👈 แสดง Filter
                  onFilterTap: _openFilterSheet, // 👈 ผูกกับ Filter
                ),
              ),
            ),
          ),

          // 8. แสดงผลลัพธ์ตามเงื่อนไข
          if (isSearching || isFiltering)
            ..._buildSearchResults() // 👈 แสดงผลการค้นหา/กรอง
          else
            ..._buildCategoryGrid(), // 👈 แสดง Grid หมวดหมู่ (เมื่อหน้าว่าง)
        ],
      ),
    );
  }

  // --- Widget Builder สำหรับ Grid หมวดหมู่ ---
  List<Widget> _buildCategoryGrid() {
    return [
      SliverPadding(
        padding: const EdgeInsets.all(16.0),
        sliver: SliverGrid(
          gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
            crossAxisCount: 2,
            mainAxisSpacing: 16,
            crossAxisSpacing: 16,
            childAspectRatio: 0.9,
          ),
          delegate: SliverChildBuilderDelegate((context, index) {
            final category = _categories[index];
            return _CategoryCard(
              title: category['name'] as String,
              imagePath: category['image'] as String,
              borderColor: category['color'] as Color,
              // 9. เมื่อกดการ์ด ให้ตั้งค่า Category และโหลด
              onTap: () {
                setState(() {
                  _search.clear(); // ล้างข้อความค้นหา
                  _selectedCategories = {category['key'] as String};
                });
                _load(); // สั่งโหลดข้อมูลใหม่
              },
            );
          }, childCount: _categories.length),
        ),
      ),
    ];
  }

  // --- Widget Builder สำหรับ ผลการค้นหา ---
  List<Widget> _buildSearchResults() {
    final width = MediaQuery.of(context).size.width;
    int crossAxisCount = (width > 600) ? 3 : 2;
    double aspectRatio = (width > 600) ? 0.75 : 0.72;

    if (_items.isEmpty) {
      return [
        const SliverFillRemaining(
          hasScrollBody: false,
          child: Center(child: Text('No results found')),
        ),
      ];
    }

    return [
      SliverPadding(
        padding: const EdgeInsets.all(16),
        sliver: SliverGrid(
          delegate: SliverChildBuilderDelegate((context, i) {
            final p = _items[i];
            return ProductCard(
              product: p,
              onTap: () => context.go('/shop/product/${p.id}', extra: p),
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
            mainAxisSpacing: 0,
            crossAxisSpacing: 0,
            childAspectRatio: aspectRatio,
          ),
        ),
      ),
    ];
  }
}

// --- Widget สำหรับการ์ดหมวดหมู่ (ตามดีไซน์ใหม่) ---
class _CategoryCard extends StatelessWidget {
  final String title;
  final String imagePath;
  final Color borderColor;
  final VoidCallback onTap;

  const _CategoryCard({
    required this.title,
    required this.imagePath,
    required this.borderColor,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Semantics(
      label: 'Category, $title',
      button: true,
      child: Container(
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(20),
          border: Border.all(color: borderColor, width: 1.5),
          boxShadow: [
            BoxShadow(
              color: Colors.grey.withValues(alpha: 0.05),
              blurRadius: 10,
              offset: const Offset(0, 5),
            ),
          ],
        ),
        child: Material(
          color: Colors.transparent,
          borderRadius: BorderRadius.circular(20),
          child: InkWell(
            onTap: onTap,
            borderRadius: BorderRadius.circular(20),
            splashColor: borderColor.withValues(alpha: 0.3),
            highlightColor: borderColor.withValues(alpha: 0.1),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Image.asset(
                  imagePath,
                  height: 80,
                  width: 100,
                  fit: BoxFit.contain,
                  errorBuilder: (context, error, stackTrace) {
                    return const Icon(
                      Icons.category_outlined,
                      size: 60,
                      color: Colors.grey,
                    );
                  },
                ),
                const SizedBox(height: 12),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: Text(
                    title,
                    textAlign: TextAlign.center,
                    style: const TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 16,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

// --- ข้อมูลหมวดหมู่ (อ้างอิงจากภาพใหม่) ---
final List<Map<String, dynamic>> _categories = const [
  {
    'name': 'Fresh Fruits & Vegetable',
    'color': Color(0xFF53B175), // Green
    'image': 'assets/images/cat_fruits.png', // 👈 ต้องมีรูปนี้
    'key': 'Fruits',
  },
  {
    'name': 'Cooking Oil & Ghee',
    'color': Color(0xFFF8A44C), // Orange
    'image': 'assets/images/cat_oil.png', // 👈 ต้องมีรูปนี้
    'key': 'Oil',
  },
  {
    'name': 'Meat & Fish',
    'color': Color(0xFFF7A593), // Red
    'image': 'assets/images/cat_meat.png', // 👈 ต้องมีรูปนี้
    'key': 'Meat',
  },
  {
    'name': 'Bakery & Snacks',
    'color': Color(0xFFD3B0E0), // Purple
    'image': 'assets/images/cat_bakery.png', // 👈 ต้องมีรูปนี้
    'key': 'Bakery',
  },
  {
    'name': 'Dairy & Eggs',
    'color': Color(0xFFFDE598), // Yellow
    'image': 'assets/images/cat_dairy.png', // 👈 ต้องมีรูปนี้
    'key': 'Dairy',
  },
  {
    'name': 'Beverages',
    'color': Color(0xFFB7DFF5), // Blue
    'image': 'assets/images/cat_beverages.png', // 👈 ต้องมีรูปนี้
    'key': 'Beverages',
  },
];
