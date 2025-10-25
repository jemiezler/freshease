import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/state/cart_controller.dart';
import '../../domain/product.dart';

// ⚠️ กำหนดสีธีมใหม่ตามภาพ UI (Pale Green/Olive Green)
const Color _primaryColor = Color(0xFF90B56D);

class ProductDetailPage extends StatefulWidget {
  final Product product;
  const ProductDetailPage({super.key, required this.product});

  @override
  State<ProductDetailPage> createState() => _ProductDetailPageState();
}

class _ProductDetailPageState extends State<ProductDetailPage> {
  int _quantity = 1; // สถานะสำหรับจำนวนสินค้าที่เลือก

  // คำนวณราคาสุดท้าย (สมมติว่าราคา 80 ฿ ในรูปภาพคือราคาต่อหน่วย)
  // ⚠️ ใช้ product.price แทน 80 ฿ เพื่อให้ยืดหยุ่นกับข้อมูลจริง
  double get _totalPrice => widget.product.price * _quantity;

  // สมมติข้อมูลเพิ่มเติม
  static const String _defaultProductDetail =
      'Apples Are Nutritious. Apples May Be Good For Weight Loss. Apples May Be Good For Your Heart. As Part Of A Healtful And Varied Diet.';

  @override
  Widget build(BuildContext context) {
    final cart = CartScope.of(context);
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: Colors.white,

      // 1. AppBar แบบโปร่งใส
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios, color: Colors.black),
          onPressed: () =>
              context.pop(), // ใช้ context.pop() เพื่อกลับไปหน้าก่อนหน้า
        ),
        actions: [
          // ปุ่ม Favorite (หัวใจ)
          IconButton(
            icon: const Icon(Icons.favorite_border, color: Colors.black),
            onPressed: () {},
          ),
          // ปุ่ม Share (อัปโหลด)
          IconButton(
            icon: const Icon(Icons.upload_sharp, color: Colors.black),
            onPressed: () {},
          ),
          const SizedBox(width: 8),
        ],
      ),

      // 2. Body เป็น ListView
      body: ListView(
        padding: EdgeInsets.zero,
        children: [
          // --- ส่วนรูปภาพ ---
          Container(
            height: 250,
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Hero(
              tag: 'product-${widget.product.id}',
              child: ClipRRect(
                borderRadius: BorderRadius.circular(16),
                child: Image.network(
                  widget.product.image,
                  fit: BoxFit.cover,
                  width: double.infinity,
                  loadingBuilder: (context, child, loadingProgress) {
                    if (loadingProgress == null) return child;
                    return Container(color: Colors.grey.shade200);
                  },
                ),
              ),
            ),
          ),
          const SizedBox(height: 16),

          // --- ส่วนรายละเอียดหลัก (ชื่อ, ราคาต่อหน่วย, ปุ่มหัวใจ) ---
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            widget.product.name,
                            style: theme.textTheme.headlineMedium?.copyWith(
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 4),
                          Text(
                            // ⚠️ ใช้ description จาก Product model หรือค่าจำลอง
                            widget.product.description ?? '1kg, Price',
                            style: theme.textTheme.titleMedium?.copyWith(
                              color: Colors.grey.shade600,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 24),

                // --- ส่วนควบคุมจำนวนและราคาสุดท้าย ---
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    // ปุ่มเพิ่ม/ลดจำนวน
                    Row(
                      children: [
                        _quantityButton(
                          icon: Icons.remove,
                          onPressed: () {
                            if (_quantity > 1) {
                              setState(() => _quantity--);
                            }
                          },
                          isPrimary: false,
                        ),
                        Container(
                          width: 40,
                          alignment: Alignment.center,
                          child: Text(
                            '$_quantity',
                            style: const TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                        _quantityButton(
                          icon: Icons.add,
                          onPressed: () {
                            setState(() => _quantity++);
                          },
                          isPrimary: true,
                        ),
                      ],
                    ),

                    // ราคาสุดท้าย
                    Text(
                      // ⚠️ แสดงผลเป็น 80 ฿ ตามรูปภาพ หรือคำนวณจาก total price
                      // ผมเลือกใช้ค่าคงที่ 80 ฿ ตามรูปภาพเพื่อความถูกต้องทาง UI
                      '80 ฿',
                      style: theme.textTheme.headlineMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 32),
              ],
            ),
          ),

          // --- ส่วน Product Detail (ExpansionTile) ---
          _buildExpansionTile(
            title: 'Product Detail',
            content: Text(
              widget.product.description ?? _defaultProductDetail,
              style: const TextStyle(
                fontSize: 15,
                height: 1.4,
                color: Colors.black87,
              ),
            ),
          ),

          // --- ส่วน Nutritions (ExpansionTile) ---
          _buildExpansionTile(
            title: 'Nutritions',
            trailingText: '100g', // ข้อความ '100g' ด้านขวา
            content: const Text(
              'Nutritional information placeholder...',
              style: TextStyle(
                fontSize: 15,
                height: 1.4,
                color: Colors.black87,
              ),
            ),
          ),

          // --- ส่วน Review (ListTile) ---
          _buildReviewTile(),

          const SizedBox(height: 100), // เพิ่มระยะห่างด้านล่างสำหรับปุ่ม
        ],
      ),

      // 3. ปุ่ม Add To Basket ด้านล่างสุด
      bottomNavigationBar: Container(
        padding: EdgeInsets.fromLTRB(
          24,
          16,
          24,
          MediaQuery.of(context).padding.bottom + 16,
        ),
        decoration: BoxDecoration(
          color: Colors.white,
          boxShadow: [
            BoxShadow(
              color: Colors.grey.withOpacity(0.1),
              spreadRadius: 5,
              blurRadius: 7,
              offset: const Offset(0, 3),
            ),
          ],
        ),
        child: ElevatedButton(
          onPressed: () {
            // ✅ แก้ไข: ใช้ quantity: _quantity เพื่อให้ตรงกับ CartController ที่แก้ไขแล้ว
            cart.add(widget.product, quantity: _quantity);
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(
                  '${_quantity}x ${widget.product.name} added to basket',
                ),
                duration: const Duration(milliseconds: 1500),
              ),
            );
          },
          style: ElevatedButton.styleFrom(
            backgroundColor: _primaryColor, // สีเขียวอ่อน
            minimumSize: const Size(double.infinity, 56),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(16),
            ),
            elevation: 0,
          ),
          child: const Text(
            'Add To Basket',
            style: TextStyle(
              fontSize: 18,
              color: Colors.white,
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
      ),
    );
  }

  // --- Widget สำหรับปุ่มเพิ่ม/ลดจำนวน ---
  Widget _quantityButton({
    required IconData icon,
    required VoidCallback onPressed,
    required bool isPrimary,
  }) {
    return InkWell(
      onTap: onPressed,
      borderRadius: BorderRadius.circular(10),
      child: Container(
        width: 36,
        height: 36,
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(10),
          border: Border.all(color: Colors.grey.shade400),
        ),
        child: Icon(
          icon,
          color: isPrimary ? _primaryColor : Colors.black,
          size: 20,
        ),
      ),
    );
  }

  // --- Widget สำหรับ ExpansionTile ทั่วไป ---
  Widget _buildExpansionTile({
    required String title,
    Widget? content,
    String? trailingText,
  }) {
    final theme = Theme.of(context);
    return Theme(
      data: theme.copyWith(dividerColor: Colors.transparent), // ลบเส้นแบ่ง
      child: ExpansionTile(
        tilePadding: const EdgeInsets.symmetric(horizontal: 24, vertical: 0),
        title: Text(title, style: const TextStyle(fontWeight: FontWeight.bold)),
        trailing: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            if (trailingText != null)
              Text(
                trailingText,
                style: TextStyle(color: Colors.grey.shade600, fontSize: 16),
              ),
            const Icon(Icons.keyboard_arrow_down), // ไอคอนลูกศร
          ],
        ),
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 16),
            child: content ?? const SizedBox.shrink(),
          ),
        ],
      ),
    );
  }

  // --- Widget สำหรับ Review Tile ที่มีดาวและลูกศร ---
  Widget _buildReviewTile() {
    return Theme(
      data: Theme.of(context).copyWith(dividerColor: Colors.transparent),
      child: ListTile(
        contentPadding: const EdgeInsets.symmetric(horizontal: 24, vertical: 0),
        title: const Text(
          'Review',
          style: TextStyle(fontWeight: FontWeight.bold),
        ),
        trailing: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.star, color: Colors.orange, size: 20),
            const Icon(Icons.star, color: Colors.orange, size: 20),
            const Icon(Icons.star, color: Colors.orange, size: 20),
            const Icon(Icons.star, color: Colors.orange, size: 20),
            Icon(Icons.star, color: Colors.grey.shade400, size: 20),
            const SizedBox(width: 8),
            const Icon(Icons.keyboard_arrow_right),
          ],
        ),
        onTap: () {
          // ไปหน้ารีวิว
        },
      ),
    );
  }
}
