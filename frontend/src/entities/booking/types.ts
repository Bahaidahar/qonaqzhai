export type BookingStatus =
  | "pending"
  | "accepted"
  | "declined"
  | "cancelled"
  | "completed"
  | "paid";

export interface Booking {
  id: string;
  customerId: string;
  vendorId: string;
  eventDate: string;
  guestCount: number;
  note: string;
  status: BookingStatus;
  amount: number;
  paymentId?: string;
  createdAt: string;
}
