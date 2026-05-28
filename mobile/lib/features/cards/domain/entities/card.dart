class PaymentCard {
  const PaymentCard({
    required this.id,
    required this.userId,
    required this.brand,
    required this.last4,
    required this.expMonth,
    required this.expYear,
    required this.holder,
    required this.isDefault,
    required this.createdAt,
  });

  final String id;
  final String userId;
  final String brand;
  final String last4;
  final int expMonth;
  final int expYear;
  final String holder;
  final bool isDefault;
  final String createdAt;

  factory PaymentCard.fromJson(Map<String, dynamic> json) => PaymentCard(
        id: json['id'] as String,
        userId: json['userId'] as String? ?? '',
        brand: json['brand'] as String? ?? 'unknown',
        last4: json['last4'] as String? ?? '',
        expMonth: json['expMonth'] as int? ?? 0,
        expYear: json['expYear'] as int? ?? 0,
        holder: json['holder'] as String? ?? '',
        isDefault: json['isDefault'] as bool? ?? false,
        createdAt: json['createdAt'] as String? ?? '',
      );
}
