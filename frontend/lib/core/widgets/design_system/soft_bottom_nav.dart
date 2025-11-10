import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/theme/design_tokens.dart';

/// A neumorphic bottom navigation bar widget
class SoftBottomNav extends StatelessWidget {
  final int currentIndex;
  final ValueChanged<int> onTap;
  final List<SoftNavItem> items;

  const SoftBottomNav({
    super.key,
    required this.currentIndex,
    required this.onTap,
    required this.items,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: DesignTokens.paddingMedium,
        vertical: DesignTokens.paddingSmall,
      ),
      decoration: BoxDecoration(
        color: AppColors.bg,
        boxShadow: [
          // Top shadow for separation
          BoxShadow(
            color: const Color(0xFFD6DFEA).withValues(alpha: 0.3),
            offset: const Offset(0, -4),
            blurRadius: 8,
            spreadRadius: 0,
          ),
        ],
      ),
      child: SafeArea(
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: List.generate(items.length, (index) {
            final item = items[index];
            final isSelected = index == currentIndex;

            return Expanded(
              child: GestureDetector(
                onTap: () => onTap(index),
                child: AnimatedContainer(
                  duration: DesignTokens.smallAnimation,
                  curve: DesignTokens.animationCurve,
                  padding: const EdgeInsets.symmetric(
                    vertical: DesignTokens.paddingSmall,
                  ),
                  margin: const EdgeInsets.symmetric(horizontal: 4),
                  decoration: BoxDecoration(
                    color: isSelected ? AppColors.surface : Colors.transparent,
                    borderRadius: BorderRadius.circular(
                      DesignTokens.radiusMedium,
                    ),
                    boxShadow: isSelected ? DesignTokens.raisedShadow : null,
                  ),
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Stack(
                        clipBehavior: Clip.none,
                        children: [
                          Icon(
                            isSelected ? item.selectedIcon : item.icon,
                            color: isSelected
                                ? AppColors.primary
                                : AppColors.textSecondary,
                            size: 24,
                          ),
                          if (item.badge != null && item.badge! > 0)
                            Positioned(
                              right: -6,
                              top: -6,
                              child: Container(
                                padding: const EdgeInsets.all(4),
                                decoration: BoxDecoration(
                                  color: AppColors.primary,
                                  shape: BoxShape.circle,
                                  boxShadow: DesignTokens.raisedShadow,
                                ),
                                constraints: const BoxConstraints(
                                  minWidth: 16,
                                  minHeight: 16,
                                ),
                                child: Text(
                                  '${item.badge! > 9 ? '9+' : item.badge}',
                                  style: const TextStyle(
                                    color: Colors.white,
                                    fontSize: 10,
                                    fontWeight: FontWeight.w600,
                                  ),
                                  textAlign: TextAlign.center,
                                ),
                              ),
                            ),
                        ],
                      ),
                      const SizedBox(height: 4),
                      Text(
                        item.label,
                        style: TextStyle(
                          color: isSelected
                              ? AppColors.primary
                              : AppColors.textSecondary,
                          fontSize: 12,
                          fontWeight: isSelected
                              ? FontWeight.w600
                              : FontWeight.w500,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            );
          }),
        ),
      ),
    );
  }
}

class SoftNavItem {
  final IconData icon;
  final IconData selectedIcon;
  final String label;
  final int? badge;

  const SoftNavItem({
    required this.icon,
    required this.selectedIcon,
    required this.label,
    this.badge,
  });
}
