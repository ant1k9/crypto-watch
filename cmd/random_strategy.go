package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/urfave/cli"

	"github.com/ant1k9/crypto-watch/internal/pkg/db"
)

func randomStrategyCommand(client *http.Client, d *db.DB, apiKey string) func(c *cli.Context) {
	return func(_ *cli.Context) {
		rand.Seed(time.Now().UnixNano())

		coins, err := d.GetCoins()
		if err != nil {
			log.Fatalf("cannot get coins from db: %s", err)
			return
		}

		var profit float64
		for _, lowerBorder := range []int{0, 0, 0, 10, 10, 10, 20, 20} {
			coin := coins[rand.Intn(10)+lowerBorder]
			log.Println("coin " + coin.Name + " (" + coin.Symbol + ")")

			currentRate, err := d.GetLastRate(coin.UUID, time.Now().Add(-time.Hour*0))
			if err != nil {
				log.Fatalf("cannot get current rate: %s", err)
				return
			}

			initialRate, err := d.GetLastRate(coin.UUID, time.Now().Add(-time.Hour*30*24-1))
			if err != nil {
				log.Fatalf("cannot get current rate: %s", err)
				return
			}

			if initialRate.Value == 0 {
				log.Printf("zero price for coin %s", coin.Name)
				continue
			}

			profit += 125.0 * (currentRate.Value - initialRate.Value) / initialRate.Value
		}

		fmt.Println("initial sum: 1000$")
		fmt.Printf("final sum:   %.2f$", 1000.0+profit)
	}
}
