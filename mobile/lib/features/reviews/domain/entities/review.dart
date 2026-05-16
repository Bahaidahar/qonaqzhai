class Review {
  const Review({
    required this.id,
    required this.bookingId,
    required this.customerId,
    required this.vendorId,
    required this.rating,
    required this.text,
    required this.createdAt,
  });

  final String id;
  final String bookingId;
  final String customerId;
  final String vendorId;
  final int rating;
  final String text;
  final String createdAt;
}
