import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';

class AppTheme {
  static ThemeData light() {
    return ThemeData(
      useMaterial3: true,
      colorSchemeSeed: AppColors.primary,
      scaffoldBackgroundColor: AppColors.background,
      appBarTheme: AppBarThemeData(
        backgroundColor: AppColors.background,
        surfaceTintColor: Colors.transparent,
        iconTheme: IconThemeData(color: Colors.black),
      ),
      visualDensity: VisualDensity.adaptivePlatformDensity,
    );
  }
}
