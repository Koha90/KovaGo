-- +goose Up

CREATE TABLE shops (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

  name VARCHAR(160) NOT NULL,
  slug VARCHAR(120) NOT NULL,
  description TEXT NOT NULL DEFAULT '',

  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,

  CONSTRAINT shops_name_not_blank
    CHECK (length(trim(name)) > 0),

  CONSTRAINT shops_slug_not_blank
    CHECK (length(trim(slug)) > 0),

  CONSTRAINT shops_slug_unique
    UNIQUE (slug)
);

-- +goose Down

DROP TABLE shops;
