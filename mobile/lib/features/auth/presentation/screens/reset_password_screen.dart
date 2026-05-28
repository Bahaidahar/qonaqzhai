import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/i18n/i18n.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../../core/ui/ui.dart';
import '../viewmodels/auth_viewmodel.dart';

class ResetPasswordScreen extends ConsumerStatefulWidget {
  const ResetPasswordScreen({super.key, this.initialToken});

  final String? initialToken;

  @override
  ConsumerState<ResetPasswordScreen> createState() => _ResetPasswordScreenState();
}

class _ResetPasswordScreenState extends ConsumerState<ResetPasswordScreen> {
  late final TextEditingController _token =
      TextEditingController(text: widget.initialToken ?? '');
  final _password = TextEditingController();
  bool _busy = false;

  @override
  void dispose() {
    _token.dispose();
    _password.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    setState(() => _busy = true);
    final ok = await ref.read(authViewModelProvider.notifier).resetPassword(
          token: _token.text.trim(),
          newPassword: _password.text,
        );
    if (!mounted) return;
    setState(() => _busy = false);
    if (ok) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(tr(ref, 'reset_done'))),
      );
      context.go('/login');
    }
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(authViewModelProvider);
    final p = AppPalette.of(context);

    return Scaffold(
      appBar: AppBar(),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.fromLTRB(24, 8, 24, 24),
          children: [
            AppPageHeader(
              title: tr(ref, 'reset_title'),
              subtitle: 'Paste the token from the email and set a new password.',
            ),
            const SizedBox(height: 24),
            _LabeledField(
              label: tr(ref, 'reset_token'),
              child: TextField(controller: _token),
            ),
            const SizedBox(height: 14),
            _LabeledField(
              label: tr(ref, 'reset_new_password'),
              child: TextField(controller: _password, obscureText: true),
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
            FilledButton(
              onPressed: _busy ? null : _submit,
              child: _busy
                  ? const SizedBox(
                      width: 18,
                      height: 18,
                      child: CircularProgressIndicator(
                          strokeWidth: 2, color: Colors.white),
                    )
                  : Text(tr(ref, 'reset_btn')),
            ),
            const SizedBox(height: 6),
            Center(
              child: TextButton(
                onPressed: () => context.go('/login'),
                child: Text(tr(ref, 'common_back')),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _LabeledField extends StatelessWidget {
  const _LabeledField({required this.label, required this.child});
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
