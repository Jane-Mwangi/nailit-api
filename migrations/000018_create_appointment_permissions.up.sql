INSERT INTO permissions (code)
VALUES
    ('appointments:read'),
    ('appointments:write')
ON CONFLICT (code) DO NOTHING;