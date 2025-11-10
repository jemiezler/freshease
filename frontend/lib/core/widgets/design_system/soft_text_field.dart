import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/theme/design_tokens.dart';

/// A neumorphic text field widget with inset effect
class SoftTextField extends StatelessWidget {
  final TextEditingController? controller;
  final String? hintText;
  final String? labelText;
  final IconData? prefixIcon;
  final IconData? suffixIcon;
  final VoidCallback? onSuffixIconTap;
  final bool obscureText;
  final TextInputType? keyboardType;
  final ValueChanged<String>? onChanged;
  final FormFieldValidator<String>? validator;
  final int? maxLines;
  final bool enabled;
  final FocusNode? focusNode;

  const SoftTextField({
    super.key,
    this.controller,
    this.hintText,
    this.labelText,
    this.prefixIcon,
    this.suffixIcon,
    this.onSuffixIconTap,
    this.obscureText = false,
    this.keyboardType,
    this.onChanged,
    this.validator,
    this.maxLines = 1,
    this.enabled = true,
    this.focusNode,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (labelText != null) ...[
          Padding(
            padding: const EdgeInsets.only(bottom: 8.0),
            child: Text(
              labelText!,
              style: const TextStyle(
                color: AppColors.textPrimary,
                fontSize: 14,
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
        ],
        Container(
          decoration: BoxDecoration(
            color: AppColors.surface,
            borderRadius: BorderRadius.circular(DesignTokens.radiusMedium),
            boxShadow: DesignTokens.insetShadow,
          ),
          child: TextFormField(
            controller: controller,
            focusNode: focusNode,
            obscureText: obscureText,
            keyboardType: keyboardType,
            onChanged: onChanged,
            validator: validator,
            maxLines: maxLines,
            enabled: enabled,
            style: const TextStyle(color: AppColors.textPrimary, fontSize: 16),
            decoration: InputDecoration(
              hintText: hintText,
              hintStyle: TextStyle(
                color: AppColors.textSecondary.withValues(alpha: 0.6),
                fontSize: 16,
              ),
              prefixIcon: prefixIcon != null
                  ? Icon(prefixIcon, color: AppColors.textSecondary, size: 20)
                  : null,
              suffixIcon: suffixIcon != null
                  ? GestureDetector(
                      onTap: onSuffixIconTap,
                      child: Icon(
                        suffixIcon,
                        color: AppColors.textSecondary,
                        size: 20,
                      ),
                    )
                  : null,
              border: InputBorder.none,
              contentPadding: const EdgeInsets.symmetric(
                horizontal: DesignTokens.paddingMedium,
                vertical: DesignTokens.paddingMedium,
              ),
              enabledBorder: InputBorder.none,
              focusedBorder: InputBorder.none,
              disabledBorder: InputBorder.none,
              errorBorder: InputBorder.none,
              focusedErrorBorder: InputBorder.none,
            ),
          ),
        ),
      ],
    );
  }
}
