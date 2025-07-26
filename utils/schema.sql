CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS deleted_pods (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pod_name        TEXT NOT NULL,
    namespace       TEXT NOT NULL,
    node_name       TEXT NOT NULL,
    owner_type      TEXT,
    owner_name      TEXT,
    deleted_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    action_type     TEXT NOT NULL, 
    deletion_reason TEXT,
    status          TEXT
);
