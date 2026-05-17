export interface BookingThread {
  id: string;
  bookingId: string;
  customerId: string;
  vendorId: string;
  createdAt: string;
  updatedAt: string;
}

export interface ThreadSummary {
  thread: BookingThread;
  bookingId: string;
  eventDate: string;
  guestCount: number;
  amount: number;
  status: string;
  vendorName: string;
  counterpart: string;
}

export interface ThreadMessage {
  id: string;
  threadId: string;
  senderId: string;
  text: string;
  createdAt: string;
}
