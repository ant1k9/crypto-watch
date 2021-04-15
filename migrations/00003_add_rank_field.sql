-- +goose Up
ALTER TABLE coins
    ADD COLUMN rank INT DEFAULT 1000000;

-- +goose Down
ALTER TABLE coins
    DROP COLUMN rank;
