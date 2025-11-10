import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/theme/design_tokens.dart';

/// A neumorphic chip widget
class SoftChip extends StatelessWidget {
  final String label;
  final bool isSelected;
  final VoidCallback? onTap;
  final IconData? icon;
  final Color? backgroundColor;
  final Color? selectedColor;

  const SoftChip({
    super.key,
    required this.label,
    this.isSelected = false,
    this.onTap,
    this.icon,
    this.backgroundColor,
    this.selectedColor,
  });

  @override
  Widget build(BuildContext context) {
    final effectiveBackgroundColor = isSelected
        ? (selectedColor ?? AppColors.primary)
        : (backgroundColor ?? AppColors.surface);

    final textColor = isSelected ? Colors.white : AppColors.textPrimary;

    final shadows = isSelected
        ? DesignTokens.raisedShadow
        : DesignTokens.insetShadow;

    Widget chip = Container(
      padding: const EdgeInsets.symmetric(
        horizontal: DesignTokens.paddingMedium,
        vertical: DesignTokens.paddingSmall,
      ),
      decoration: BoxDecoration(
        color: effectiveBackgroundColor,
        borderRadius: BorderRadius.circular(DesignTokens.radiusLarge),
        boxShadow: shadows,
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (icon != null) ...[
            Icon(icon, color: textColor, size: 16),
            const SizedBox(width: 6),
          ],
          Text(
            label,
            style: TextStyle(
              color: textColor,
              fontSize: 14,
              fontWeight: isSelected ? FontWeight.w600 : FontWeight.w500,
            ),
          ),
        ],
      ),
    );

    if (onTap != null) {
      return GestureDetector(onTap: onTap, child: chip);
    }

    return chip;
  }
}
