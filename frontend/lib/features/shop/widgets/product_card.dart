import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/theme/design_tokens.dart';
import 'package:frontend/core/widgets/design_system/soft_card.dart';
import 'package:frontend/core/widgets/design_system/soft_icon_button.dart';
import 'package:frontend/features/shop/domain/product.dart';

class ProductCard extends StatelessWidget {
  final Product product;
  final VoidCallback onAdd;
  final VoidCallback onTap;
  const ProductCard({
    super.key,
    required this.product,
    required this.onAdd,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return SoftCard(
      onTap: onTap,
      padding: EdgeInsets.zero,
      margin: const EdgeInsets.all(8),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: ClipRRect(
              borderRadius: BorderRadius.vertical(
                top: Radius.circular(DesignTokens.radiusMedium),
              ),
              child: Padding(
                padding: const EdgeInsets.all(6),
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(DesignTokens.radiusSmall),
                  child: Image.network(
                    product.image,
                    fit: BoxFit.cover,
                    width: double.infinity,
                    errorBuilder: (context, error, stackTrace) {
                      return Container(
                        color: AppColors.surface,
                        child: const Icon(Icons.image_not_supported),
                      );
                    },
                  ),
                ),
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(DesignTokens.paddingSmall),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  product.name,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                  style: const TextStyle(
                    fontWeight: FontWeight.w700,
                    fontSize: 14,
                    color: AppColors.textPrimary,
                  ),
                ),
                const SizedBox(height: 8),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    Text(
                      '฿${product.price.toStringAsFixed(2)}',
                      style: const TextStyle(
                        color: AppColors.primary,
                        fontWeight: FontWeight.w600,
                        fontSize: 16,
                      ),
                    ),
                    SoftIconButton(
                      icon: Icons.add_shopping_cart,
                      onPressed: onAdd,
                      size: 40,
                      tooltip: 'Add to cart',
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
