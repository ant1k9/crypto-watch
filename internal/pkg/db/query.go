package db

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type (
	DB struct {
		inner *sqlx.DB
	}
)

func NewDb(inner *sqlx.DB) *DB {
	return &DB{inner: inner}
}

func (d *DB) SaveCoin(coin Coin) error {
	query, args, _ := sq.Insert("coins").
		Columns("uuid", "name", "symbol", "rank").
		Values(coin.UUID, coin.Name, coin.Symbol, coin.Rank).
		Suffix(
			`ON CONFLICT(uuid) DO UPDATE SET
				symbol = excluded.symbol,
				name = excluded.name,
				rank = excluded.rank`,
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err := d.inner.DB.Exec(query, args...)
	return err
}

func (d *DB) GetCoins() ([]Coin, error) {
	query, args, _ := sq.Select("id", "uuid", "name", "symbol", "rank").
		From("coins").
		OrderBy("rank NULLS LAST").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var coins []Coin
	err := d.inner.Select(&coins, query, args...)
	return coins, err
}

func (d *DB) SaveRates(rates []Rate) error {
	q := sq.Insert("rates").
		Columns("coin_uuid", "value", "ts")

	for _, rate := range rates {
		q = q.Values(rate.CoinUUID, rate.Value, rate.Ts)
	}
	query, args, _ := q.Suffix("ON CONFLICT(coin_uuid, ts) DO NOTHING").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err := d.inner.DB.Exec(query, args...)
	return err
}

func (d *DB) GetLastRate(uuid string, ts time.Time) (Rate, error) {
	query, args, _ := sq.Select("id", "coin_uuid", "value", "ts").
		From("rates").
		Where(sq.Eq{"coin_uuid": uuid}).
		Where(sq.Lt{"ts": ts}).
		OrderBy("ts desc").
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var rate Rate
	err := d.inner.Get(&rate, query, args...)
	return rate, err
}
