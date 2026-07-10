.PHONY: help env up down

DEP_DIR := dep
IMAGE_ENV := image_service/.env
COMPOSE := docker compose

-include $(DEP_DIR)/.env
export

help:
	@echo "Команды:"
	@echo "  make help  — эта справка"
	@echo "  make env   — создать и дополнить .env файлы"
	@echo "  make up    — запустить docker compose"
	@echo "  make down  — остановить docker compose"

env:
	@test -f $(DEP_DIR)/.env || cp $(DEP_DIR)/.env.example $(DEP_DIR)/.env
	@grep -v '^#' $(DEP_DIR)/.env.example | grep '=' | while IFS= read -r line; do \
		key=$${line%%=*}; \
		grep -qE "^$$key=" $(DEP_DIR)/.env || echo "$$line" >> $(DEP_DIR)/.env; \
	done
	@touch $(IMAGE_ENV)
	@grep -qE '^DB_PASSWORD=' $(IMAGE_ENV) || echo 'DB_PASSWORD=12345678' >> $(IMAGE_ENV)
	@grep -qE '^JWT_SECRET=' $(IMAGE_ENV) || echo 'JWT_SECRET=changethiskey!!!' >> $(IMAGE_ENV)
	@echo "Готово: $(DEP_DIR)/.env, $(IMAGE_ENV)"

up: env
	cd $(DEP_DIR) && $(COMPOSE) up --build -d

down:
	cd $(DEP_DIR) && $(COMPOSE) down
