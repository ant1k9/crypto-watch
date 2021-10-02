DB_BACKUP_PATH=/tmp/$(shell date +'%Y%m%d_%H%M%S').dump
DB_BACKUP_ARCHIVE=${DB_BACKUP_PATH}.zip

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

.PHONY: backup
backup:
	@pg_dump --data-only -d "$$DB_DSN" > "${DB_BACKUP_PATH}" &>/dev/null
	@zip "${DB_BACKUP_ARCHIVE}" "${DB_BACKUP_PATH}" &>/dev/null
	@rm "${DB_BACKUP_PATH}"
	@echo "${DB_BACKUP_ARCHIVE}"

.PHONY: load
load: get_coins get_rates
