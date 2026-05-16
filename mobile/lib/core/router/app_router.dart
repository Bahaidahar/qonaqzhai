import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../features/ai_chat/presentation/screens/chat_screen.dart';
import '../../features/auth/presentation/screens/forgot_password_screen.dart';
import '../../features/auth/presentation/screens/login_screen.dart';
import '../../features/auth/presentation/screens/reset_password_screen.dart';
import '../../features/auth/presentation/screens/signup_screen.dart';
import '../../features/auth/presentation/viewmodels/auth_viewmodel.dart';
import '../../features/booking/presentation/screens/booking_form_screen.dart';
import '../../features/booking/presentation/screens/bookings_screen.dart';
import '../../features/notifications/presentation/screens/notifications_screen.dart';
import '../../features/reviews/presentation/screens/review_submit_screen.dart';
import '../../features/vendor_catalog/presentation/screens/vendor_catalog_screen.dart';
import '../../features/vendor_catalog/presentation/screens/vendor_detail_screen.dart';

final appRouterProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: '/',
    refreshListenable: _AuthListenable(ref),
    redirect: (context, state) {
      final auth = ref.read(authViewModelProvider);
      if (auth.loading) return null;
      final loggedIn = auth.user != null;
      final publicPaths = {'/login', '/signup', '/forgot', '/reset'};
      final goingToAuth = publicPaths.contains(state.matchedLocation);
      if (!loggedIn && !goingToAuth) return '/login';
      if (loggedIn && (state.matchedLocation == '/login' || state.matchedLocation == '/signup')) return '/';
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
        ],
      ),
      GoRoute(
        path: '/vendors/:id',
        builder: (_, state) => VendorDetailScreen(id: state.pathParameters['id']!),
      ),
      GoRoute(
        path: '/bookings/new',
        builder: (_, state) {
          final vendor = state.uri.queryParameters['vendor'] ?? '';
          final price = int.tryParse(state.uri.queryParameters['price'] ?? '') ?? 0;
          return BookingFormScreen(vendorId: vendor, priceFrom: price);
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
  }
  // ignore: unused_field
  final Ref _ref;
}

class _HomeShell extends StatelessWidget {
  const _HomeShell({required this.child});
  final Widget child;

  int _indexFor(BuildContext context) {
    final loc = GoRouterState.of(context).matchedLocation;
    if (loc.startsWith('/vendors')) return 1;
    if (loc.startsWith('/bookings')) return 2;
    if (loc.startsWith('/notifications')) return 3;
    return 0;
  }

  @override
  Widget build(BuildContext context) {
    final index = _indexFor(context);
    return Scaffold(
      body: child,
      bottomNavigationBar: NavigationBar(
        selectedIndex: index,
        onDestinationSelected: (i) {
          switch (i) {
            case 0:
              context.go('/');
              break;
            case 1:
              context.go('/vendors');
              break;
            case 2:
              context.go('/bookings');
              break;
            case 3:
              context.go('/notifications');
              break;
          }
        },
        destinations: const [
          NavigationDestination(icon: Icon(Icons.chat_bubble_outline), label: 'AI'),
          NavigationDestination(icon: Icon(Icons.store_outlined), label: 'Vendors'),
          NavigationDestination(icon: Icon(Icons.calendar_month_outlined), label: 'Bookings'),
          NavigationDestination(icon: Icon(Icons.notifications_outlined), label: 'Inbox'),
        ],
      ),
    );
  }
}
