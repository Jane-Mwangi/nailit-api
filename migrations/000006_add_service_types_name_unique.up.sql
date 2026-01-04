ALTER TABLE service_types
ADD CONSTRAINT service_types_service_name_unique
UNIQUE (service_id, name);
