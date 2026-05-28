import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

import '../theme/app_theme.dart';

enum AppBadgeTone { neutral, success, warning, danger, info }

/// Pill matching the web `StatusBadge` look — soft tint + thin colored border.
class AppBadge extends StatelessWidget {
  const AppBadge({
    super.key,
    required this.label,
    this.tone = AppBadgeTone.neutral,
    this.icon,
  });

  final String label;
  final AppBadgeTone tone;
  final IconData? icon;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    final (Color border, Color bg, Color fg) = switch (tone) {
      AppBadgeTone.neutral => (p.border, p.muted, p.mutedFg),
      AppBadgeTone.success => (
          const Color(0x4D34C759), // 30% green
          const Color(0x1A34C759),
          const Color(0xFF1B7A3A),
        ),
      AppBadgeTone.warning => (
          const Color(0x4DF59E0B),
          const Color(0x1AF59E0B),
          const Color(0xFFB45309),
        ),
      AppBadgeTone.danger => (
          const Color(0x4DDC2626),
          const Color(0x1ADC2626),
          const Color(0xFFB91C1C),
        ),
      AppBadgeTone.info => (
          p.primary.withValues(alpha: 0.3),
          p.primary.withValues(alpha: 0.1),
          p.primary,
        ),
    };

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
      decoration: BoxDecoration(
        color: bg,
        borderRadius: BorderRadius.circular(999),
        border: Border.all(color: border),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (icon != null) ...[
            Icon(icon, size: 12, color: fg),
            const SizedBox(width: 5),
          ],
          Text(
            label,
            style: GoogleFonts.manrope(
              fontSize: 11,
              fontWeight: FontWeight.w600,
              color: fg,
            ),
          ),
        ],
      ),
    );
  }
}
