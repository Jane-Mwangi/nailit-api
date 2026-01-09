CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at timestamptz NOT NULL DEFAULT now(),
  email citext UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  role text NOT NULL CHECK (role IN ('admin', 'customer')),
  activated bool NOT NULL DEFAULT true,
  version integer NOT NULL DEFAULT 1
);
