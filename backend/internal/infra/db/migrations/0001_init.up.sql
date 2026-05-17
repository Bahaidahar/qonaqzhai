CREATE TABLE IF NOT EXISTS users (
  id            TEXT PRIMARY KEY,
  email         TEXT NOT NULL UNIQUE,
  name          TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  role          TEXT NOT NULL CHECK (role IN ('customer','vendor','admin')),
  status        TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','suspended')),
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS vendors (
  id           TEXT PRIMARY KEY,
  user_id      TEXT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  name         TEXT NOT NULL,
  category     TEXT NOT NULL,
  city         TEXT NOT NULL,
  description  TEXT NOT NULL DEFAULT '',
  price_from   BIGINT NOT NULL DEFAULT 0,
  status       TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','approved','rejected')),
  rating_avg   DOUBLE PRECISION NOT NULL DEFAULT 0,
  rating_count INTEGER NOT NULL DEFAULT 0,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_vendors_category ON vendors(category);
CREATE INDEX IF NOT EXISTS idx_vendors_city     ON vendors(city);
CREATE INDEX IF NOT EXISTS idx_vendors_status   ON vendors(status);
CREATE INDEX IF NOT EXISTS idx_vendors_price    ON vendors(price_from);
CREATE INDEX IF NOT EXISTS idx_vendors_rating   ON vendors(rating_avg);

CREATE TABLE IF NOT EXISTS photos (
  id         TEXT PRIMARY KEY,
  vendor_id  TEXT NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
  mime       TEXT NOT NULL,
  size       BIGINT NOT NULL,
  data       BYTEA NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_photos_vendor ON photos(vendor_id);

CREATE TABLE IF NOT EXISTS bookings (
  id           TEXT PRIMARY KEY,
  customer_id  TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  vendor_id    TEXT NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
  event_date   TEXT NOT NULL,
  guest_count  INTEGER NOT NULL DEFAULT 0,
  note         TEXT NOT NULL DEFAULT '',
  status       TEXT NOT NULL DEFAULT 'pending'
               CHECK (status IN ('pending','accepted','declined','cancelled','completed','paid')),
  amount       BIGINT NOT NULL DEFAULT 0,
  payment_id   TEXT NOT NULL DEFAULT '',
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_bookings_customer ON bookings(customer_id);
CREATE INDEX IF NOT EXISTS idx_bookings_vendor   ON bookings(vendor_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status   ON bookings(status);

CREATE TABLE IF NOT EXISTS reviews (
  id          TEXT PRIMARY KEY,
  booking_id  TEXT NOT NULL UNIQUE REFERENCES bookings(id) ON DELETE CASCADE,
  customer_id TEXT NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
  vendor_id   TEXT NOT NULL REFERENCES vendors(id)  ON DELETE CASCADE,
  rating      INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
  text        TEXT NOT NULL DEFAULT '',
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_reviews_vendor ON reviews(vendor_id);

CREATE TABLE IF NOT EXISTS refresh_tokens (
  id          TEXT PRIMARY KEY,
  user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash  TEXT NOT NULL UNIQUE,
  expires_at  TIMESTAMPTZ NOT NULL,
  revoked_at  TIMESTAMPTZ,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_refresh_user ON refresh_tokens(user_id);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
  id          TEXT PRIMARY KEY,
  user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash  TEXT NOT NULL UNIQUE,
  expires_at  TIMESTAMPTZ NOT NULL,
  used_at     TIMESTAMPTZ,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_password_reset_user ON password_reset_tokens(user_id);

CREATE TABLE IF NOT EXISTS notifications (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type       TEXT NOT NULL,
  channel    TEXT NOT NULL,
  title      TEXT NOT NULL,
  body       TEXT NOT NULL,
  status     TEXT NOT NULL DEFAULT 'queued',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);

CREATE TABLE IF NOT EXISTS fcm_tokens (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token      TEXT NOT NULL UNIQUE,
  platform   TEXT NOT NULL DEFAULT 'unknown',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_fcm_user ON fcm_tokens(user_id);

-- Full-text search on vendors via tsvector generated column + GIN index.
ALTER TABLE vendors
  ADD COLUMN IF NOT EXISTS search_tsv tsvector
  GENERATED ALWAYS AS (
    setweight(to_tsvector('simple', coalesce(name, '')), 'A') ||
    setweight(to_tsvector('simple', coalesce(description, '')), 'B')
  ) STORED;

CREATE INDEX IF NOT EXISTS idx_vendors_search ON vendors USING GIN (search_tsv);
