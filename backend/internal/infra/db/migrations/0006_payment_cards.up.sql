CREATE TABLE IF NOT EXISTS payment_cards (
  id          TEXT PRIMARY KEY,
  user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  brand       TEXT NOT NULL DEFAULT 'unknown',
  last4       TEXT NOT NULL,
  exp_month   INTEGER NOT NULL,
  exp_year    INTEGER NOT NULL,
  holder      TEXT NOT NULL DEFAULT '',
  is_default  BOOLEAN NOT NULL DEFAULT FALSE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_cards_user ON payment_cards(user_id);
