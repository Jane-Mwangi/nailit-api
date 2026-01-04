CREATE TABLE IF NOT EXISTS service_types (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id uuid NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    name varchar NOT NULL,
    price integer NOT NULL,
    duration_minutes integer NOT NULL,
    image_url varchar NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    version integer NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS service_types_service_id_idx
ON service_types (service_id);
