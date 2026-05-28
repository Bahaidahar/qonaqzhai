import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../../core/i18n/i18n.dart';
import '../../../core/router/app_router.dart';
import '../../../core/theme/app_theme.dart';

const String onboardingSeenKey = 'qz_onboarding_seen';

class OnboardingScreen extends ConsumerStatefulWidget {
  const OnboardingScreen({super.key});

  @override
  ConsumerState<OnboardingScreen> createState() => _OnboardingScreenState();
}

class _OnboardingScreenState extends ConsumerState<OnboardingScreen> {
  final _page = PageController();
  int _index = 0;

  @override
  void dispose() {
    _page.dispose();
    super.dispose();
  }

  Future<void> _finish() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool(onboardingSeenKey, true);
    ref.read(onboardingSeenStateProvider.notifier).state = true;
    if (!mounted) return;
    context.go('/');
  }

  @override
  Widget build(BuildContext context) {
    final locale = ref.watch(localeProvider);
    final pages = <_OnboardPage>[
      _OnboardPage(
        icon: CupertinoIcons.globe,
        title: trFor(locale, 'onb_lang_title'),
        body: trFor(locale, 'onb_lang_body'),
        widget: _LanguagePicker(active: locale),
      ),
      _OnboardPage(
        icon: CupertinoIcons.sparkles,
        title: trFor(locale, 'onb_ai_title'),
        body: trFor(locale, 'onb_ai_body'),
      ),
      _OnboardPage(
        icon: CupertinoIcons.square_grid_2x2,
        title: trFor(locale, 'onb_vendors_title'),
        body: trFor(locale, 'onb_vendors_body'),
      ),
      _OnboardPage(
        icon: CupertinoIcons.calendar_badge_plus,
        title: trFor(locale, 'onb_book_title'),
        body: trFor(locale, 'onb_book_body'),
      ),
    ];
    final isLast = _index == pages.length - 1;
    return Scaffold(
      body: SafeArea(
        child: Column(children: [
          Align(
            alignment: Alignment.centerRight,
            child: TextButton(
              onPressed: _finish,
              child: Text(trFor(locale, 'onb_skip')),
            ),
          ),
          Expanded(
            child: PageView.builder(
              controller: _page,
              itemCount: pages.length,
              onPageChanged: (i) => setState(() => _index = i),
              itemBuilder: (_, i) => _OnboardSlide(page: pages[i]),
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 16),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: List.generate(pages.length, (i) {
                final active = i == _index;
                return AnimatedContainer(
                  duration: const Duration(milliseconds: 200),
                  margin: const EdgeInsets.symmetric(horizontal: 4),
                  height: 8,
                  width: active ? 24 : 8,
                  decoration: BoxDecoration(
                    color: active
                        ? Theme.of(context).colorScheme.primary
                        : Theme.of(context).colorScheme.outline,
                    borderRadius: BorderRadius.circular(999),
                  ),
                );
              }),
            ),
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 24),
            child: SizedBox(
              width: double.infinity,
              child: FilledButton(
                onPressed: () {
                  if (isLast) {
                    _finish();
                  } else {
                    _page.nextPage(duration: const Duration(milliseconds: 250), curve: Curves.easeOut);
                  }
                },
                child: Text(isLast ? trFor(locale, 'onb_start') : trFor(locale, 'onb_next')),
              ),
            ),
          ),
        ]),
      ),
    );
  }
}

class _OnboardPage {
  const _OnboardPage({required this.icon, required this.title, required this.body, this.widget});
  final IconData icon;
  final String title;
  final String body;
  final Widget? widget;
}

class _OnboardSlide extends StatelessWidget {
  const _OnboardSlide({required this.page});
  final _OnboardPage page;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 28),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            padding: const EdgeInsets.all(22),
            decoration: BoxDecoration(
              color: p.primary.withValues(alpha: 0.1),
              borderRadius: BorderRadius.circular(20),
              border: Border.all(color: p.primary.withValues(alpha: 0.25)),
            ),
            child: Icon(page.icon, size: 36, color: p.primary),
          ),
          const SizedBox(height: 28),
          Text(
            page.title,
            textAlign: TextAlign.center,
            style: GoogleFonts.manrope(
              fontSize: 28,
              fontWeight: FontWeight.w700,
              color: p.fg,
              letterSpacing: -0.8,
              height: 1.1,
            ),
          ),
          const SizedBox(height: 12),
          Text(
            page.body,
            textAlign: TextAlign.center,
            style: GoogleFonts.manrope(
              fontSize: 14,
              color: p.mutedFg,
              height: 1.55,
            ),
          ),
          if (page.widget != null) ...[
            const SizedBox(height: 24),
            page.widget!,
          ],
        ],
      ),
    );
  }
}

class _LanguagePicker extends ConsumerWidget {
  const _LanguagePicker({required this.active});
  final AppLocale active;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final p = AppPalette.of(context);
    return Wrap(
      spacing: 8,
      runSpacing: 8,
      alignment: WrapAlignment.center,
      children: [
        for (final l in AppLocale.values)
          GestureDetector(
            onTap: () => ref.read(localeProvider.notifier).set(l),
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
              decoration: BoxDecoration(
                color: l == active ? p.primary : p.muted,
                borderRadius: BorderRadius.circular(999),
                border: Border.all(color: l == active ? p.primary : p.border),
              ),
              child: Text(
                l.label,
                style: GoogleFonts.manrope(
                  color: l == active ? p.onPrimary : p.fg,
                  fontWeight: FontWeight.w600,
                  fontSize: 12.5,
                ),
              ),
            ),
          ),
      ],
    );
  }
}
