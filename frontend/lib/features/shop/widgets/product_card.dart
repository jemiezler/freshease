// ไฟล์: product_card.dart

import 'package:flutter/material.dart';
import 'package:frontend/features/shop/domain/product.dart';

// ‼️ ดึงสีหลักมาจาก shoppage
const Color _primaryColor = Color(0xFF90B56D);

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
    return InkWell(
      // ทำให้ทั้งการ์ดกดได้
      onTap: onTap,
      borderRadius: BorderRadius.circular(18),
      child: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(18),
          // ‼️ เพิ่มเส้นขอบสีเทาอ่อนตามภาพ
          border: Border.all(color: Colors.grey.shade200, width: 1.5),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // --- รูปสินค้า ---
            Expanded(
              flex: 3,
              child: Center(
                // ‼️ ตรวจสอบว่า product.image เป็น Asset Path ที่ถูกต้อง
                child: Image.asset(
                  product.image,
                  fit: BoxFit.contain,
                  errorBuilder: (ctx, err, stack) => const Icon(
                    Icons.image_not_supported_outlined,
                    color: Colors.grey,
                    size: 40,
                  ),
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

            // --- รายละเอียด (เช่น "7pcs, Priceg") ---
            Text(
              productDetail ?? '${product.category} item', // Fallback
              style: TextStyle(color: Colors.grey.shade600, fontSize: 14),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const Spacer(),

            // --- ราคา และ ปุ่มบวก ---
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                Text(
                  // ‼️ แสดงราคา (สมมติว่าเป็น $ ตามภาพ)
                  '\$${product.price.toStringAsFixed(2)}',
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 18,
                  ),
                ),

                // ‼️ นี่คือ "ปุ่ม" ที่แก้ไขแล้ว
                GestureDetector(
                  onTap: onAdd,
                  child: Container(
                    width: 45,
                    height: 45,
                    decoration: const BoxDecoration(
                      color: _primaryColor, // 👈 สีเขียว
                      shape: BoxShape.circle, // 👈 ทรงกลม
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
