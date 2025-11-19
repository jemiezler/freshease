import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/theme/design_tokens.dart';
import 'package:frontend/core/widgets/design_system/soft_icon_button.dart';

class SearchPill extends StatelessWidget {
  const SearchPill({
    super.key,
    required this.controller,
    this.onFilterTap,
    this.onChanged,
    this.onSubmitted,
    this.hintText = 'Search…',
    this.readOnly = false,
    this.onTap,
    this.padding = const EdgeInsets.symmetric(horizontal: 14),
    this.backgroundColor,
    this.showClear = true,
    this.leading,
    this.trailing,
  });

  final TextEditingController controller;
  final VoidCallback? onFilterTap;
  final ValueChanged<String>? onChanged;
  final ValueChanged<String>? onSubmitted;
  final String hintText;

  final bool readOnly;
  final VoidCallback? onTap;

  final EdgeInsetsGeometry padding;
  final Color? backgroundColor;

  /// show an "X" to clear the field when there's text
  final bool showClear;

  /// Optional overrides for leading/trailing widgets
  final Widget? leading;
  final Widget? trailing;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: backgroundColor ?? AppColors.surface,
        borderRadius: BorderRadius.circular(DesignTokens.radiusMedium),
      ),
      padding: padding,
      child: Row(
        children: [
          leading ??
              Padding(
                padding: const EdgeInsets.only(right: 8),
                child: Icon(
                  Icons.search,
                  size: 22,
                  color: AppColors.textSecondary,
                ),
              ),

          /// TextField with live change + submit callbacks
          Expanded(
            child: TextField(
              controller: controller,
              readOnly: readOnly,
              onTap: onTap,
              style: const TextStyle(
                color: AppColors.textPrimary,
                fontSize: 16,
              ),
              decoration: InputDecoration(
                border: InputBorder.none,
                hintText: hintText,
                hintStyle: TextStyle(
                  color: AppColors.textSecondary.withValues(alpha: 0.6),
                  fontSize: 16,
                ),
                isDense: true,
                contentPadding: EdgeInsets.zero,
              ),
              textInputAction: TextInputAction.search,
              onChanged: onChanged,
              onSubmitted: onSubmitted,
            ),
          ),

          // clear button (shows only when there's text)
          if (showClear)
            ValueListenableBuilder<TextEditingValue>(
              valueListenable: controller,
              builder: (_, value, _) => AnimatedSwitcher(
                duration: DesignTokens.microAnimation,
                child: value.text.isEmpty
                    ? const SizedBox.shrink()
                    : Padding(
                        key: const ValueKey('clear'),
                        padding: const EdgeInsets.only(left: 4),
                        child: SoftIconButton(
                          icon: Icons.close_rounded,
                          onPressed: () {
                            controller.clear();
                            onChanged?.call('');
                          },
                          size: 32,
                          tooltip: 'Clear',
                        ),
                      ),
              ),
            ),

          // trailing: filter or a custom widget
          Padding(
            padding: const EdgeInsets.only(left: 4),
            child:
                trailing ??
                SoftIconButton(
                  icon: Icons.tune_rounded,
                  onPressed: onFilterTap,
                  size: 32,
                  tooltip: 'Filters',
                ),
          ),
        ],
      ),
    );
  }
}
