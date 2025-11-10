import 'package:flutter/material.dart';

/// Design tokens for Soft UI (Neumorphism) design system
class DesignTokens {
  // Border radius
  static const double radiusSmall = 20.0;
  static const double radiusMedium = 24.0;
  static const double radiusLarge = 28.0;

  // Padding
  static const double paddingSmall = 12.0;
  static const double paddingMedium = 16.0;
  static const double paddingLarge = 20.0;

  // Shadows for raised effect (neumorphic)
  static List<BoxShadow> get raisedShadow => [
    // Light shadow (top-left)
    BoxShadow(
      color: Colors.white.withValues(alpha: 0.9),
      offset: const Offset(-8, -8),
      blurRadius: 16,
      spreadRadius: 0,
    ),
    // Dark shadow (bottom-right)
    BoxShadow(
      color: const Color(0xFFD6DFEA).withValues(alpha: 0.9),
      offset: const Offset(8, 8),
      blurRadius: 16,
      spreadRadius: 0,
    ),
  ];

  // Shadows for inset effect
  static List<BoxShadow> get insetShadow => [
    // Inner highlight (top-left)
    BoxShadow(
      color: Colors.white.withValues(alpha: 0.7),
      offset: const Offset(4, 4),
      blurRadius: 8,
      spreadRadius: -2,
    ),
    // Inner shade (bottom-right)
    BoxShadow(
      color: const Color(0xFFD6DFEA).withValues(alpha: 0.7),
      offset: const Offset(-4, -4),
      blurRadius: 8,
      spreadRadius: -2,
    ),
  ];

  // Shadows for pressed/active state (deeper inset)
  static List<BoxShadow> get pressedShadow => [
    BoxShadow(
      color: Colors.white.withValues(alpha: 0.5),
      offset: const Offset(2, 2),
      blurRadius: 4,
      spreadRadius: -1,
    ),
    BoxShadow(
      color: const Color(0xFFD6DFEA).withValues(alpha: 0.8),
      offset: const Offset(-2, -2),
      blurRadius: 4,
      spreadRadius: -1,
    ),
  ];

  // Subtle shadow for avatars
  static List<BoxShadow> get avatarGlow => [
    BoxShadow(
      color: Colors.white.withValues(alpha: 0.6),
      offset: const Offset(-2, -2),
      blurRadius: 4,
    ),
    BoxShadow(
      color: const Color(0xFFD6DFEA).withValues(alpha: 0.6),
      offset: const Offset(2, 2),
      blurRadius: 4,
    ),
  ];

  // Animation durations
  static const Duration microAnimation = Duration(milliseconds: 100);
  static const Duration smallAnimation = Duration(milliseconds: 200);

  // Animation curves
  static const Curve animationCurve = Curves.easeOut;

  // Scale for hover/press effects
  static const double hoverScale = 1.02;
  static const double pressScale = 0.98;
}
