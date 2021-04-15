package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/urfave/cli"

	"github.com/ant1k9/crypto-watch/internal/pkg/db"
)

func getCoinsCommand(client *http.Client, d *db.DB, apiKey string) func(c *cli.Context) {
	return func(_ *cli.Context) {
		req, err := http.NewRequest(http.MethodGet, getCoinsTemplate, nil)
		if err != nil {
			log.Fatalf("creating http request: %s", err)
			return
		}

		req.Header.Set("X-Access-Token", apiKey)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("getting coins: %s", err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("reading response: %s", err)
			return
		}

		var coinsData CoinsData
		err = json.Unmarshal(body, &coinsData)
		if err != nil {
			log.Fatalf("unmarshalling coins: %s", err)
			return
		}

		for _, coin := range coinsData.Data.Coins {
			err = d.SaveCoin(coin)
			if err != nil {
				log.Fatalf("save coin: %s", err)
				return
			}
		}
	}
}
