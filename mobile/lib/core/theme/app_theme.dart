import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

/// Mirrors `frontend/src/app/globals.css` so mobile and web read as the same
/// product. Tokens here are the OKLCH-source palette converted to sRGB:
/// crisp white surface in light, deep indigo-black in dark, single electric
/// indigo accent, generous radii, and Manrope across the type ramp.
class AppTheme {
  AppTheme._();

  // --- Light surface palette ---
  // oklch(1 0 0) → pure white
  static const _bgLight = Color(0xFFFFFFFF);
  // oklch(0.14 0.018 275) → near-black with a violet tint
  static const _fgLight = Color(0xFF1A1A2C);
  static const _cardLight = Color(0xFFFFFFFF);
  // oklch(0.91 0.008 275)
  static const _borderLight = Color(0xFFE4E3EB);
  // oklch(0.96 0.008 275)
  static const _inputLight = Color(0xFFEFEEF4);
  // oklch(0.97 0.006 275)
  static const _mutedLight = Color(0xFFF3F3F7);
  // oklch(0.46 0.015 275)
  static const _mutedFgLight = Color(0xFF6E6D80);

  // --- Dark surface palette ---
  // oklch(0.16 0.018 275)
  static const _bgDark = Color(0xFF1D1D2E);
  // oklch(0.20 0.020 275)
  static const _cardDark = Color(0xFF26263A);
  // oklch(0.22 0.020 275)
  static const _mutedDark = Color(0xFF2A2A40);
  // oklch(0.24 0.025 275)
  static const _inputDark = Color(0xFF2F2F46);
  // oklch(0.28 0.025 275)
  static const _borderDark = Color(0xFF3A3A55);
  // oklch(0.65 0.015 275)
  static const _mutedFgDark = Color(0xFF9B9AAA);
  // oklch(0.95 0.012 90)
  static const _fgDark = Color(0xFFF2F0E6);

  // --- Accent (single bold) ---
  // oklch(0.48 0.22 280) — electric indigo
  static const _primaryLight = Color(0xFF4F2FE3);
  // oklch(0.72 0.20 280) — softened indigo for dark surfaces
  static const _primaryDark = Color(0xFFA294FF);

  static const _destructive = Color(0xFFDC2626);

  static const _radiusInput = 12.0;
  static const _radiusCard = 14.0;

  static ThemeData light() => _build(
        brightness: Brightness.light,
        primary: _primaryLight,
        onPrimary: Colors.white,
        bg: _bgLight,
        card: _cardLight,
        border: _borderLight,
        input: _inputLight,
        muted: _mutedLight,
        mutedFg: _mutedFgLight,
        fg: _fgLight,
      );

  static ThemeData dark() => _build(
        brightness: Brightness.dark,
        primary: _primaryDark,
        onPrimary: _bgDark,
        bg: _bgDark,
        card: _cardDark,
        border: _borderDark,
        input: _inputDark,
        muted: _mutedDark,
        mutedFg: _mutedFgDark,
        fg: _fgDark,
      );

  static ThemeData _build({
    required Brightness brightness,
    required Color primary,
    required Color onPrimary,
    required Color bg,
    required Color card,
    required Color border,
    required Color input,
    required Color muted,
    required Color mutedFg,
    required Color fg,
  }) {
    final colorScheme = ColorScheme(
      brightness: brightness,
      primary: primary,
      onPrimary: onPrimary,
      secondary: primary,
      onSecondary: onPrimary,
      error: _destructive,
      onError: Colors.white,
      surface: bg,
      onSurface: fg,
      surfaceContainerLowest: bg,
      surfaceContainerLow: card,
      surfaceContainer: card,
      surfaceContainerHigh: muted,
      surfaceContainerHighest: muted,
      outline: border,
      outlineVariant: border,
      onSurfaceVariant: mutedFg,
    );

    final base = ThemeData(brightness: brightness);
    final textTheme = GoogleFonts.manropeTextTheme(base.textTheme).apply(
      bodyColor: fg,
      displayColor: fg,
    );

    return ThemeData(
      useMaterial3: true,
      brightness: brightness,
      colorScheme: colorScheme,
      scaffoldBackgroundColor: bg,
      canvasColor: bg,
      textTheme: textTheme,
      primaryTextTheme: textTheme,
      extensions: [
        AppPalette(
          bg: bg,
          card: card,
          border: border,
          input: input,
          muted: muted,
          mutedFg: mutedFg,
          fg: fg,
          primary: primary,
          onPrimary: onPrimary,
          destructive: _destructive,
        ),
      ],
      appBarTheme: AppBarTheme(
        backgroundColor: bg,
        foregroundColor: fg,
        elevation: 0,
        scrolledUnderElevation: 0,
        centerTitle: false,
        titleTextStyle: GoogleFonts.manrope(
          fontSize: 17,
          fontWeight: FontWeight.w700,
          color: fg,
          letterSpacing: -0.3,
        ),
        iconTheme: IconThemeData(color: fg),
        surfaceTintColor: Colors.transparent,
        shape: Border(bottom: BorderSide(color: border, width: 0.5)),
      ),
      cardTheme: CardThemeData(
        color: card,
        elevation: 0,
        margin: EdgeInsets.zero,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(_radiusCard),
          side: BorderSide(color: border),
        ),
        clipBehavior: Clip.antiAlias,
      ),
      dividerTheme: DividerThemeData(color: border, space: 1, thickness: 0.5),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: input,
        contentPadding: const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(_radiusInput),
          borderSide: BorderSide(color: border),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(_radiusInput),
          borderSide: BorderSide(color: border),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(_radiusInput),
          borderSide: BorderSide(color: primary, width: 1.5),
        ),
        errorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(_radiusInput),
          borderSide: const BorderSide(color: _destructive),
        ),
        hintStyle: TextStyle(color: mutedFg),
        labelStyle: TextStyle(color: mutedFg),
      ),
      filledButtonTheme: FilledButtonThemeData(
        style: FilledButton.styleFrom(
          backgroundColor: primary,
          foregroundColor: onPrimary,
          padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 14),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(_radiusInput)),
          textStyle: GoogleFonts.manrope(fontWeight: FontWeight.w600, fontSize: 14),
          minimumSize: const Size(0, 44),
        ),
      ),
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          foregroundColor: fg,
          side: BorderSide(color: border),
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(_radiusInput)),
          textStyle: GoogleFonts.manrope(fontWeight: FontWeight.w600, fontSize: 14),
          minimumSize: const Size(0, 44),
        ),
      ),
      textButtonTheme: TextButtonThemeData(
        style: TextButton.styleFrom(
          foregroundColor: primary,
          textStyle: GoogleFonts.manrope(fontWeight: FontWeight.w600),
        ),
      ),
      iconButtonTheme: IconButtonThemeData(
        style: IconButton.styleFrom(foregroundColor: fg),
      ),
      navigationBarTheme: NavigationBarThemeData(
        backgroundColor: bg,
        indicatorColor: primary.withValues(alpha: 0.12),
        labelTextStyle: WidgetStatePropertyAll(
          GoogleFonts.manrope(fontSize: 11, fontWeight: FontWeight.w600, letterSpacing: 0.1),
        ),
        iconTheme: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) return IconThemeData(color: primary, size: 22);
          return IconThemeData(color: mutedFg, size: 22);
        }),
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        height: 68,
      ),
      bottomSheetTheme: BottomSheetThemeData(
        backgroundColor: card,
        surfaceTintColor: Colors.transparent,
        shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
        ),
      ),
      chipTheme: ChipThemeData(
        backgroundColor: muted,
        side: BorderSide(color: border),
        labelStyle: GoogleFonts.manrope(fontSize: 12, fontWeight: FontWeight.w500, color: fg),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(999)),
      ),
      listTileTheme: ListTileThemeData(
        iconColor: mutedFg,
        textColor: fg,
        tileColor: card,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(_radiusCard)),
      ),
      dialogTheme: DialogThemeData(
        backgroundColor: card,
        surfaceTintColor: Colors.transparent,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      ),
      snackBarTheme: SnackBarThemeData(
        backgroundColor: fg,
        contentTextStyle: GoogleFonts.manrope(color: bg, fontWeight: FontWeight.w500),
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(_radiusInput)),
      ),
    );
  }
}

/// Theme extension exposing the raw token names from globals.css so widgets can
/// pull semantic colors (border, mutedFg) without reaching into the M3 scheme
/// and praying onSurfaceVariant happens to mean the same thing.
@immutable
class AppPalette extends ThemeExtension<AppPalette> {
  const AppPalette({
    required this.bg,
    required this.card,
    required this.border,
    required this.input,
    required this.muted,
    required this.mutedFg,
    required this.fg,
    required this.primary,
    required this.onPrimary,
    required this.destructive,
  });

  final Color bg;
  final Color card;
  final Color border;
  final Color input;
  final Color muted;
  final Color mutedFg;
  final Color fg;
  final Color primary;
  final Color onPrimary;
  final Color destructive;

  static AppPalette of(BuildContext context) {
    final ext = Theme.of(context).extension<AppPalette>();
    if (ext == null) {
      throw FlutterError('AppPalette missing from ThemeData.extensions');
    }
    return ext;
  }

  @override
  AppPalette copyWith({
    Color? bg,
    Color? card,
    Color? border,
    Color? input,
    Color? muted,
    Color? mutedFg,
    Color? fg,
    Color? primary,
    Color? onPrimary,
    Color? destructive,
  }) {
    return AppPalette(
      bg: bg ?? this.bg,
      card: card ?? this.card,
      border: border ?? this.border,
      input: input ?? this.input,
      muted: muted ?? this.muted,
      mutedFg: mutedFg ?? this.mutedFg,
      fg: fg ?? this.fg,
      primary: primary ?? this.primary,
      onPrimary: onPrimary ?? this.onPrimary,
      destructive: destructive ?? this.destructive,
    );
  }

  @override
  AppPalette lerp(ThemeExtension<AppPalette>? other, double t) {
    if (other is! AppPalette) return this;
    return AppPalette(
      bg: Color.lerp(bg, other.bg, t)!,
      card: Color.lerp(card, other.card, t)!,
      border: Color.lerp(border, other.border, t)!,
      input: Color.lerp(input, other.input, t)!,
      muted: Color.lerp(muted, other.muted, t)!,
      mutedFg: Color.lerp(mutedFg, other.mutedFg, t)!,
      fg: Color.lerp(fg, other.fg, t)!,
      primary: Color.lerp(primary, other.primary, t)!,
      onPrimary: Color.lerp(onPrimary, other.onPrimary, t)!,
      destructive: Color.lerp(destructive, other.destructive, t)!,
    );
  }
}
