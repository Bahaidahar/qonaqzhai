DROP INDEX IF EXISTS idx_bookings_service;
ALTER TABLE bookings DROP COLUMN service_id;
DROP INDEX IF EXISTS idx_services_active;
DROP INDEX IF EXISTS idx_services_vendor;
DROP TABLE IF EXISTS services;
