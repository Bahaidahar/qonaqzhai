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
      // Persistent bottom-nav shell. Each tab is its own branch kept alive in an
      // IndexedStack, so switching tabs is instant and preserves scroll/state —
      // it feels like tab switching, not page navigation. Branch order is fixed;
      // the nav bar shows a role-specific subset (see _HomeShell).
      StatefulShellRoute.indexedStack(
        builder: (context, state, navigationShell) =>
            _HomeShell(navigationShell: navigationShell),
        branches: [
          StatefulShellBranch(routes: [
            GoRoute(path: '/', builder: (_, __) => const ChatScreen()),
          ]),
          StatefulShellBranch(routes: [
            GoRoute(path: '/vendors', builder: (_, __) => const VendorCatalogScreen()),
          ]),
          StatefulShellBranch(routes: [
            GoRoute(path: '/bookings', builder: (_, __) => const BookingsScreen()),
          ]),
          StatefulShellBranch(routes: [
            GoRoute(path: '/notifications', builder: (_, __) => const NotificationsScreen()),
          ]),
          StatefulShellBranch(routes: [
            GoRoute(path: '/threads', builder: (_, __) => const ThreadsScreen()),
          ]),
          StatefulShellBranch(routes: [
            GoRoute(path: '/cards', builder: (_, __) => const CardsScreen()),
          ]),
          StatefulShellBranch(routes: [
            GoRoute(path: '/vendor-profile', builder: (_, __) => const VendorSelfScreen()),
          ]),
          StatefulShellBranch(routes: [
            GoRoute(path: '/settings', builder: (_, __) => const SettingsScreen()),
          ]),
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

// Branch indices into StatefulShellRoute (must match the branch order above).
const int _bChat = 0;
const int _bVendors = 1;
const int _bBookings = 2;
const int _bNotifications = 3;
const int _bThreads = 4;
// ignore: unused_element
const int _bCards = 5;
const int _bVendorProfile = 6;
const int _bSettings = 7;

class _HomeShell extends ConsumerWidget {
  const _HomeShell({required this.navigationShell});
  final StatefulNavigationShell navigationShell;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(authViewModelProvider).user;
    final role = user?.role ?? UserRole.customer;
    final tabs = _tabsFor(role, ref);

    int selected = tabs.indexWhere((t) => t.branch == navigationShell.currentIndex);
    if (selected < 0) selected = 0;

    return Scaffold(
      // IndexedStack keeps every visited tab mounted -> instant switch, no
      // page transition, state preserved per tab.
      body: navigationShell,
      bottomNavigationBar: NavigationBar(
        selectedIndex: selected,
        onDestinationSelected: (i) => navigationShell.goBranch(
          tabs[i].branch,
          // tapping the active tab resets it to the branch root
          initialLocation: tabs[i].branch == navigationShell.currentIndex,
        ),
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
          _NavTab(_bChat, CupertinoIcons.chat_bubble_text, tr(ref, 'nav_chat')),
          _NavTab(_bVendors, CupertinoIcons.square_grid_2x2, tr(ref, 'nav_vendors')),
          _NavTab(_bBookings, CupertinoIcons.calendar, tr(ref, 'nav_bookings')),
          _NavTab(_bThreads, CupertinoIcons.chat_bubble_2, tr(ref, 'nav_threads')),
          _NavTab(_bSettings, CupertinoIcons.settings, tr(ref, 'nav_settings')),
        ];
      case UserRole.vendor:
        return [
          _NavTab(_bVendorProfile, CupertinoIcons.house, tr(ref, 'nav_vendor_profile')),
          _NavTab(_bBookings, CupertinoIcons.calendar, tr(ref, 'nav_bookings')),
          _NavTab(_bThreads, CupertinoIcons.chat_bubble_2, tr(ref, 'nav_threads')),
          _NavTab(_bNotifications, CupertinoIcons.bell, tr(ref, 'nav_notifications')),
          _NavTab(_bSettings, CupertinoIcons.settings, tr(ref, 'nav_settings')),
        ];
    }
  }
}

class _NavTab {
  const _NavTab(this.branch, this.icon, this.label);
  final int branch;
  final IconData icon;
  final String label;
}
