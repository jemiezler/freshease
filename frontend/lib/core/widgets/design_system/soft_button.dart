import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/theme/design_tokens.dart';

/// A neumorphic button widget with raised/inset states
class SoftButton extends StatefulWidget {
  final String label;
  final VoidCallback? onPressed;
  final bool isPrimary;
  final bool isFullWidth;
  final IconData? icon;
  final double? height;
  final EdgeInsetsGeometry? padding;

  const SoftButton({
    super.key,
    required this.label,
    this.onPressed,
    this.isPrimary = false,
    this.isFullWidth = true,
    this.icon,
    this.height,
    this.padding,
  });

  @override
  State<SoftButton> createState() => _SoftButtonState();
}

class _SoftButtonState extends State<SoftButton>
    with SingleTickerProviderStateMixin {
  bool _isPressed = false;
  late AnimationController _controller;
  late Animation<double> _scaleAnimation;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      duration: DesignTokens.microAnimation,
      vsync: this,
    );
    _scaleAnimation = Tween<double>(begin: 1.0, end: DesignTokens.pressScale)
        .animate(
          CurvedAnimation(
            parent: _controller,
            curve: DesignTokens.animationCurve,
          ),
        );
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _handleTapDown(TapDownDetails details) {
    setState(() => _isPressed = true);
    _controller.forward();
  }

  void _handleTapUp(TapUpDetails details) {
    setState(() => _isPressed = false);
    _controller.reverse();
    widget.onPressed?.call();
  }

  void _handleTapCancel() {
    setState(() => _isPressed = false);
    _controller.reverse();
  }

  @override
  Widget build(BuildContext context) {
    final effectiveHeight = widget.height ?? 56.0;
    final effectivePadding =
        widget.padding ??
        const EdgeInsets.symmetric(
          horizontal: DesignTokens.paddingLarge,
          vertical: DesignTokens.paddingMedium,
        );

    final backgroundColor = widget.isPrimary
        ? AppColors.primary
        : AppColors.surface;

    final textColor = widget.isPrimary ? Colors.white : AppColors.textPrimary;

    final shadows = _isPressed
        ? DesignTokens.pressedShadow
        : DesignTokens.raisedShadow;

    Widget buttonContent = Container(
      height: effectiveHeight,
      padding: effectivePadding,
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(effectiveHeight / 2), // Pill shape
        boxShadow: shadows,
      ),
      child: Row(
        mainAxisSize: widget.isFullWidth ? MainAxisSize.max : MainAxisSize.min,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          if (widget.icon != null) ...[
            Icon(widget.icon, color: textColor, size: 20),
            const SizedBox(width: 8),
          ],
          Text(
            widget.label,
            style: TextStyle(
              color: textColor,
              fontSize: 16,
              fontWeight: FontWeight.w600,
            ),
          ),
        ],
      ),
    );

    if (widget.onPressed == null) {
      return Opacity(opacity: 0.5, child: buttonContent);
    }

    return GestureDetector(
      onTapDown: _handleTapDown,
      onTapUp: _handleTapUp,
      onTapCancel: _handleTapCancel,
      child: ScaleTransition(scale: _scaleAnimation, child: buttonContent),
    );
  }
}
