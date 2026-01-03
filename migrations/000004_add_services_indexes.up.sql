CREATE INDEX IF NOT EXISTS services_name_fts_idx
ON services
USING GIN (to_tsvector('simple', name));
