push:
	@bash push.sh

# make start
docker-start:
	@docker compose up -d

# 只重建后端代码，不改mysql,redis
docker-rebuild:
	@docker compose up -d --build backend

# make remove
docker-remove:
	@docker compose down
	@docker image rm go-alpha:v1

# make restart
docker-restart:
	@docker restart go-alpha mysql redis

# make stop
docker-stop:
	@docker stop go-alpha mysql redis
