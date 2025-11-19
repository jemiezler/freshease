import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/theme/design_tokens.dart';

/// A neumorphic icon button widget
class SoftIconButton extends StatefulWidget {
  final IconData icon;
  final VoidCallback? onPressed;
  final double? size;
  final Color? iconColor;
  final Color? backgroundColor;
  final String? tooltip;

  const SoftIconButton({
    super.key,
    required this.icon,
    this.onPressed,
    this.size,
    this.iconColor,
    this.backgroundColor,
    this.tooltip,
  });

  @override
  State<SoftIconButton> createState() => _SoftIconButtonState();
}

class _SoftIconButtonState extends State<SoftIconButton>
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
    final effectiveSize = widget.size ?? 48.0;
    final effectiveIconColor = widget.iconColor ?? AppColors.textPrimary;
    final effectiveBackgroundColor =
        widget.backgroundColor ?? AppColors.surface;

    Widget button = Material(
      shape: const CircleBorder(),
      child: Container(
        width: effectiveSize,
        height: effectiveSize,
        decoration: BoxDecoration(
          color: effectiveBackgroundColor,
          shape: BoxShape.circle,
        ),
        child: Center(
          child: Icon(
            widget.icon,
            color: effectiveIconColor,
            size: effectiveSize * 0.5,
          ),
        ),
      ),
    );

    if (widget.onPressed == null) {
      button = Opacity(opacity: 0.5, child: button);
    } else {
      button = GestureDetector(
        onTapDown: _handleTapDown,
        onTapUp: _handleTapUp,
        onTapCancel: _handleTapCancel,
        child: ScaleTransition(scale: _scaleAnimation, child: button),
      );
    }

    if (widget.tooltip != null) {
      return Tooltip(message: widget.tooltip!, child: button);
    }

    return button;
  }
}
