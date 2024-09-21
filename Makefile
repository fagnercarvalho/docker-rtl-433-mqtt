test:
	docker compose --profile test up -d

run:
	docker compose --profile app up -d --build --force-recreate