import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/theme/design_tokens.dart';

/// A neumorphic avatar widget with outer glow
class SoftAvatar extends StatelessWidget {
  final String? imageUrl;
  final String? initials;
  final double? radius;
  final Color? backgroundColor;
  final IconData? icon;

  const SoftAvatar({
    super.key,
    this.imageUrl,
    this.initials,
    this.radius,
    this.backgroundColor,
    this.icon,
  });

  @override
  Widget build(BuildContext context) {
    final effectiveRadius = radius ?? 24.0;
    final effectiveBackgroundColor = backgroundColor ?? AppColors.surface;

    Widget avatarContent;

    if (imageUrl != null) {
      avatarContent = CircleAvatar(
        radius: effectiveRadius,
        backgroundImage: NetworkImage(imageUrl!),
        backgroundColor: effectiveBackgroundColor,
      );
    } else if (initials != null) {
      avatarContent = CircleAvatar(
        radius: effectiveRadius,
        backgroundColor: effectiveBackgroundColor,
        child: Text(
          initials!,
          style: const TextStyle(
            color: AppColors.textPrimary,
            fontWeight: FontWeight.w600,
          ),
        ),
      );
    } else if (icon != null) {
      avatarContent = CircleAvatar(
        radius: effectiveRadius,
        backgroundColor: effectiveBackgroundColor,
        child: Icon(icon, color: AppColors.textPrimary, size: effectiveRadius),
      );
    } else {
      avatarContent = CircleAvatar(
        radius: effectiveRadius,
        backgroundColor: effectiveBackgroundColor,
        child: const Icon(Icons.person, color: AppColors.textPrimary),
      );
    }

    return Container(
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        boxShadow: DesignTokens.avatarGlow,
      ),
      child: avatarContent,
    );
  }
}
