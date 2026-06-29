# shell: make push
push:
	@bash push.sh

# go-alpha, mysql, redis
docker-build:
	@docker compose up -d

# 只重建后端代码，不重建mysql, redis
docker-rebuild:
	@docker compose up -d --build backend
