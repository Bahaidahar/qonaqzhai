import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/i18n/i18n.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../../core/ui/ui.dart';
import '../../domain/entities/user.dart';
import '../viewmodels/auth_viewmodel.dart';

class SignupScreen extends ConsumerStatefulWidget {
  const SignupScreen({super.key});

  @override
  ConsumerState<SignupScreen> createState() => _SignupScreenState();
}

class _SignupScreenState extends ConsumerState<SignupScreen> {
  final _email = TextEditingController();
  final _password = TextEditingController();
  final _name = TextEditingController();
  UserRole _role = UserRole.customer;
  bool _showPw = false;

  @override
  void dispose() {
    _email.dispose();
    _password.dispose();
    _name.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(authViewModelProvider);
    final p = AppPalette.of(context);

    return Scaffold(
      appBar: AppBar(),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.fromLTRB(24, 8, 24, 32),
          children: [
            AppPageHeader(
              title: tr(ref, 'signup_title'),
              subtitle: 'Choose your role and create your account.',
            ),
            const SizedBox(height: 24),
            const AppSectionHeader('I am a'),
            const SizedBox(height: 10),
            Row(
              children: [
                Expanded(
                  child: _RoleTile(
                    icon: CupertinoIcons.person_2,
                    label: tr(ref, 'signup_role_customer'),
                    selected: _role == UserRole.customer,
                    onTap: () => setState(() => _role = UserRole.customer),
                  ),
                ),
                const SizedBox(width: 10),
                Expanded(
                  child: _RoleTile(
                    icon: CupertinoIcons.briefcase,
                    label: tr(ref, 'signup_role_vendor'),
                    selected: _role == UserRole.vendor,
                    onTap: () => setState(() => _role = UserRole.vendor),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 22),
            _Field(
              label: tr(ref, 'signup_name'),
              child: TextField(controller: _name),
            ),
            const SizedBox(height: 14),
            _Field(
              label: tr(ref, 'login_email'),
              child: TextField(
                controller: _email,
                keyboardType: TextInputType.emailAddress,
                autocorrect: false,
              ),
            ),
            const SizedBox(height: 14),
            _Field(
              label: tr(ref, 'login_password'),
              child: TextField(
                controller: _password,
                obscureText: !_showPw,
                decoration: InputDecoration(
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
              onPressed: state.loading
                  ? null
                  : () => ref.read(authViewModelProvider.notifier).signup(
                        email: _email.text.trim(),
                        password: _password.text,
                        role: _role,
                        name: _name.text.trim().isEmpty ? null : _name.text.trim(),
                      ),
              child: state.loading
                  ? const SizedBox(
                      width: 18,
                      height: 18,
                      child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                    )
                  : Text(tr(ref, 'signup_btn')),
            ),
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

class _RoleTile extends StatelessWidget {
  const _RoleTile({
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
          padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 16),
          decoration: BoxDecoration(
            color: selected ? p.primary.withValues(alpha: 0.08) : p.card,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: selected ? p.primary : p.border,
              width: selected ? 1.4 : 1,
            ),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Icon(icon, color: selected ? p.primary : p.mutedFg, size: 22),
              const SizedBox(height: 8),
              Text(
                label,
                style: GoogleFonts.manrope(
                  fontSize: 13.5,
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
