-- +goose Up
CREATE TABLE IF NOT EXISTS water (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR (100) NOT NULL,
    model VARCHAR (100) NOT NULL,
    manufacturer VARCHAR (100) NOT NULL,
    material VARCHAR (100) NOT NULL,
    speed INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TYPE water_events_type AS ENUM ('created', 'updated', 'removed');
CREATE TYPE water_events_status AS ENUM ('lock', 'unlock');
CREATE TABLE IF NOT EXISTS water_events (
    id BIGSERIAL PRIMARY KEY,
    water_id BIGINT NOT NULL,
    type water_events_type NOT NULL,
    status water_events_status NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);

CREATE INDEX idx_water_id ON water_events(water_id);
CREATE INDEX idx_status ON water_events(status);

-- +goose Down
DROP TABLE water, water_events;
DROP TYPE water_events_type, water_events_status;