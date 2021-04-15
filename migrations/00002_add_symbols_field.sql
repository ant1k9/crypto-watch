-- +goose Up
ALTER TABLE coins
    ADD COLUMN symbol text NOT NULL default '';

-- +goose Down
ALTER TABLE coins
    DROP COLUMN symbol;
