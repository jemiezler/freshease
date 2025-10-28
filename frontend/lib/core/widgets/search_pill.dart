import 'package:flutter/material.dart';

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

  /// show an “X” to clear the field when there’s text
  final bool showClear;

  /// Optional overrides for leading/trailing widgets
  final Widget? leading;
  final Widget? trailing;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final surface =
        backgroundColor ??
        // use withOpacity for wider SDK compatibility; swap to withValues if you prefer
        theme.colorScheme.surfaceContainerHighest.withOpacity(0.95);

    return Container(
      decoration: BoxDecoration(
        color: surface,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            blurRadius: 12,
            spreadRadius: -2,
            color: Colors.black.withOpacity(0.06),
            offset: const Offset(0, 6),
          ),
        ],
      ),
      padding: padding,
      child: Row(
        children: [
          leading ??
              const Padding(
                padding: EdgeInsets.only(right: 8),
                child: Icon(Icons.search, size: 22),
              ),

          /// TextField with live change + submit callbacks
          Expanded(
            child: TextField(
              controller: controller,
              readOnly: readOnly,
              onTap: onTap,
              decoration: InputDecoration(
                border: InputBorder.none,
                hintText: hintText,
                isDense: true,
              ),
              textInputAction: TextInputAction.search,
              onChanged: onChanged,
              onSubmitted: onSubmitted,
            ),
          ),

          // clear button (shows only when there’s text)
          if (showClear)
            ValueListenableBuilder<TextEditingValue>(
              valueListenable: controller,
              builder: (_, value, __) => AnimatedSwitcher(
                duration: const Duration(milliseconds: 150),
                child: value.text.isEmpty
                    ? const SizedBox.shrink()
                    : IconButton(
                        key: const ValueKey('clear'),
                        tooltip: 'Clear',
                        icon: const Icon(Icons.close_rounded),
                        onPressed: () {
                          controller.clear();
                          onChanged?.call('');
                        },
                      ),
              ),
            ),

          // trailing: filter or a custom widget
          trailing ??
              IconButton(
                tooltip: 'Filters',
                icon: const Icon(Icons.tune_rounded),
                onPressed: onFilterTap,
              ),
        ],
      ),
    );
  }
}
