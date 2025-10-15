import 'package:flutter/material.dart';

class AppTheme {
  static ThemeData light() {
    return ThemeData(
      useMaterial3: true,
      colorSchemeSeed: Colors.teal,
      scaffoldBackgroundColor: const Color(0xFFF9FAFB),
      visualDensity: VisualDensity.adaptivePlatformDensity,
    );
  }
}
