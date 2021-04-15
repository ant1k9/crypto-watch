package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli"

	"github.com/ant1k9/crypto-watch/internal/pkg/db"
)

const (
	getRatesTemplate = "https://api.coinranking.com/v2/coin/%s/history?timePeriod=3m"
	getCoinsTemplate = "https://api.coinranking.com/v2/coins?limit=60"
)

type (
	CoinsData struct {
		Data struct {
			Coins []db.Coin `json:"coins"`
		} `json:"data"`
	}

	RatesData struct {
		Data struct {
			Rates []struct {
				Price     string `json:"price"`
				Timestamp int64  `json:"timestamp"`
			} `json:"history"`
		} `json:"data"`
	}
)

var app = cli.NewApp()

func commands(d *db.DB, apiKey string) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	app.Commands = []cli.Command{
		{
			Name:   "get_coins",
			Usage:  "Updates coins in database",
			Action: getCoinsCommand(client, d, apiKey),
		},
		{
			Name:   "get_rates",
			Usage:  "Update rates in database",
			Action: getRatesCommand(client, d, apiKey),
		},
		{
			Name:   "random_strategy",
			Usage:  "Update rates in database",
			Action: randomStrategyCommand(client, d, apiKey),
		},
	}
}

func main() {
	d := sqlx.MustOpen("postgres", os.Getenv("DB_DSN"))
	commands(db.NewDb(d), os.Getenv("API_KEY"))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
