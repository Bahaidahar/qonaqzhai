/// Central catalog of HTTP endpoints exposed by the qonaqzhai gateway.
///
/// Mobile only ships customer + vendor flows — admin endpoints live on the
/// web app and are intentionally absent here so a stray reference fails the
/// build instead of silently shipping unreachable code.
class ApiEndpoints {
  ApiEndpoints._();

  static const String baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );

  // auth
  static const String signup = '/api/signup';
  static const String login = '/api/login';
  static const String refresh = '/api/refresh';
  static const String logout = '/api/logout';
  static const String forgotPassword = '/api/forgot-password';
  static const String resetPassword = '/api/reset-password';
  static const String me = '/api/me';

  // ai chat
  static const String chat = '/api/chat';
  static const String chats = '/api/chats';
  static String chatById(String id) => '/api/chats/$id';

  // vendors (public)
  static const String vendors = '/api/vendors';
  static String vendor(String id) => '/api/vendors/$id';
  static String vendorReviews(String id) => '/api/vendors/$id/reviews';
  static String vendorPublicServices(String id) => '/api/vendors/$id/services';

  // vendor self
  static const String myVendor = '/api/me/vendor';
  static const String vendorPhotos = '/api/me/vendor/photos';
  static String vendorPhoto(String id) => '/api/me/vendor/photos/$id';
  static const String vendorServices = '/api/me/vendor/services';
  static String vendorService(String id) => '/api/me/vendor/services/$id';

  // bookings
  static const String bookings = '/api/bookings';
  static String booking(String id) => '/api/bookings/$id';
  static String startPayment(String id) => '/api/bookings/$id/pay';
  static String mockPayment(String id) => '/api/bookings/$id/pay/mock';

  // reviews
  static const String reviews = '/api/reviews';

  // notifications
  static const String notifications = '/api/notifications';
  static const String fcmTokens = '/api/notifications/tokens';

  // threads
  static const String threads = '/api/threads';
  static String thread(String id) => '/api/threads/$id';
  static String threadMessages(String id) => '/api/threads/$id/messages';

  // saved cards
  static const String cards = '/api/cards';
  static String card(String id) => '/api/cards/$id';
  static String cardDefault(String id) => '/api/cards/$id/default';

  // photos
  static String photo(String id) => '/api/photos/$id';
}
