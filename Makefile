test:
	docker compose --profile test up -d

run:
	set "DOCKER_BUILDKIT=0" && docker compose --profile app up -d --build --force-recreate