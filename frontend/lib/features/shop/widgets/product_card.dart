// ไฟล์: product_card.dart

import 'package:flutter/material.dart';
import 'package:frontend/features/shop/domain/product.dart'; // ตรวจสอบ import นี้

// ‼️ ดึงสีหลัก (อาจต้อง import จาก shoppage หรือประกาศใหม่)
// const Color _primaryColor = Color(0xFF53B175); // ใช้สีเขียวใหม่จาก shoppage
const Color _primaryColor = Color(0xFF90B56D); // หรือใช้สีเขียวเดิมที่คุณมี

class ProductCard extends StatelessWidget {
  final Product product;
  final String? productDetail;
  final VoidCallback onTap;
  final VoidCallback onAdd;

  const ProductCard({
    super.key,
    required this.product,
    this.productDetail,
    required this.onTap,
    required this.onAdd,
  });

  @override
  Widget build(BuildContext context) {
    // ‼️ 1. ตรวจสอบว่า image เป็น URL หรือ Asset Path
    final bool isNetworkImage = product.image.startsWith('http');

    // ‼️ 2. สร้าง ImageProvider ที่ถูกต้อง
    final ImageProvider imageProvider;
    if (isNetworkImage) {
      imageProvider = NetworkImage(product.image); // 👈 ถ้าเป็น URL
    } else {
      // 👈 ถ้าเป็น Asset Path
      // (ต้องแน่ใจว่าประกาศ Assets ใน pubspec.yaml ถูกต้อง)
      imageProvider = AssetImage(product.image);
    }

    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(18),
      child: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(18),
          border: Border.all(color: Colors.grey.shade200, width: 1.5),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // --- รูปสินค้า ---
            Expanded(
              flex: 3,
              child: Center(
                // ‼️ 3. ใช้ Image.provider แทน Image.asset
                child: Image(
                  image: imageProvider, // 👈 ใช้ Provider ที่สร้างไว้
                  fit: BoxFit.contain, // หรือ BoxFit.cover ตามต้องการ
                  // ‼️ 4. เพิ่ม Loading Builder (สำคัญสำหรับ NetworkImage)
                  loadingBuilder: (context, child, loadingProgress) {
                    if (loadingProgress == null) return child; // โหลดเสร็จ
                    return Center(
                      // กำลังโหลด...
                      child: CircularProgressIndicator(
                        valueColor: const AlwaysStoppedAnimation<Color>(
                          _primaryColor,
                        ),
                        strokeWidth: 2,
                        value: loadingProgress.expectedTotalBytes != null
                            ? loadingProgress.cumulativeBytesLoaded /
                                  loadingProgress.expectedTotalBytes!
                            : null,
                      ),
                    );
                  },

                  // ‼️ 5. errorBuilder (สำคัญมาก)
                  errorBuilder: (ctx, err, stack) {
                    debugPrint(
                      '!!! ProductCard image error: ${product.image}\nError: $err',
                    );
                    return const Icon(
                      Icons.broken_image_outlined, // Icon รูปเสีย
                      color: Colors.grey,
                      size: 40,
                    );
                  },
                ),
              ),
            ),
            const SizedBox(height: 12),

            // --- ชิ่อสินค้า ---
            Text(
              product.name,
              style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const SizedBox(height: 4),

            // --- รายละเอียด ---
            Text(
              productDetail ?? '${product.category} item',
              style: TextStyle(color: Colors.grey.shade600, fontSize: 14),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const Spacer(), // ดันราคา/ปุ่มไปล่างสุด
            // --- ราคา และ ปุ่มบวก ---
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                Text(
                  '\$${product.price.toStringAsFixed(2)}',
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 18,
                  ),
                ),
                GestureDetector(
                  onTap: onAdd,
                  child: Container(
                    width: 45,
                    height: 45,
                    decoration: const BoxDecoration(
                      color: _primaryColor, // สีเขียว
                      shape: BoxShape.circle, // ทรงกลม
                    ),
                    child: const Icon(Icons.add, color: Colors.white, size: 28),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
