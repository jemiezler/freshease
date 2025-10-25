// ไฟล์: shoppage.dart

import 'package:flutter/material.dart';
import 'package:frontend/core/state/cart_controller.dart'; // ตรวจสอบว่า import ถูกต้อง
import 'package:frontend/features/shop/widgets/product_card.dart'; // ตรวจสอบว่า import ถูกต้อง
import 'package:go_router/go_router.dart';

// ตรวจสอบว่า path เหล่านี้ถูกต้องตามโครงสร้างโปรเจกต์ของคุณ
import '../../../shop/data/product_repository.dart';
import '../../../shop/domain/product.dart';
import '../../widgets/search_pill_widget.dart';

// ⚠️ กำหนดสีธีมใหม่ตามภาพ UI
const Color _primaryColor = Color(0xFF53B175);
const Color _selectedChipColor = Color(0xFFE8F5E9);
const Color _scaffoldBgColor = Color(0xFFFAFAFA);
const Color AppColors_primary = _primaryColor;

class ShopPage extends StatefulWidget {
  const ShopPage({super.key});
  @override
  State<ShopPage> createState() => _ShopPageState();
}

class _ShopPageState extends State<ShopPage> {
  // --- สถานะ (State) ---
  final _repo = MockProductRepository(); // หรือ Repository จริงของคุณ
  String _category = 'All'; // สำหรับ Chip Filter
  List<Product> _items = []; // รายการสินค้าที่จะแสดงใน Grid

  // สถานะสำหรับ Banner
  int _bannerCurrentPage = 0;
  final PageController _bannerController = PageController();
  // ใช้ URL สำหรับ Banner
  final List<String> _bannerImages = [
    'https://d1csarkz8obe9u.cloudfront.net/posterpreviews/healthy-vegetables-banner-design-template-21a9d6f7102f16ddd973d540d30bbe83_screen.jpg?ts=1758505842',
    'https://d1csarkz8obe9u.cloudfront.net/posterpreviews/healthy-food-restaurant-banner-design-template-5d8526f015d6a01027536b17714b98d3_screen.jpg?ts=1662349433',
    'https://d1csarkz8obe9u.cloudfront.net/posterpreviews/fresh-vegetables-flyer-design-template-0396f1a5981cef834fe21743c77d8dfe_screen.jpg?ts=1621616325',
  ];

  // ข้อมูลจำลอง Exclusive Offers (ใช้ URL)
  final List<Map<String, dynamic>> _exclusiveOffersData = const [
    {
      'id': 6,
      'name': 'Organic Bananas',
      'category': 'Fruits',
      'description': '7pcs, Priceg',
      'price': 4.99,
      'image':
          'https://i.pinimg.com/736x/02/49/5f/02495fb1b8bd32a24fb8eb483a18a074.jpg',
    },
    {
      'id': 7,
      'name': 'Red Apple',
      'category': 'Fruits',
      'description': '1kg, Priceg',
      'price': 4.99,
      'image':
          'https://i.pinimg.com/736x/cc/4c/8e/cc4c8e9e9c9ee1bab48b41f1863e971e.jpg',
    },
  ];

  // ข้อมูลจำลอง Best Selling (ใช้ URL)
  final List<Map<String, dynamic>> _bestSellingData = const [
    {
      'id': 8,
      'name': 'Red Pepper',
      'category': 'Veggies',
      'description': '1kg, Priceg',
      'price': 4.99,
      'image':
          'https://i.pinimg.com/736x/41/8a/4d/418a4dba2668bf8446094fdaf94fe85e.jpg',
    },
    {
      'id': 6,
      'name': 'Organic Bananas',
      'category': 'Fruits',
      'description': '7pcs, Priceg',
      'price': 4.99,
      'image':
          'https://i.pinimg.com/736x/02/49/5f/02495fb1b8bd32a24fb8eb483a18a074.jpg',
    },
  ];

  // --- เมธอด (Methods) ---

  // แปลง Map เป็น Object Product
  Product _mapToProduct(Map<String, dynamic> data) {
    return Product(
      id: data['id'] as int,
      name: data['name'] as String,
      category: data['category'] as String,
      price: data['price'] as double,
      image: data['image'] as String? ?? '', // รับ URL หรือ Asset Path
    );
  }

  // รายการ Chip Filter
  List<String> get _chips => const ['All', 'Fruits', 'Veggies', 'Herbs'];

  // โหลดข้อมูลสินค้าตาม Filter ที่เลือก
  Future<void> _load() async {
    // สร้าง filters map สำหรับส่งให้ Repository
    final Map<String, List<String>> filters = {
      'categories': _category == 'All'
          ? []
          : [_category], // ส่ง category ที่เลือก (ถ้าไม่ใช่ 'All')
      'brands': [], // หน้านี้ไม่มี filter brand
    };
    // เรียก Repository (ส่ง q ว่างเปล่า และ filters)
    final list = await _repo.list(q: '', filters: filters);
    // อัปเดต UI ถ้ายังอยู่ในหน้านี้
    if (mounted) {
      setState(() => _items = list);
    }
  }

  // --- Lifecycle ---

  @override
  void initState() {
    super.initState();
    _load(); // โหลดข้อมูลครั้งแรกเมื่อหน้าเปิด
  }

  // --- UI Build ---

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: _scaffoldBgColor, // ใช้สีพื้นหลังที่กำหนด
      body: RefreshIndicator(
        onRefresh: _load, // ดึงลงเพื่อโหลดใหม่
        child: LayoutBuilder(
          builder: (context, constraints) {
            // คำนวณ Grid layout ตามความกว้างหน้าจอ
            final width = constraints.maxWidth;
            int crossAxisCount = (width > 600) ? 3 : 2; // จำนวนคอลัมน์
            double aspectRatio = (width > 600) ? 0.75 : 0.72; // สัดส่วนการ์ด

            return CustomScrollView(
              slivers: [
                // --- ส่วนหัว: Title ---
                SliverToBoxAdapter(
                  child: SafeArea(
                    bottom: false,
                    child: Padding(
                      padding: const EdgeInsets.only(top: 16.0),
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
                // --- ส่วนหัว: ช่องค้นหา ---
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(16, 16, 16, 16),
                    child: SearchPill(
                      readOnly: true, // กดได้อย่างเดียว พิมพ์ไม่ได้
                      showFilter: false, // ไม่มีปุ่ม Filter
                      onTap: () {
                        // กดแล้วไปหน้า Explore
                        context.go('/explore');
                      },
                    ),
                  ),
                ),

                // --- แสดง Banner, Exclusive, Best Selling เฉพาะเมื่อเลือก 'All' ---
                if (_category == 'All') ...[
                  // --- แบนเนอร์ ---
                  SliverToBoxAdapter(child: _buildBanner()),

                  // --- Exclusive Offer Section ---
                  SliverToBoxAdapter(
                    child: _buildSectionHeader(
                      title: 'Exclusive Offer',
                      onSeeAllTap: () {
                        // ไปหน้า Explore พร้อมกรอง Category 'Fruits'
                        context.go('/explore?category=Fruits');
                      },
                    ),
                  ),
                  SliverToBoxAdapter(
                    child: _buildHorizontalProductList(_exclusiveOffersData),
                  ),

                  // --- Best Selling Section ---
                  SliverToBoxAdapter(
                    child: _buildSectionHeader(
                      title: 'Best Selling',
                      onSeeAllTap: () {
                        // ไปหน้า Explore พร้อมกรอง Category 'Veggies' (ตัวอย่าง)
                        context.go('/explore?category=Veggies');
                      },
                    ),
                  ),
                  SliverToBoxAdapter(
                    child: _buildHorizontalProductList(_bestSellingData),
                  ),
                ], // --- สิ้นสุด if (_category == 'All') ---
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
                            // เมื่อเลือก Chip ให้เปลี่ยน _category และโหลดใหม่
                            setState(() => _category = label);
                            _load();
                          },
                          backgroundColor: Colors.white,
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
                                ? _primaryColor
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

                // --- Product Grid (แสดงผลลัพธ์จาก _load) ---
                SliverPadding(
                  padding: const EdgeInsets.all(16),
                  sliver: SliverGrid(
                    delegate: SliverChildBuilderDelegate((context, i) {
                      final p = _items[i]; // สินค้าจาก State _items
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
                      mainAxisSpacing: 12, // ระยะห่างแนวตั้ง
                      crossAxisSpacing: 12, // ระยะห่างแนวนอน
                      childAspectRatio: aspectRatio, // สัดส่วนการ์ด
                    ),
                  ),
                ),

                // --- แสดงข้อความเมื่อไม่มีสินค้า ---
                if (_items.isEmpty)
                  const SliverFillRemaining(
                    hasScrollBody: false, // ไม่ต้องเลื่อนถ้าเนื้อหาไม่พอ
                    child: Center(child: Text('No results')),
                  ),
              ],
            );
          },
        ),
      ),
    );
  }

  // --- Helper Widgets ---

  // Widget สำหรับสร้างส่วนหัว Section (เช่น "Exclusive Offer")
  Widget _buildSectionHeader({
    required String title,
    VoidCallback? onSeeAllTap,
  }) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 0), // เพิ่มระยะห่างด้านบน
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            title,
            style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          if (onSeeAllTap != null) // แสดงปุ่ม "See all" ถ้ามีฟังก์ชัน onTap
            GestureDetector(
              onTap: onSeeAllTap,
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
    );
  }

  // Widget สำหรับสร้าง List สินค้าแนวนอน
  Widget _buildHorizontalProductList(List<Map<String, dynamic>> productData) {
    return SizedBox(
      height: 270, // ความสูงของ List
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        padding: const EdgeInsets.all(16),
        itemCount: productData.length,
        separatorBuilder: (context, index) => const SizedBox(width: 12),
        itemBuilder: (context, index) {
          final data = productData[index];
          final p = _mapToProduct(data);
          return SizedBox(
            width: 150, // ความกว้างของการ์ด
            child: ProductCard(
              product: p,
              productDetail: data['description'],
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
            ),
          );
        },
      ),
    );
  }

  // Widget สำหรับสร้าง Banner
  Widget _buildBanner() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
      child: AspectRatio(
        aspectRatio: 16 / 7, // สัดส่วนแบนเนอร์
        child: Stack(
          children: [
            PageView.builder(
              controller: _bannerController,
              itemCount: _bannerImages.length,
              onPageChanged: (index) {
                // อัปเดต State เมื่อเปลี่ยนหน้า Banner
                setState(() => _bannerCurrentPage = index);
              },
              itemBuilder: (context, index) {
                final imageUrl = _bannerImages[index];
                return Container(
                  clipBehavior: Clip.antiAlias, // ทำให้ขอบมนมีผลกับ Image
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Image.network(
                    // ใช้ Image.network โดยตรง
                    imageUrl,
                    fit: BoxFit.cover,
                    // แสดง Loading Indicator
                    loadingBuilder: (context, child, loadingProgress) {
                      if (loadingProgress == null) return child;
                      return const Center(
                        child: CircularProgressIndicator(
                          valueColor: AlwaysStoppedAnimation<Color>(
                            _primaryColor,
                          ),
                          strokeWidth: 2,
                        ),
                      );
                    },
                    // แสดง Icon รูปเสีย
                    errorBuilder: (context, error, stackTrace) {
                      debugPrint(
                        'Banner image failed to load: $imageUrl\nError: $error',
                      );
                      return const Center(
                        child: Icon(
                          Icons.broken_image_outlined,
                          color: Colors.grey,
                          size: 40,
                        ),
                      );
                    },
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
                          ? _primaryColor
                          : Colors.grey.shade400,
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
} // --- สิ้นสุดคลาส _ShopPageState ---
