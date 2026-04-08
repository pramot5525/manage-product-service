COMPOSE = docker compose

.PHONY: start down restart logs ps build test

start:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down

restart: down start

logs:
	$(COMPOSE) logs -f app

ps:
	$(COMPOSE) ps

build:
	$(COMPOSE) build --no-cache

test:
	go test ./... -cover