import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/i18n/i18n.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../../core/ui/ui.dart';
import '../viewmodels/auth_viewmodel.dart';

class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});

  @override
  ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  final _email = TextEditingController();
  final _password = TextEditingController();
  bool _showPw = false;

  @override
  void dispose() {
    _email.dispose();
    _password.dispose();
    super.dispose();
  }

  void _fill(String email, String password) {
    _email.text = email;
    _password.text = password;
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(authViewModelProvider);
    final p = AppPalette.of(context);

    return Scaffold(
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.fromLTRB(24, 40, 24, 32),
          children: [
            Row(
              children: [
                Container(
                  width: 36,
                  height: 36,
                  decoration: BoxDecoration(
                    color: p.primary,
                    borderRadius: BorderRadius.circular(10),
                  ),
                  alignment: Alignment.center,
                  child: Text(
                    'q',
                    style: GoogleFonts.manrope(
                      color: p.onPrimary,
                      fontWeight: FontWeight.w800,
                      fontSize: 18,
                    ),
                  ),
                ),
                const SizedBox(width: 10),
                Text(
                  'qonaqzhai',
                  style: GoogleFonts.manrope(
                    fontSize: 16,
                    fontWeight: FontWeight.w700,
                    color: p.fg,
                    letterSpacing: -0.3,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 36),
            AppPageHeader(
              title: tr(ref, 'login_title'),
              subtitle: 'Welcome back. Sign in to plan, manage, and message.',
            ),
            const SizedBox(height: 28),
            _Field(
              label: tr(ref, 'login_email'),
              child: Semantics(
                identifier: 'login-email',
                child: TextField(
                  key: const Key('login-email'),
                  controller: _email,
                  keyboardType: TextInputType.emailAddress,
                  autocorrect: false,
                  decoration: const InputDecoration(hintText: 'you@example.com'),
                ),
              ),
            ),
            const SizedBox(height: 14),
            _Field(
              label: tr(ref, 'login_password'),
              child: Semantics(
                identifier: 'login-password',
                child: TextField(
                  key: const Key('login-password'),
                  controller: _password,
                  obscureText: !_showPw,
                  decoration: InputDecoration(
                    hintText: '••••••••',
                    suffixIcon: IconButton(
                      icon: Icon(
                        _showPw ? CupertinoIcons.eye_slash : CupertinoIcons.eye,
                        size: 18,
                        color: p.mutedFg,
                      ),
                      onPressed: () => setState(() => _showPw = !_showPw),
                    ),
                  ),
                ),
              ),
            ),
            if (state.error != null) ...[
              const SizedBox(height: 14),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                decoration: BoxDecoration(
                  color: p.destructive.withValues(alpha: 0.08),
                  borderRadius: BorderRadius.circular(10),
                  border: Border.all(color: p.destructive.withValues(alpha: 0.3)),
                ),
                child: Text(state.error!,
                    style: GoogleFonts.manrope(fontSize: 12, color: p.destructive)),
              ),
            ],
            const SizedBox(height: 22),
            Semantics(
              identifier: 'login-submit',
              child: FilledButton(
                key: const Key('login-submit'),
                onPressed: state.loading
                    ? null
                    : () => ref
                        .read(authViewModelProvider.notifier)
                        .login(_email.text.trim(), _password.text),
                child: state.loading
                    ? const SizedBox(
                        width: 18,
                        height: 18,
                        child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                      )
                    : Text(tr(ref, 'login_btn')),
              ),
            ),
            const SizedBox(height: 6),
            Align(
              alignment: Alignment.centerRight,
              child: TextButton(
                onPressed: () => context.push('/forgot'),
                child: Text(tr(ref, 'login_forgot')),
              ),
            ),
            const SizedBox(height: 16),
            Row(children: [
              Expanded(child: Divider(color: p.border)),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10),
                child: Text(
                  'or',
                  style: GoogleFonts.manrope(fontSize: 11, color: p.mutedFg),
                ),
              ),
              Expanded(child: Divider(color: p.border)),
            ]),
            const SizedBox(height: 16),
            OutlinedButton.icon(
              onPressed: () => context.push('/signup'),
              icon: const Icon(CupertinoIcons.person_add, size: 16),
              label: Text('${tr(ref, 'login_no_account')} ${tr(ref, 'login_signup_link')}'),
            ),
            const SizedBox(height: 32),
            _DemoAccountsCard(onPick: _fill),
          ],
        ),
      ),
    );
  }
}

class _Field extends StatelessWidget {
  const _Field({required this.label, required this.child});
  final String label;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: GoogleFonts.manrope(
            fontSize: 12.5,
            fontWeight: FontWeight.w600,
            color: p.fg,
          ),
        ),
        const SizedBox(height: 6),
        child,
      ],
    );
  }
}

class _DemoAccountsCard extends ConsumerWidget {
  const _DemoAccountsCard({required this.onPick});
  final void Function(String email, String password) onPick;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final p = AppPalette.of(context);
    // Mobile only ships customer + vendor; admin lives on the web app.
    final accounts = const [
      _Acct('customer1@demo.kz', 'demo12345', 'demo_acc_customer'),
      _Acct('vendor1@demo.kz', 'demo12345', 'demo_acc_vendor'),
    ];
    return AppCard(
      padding: const EdgeInsets.fromLTRB(16, 14, 16, 10),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const AppSectionHeader('Demo accounts'),
          const SizedBox(height: 10),
          for (final a in accounts)
            Semantics(
              identifier: 'demo-${a.labelKey}',
              child: InkWell(
                onTap: () => onPick(a.email, a.password),
                borderRadius: BorderRadius.circular(8),
                child: Padding(
                padding: const EdgeInsets.symmetric(vertical: 10),
                child: Row(
                  children: [
                    Container(
                      width: 32,
                      height: 32,
                      decoration: BoxDecoration(
                        color: p.muted,
                        shape: BoxShape.circle,
                        border: Border.all(color: p.border),
                      ),
                      alignment: Alignment.center,
                      child: Icon(CupertinoIcons.person, size: 14, color: p.mutedFg),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            tr(ref, a.labelKey),
                            style: GoogleFonts.manrope(
                              fontWeight: FontWeight.w600,
                              fontSize: 13,
                              color: p.fg,
                            ),
                          ),
                          Text(
                            '${a.email} · ${a.password}',
                            style: GoogleFonts.manrope(
                              fontSize: 11.5,
                              color: p.mutedFg,
                            ),
                          ),
                        ],
                      ),
                    ),
                    Icon(CupertinoIcons.arrow_right, size: 14, color: p.mutedFg),
                  ],
                ),
              ),
            ),
            ),
        ],
      ),
    );
  }
}

class _Acct {
  const _Acct(this.email, this.password, this.labelKey);
  final String email;
  final String password;
  final String labelKey;
}
