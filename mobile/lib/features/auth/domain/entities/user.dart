enum UserRole { customer, vendor, admin }

extension UserRoleX on UserRole {
  static UserRole parse(String s) {
    switch (s) {
      case 'vendor':
        return UserRole.vendor;
      case 'admin':
        return UserRole.admin;
      default:
        return UserRole.customer;
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
