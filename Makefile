get_coins:
	./bin/crypto-watch get_coins

get_rates:
	./bin/crypto-watch get_rates

random_strategy:
	./bin/crypto-watch random_strategy

build:
	go build -o ./bin/crypto-watch cmd/*.go

migrate:
	goose -dir ./migrations postgres $$DB_DSN up
