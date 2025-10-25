import 'package:flutter/material.dart';
import 'package:frontend/features/shop/domain/product.dart';

class ProductCard extends StatelessWidget {
  final Product product;
  final VoidCallback onAdd;
  final VoidCallback onTap;
  final String? productDetail; // รายละเอียดสินค้าเสริม

  const ProductCard({
    super.key,
    required this.product,
    required this.onAdd,
    required this.onTap,
    this.productDetail,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    // ⚠️ ใช้ theme.colorScheme.primary สำหรับสีเขียวอ่อนของปุ่ม
    final primaryColor = theme.colorScheme.primary;

    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: Colors.grey.shade200),
        ),
        // ⚠️ ลบ padding ด้านล่างออก และใช้ Column.mainAxisSize.min
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          // ⚠️ ใช้ min เพื่อให้ Column มีความสูงตามเนื้อหาจริงเท่านั้น (สำคัญสำหรับ ListView แนวนอน)
          mainAxisSize: MainAxisSize.min,
          children: [
            // --- ส่วนรูปภาพ (120px) ---
            Container(
              height: 120,
              width: double.infinity,
              child: ClipRRect(
                borderRadius: const BorderRadius.vertical(
                  top: Radius.circular(16),
                ),
                child: Image.network(
                  product.image,
                  fit: BoxFit.cover,
                  width: double.infinity,
                  errorBuilder: (context, error, stackTrace) {
                    return Container(
                      color: Colors.grey.shade300,
                      child: const Center(child: Icon(Icons.broken_image)),
                    );
                  },
                ),
              ),
            ),

            // --- ส่วนรายละเอียดสินค้า ---
            Padding(
              padding: const EdgeInsets.symmetric(
                horizontal: 8.0,
                vertical: 4.0,
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    // ชื่อสินค้า
                    product.name,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: theme.textTheme.titleSmall?.copyWith(
                      fontWeight: FontWeight.w700,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    // รายละเอียดสินค้า
                    productDetail ?? product.category,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: Colors.grey.shade600,
                    ),
                  ),
                ],
              ),
            ),

            // ⚠️ แทนที่ Spacer ด้วย SizedBox เพื่อควบคุมระยะห่าง (UX/UI spacing)
            const SizedBox(height: 16),

            // --- ส่วนราคาและปุ่มเพิ่ม (ตามรูปภาพ) ---
            Padding(
              // ⚠️ ปรับ padding ด้านล่างให้เหมาะสมกับดีไซน์
              padding: const EdgeInsets.fromLTRB(8.0, 0, 8.0, 12.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  Text(
                    // ราคา
                    '\$${product.price.toStringAsFixed(2)}',
                    style: theme.textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.w700,
                      color: Colors.black,
                    ),
                  ),

                  // ปุ่มเพิ่มลงตะกร้า (ปุ่มเขียวอ่อน)
                  SizedBox(
                    width: 36,
                    height: 36,
                    child: FloatingActionButton(
                      heroTag: null,
                      onPressed: onAdd,
                      mini: true,
                      elevation: 0,
                      // ⚠️ ใช้สีเขียวอ่อนที่ได้จาก Theme
                      backgroundColor: primaryColor,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: const Icon(
                        Icons.add,
                        color: Colors.white,
                        size: 20,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
