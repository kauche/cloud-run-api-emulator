CREATE TABLE IF NOT EXISTS services (
    parent TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    uid TEXT NOT NULL,
    generation BIGINT NOT NULL,
    uri TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (parent, name)
);

CREATE INDEX IF NOT EXISTS created_at_desc ON services(created_at DESC);

CREATE TABLE IF NOT EXISTS service_labels (
    service_parent TEXT NOT NULL,
    service_name TEXT NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY (service_parent, service_name, key)
);

CREATE TABLE IF NOT EXISTS service_annotations (
    service_parent TEXT NOT NULL,
    service_name TEXT NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY (service_parent, service_name, key)
);
