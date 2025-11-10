import 'package:flutter/material.dart';

/// Soft UI (Neumorphism) Design System Colors
class AppColors {
  // Primary brand color
  static const primary = Color(0xFF53B175);

  // Background colors
  static const bg = Color(0xFFF2F5F9); // soft off-white base
  static const surface = Color(0xFFF6F9FC); // slightly lighter than bg

  // Text colors
  static const textPrimary = Color(0xFF1F2937); // gray-900
  static const textSecondary = Color(0xFF6B7280); // gray-500

  // UI element colors
  static const divider = Color(0xFFE5E7EB);
  static const success = Color(0xFF34D399);
  static const warning = Color(0xFFF59E0B);

  // Legacy support (deprecated, use above)
  @Deprecated('Use bg instead')
  static const background = bg;

  @Deprecated('Use surface instead')
  static const foreground = surface;

  @Deprecated('Use primary instead')
  static const secondary = primary;

  @Deprecated('Use primary instead')
  static const accent = primary;
}
