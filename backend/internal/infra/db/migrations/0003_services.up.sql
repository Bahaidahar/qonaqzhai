CREATE TABLE IF NOT EXISTS services (
  id          TEXT PRIMARY KEY,
  vendor_id   TEXT NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
  name        TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  price       BIGINT NOT NULL DEFAULT 0,
  unit        TEXT NOT NULL DEFAULT 'fixed'
              CHECK (unit IN ('fixed','hour','item','person','day')),
  is_active   BOOLEAN NOT NULL DEFAULT TRUE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_services_vendor ON services(vendor_id);
CREATE INDEX IF NOT EXISTS idx_services_active ON services(is_active);

-- Booking now optionally references a single service for transparency.
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS service_id TEXT NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_bookings_service ON bookings(service_id);
