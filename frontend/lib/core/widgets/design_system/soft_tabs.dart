import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/theme/design_tokens.dart';

/// A neumorphic tab bar widget with floating pills
class SoftTabs extends StatelessWidget {
  final List<String> tabs;
  final int selectedIndex;
  final ValueChanged<int> onTabSelected;
  final bool isScrollable;

  const SoftTabs({
    super.key,
    required this.tabs,
    required this.selectedIndex,
    required this.onTabSelected,
    this.isScrollable = false,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(4),
      decoration: BoxDecoration(
        color: AppColors.bg,
        borderRadius: BorderRadius.circular(DesignTokens.radiusLarge),
        boxShadow: DesignTokens.insetShadow,
      ),
      child: Row(
        mainAxisSize: isScrollable ? MainAxisSize.min : MainAxisSize.max,
        children: List.generate(tabs.length, (index) {
          final isSelected = index == selectedIndex;
          return Expanded(
            child: GestureDetector(
              onTap: () => onTabSelected(index),
              child: AnimatedContainer(
                duration: DesignTokens.smallAnimation,
                curve: DesignTokens.animationCurve,
                padding: const EdgeInsets.symmetric(
                  horizontal: DesignTokens.paddingMedium,
                  vertical: DesignTokens.paddingSmall,
                ),
                decoration: BoxDecoration(
                  color: isSelected ? AppColors.surface : Colors.transparent,
                  borderRadius: BorderRadius.circular(
                    DesignTokens.radiusMedium,
                  ),
                  boxShadow: isSelected ? DesignTokens.raisedShadow : null,
                ),
                child: Center(
                  child: Text(
                    tabs[index],
                    style: TextStyle(
                      color: isSelected
                          ? AppColors.textPrimary
                          : AppColors.textSecondary,
                      fontSize: 14,
                      fontWeight: isSelected
                          ? FontWeight.w600
                          : FontWeight.w500,
                    ),
                  ),
                ),
              ),
            ),
          );
        }),
      ),
    );
  }
}
