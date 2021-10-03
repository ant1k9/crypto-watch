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

func (d *DB) GetTrendingCoins(ts time.Time) ([]Coin, error) {
	return d.getTrendingCoins(ts, 10, "cs2.value / cs1.value DESC")
}

func (d *DB) GetDescendingCoins(ts time.Time) ([]Coin, error) {
	return d.getTrendingCoins(ts, 30, "cs1.value / cs2.value DESC")
}

func (d *DB) getTrendingCoins(ts time.Time, window int, orderBy string) ([]Coin, error) {
	query, args, _ := sq.Select("coins.name", "coins.uuid", "coins.id", "coins.symbol").
		From("coin_stats cs1").
		InnerJoin("coin_stats cs2 USING(coin_uuid)").
		InnerJoin("coins ON cs1.coin_uuid = coins.uuid").
		Where(sq.Eq{"cs1.rn": 1, "cs2.rn": window}).
		Where(sq.NotEq{"cs1.value": 0, "cs2.value": 0}).
		Where(sq.LtOrEq{"coins.rank": 40}).
		OrderBy(orderBy).
		Limit(10).
		Prefix(
			`WITH coin_stats AS (
		       SELECT
				   coin_uuid, value, 
		           ROW_NUMBER() OVER (PARTITION BY coin_uuid ORDER BY ts DESC) rn
		       FROM rates
			   WHERE ts <= '` + ts.Format("2006-01-02 15:04:05") + `'
		   )
		   `).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var coins []Coin
	err := d.inner.Select(&coins, query, args...)
	return coins, err
}
