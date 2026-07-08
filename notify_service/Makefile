include ./config/data.env
export DATABASE_CONNECT

run:
	docker compose up -d --build

up:
	migrate -source file://./Backend/migrations -database $$DATABASE_CONNECT up

up-1:
	migrate -source file://./Backend/migrations -database $$DATABASE_CONNECT up 1

down:
	migrate -source file://./Backend/migrations -database $$DATABASE_CONNECT down

down-1:
	migrate -source file://./Backend/migrations -database $$DATABASE_CONNECT down 1