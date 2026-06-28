# make push
push:
	@bash push.sh

# make start [ `docker-compose` ]
docker-start:
	# 后台运行
	@docker compose up -d

# make remove  [ `docker-compose` ]
docker-remove:
	# 删除当前所有容器
	@docker compose down
	@docker image rm gin-admin


# make restart  [ run by `docker-compose` ]
docker-restart:
	# 重启当前容器
	@docker restart gin-admin mysql redis


# make stop [ stop by `docker-compose` ]
docker-stop:
	@docker stop gin-admin mysql redis
