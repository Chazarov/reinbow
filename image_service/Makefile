.PHONY: help env up updev down drop test test-upload logs swagger

COMPOSE := docker compose
export DOCKER_BUILDKIT := 1

-include .env
export

TESTS_DIR := tests
FRONTEND_DIR := web

help:
		@echo "Доступные команды: "
		@echo "  make env          — создать .env (если отсутствует)"
		@echo "  make up           — build + up docker compose"
		@echo "  make updev        — backend в docker + frontend dev-сервер"
		@echo "  make down         — остановить контейнеры"
		@echo "  make swagger      — сгенерировать Swagger (перед коммитом при изменении API)"
		@echo "  make test-upload  — тест загрузки изображения"
		@echo "  make test         — все тесты из api.sh"

env:
		@touch .env
		@grep -qE '^DB_PASSWORD=' .env || echo 'DB_PASSWORD=12345678' >> .env
		@grep -qE '^JWT_SECRET=' .env || echo 'JWT_SECRET=changethiskey!!!' >> .env

swagger:
		go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g spec.go -o api -d api,internal/app/ports/api/v1 --parseInternal --instanceName internal

up: env
		$(COMPOSE) build
		$(COMPOSE) up -d

updev: env
		$(COMPOSE) build app
		$(COMPOSE) up -d app
		cd $(FRONTEND_DIR) && npm install && npm run dev

down:
		$(COMPOSE) down

drop:
		$(COMPOSE) down -v

test-upload:
		cd $(TESTS_DIR) && bash api.sh upload

test:
		cd $(TESTS_DIR) && bash api.sh all

logs:
		$(COMPOSE) logs -f
