import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

import '../theme/app_theme.dart';

/// Mirrors the `/ section` mono caption used across the web layout.
class AppSectionHeader extends StatelessWidget {
  const AppSectionHeader(this.label, {super.key});
  final String label;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Text(
      '/ ${label.toUpperCase()}',
      style: GoogleFonts.jetBrainsMono(
        fontSize: 10,
        fontWeight: FontWeight.w500,
        letterSpacing: 1.4,
        color: p.mutedFg,
      ),
    );
  }
}

/// Big editorial heading + subdued subtitle, matching `font-display
/// tracking-[-0.045em]` on the web hero blocks.
class AppPageHeader extends StatelessWidget {
  const AppPageHeader({
    super.key,
    required this.title,
    this.subtitle,
  });

  final String title;
  final String? subtitle;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          title,
          style: GoogleFonts.manrope(
            fontSize: 32,
            fontWeight: FontWeight.w700,
            color: p.fg,
            letterSpacing: -1.3,
            height: 1.05,
          ),
        ),
        if (subtitle != null) ...[
          const SizedBox(height: 6),
          Text(
            subtitle!,
            style: GoogleFonts.manrope(
              fontSize: 13,
              color: p.mutedFg,
            ),
          ),
        ],
      ],
    );
  }
}
