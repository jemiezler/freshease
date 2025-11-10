import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/theme/design_tokens.dart';

/// A neumorphic progress bar widget
class SoftProgressBar extends StatelessWidget {
  final double progress; // 0.0 to 1.0
  final double? height;
  final Color? backgroundColor;
  final Color? progressColor;
  final bool showLabel;

  const SoftProgressBar({
    super.key,
    required this.progress,
    this.height,
    this.backgroundColor,
    this.progressColor,
    this.showLabel = false,
  });

  @override
  Widget build(BuildContext context) {
    final effectiveHeight = height ?? 12.0;
    final effectiveBackgroundColor = backgroundColor ?? AppColors.surface;
    final effectiveProgressColor = progressColor ?? AppColors.primary;
    final clampedProgress = progress.clamp(0.0, 1.0);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          height: effectiveHeight,
          decoration: BoxDecoration(
            color: effectiveBackgroundColor,
            borderRadius: BorderRadius.circular(effectiveHeight / 2),
            boxShadow: DesignTokens.insetShadow,
          ),
          child: Stack(
            children: [
              // Progress fill with gradient
              FractionallySizedBox(
                widthFactor: clampedProgress,
                child: Container(
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      colors: [
                        effectiveProgressColor,
                        effectiveProgressColor.withValues(alpha: 0.8),
                      ],
                      begin: Alignment.centerLeft,
                      end: Alignment.centerRight,
                    ),
                    borderRadius: BorderRadius.circular(effectiveHeight / 2),
                  ),
                ),
              ),
            ],
          ),
        ),
        if (showLabel) ...[
          const SizedBox(height: 4),
          Text(
            '${(clampedProgress * 100).toInt()}%',
            style: const TextStyle(
              color: AppColors.textSecondary,
              fontSize: 12,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ],
    );
  }
}
