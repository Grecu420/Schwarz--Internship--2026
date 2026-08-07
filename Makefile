up: docker-compose.yml
	docker compose up -d --build

down:
	docker compose down


rebuild:
	docker compose up --build --no-cache -d


fresh:
	docker compose down -v
	docker compose build --no-cache
	docker compose up -d