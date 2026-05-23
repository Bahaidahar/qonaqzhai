CREATE TABLE IF NOT EXISTS threads (
  id          TEXT PRIMARY KEY,
  booking_id  TEXT NOT NULL UNIQUE,
  customer_id TEXT NOT NULL,
  vendor_id   TEXT NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_threads_customer ON threads(customer_id);
CREATE INDEX IF NOT EXISTS idx_threads_vendor   ON threads(vendor_id);

CREATE TABLE IF NOT EXISTS thread_messages (
  id         TEXT PRIMARY KEY,
  thread_id  TEXT NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
  sender_id  TEXT NOT NULL,
  text       TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_thread_messages_thread ON thread_messages(thread_id, created_at);
