DELETE FROM permissions
WHERE code IN (
    'appointments:read',
    'appointments:write'
);