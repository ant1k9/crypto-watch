DB_BACKUP_PATH=/tmp/$(shell date +'%Y%m%d_%H%M%S').dump
DB_BACKUP_ARCHIVE=${DB_BACKUP_PATH}.zip

STRATEGIES = random trend descend
STRATEGIES_COMMANDS = $(addsuffix _strategy, ${STRATEGIES})

print:
	@echo ${STRATEGIES_COMMANDS}

get_coins: build
	./bin/crypto-watch get_coins

get_rates: build
	./bin/crypto-watch get_rates

$(STRATEGIES_COMMANDS):
	./bin/crypto-watch $@

build:
	go build -o ./bin/crypto-watch cmd/*.go

migrate:
	goose -dir ./migrations postgres $$DB_DSN up

.PHONY: backup
backup:
	@pg_dump --data-only -d "$$DB_DSN" > "${DB_BACKUP_PATH}" 2>/dev/null
	@zip "${DB_BACKUP_ARCHIVE}" "${DB_BACKUP_PATH}" &>/dev/null
	@rm "${DB_BACKUP_PATH}"
	@echo "${DB_BACKUP_ARCHIVE}"

.PHONY: load
load: get_coins get_rates
