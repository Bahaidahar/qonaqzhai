/// Mobile is a two-role product: customers plan events, vendors run listings.
/// Admin moderation lives on the web app only — if the API ever returns
/// `role: "admin"` we fall back to customer so the user is not stranded with
/// an unknown role.
enum UserRole { customer, vendor }

extension UserRoleX on UserRole {
  static UserRole parse(String s) {
    switch (s) {
      case 'vendor':
        return UserRole.vendor;
      default:
        return UserRole.customer;
    }
  }

  String get wire {
    switch (this) {
      case UserRole.customer:
        return 'customer';
      case UserRole.vendor:
        return 'vendor';
    }
  }
}

class User {
  const User({
    required this.id,
    required this.email,
    required this.name,
    required this.role,
    required this.status,
  });

  final String id;
  final String email;
  final String name;
  final UserRole role;
  final String status;

  bool get isActive => status == 'active';
}
