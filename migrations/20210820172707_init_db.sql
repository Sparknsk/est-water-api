-- +goose Up
CREATE TABLE IF NOT EXISTS water (
     id BIGSERIAL PRIMARY KEY,
     name VARCHAR (100) NOT NULL,
     model VARCHAR (100) NOT NULL,
     manufacturer VARCHAR (100) NOT NULL,
     material VARCHAR (100) NOT NULL,
     speed INTEGER NOT NULL,
     created_at TIMESTAMP NOT NULL,
     delete_status BOOL NOT NULL DEFAULT FALSE
) PARTITION BY RANGE (id);

CREATE TABLE IF NOT EXISTS water_p1 PARTITION OF water FOR VALUES FROM (1) TO (6);
CREATE TABLE IF NOT EXISTS water_p2 PARTITION OF water FOR VALUES FROM (6) TO (11);
CREATE TABLE IF NOT EXISTS water_p3 PARTITION OF water FOR VALUES FROM (11) TO (21);
CREATE TABLE IF NOT EXISTS water_p4 PARTITION OF water FOR VALUES FROM (21) TO (31);
CREATE TABLE IF NOT EXISTS water_p5 PARTITION OF water FOR VALUES FROM (31) TO (MAXVALUE);

CREATE INDEX idx_delete_status ON water(delete_status);

CREATE TYPE water_events_type AS ENUM ('created', 'updated', 'removed');
CREATE TYPE water_events_status AS ENUM ('lock', 'unlock');
CREATE TABLE IF NOT EXISTS water_events (
    id BIGSERIAL PRIMARY KEY,
    water_id BIGINT NOT NULL,
    type water_events_type NOT NULL,
    status water_events_status NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    FOREIGN KEY (water_id) REFERENCES water (id) ON DELETE CASCADE
) PARTITION BY RANGE (id);

CREATE TABLE IF NOT EXISTS water_events_p1 PARTITION OF water_events FOR VALUES FROM (1) TO (11);
CREATE TABLE IF NOT EXISTS water_events_p2 PARTITION OF water_events FOR VALUES FROM (11) TO (21);
CREATE TABLE IF NOT EXISTS water_events_p3 PARTITION OF water_events FOR VALUES FROM (21) TO (31);
CREATE TABLE IF NOT EXISTS water_events_p4 PARTITION OF water_events FOR VALUES FROM (31) TO (41);
CREATE TABLE IF NOT EXISTS water_events_p5 PARTITION OF water_events FOR VALUES FROM (41) TO (MAXVALUE);

CREATE INDEX idx_status ON water_events(status);

-- +goose Down
DROP TABLE water, water_events;
DROP TYPE water_events_type, water_events_status;