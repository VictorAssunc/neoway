run:
	@docker compose -f deploy/docker-compose.yaml up -d --force-recreate --build
	@docker compose -f deploy/docker-compose.yaml logs app -f

stop:
	@docker compose -f deploy/docker-compose.yaml down -v

test:
	@go test -covermode=atomic -cover -race ./...
