up: up
up-build: up-build
down: down

up:
	@ docker compose up 

up-build:
	@ docker compose up --build

down:
	@ docker compose down