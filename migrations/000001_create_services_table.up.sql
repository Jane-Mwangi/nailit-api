CREATE TABLE IF NOT EXISTS services (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
     version integer NOT NULL DEFAULT 1
);
