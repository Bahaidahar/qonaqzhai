import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../features/onboarding/presentation/onboarding_screen.dart';
import '../../features/ai_chat/presentation/screens/chat_screen.dart';
import '../../features/auth/domain/entities/user.dart';
import '../../features/auth/presentation/screens/forgot_password_screen.dart';
import '../../features/auth/presentation/screens/login_screen.dart';
import '../../features/auth/presentation/screens/reset_password_screen.dart';
import '../../features/auth/presentation/screens/signup_screen.dart';
import '../../features/auth/presentation/viewmodels/auth_viewmodel.dart';
import '../../features/booking/presentation/screens/booking_form_screen.dart';
import '../../features/booking/presentation/screens/bookings_screen.dart';
import '../../features/cards/presentation/cards_screen.dart';
import '../../features/messaging/presentation/screens/thread_chat_screen.dart';
import '../../features/messaging/presentation/screens/threads_screen.dart';
import '../../features/notifications/presentation/screens/notifications_screen.dart';
import '../../features/payment/presentation/payment_screen.dart';
import '../../features/reviews/presentation/screens/review_submit_screen.dart';
import '../../features/settings/presentation/settings_screen.dart';
import '../../features/vendor_catalog/presentation/screens/vendor_catalog_screen.dart';
import '../../features/vendor_catalog/presentation/screens/vendor_detail_screen.dart';
import '../../features/vendor_self/presentation/vendor_self_screen.dart';
import '../i18n/i18n.dart';

final onboardingSeenStateProvider = StateProvider<bool?>((_) => null);

final appRouterProvider = Provider<GoRouter>((ref) {
  // Resolve onboarding flag on app start.
  Future(() async {
    final p = await SharedPreferences.getInstance();
    ref.read(onboardingSeenStateProvider.notifier).state = p.getBool(onboardingSeenKey) ?? false;
  });
  return GoRouter(
    initialLocation: '/',
    refreshListenable: _AuthListenable(ref),
    redirect: (context, state) {
      final auth = ref.read(authViewModelProvider);
      if (auth.loading) return null;
      final onboardingSeen = ref.read(onboardingSeenStateProvider);
      if (onboardingSeen == null) return null;
      final loggedIn = auth.user != null;
      final loc = state.matchedLocation;
      if (!onboardingSeen && loc != '/onboarding') return '/onboarding';
      if (onboardingSeen && loc == '/onboarding') return loggedIn ? '/' : '/login';
      final publicPaths = {'/login', '/signup', '/forgot', '/reset', '/onboarding'};
      final goingToAuth = publicPaths.contains(loc);
      if (!loggedIn && !goingToAuth) return '/login';
      if (loggedIn && (loc == '/login' || loc == '/signup')) return '/';
      return null;
    },
    routes: [
      ShellRoute(
        builder: (context, state, child) => _HomeShell(child: child),
        routes: [
          GoRoute(path: '/', builder: (_, __) => const ChatScreen()),
          GoRoute(path: '/vendors', builder: (_, __) => const VendorCatalogScreen()),
          GoRoute(path: '/bookings', builder: (_, __) => const BookingsScreen()),
          GoRoute(path: '/notifications', builder: (_, __) => const NotificationsScreen()),
          GoRoute(path: '/threads', builder: (_, __) => const ThreadsScreen()),
          GoRoute(path: '/cards', builder: (_, __) => const CardsScreen()),
          GoRoute(path: '/vendor-profile', builder: (_, __) => const VendorSelfScreen()),
          GoRoute(path: '/settings', builder: (_, __) => const SettingsScreen()),
        ],
      ),
      GoRoute(
        path: '/threads/:id',
        builder: (_, state) => ThreadChatScreen(id: state.pathParameters['id']!),
      ),
      GoRoute(
        path: '/vendors/:id',
        builder: (_, state) => VendorDetailScreen(id: state.pathParameters['id']!),
      ),
      GoRoute(
        path: '/bookings/new',
        builder: (_, state) {
          final qp = state.uri.queryParameters;
          final vendor = qp['vendor'] ?? '';
          final price = int.tryParse(qp['price'] ?? '') ?? 0;
          return BookingFormScreen(
            vendorId: vendor,
            priceFrom: price,
            serviceId: qp['service'],
            serviceUnit: qp['unit'],
          );
        },
      ),
      GoRoute(
        path: '/reviews/new',
        builder: (_, state) {
          final booking = state.uri.queryParameters['booking'] ?? '';
          final vendor = state.uri.queryParameters['vendor'];
          return ReviewSubmitScreen(bookingId: booking, vendorId: vendor);
        },
      ),
      GoRoute(
        path: '/pay',
        builder: (_, state) {
          final id = state.uri.queryParameters['booking'] ?? '';
          final amount = int.tryParse(state.uri.queryParameters['amount'] ?? '') ?? 0;
          return PaymentScreen(bookingId: id, amount: amount);
        },
      ),
      GoRoute(path: '/onboarding', builder: (_, __) => const OnboardingScreen()),
      GoRoute(path: '/login', builder: (_, __) => const LoginScreen()),
      GoRoute(path: '/signup', builder: (_, __) => const SignupScreen()),
      GoRoute(path: '/forgot', builder: (_, __) => const ForgotPasswordScreen()),
      GoRoute(
        path: '/reset',
        builder: (_, state) => ResetPasswordScreen(
          initialToken: state.uri.queryParameters['token'],
        ),
      ),
    ],
  );
});

class _AuthListenable extends ChangeNotifier {
  _AuthListenable(this._ref) {
    _ref.listen(authViewModelProvider, (_, __) => notifyListeners());
    _ref.listen(onboardingSeenStateProvider, (_, __) => notifyListeners());
  }
  // ignore: unused_field
  final Ref _ref;
}

class _HomeShell extends ConsumerWidget {
  const _HomeShell({required this.child});
  final Widget child;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(authViewModelProvider).user;
    final role = user?.role ?? UserRole.customer;
    final tabs = _tabsFor(role, ref);
    final loc = GoRouterState.of(context).matchedLocation;
    int index = tabs.indexWhere((t) => t.match(loc));
    if (index < 0) index = 0;

    return Scaffold(
      body: child,
      bottomNavigationBar: NavigationBar(
        selectedIndex: index,
        onDestinationSelected: (i) => context.go(tabs[i].path),
        destinations: [
          for (final t in tabs)
            NavigationDestination(icon: Icon(t.icon), label: t.label),
        ],
      ),
    );
  }

  List<_NavTab> _tabsFor(UserRole role, WidgetRef ref) {
    switch (role) {
      case UserRole.customer:
        return [
          _NavTab('/', CupertinoIcons.chat_bubble_text, tr(ref, 'nav_chat'), exact: true),
          _NavTab('/vendors', CupertinoIcons.square_grid_2x2, tr(ref, 'nav_vendors')),
          _NavTab('/bookings', CupertinoIcons.calendar, tr(ref, 'nav_bookings')),
          _NavTab('/threads', CupertinoIcons.chat_bubble_2, tr(ref, 'nav_threads')),
          _NavTab('/settings', CupertinoIcons.settings, tr(ref, 'nav_settings')),
        ];
      case UserRole.vendor:
        return [
          _NavTab('/vendor-profile', CupertinoIcons.house, tr(ref, 'nav_vendor_profile'), exact: true),
          _NavTab('/bookings', CupertinoIcons.calendar, tr(ref, 'nav_bookings')),
          _NavTab('/threads', CupertinoIcons.chat_bubble_2, tr(ref, 'nav_threads')),
          _NavTab('/notifications', CupertinoIcons.bell, tr(ref, 'nav_notifications')),
          _NavTab('/settings', CupertinoIcons.settings, tr(ref, 'nav_settings')),
        ];
    }
  }
}

class _NavTab {
  const _NavTab(this.path, this.icon, this.label, {this.exact = false});
  final String path;
  final IconData icon;
  final String label;
  final bool exact;
  bool match(String loc) => exact ? loc == path : loc.startsWith(path);
}
