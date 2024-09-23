up: up
up-build: up-build
down: down
test: test
generate-doc: generate-doc

up:
	@ docker compose up 

up-build:
	@ docker compose up --build

down:
	@ docker compose down

test:
	@ find . -type d -name 'test*' -exec go test {}/... \;

generate-doc:
	@ swag init -g ./cmd/main.go 