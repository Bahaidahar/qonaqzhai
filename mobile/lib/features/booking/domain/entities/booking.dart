class Booking {
  const Booking({
    required this.id,
    required this.customerId,
    required this.vendorId,
    required this.eventDate,
    required this.guestCount,
    required this.note,
    required this.status,
    required this.amount,
    required this.paymentId,
  });

  final String id;
  final String customerId;
  final String vendorId;
  final String eventDate;
  final int guestCount;
  final String note;
  final String status;
  final int amount;
  final String paymentId;
}
