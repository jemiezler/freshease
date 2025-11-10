import 'dart:async';
import 'package:flutter/material.dart';
import 'package:frontend/core/constants/app_colors.dart';
import 'package:frontend/core/theme/design_tokens.dart';
import 'package:go_router/go_router.dart';

class BannerItem {
  const BannerItem({
    required this.imageUrl,
    this.isAsset = false,
    this.route,
    this.onTap,
    this.semanticLabel,
  });

  final String imageUrl;
  final bool isAsset;
  final String? route;
  final VoidCallback? onTap;
  final String? semanticLabel;
}

class BannerCarousel extends StatefulWidget {
  const BannerCarousel({
    super.key,
    required this.items,
    this.height = 180,
    this.borderRadius = 16,
    this.autoPlay = true,
    this.autoPlayInterval = const Duration(seconds: 4),
    this.transitionDuration = const Duration(milliseconds: 420),
    this.curve = Curves.easeOut,
    this.showDots = true,
    this.gradientFade = true,
    this.onPageChanged,
  });

  final List<BannerItem> items;

  final double height;
  final double borderRadius;

  final bool autoPlay;
  final Duration autoPlayInterval;
  final Duration transitionDuration;
  final Curve curve;

  final bool showDots;
  final bool gradientFade;
  final ValueChanged<int>? onPageChanged;

  @override
  State<BannerCarousel> createState() => _BannerCarouselState();
}

class _BannerCarouselState extends State<BannerCarousel> {
  late final PageController _ctrl = PageController(viewportFraction: 1.0);
  Timer? _timer;
  int _index = 0;

  @override
  void initState() {
    super.initState();
    if (widget.autoPlay && widget.items.isNotEmpty) {
      _timer = Timer.periodic(widget.autoPlayInterval, (_) {
        if (!mounted || widget.items.isEmpty) return;
        final next = (_index + 1) % widget.items.length;
        _ctrl.animateToPage(
          next,
          duration: widget.transitionDuration,
          curve: widget.curve,
        );
      });
    }
  }

  @override
  void dispose() {
    _timer?.cancel();
    _ctrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (widget.items.isEmpty) return const SizedBox.shrink();

    return SizedBox(
      height: widget.height,
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(widget.borderRadius),
          boxShadow: DesignTokens.raisedShadow,
        ),
        child: ClipRRect(
          borderRadius: BorderRadius.circular(widget.borderRadius),
          child: Stack(
            fit: StackFit.expand,
            children: [
              PageView.builder(
                controller: _ctrl,
                clipBehavior: Clip.hardEdge, // ✅ prevents overlap/bleed
                itemCount: widget.items.length,
                onPageChanged: (i) {
                  setState(() => _index = i);
                  widget.onPageChanged?.call(i);
                },
                itemBuilder: (context, i) {
                  final item = widget.items[i];

                  void handleTap() {
                    if (item.onTap != null) {
                      item.onTap!();
                    } else if (item.route != null) {
                      context.go(item.route!);
                    }
                  }

                  // ✅ use cover to maintain ratio and fill frame
                  final image = item.isAsset
                      ? Image.asset(
                          item.imageUrl,
                          fit: BoxFit.cover,
                          width: double.infinity,
                          height: double.infinity,
                          semanticLabel: item.semanticLabel,
                        )
                      : Image.network(
                          item.imageUrl,
                          fit: BoxFit.cover,
                          width: double.infinity,
                          height: double.infinity,
                          semanticLabel: item.semanticLabel,
                        );

                  return GestureDetector(
                    onTap: (item.onTap != null || item.route != null)
                        ? handleTap
                        : null,
                    child: image,
                  );
                },
              ),

              // optional gradient fade
              if (widget.gradientFade)
                Positioned.fill(
                  child: IgnorePointer(
                    child: DecoratedBox(
                      decoration: BoxDecoration(
                        gradient: LinearGradient(
                          begin: Alignment.bottomCenter,
                          end: Alignment.center,
                          colors: [
                            Colors.black.withValues(alpha: 0.28),
                            Colors.transparent,
                          ],
                        ),
                      ),
                    ),
                  ),
                ),

              // dots indicator
              if (widget.showDots)
                Positioned(
                  bottom: 10,
                  left: 0,
                  right: 0,
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: List.generate(widget.items.length, (i) {
                      final active = i == _index;
                      return AnimatedContainer(
                        duration: DesignTokens.smallAnimation,
                        margin: const EdgeInsets.symmetric(horizontal: 3),
                        height: 8,
                        width: active ? 18 : 8,
                        decoration: BoxDecoration(
                          color: active
                              ? AppColors.primary
                              : Colors.white.withValues(alpha: 0.6),
                          borderRadius: BorderRadius.circular(8),
                          boxShadow: active ? DesignTokens.raisedShadow : null,
                        ),
                      );
                    }),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
