CREATE TABLE IF NOT EXISTS cards (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL,
  brand      TEXT NOT NULL,
  last4      TEXT NOT NULL,
  exp_month  INTEGER NOT NULL,
  exp_year   INTEGER NOT NULL,
  holder     TEXT NOT NULL DEFAULT '',
  is_default BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_cards_user ON cards(user_id);

CREATE TABLE IF NOT EXISTS payments (
  id            TEXT PRIMARY KEY,
  booking_id    TEXT NOT NULL UNIQUE,
  user_id       TEXT NOT NULL,
  card_id       TEXT NOT NULL,
  amount        BIGINT NOT NULL,
  currency      TEXT NOT NULL DEFAULT 'KZT',
  status        TEXT NOT NULL CHECK (status IN ('pending','captured','failed','refunded')),
  provider_ref  TEXT NOT NULL DEFAULT '',
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_payments_user    ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_booking ON payments(booking_id);
