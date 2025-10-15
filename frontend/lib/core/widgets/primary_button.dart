import 'package:flutter/material.dart';

class PrimaryButton extends StatelessWidget {
  final VoidCallback? onPressed;
  final Widget child;
  final bool expanded;
  const PrimaryButton({
    super.key,
    required this.onPressed,
    required this.child,
    this.expanded = true,
  });

  @override
  Widget build(BuildContext context) {
    final btn = FilledButton(
      onPressed: onPressed,
      style: FilledButton.styleFrom(
        shape: const StadiumBorder(),
        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 14),
      ),
      child: child,
    );
    return expanded ? SizedBox(width: double.infinity, child: btn) : btn;
  }
}
