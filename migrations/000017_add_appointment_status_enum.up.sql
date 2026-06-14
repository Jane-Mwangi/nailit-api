
UPDATE appointments
SET status = 'booked'
WHERE status = 'pending';


CREATE TYPE appointment_status AS ENUM (
  'booked',
  'cancelled',
  'completed'
);


ALTER TABLE appointments
ALTER COLUMN status TYPE appointment_status
USING status::appointment_status;