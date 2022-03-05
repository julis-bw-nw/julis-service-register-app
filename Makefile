test:
	go install github.com/golang/mock/mockgen@v1.6.0
	go generate ./...
	go test -race -timeout 30s ./...
	
dev:
	docker-compose -p register-app -f deployments/docker-compose.dev.yml up --force-recreate --remove-orphans

run:
	JULIS_REGISTER_APP_CONFIG_PATH=config.dev.yml \
	go run -race ./cmd/register/.

test-docker:
	docker-compose -p register-app -f deployments/docker-compose.test.yml build --no-cache --force-rm
	docker-compose -p register-app -f deployments/docker-compose.test.yml up --force-recreate --remove-orphans
