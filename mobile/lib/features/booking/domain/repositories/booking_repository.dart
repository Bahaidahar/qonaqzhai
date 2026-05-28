import '../entities/booking.dart';

abstract class BookingRepository {
  Future<List<Booking>> myBookings();
  Future<Booking> create({
    required String vendorId,
    required String eventDate,
    required int guestCount,
    String? note,
    int? amount,
    String? serviceId,
  });
  Future<Booking> updateStatus(String id, String status);
  Future<String> startPayment(String bookingId);
}
