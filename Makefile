test:
	go install github.com/golang/mock/mockgen@v1.6.0
	go generate ./...
	go test -race -timeout 30s ./...

dev-build-docker:
	docker-compose -p register-app -f deployments/docker-compose.dev.yml build --no-cache --force-rm

dev-run-docker:
	docker-compose -p register-app -f deployments/docker-compose.dev.yml up --force-recreate --remove-orphans

dev-docker:
	make dev-build-docker
	make dev-run-docker