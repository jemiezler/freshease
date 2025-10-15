import 'package:flutter/material.dart';

class RoundedIconButton extends StatelessWidget {
  final IconData icon;
  final VoidCallback onPressed;
  final Color? backgroundColor;
  final Color? iconColor;

  const RoundedIconButton({
    super.key,
    required this.icon,
    required this.onPressed,
    this.backgroundColor,
    this.iconColor,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return IconButton(
      onPressed: onPressed,
      icon: Container(
        padding: const EdgeInsets.all(6),
        decoration: BoxDecoration(
          color:
              backgroundColor ??
              (theme.brightness == Brightness.dark
                  ? Colors.white10
                  : Colors.black12),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Icon(
          icon,
          color: iconColor ?? theme.iconTheme.color,
          size: theme.iconTheme.size,
        ),
      ),
    );
  }
}

class GlobalAppBar extends StatelessWidget implements PreferredSizeWidget {
  final String title;
  final List<Widget>? actions;
  final PreferredSizeWidget? bottom;
  final bool showBackButton;

  const GlobalAppBar({
    super.key,
    required this.title,
    this.actions,
    this.bottom,
    this.showBackButton = true,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppBar(
      backgroundColor: theme.appBarTheme.backgroundColor,
      centerTitle: true,
      elevation: 0,
      title: Text(title, style: theme.appBarTheme.titleTextStyle),
      leading: showBackButton && Navigator.canPop(context)
          ? Padding(
              padding: const EdgeInsets.only(
                left: 6,
              ), // optional spacing from screen edge
              child: Center(
                child: GestureDetector(
                  onTap: () => Navigator.of(context).maybePop(),
                  child: Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: Colors.grey.shade300,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(
                      Icons.arrow_back_ios_new_rounded,
                      color: theme.iconTheme.color,
                      size: theme.iconTheme.size ?? 20,
                    ),
                  ),
                ),
              ),
            )
          : null,
      actions: actions,
      bottom: bottom,
    );
  }

  @override
  Size get preferredSize {
    final bottomHeight = bottom?.preferredSize.height ?? 0;
    return Size.fromHeight(kToolbarHeight + bottomHeight);
  }
}
