test:
	go install github.com/golang/mock/mockgen@v1.6.0
	go generate ./...
	go test -race -timeout 30s ./...

dev-run:
	JULIS_REGISTER_APP_CONFIG_PATH=config.dev.yml \
	go run -race ./cmd/register/.

dev-build-docker:
	docker-compose -p register-app -f deployments/docker-compose.dev.yml build --no-cache --force-rm

dev-run-docker:
	docker-compose -p register-app -f deployments/docker-compose.dev.yml up --force-recreate --remove-orphans

dev-docker:
	make dev-build-docker
	make dev-run-docker