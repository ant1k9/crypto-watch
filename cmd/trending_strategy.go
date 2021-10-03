package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/urfave/cli"

	"github.com/ant1k9/crypto-watch/internal/pkg/db"
)

func trendStrategyCommand(d *db.DB) func(c *cli.Context) {
	return trendStrategy(d, d.GetTrendingCoins)
}

func descendStrategyCommand(d *db.DB) func(c *cli.Context) {
	return trendStrategy(d, d.GetDescendingCoins)
}

func trendStrategy(d *db.DB, getCoins func(time.Time) ([]db.Coin, error)) func(c *cli.Context) {
	return func(_ *cli.Context) {
		rand.Seed(time.Now().UnixNano())

		// two weeks ago
		offset := -time.Hour * 14 * 24
		coins, err := getCoins(time.Now().Add(offset))
		if err != nil {
			log.Fatalf("cannot get coins from db: %s", err)
			return
		}

		var profit float64
		for _, coin := range coins {
			log.Println("coin " + coin.Name + " (" + coin.Symbol + ")")

			currentRate, err := d.GetLastRate(coin.UUID, time.Now().Add(-time.Hour*0))
			if err != nil {
				log.Fatalf("cannot get current rate: %s", err)
				return
			}

			initialRate, err := d.GetLastRate(coin.UUID, time.Now().Add(offset-time.Hour))
			if err != nil {
				log.Fatalf("cannot get current rate: %s", err)
				return
			}

			if initialRate.Value == 0 {
				profit += 100.0 // back profit back because it seems to be a mistake
				log.Printf("zero price for coin %s", coin.Name)
				continue
			}

			profit += 100.0 * (currentRate.Value - initialRate.Value) / initialRate.Value
		}

		fmt.Println("initial sum: 1000$")
		fmt.Printf("final sum:   %.2f$", 1000.0+profit)
	}
}
