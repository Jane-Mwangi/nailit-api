ALTER TABLE services
ADD CONSTRAINT services_name_unique UNIQUE (name);
