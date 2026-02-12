CREATE TABLE IF NOT EXISTS appointments (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    customer_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    service_id uuid NOT NULL REFERENCES services(id),
    service_type_id uuid NOT NULL REFERENCES service_types(id),

    starts_at timestamp with time zone NOT NULL,
    ends_at timestamp with time zone NOT NULL,

    status text NOT NULL CHECK (
        status IN ('pending', 'booked', 'cancelled', 'completed')
    ),

    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);
