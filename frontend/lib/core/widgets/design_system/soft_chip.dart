import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/theme/design_tokens.dart';

/// A soft neumorphic chip widget
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
    // Determine background color based on selection
    final effectiveBackgroundColor = isSelected
        ? (selectedColor ?? AppColors.primary)
        : (backgroundColor ?? AppColors.surface);

    // Determine text color
    final textColor = isSelected
        ? Colors.white
        : Theme.of(context).textTheme.bodyMedium?.color ??
              const Color(0xFF111827);

    // Determine shadows for neumorphic effect
    final shadows = isSelected
        ? [
            BoxShadow(
              color: Colors.black.withValues(alpha: 0.2),
              offset: const Offset(2, 2),
              blurRadius: 4,
            ),
            BoxShadow(
              color: Colors.white.withValues(alpha: 0.7),
              offset: const Offset(-2, -2),
              blurRadius: 4,
            ),
          ]
        : [
            BoxShadow(
              color: Colors.black.withValues(alpha: 0.15),
              offset: const Offset(2, 2),
              blurRadius: 6,
            ),
            BoxShadow(
              color: Colors.white.withOpacity(0.7),
              offset: const Offset(-2, -2),
              blurRadius: 6,
            ),
          ];

    Widget chipContent = Container(
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
              fontWeight: FontWeight.w600,
              letterSpacing: 0.2,
              height: 1, // Proper text height
            ),
          ),
        ],
      ),
    );

    // Wrap with InkWell for ripple effect if onTap is provided
    if (onTap != null) {
      return Material(
        color: Colors.transparent,
        child: InkWell(
          borderRadius: BorderRadius.circular(DesignTokens.radiusLarge),
          onTap: onTap,
          child: chipContent,
        ),
      );
    }

    return chipContent;
  }
}
