class BookingThread {
  const BookingThread({
    required this.id,
    required this.bookingId,
    required this.customerId,
    required this.vendorId,
    required this.createdAt,
    required this.updatedAt,
  });

  final String id;
  final String bookingId;
  final String customerId;
  final String vendorId;
  final String createdAt;
  final String updatedAt;
}

class ThreadMessage {
  const ThreadMessage({
    required this.id,
    required this.threadId,
    required this.senderId,
    required this.text,
    required this.createdAt,
  });

  final String id;
  final String threadId;
  final String senderId;
  final String text;
  final String createdAt;
}
