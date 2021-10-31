-- +goose Up
CREATE TABLE water (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR (100) NOT NULL,
  model VARCHAR (100) NOT NULL,
  manufacturer VARCHAR (100) NOT NULL,
  material VARCHAR (100) NOT NULL,
  speed INTEGER NOT NULL
);

-- +goose Down
DROP TABLE water;
