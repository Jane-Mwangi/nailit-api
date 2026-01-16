-- Remove incorrect permissions added earlier
DELETE FROM permissions
WHERE code LIKE 'movies:%';

-- Insert correct Nailit permissions
INSERT INTO permissions (code)
VALUES
    ('services:read'),
    ('services:write'),
    ('service-types:read'),
    ('service-types:write'),
    ('staff:read'),
    ('staff:write')
ON CONFLICT (code) DO NOTHING;
