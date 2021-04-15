package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/urfave/cli"

	"github.com/ant1k9/crypto-watch/internal/pkg/db"
)

const Attempts = 10

func getRatesCommand(client *http.Client, d *db.DB, apiKey string) func(c *cli.Context) {
	return func(_ *cli.Context) {
		coins, err := d.GetCoins()
		if err != nil {
			log.Fatalf("cannot get coins from db: %s", err)
			return
		}

		for _, coin := range coins {
			log.Printf("getting rates for %s", coin.UUID)
			for attempt := range [Attempts]struct{}{} {
				req, err := http.NewRequest(
					http.MethodGet, fmt.Sprintf(getRatesTemplate, coin.UUID), nil,
				)
				if err != nil {
					log.Fatalf("creating http request: %s", err)
					return
				}

				req.Header.Set("X-Access-Token", apiKey)
				resp, err := client.Do(req)
				if err != nil {
					log.Fatalf("getting rates: %s", err)
					return
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatalf("reading response: %s", err)
					return
				}

				var rates RatesData
				err = json.Unmarshal(body, &rates)
				if err != nil {
					log.Fatalf("unmarshalling rates: %s", err)
					return
				}

				if len(rates.Data.Rates) == 0 {
					log.Printf(
						"attempt %d unexpectedly empty rates for coin %s",
						attempt+1, coin.Name,
					)
					continue
				}

				dbRates := make([]db.Rate, 0, len(rates.Data.Rates))
				for _, rate := range rates.Data.Rates {
					price, _ := strconv.ParseFloat(rate.Price, 64)
					dbRates = append(dbRates, db.Rate{
						Value:    price,
						CoinUUID: coin.UUID,
						Ts:       time.Unix(rate.Timestamp, 0),
					})
				}

				err = d.SaveRates(dbRates)
				if err != nil {
					log.Fatalf("saving rates to db: %s", err)
					return
				}
				break
			}
		}
	}
}
