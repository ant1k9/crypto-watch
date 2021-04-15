package db

import "time"

type (
	Coin struct {
		ID     uint64 `db:"id" json:"id"`
		Name   string `db:"name" json:"name"`
		UUID   string `db:"uuid" json:"uuid"`
		Symbol string `json:"symbol" db:"symbol"`
		Rank   uint64 `json:"rank" db:"rank"`
	}

	Rate struct {
		ID       uint64    `db:"id" json:"id"`
		CoinUUID string    `db:"coin_uuid" json:"-"`
		Value    float64   `db:"value" json:"price"`
		Ts       time.Time `db:"ts" json:"ts"`
	}
)
