class ApiEndpoints {
  ApiEndpoints._();

  static const String baseUrl =
      String.fromEnvironment('API_BASE_URL', defaultValue: 'http://localhost:8080');

  static const String signup = '/api/signup';
  static const String login = '/api/login';
  static const String refresh = '/api/auth/refresh';
  static const String logout = '/api/auth/logout';
  static const String forgotPassword = '/api/auth/forgot-password';
  static const String resetPassword = '/api/auth/reset-password';
  static const String me = '/api/me';
  static const String chat = '/api/chat';
  static const String vendors = '/api/vendors';
  static const String myVendor = '/api/vendor';
  static const String vendorPhotos = '/api/vendor/photos';
  static const String bookings = '/api/bookings';
  static const String reviews = '/api/reviews';
  static const String notifications = '/api/notifications';
  static const String fcmTokens = '/api/notifications/tokens';

  static String vendor(String id) => '/api/vendors/$id';
  static String vendorReviews(String id) => '/api/vendors/$id/reviews';
  static String booking(String id) => '/api/bookings/$id';
  static String startPayment(String id) => '/api/bookings/$id/pay';
  static String photo(String id) => '/api/photos/$id';
}
