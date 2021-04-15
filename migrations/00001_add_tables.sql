-- +goose Up
CREATE TABLE coins (
    id   SERIAL PRIMARY KEY NOT NULL,
    uuid TEXT   UNIQUE      NOT NULL,
    name TEXT               NOT NULL
);

CREATE TABLE rates (
    id          SERIAL  PRIMARY KEY NOT NULL,
    coin_uuid   TEXT    REFERENCES coins(uuid),
    value       NUMERIC(10, 2)      NOT NULL,
    ts          TIMESTAMP           NOT NULL,
    UNIQUE(coin_uuid, ts)
);

-- +goose Down
DROP TABLE coins;
DROP TABLE rates;
