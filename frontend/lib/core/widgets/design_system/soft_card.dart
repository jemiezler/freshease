import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/theme/design_tokens.dart';

/// A neumorphic card widget with raised shadow effect
class SoftCard extends StatelessWidget {
  final Widget child;
  final EdgeInsetsGeometry? padding;
  final EdgeInsetsGeometry? margin;
  final double? borderRadius;
  final Color? backgroundColor;
  final VoidCallback? onTap;
  final bool isPressed;

  const SoftCard({
    super.key,
    required this.child,
    this.padding,
    this.margin,
    this.borderRadius,
    this.backgroundColor,
    this.onTap,
    this.isPressed = false,
  });

  @override
  Widget build(BuildContext context) {
    final effectivePadding =
        padding ?? const EdgeInsets.all(DesignTokens.paddingMedium);
    final effectiveBorderRadius = borderRadius ?? DesignTokens.radiusMedium;
    final effectiveBackgroundColor = backgroundColor ?? AppColors.surface;
    final effectiveShadows = isPressed
        ? DesignTokens.pressedShadow
        : DesignTokens.raisedShadow;

    Widget content = Container(
      padding: effectivePadding,
      decoration: BoxDecoration(
        color: effectiveBackgroundColor,
        borderRadius: BorderRadius.circular(effectiveBorderRadius),
        boxShadow: effectiveShadows,
      ),
      child: child,
    );

    if (onTap != null) {
      return GestureDetector(onTap: onTap, child: content);
    }

    return content;
  }
}
