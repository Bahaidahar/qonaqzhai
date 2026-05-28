import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../core/i18n/i18n.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/theme/theme_controller.dart';
import '../../../core/ui/ui.dart';
import '../../auth/presentation/viewmodels/auth_viewmodel.dart';

class SettingsScreen extends ConsumerWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final locale = ref.watch(localeProvider);
    final theme = ref.watch(themeControllerProvider);
    final auth = ref.watch(authViewModelProvider);
    final user = auth.user;
    final p = AppPalette.of(context);

    return Scaffold(
      appBar: AppBar(title: Text(tr(ref, 'settings_title'))),
      body: ListView(
        padding: const EdgeInsets.fromLTRB(20, 16, 20, 32),
        children: [
          const AppPageHeader(title: 'Settings'),
          const SizedBox(height: 24),
          if (user != null) ...[
            AppCard(
              padding: const EdgeInsets.all(18),
              child: Row(
                children: [
                  Container(
                    width: 44,
                    height: 44,
                    decoration: BoxDecoration(
                      color: p.primary.withValues(alpha: 0.1),
                      shape: BoxShape.circle,
                      border: Border.all(color: p.primary.withValues(alpha: 0.25)),
                    ),
                    alignment: Alignment.center,
                    child: Text(
                      _initial(user.name, user.email),
                      style: GoogleFonts.manrope(
                          fontWeight: FontWeight.w700, color: p.primary),
                    ),
                  ),
                  const SizedBox(width: 14),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          user.name.isNotEmpty ? user.name : user.email,
                          style: GoogleFonts.manrope(
                            fontWeight: FontWeight.w700,
                            fontSize: 15,
                            color: p.fg,
                          ),
                        ),
                        const SizedBox(height: 2),
                        Text(
                          user.email,
                          style: GoogleFonts.manrope(
                              fontSize: 12, color: p.mutedFg),
                        ),
                      ],
                    ),
                  ),
                  AppBadge(label: user.role.name, tone: AppBadgeTone.info),
                ],
              ),
            ),
            const SizedBox(height: 24),
          ],

          // Account
          const AppSectionHeader('Account'),
          const SizedBox(height: 10),
          AppCard(
            padding: EdgeInsets.zero,
            child: Column(
              children: [
                _Row(label: tr(ref, 'settings_name'), value: user?.name ?? '—'),
                Divider(height: 1, color: p.border),
                _Row(label: tr(ref, 'settings_email'), value: user?.email ?? '—'),
                Divider(height: 1, color: p.border),
                const _Row(label: 'Plan', value: 'Free · MVP'),
              ],
            ),
          ),

          const SizedBox(height: 24),
          // Preferences — language
          const AppSectionHeader('Language'),
          const SizedBox(height: 10),
          AppCard(
            padding: EdgeInsets.zero,
            child: Column(
              children: [
                for (var i = 0; i < AppLocale.values.length; i++) ...[
                  _LocaleRow(
                    locale: AppLocale.values[i],
                    selected: AppLocale.values[i] == locale,
                    onTap: () =>
                        ref.read(localeProvider.notifier).set(AppLocale.values[i]),
                  ),
                  if (i < AppLocale.values.length - 1)
                    Divider(height: 1, color: p.border),
                ],
              ],
            ),
          ),

          const SizedBox(height: 24),
          // Appearance
          const AppSectionHeader('Appearance'),
          const SizedBox(height: 10),
          Row(
            children: [
              Expanded(
                child: _ThemeTile(
                  icon: CupertinoIcons.sun_max,
                  label: 'Light',
                  selected: theme == AppThemePref.light,
                  onTap: () => ref
                      .read(themeControllerProvider.notifier)
                      .set(AppThemePref.light),
                ),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: _ThemeTile(
                  icon: CupertinoIcons.moon,
                  label: 'Dark',
                  selected: theme == AppThemePref.dark,
                  onTap: () => ref
                      .read(themeControllerProvider.notifier)
                      .set(AppThemePref.dark),
                ),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: _ThemeTile(
                  icon: CupertinoIcons.device_phone_portrait,
                  label: 'System',
                  selected: theme == AppThemePref.system,
                  onTap: () => ref
                      .read(themeControllerProvider.notifier)
                      .set(AppThemePref.system),
                ),
              ),
            ],
          ),

          const SizedBox(height: 24),
          // Payment shortcuts (cards screen, accessible only to customers in
          // the web app — vendors don't pay through the platform).
          if (user?.role.name == 'customer') ...[
            const AppSectionHeader('Payments'),
            const SizedBox(height: 10),
            AppCard(
              padding: EdgeInsets.zero,
              child: InkWell(
                onTap: () => context.push('/cards'),
                borderRadius: BorderRadius.circular(14),
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                  child: Row(
                    children: [
                      Icon(CupertinoIcons.creditcard, color: p.mutedFg, size: 18),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Text(
                          'Saved cards',
                          style: GoogleFonts.manrope(
                              fontSize: 14, fontWeight: FontWeight.w500, color: p.fg),
                        ),
                      ),
                      Icon(CupertinoIcons.chevron_forward,
                          size: 16, color: p.mutedFg),
                    ],
                  ),
                ),
              ),
            ),
            const SizedBox(height: 24),
          ],

          // About
          const AppSectionHeader('About'),
          const SizedBox(height: 10),
          AppCard(
            padding: EdgeInsets.zero,
            child: Column(
              children: [
                const _Row(label: 'Version', value: '0.1.0 · MVP'),
                Divider(height: 1, color: p.border),
                const _Row(label: 'Model', value: 'Gemini 2.5 Flash'),
                Divider(height: 1, color: p.border),
                const _Row(label: 'Built in', value: 'Almaty, KZ'),
              ],
            ),
          ),

          const SizedBox(height: 32),
          OutlinedButton.icon(
            onPressed: () async {
              await ref.read(authViewModelProvider.notifier).logout();
              if (context.mounted) context.go('/login');
            },
            icon: Icon(CupertinoIcons.square_arrow_right, size: 18, color: p.destructive),
            label: Text(
              tr(ref, 'logout'),
              style: GoogleFonts.manrope(
                  color: p.destructive, fontWeight: FontWeight.w600),
            ),
            style: OutlinedButton.styleFrom(
              side: BorderSide(color: p.destructive.withValues(alpha: 0.3)),
            ),
          ),
        ],
      ),
    );
  }

  static String _initial(String name, String email) {
    final src = name.trim().isNotEmpty ? name.trim() : email.trim();
    if (src.isEmpty) return '?';
    return src.characters.first.toUpperCase();
  }
}

class _Row extends StatelessWidget {
  const _Row({required this.label, required this.value});
  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 13),
      child: Row(
        children: [
          Text(
            label,
            style: GoogleFonts.manrope(
                fontSize: 13, color: p.mutedFg, fontWeight: FontWeight.w500),
          ),
          const Spacer(),
          Flexible(
            child: Text(
              value,
              textAlign: TextAlign.end,
              overflow: TextOverflow.ellipsis,
              style: GoogleFonts.manrope(
                  fontSize: 13.5, color: p.fg, fontWeight: FontWeight.w600),
            ),
          ),
        ],
      ),
    );
  }
}

class _LocaleRow extends StatelessWidget {
  const _LocaleRow({required this.locale, required this.selected, required this.onTap});

  final AppLocale locale;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        child: Row(
          children: [
            Expanded(
              child: Text(
                locale.label,
                style: GoogleFonts.manrope(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: p.fg,
                ),
              ),
            ),
            if (selected) Icon(CupertinoIcons.checkmark_alt, color: p.primary, size: 18),
          ],
        ),
      ),
    );
  }
}

class _ThemeTile extends StatelessWidget {
  const _ThemeTile({
    required this.icon,
    required this.label,
    required this.selected,
    required this.onTap,
  });

  final IconData icon;
  final String label;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 14),
          decoration: BoxDecoration(
            color: selected ? p.primary.withValues(alpha: 0.08) : p.card,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: selected ? p.primary : p.border,
              width: selected ? 1.4 : 1,
            ),
          ),
          child: Column(
            children: [
              Icon(icon, color: selected ? p.primary : p.mutedFg, size: 20),
              const SizedBox(height: 6),
              Text(
                label,
                style: GoogleFonts.manrope(
                  fontSize: 12.5,
                  fontWeight: FontWeight.w700,
                  color: p.fg,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
