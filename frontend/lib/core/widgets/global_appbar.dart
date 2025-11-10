import 'package:flutter/material.dart';
import 'package:frontend/core/widgets/design_system/soft_icon_button.dart';

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
    return SoftIconButton(
      icon: icon,
      onPressed: onPressed,
      iconColor: iconColor,
      backgroundColor: backgroundColor,
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
      backgroundColor: Colors.transparent,
      centerTitle: true,
      elevation: 0,
      title: Text(
        title,
        style: theme.appBarTheme.titleTextStyle,
      ),
      leading: showBackButton && Navigator.canPop(context)
          ? Padding(
              padding: const EdgeInsets.only(left: 8),
              child: Center(
                child: SoftIconButton(
                  icon: Icons.arrow_back_ios_new_rounded,
                  onPressed: () => Navigator.of(context).maybePop(),
                  size: 40,
                  tooltip: 'Back',
                ),
              ),
            )
          : null,
      actions: actions?.map((action) {
        if (action is IconButton) {
          IconData? iconData;
          if (action.icon is Icon) {
            iconData = (action.icon as Icon).icon;
          }
          if (iconData != null) {
            return Padding(
              padding: const EdgeInsets.only(right: 8),
              child: Center(
                child: SoftIconButton(
                  icon: iconData,
                  onPressed: action.onPressed ?? () {},
                  size: 40,
                ),
              ),
            );
          }
        }
        return Padding(
          padding: const EdgeInsets.only(right: 8),
          child: Center(child: action),
        );
      }).toList(),
      bottom: bottom,
    );
  }

  @override
  Size get preferredSize {
    final bottomHeight = bottom?.preferredSize.height ?? 0;
    return Size.fromHeight(kToolbarHeight + bottomHeight);
  }
}
