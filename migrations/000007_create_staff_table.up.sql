CREATE TABLE staff (
    id uuid PRIMARY KEY,
    name varchar NOT NULL,
    email varchar NOT NULL UNIQUE,
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamp NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX staff_email_unique
ON staff (email)
WHERE email IS NOT NULL;
