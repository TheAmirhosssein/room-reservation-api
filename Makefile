up: up
up-build: up-build
down: down
test: test

up:
	@ docker compose up 

up-build:
	@ docker compose up --build

down:
	@ docker compose down

test:
	@ go test ./...