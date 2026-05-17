CREATE TABLE IF NOT EXISTS booking_threads (
  id          TEXT PRIMARY KEY,
  booking_id  TEXT NOT NULL UNIQUE REFERENCES bookings(id) ON DELETE CASCADE,
  customer_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  vendor_id   TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_threads_customer ON booking_threads(customer_id, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_threads_vendor ON booking_threads(vendor_id, updated_at DESC);

CREATE TABLE IF NOT EXISTS thread_messages (
  id         TEXT PRIMARY KEY,
  thread_id  TEXT NOT NULL REFERENCES booking_threads(id) ON DELETE CASCADE,
  sender_id  TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  text       TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_thread_messages_thread ON thread_messages(thread_id, created_at);
