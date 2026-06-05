-- Runs once on first container start (role qonaqzhai is created by
-- POSTGRES_USER). Each service migrates its own schema on boot, so we only
-- need the four empty databases here — one per service, no cross-DB joins.
CREATE DATABASE qonaqzhai_auth     OWNER qonaqzhai;
CREATE DATABASE qonaqzhai_core     OWNER qonaqzhai;
CREATE DATABASE qonaqzhai_payment  OWNER qonaqzhai;
CREATE DATABASE qonaqzhai_realtime OWNER qonaqzhai;
