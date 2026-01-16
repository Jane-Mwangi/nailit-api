DELETE FROM permissions
WHERE code IN (
    'services:read',
    'services:write',
    'service-types:read',
    'service-types:write',
    'staff:read',
    'staff:write'
);
