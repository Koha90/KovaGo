-- +goose Up

ALTER TABLE shops
  ADD COLUMN status VARCHAR(16) NOT NULL DEFAULT 'active';

ALTER TABLE shops
  ADD CONSTRAINT shops_status_valid
  CHECK (status IN ('active', 'inactive'));

-- +goose Down

ALTER TABLE shops
  DROP CONSTRAINT shops_status_valid;

ALTER TABLE shops
  DROP COLUMN status;
