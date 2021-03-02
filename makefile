SHELL := /bin/bash

# ==============================================================================
# Building containers

all: quote

quote:
	docker build \
		-f dockerfile.quote-api \
		-t quote-api-amd64:1.0 \
		.

# ==============================================================================
# Running from within docker compose

up:
	docker-compose -f docker-compose.yaml up --detach --remove-orphans

down:
	docker-compose -f docker-compose.yaml down --remove-orphans

logs:
	docker-compose -f docker-compose.yaml logs -f

migrate:
	docker exec -it quote-api /service/admin migrate

seed:
	docker exec -it quote-api /service/admin seed

