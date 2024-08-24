up-windows: docker-up-window
up-linux: docker-up-linux

docker-up-window:
	@ docker compose up --build

docker-up-linux:
	@ sudo docker compose up --build