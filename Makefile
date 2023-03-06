# ===================================================== MongoDB ===================================================== #
MONGO_CONTAINER=sportsnews-mongo
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_ADMIN=admin
MONGO_SECRET=secret
MONGO_DB=sportsnews

.PHONY: run-mongo start-mongo stop-mongo info-mongo remove-mongo
run-mongo:
	@echo "Starting mongodb container $(MONGO_CONTAINER)"
	@docker run --name $(MONGO_CONTAINER) \
			-p $(MONGO_PORT):27017 \
			-e MONGO_INITDB_ROOT_USERNAME=$(MONGO_ADMIN) \
			-e MONGO_INITDB_ROOT_PASSWORD=$(MONGO_SECRET) \
			-e MONGO_INITDB_DATABASE=$(MONGO_DB) \
			-d mongo

start-mongo:
	@echo "Starting mongo container $(MONGO_CONTAINER)"
	@docker start $(MONGO_CONTAINER)

stop-mongo:
	@echo "Stopping mongo container $(MONGO_CONTAINER)"
	@docker stop $(MONGO_CONTAINER)

remove-mongo: stop-mongo
	@echo "Removing mongo container $(MONGO_CONTAINER)"
	@docker rm $(MONGO_CONTAINER)

info-mongo:
	@echo "Get logs from mongo container $(MONGO_CONTAINER)"
	@docker logs -f $(MONGO_CONTAINER)

# ===================================================== MongoDB ===================================================== #

# ====================================================== Redis ====================================================== #
REDIS_CONTAINER=sportsnews-redis
REDIS_PORT=6379

.PHONY: run-redis stop-redis info-redis remove-redis
run-redis:
	@echo "Starting redis container $(REDIS_CONTAINER)"
	@docker run --name $(REDIS_CONTAINER) -p $(REDIS_PORT):6379 -d redis

start-redis:
	@echo "Starting redis container $(REDIS_CONTAINER)"
	@docker start $(REDIS_CONTAINER)

stop-redis:
	@echo "Stopping redis container $(REDIS_CONTAINER)"
	@docker stop $(REDIS_CONTAINER)

remove-redis: stop-redis
	@echo "Removing redis container $(REDIS_CONTAINER)"
	@docker rm $(REDIS_CONTAINER)

info-redis:
	@echo "Get logs from redis container $(REDIS_CONTAINER)"
	@docker logs -f $(REDIS_CONTAINER)

# Redis commander
REDIS_COMMANDER_CONTAINER=sportsnews-redis-commander
REDIS_COMMANDER_PORT=1235

.PHONY: get-redis-container-ip run-redis-commander stop-redis-commander remove-redis-commander
get-redis-container-ip:
	@echo $(eval REDIS_HOST=$(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(REDIS_CONTAINER)))

run-redis-commander: get-redis-container-ip
	@echo "Starting redis commander $(REDIS_HOST)"
	@docker run --name $(REDIS_COMMANDER_CONTAINER) \
			--env REDIS_HOSTS=local:$(REDIS_HOST) \
			-p $(REDIS_COMMANDER_PORT):8081 \
			-d rediscommander/redis-commander:latest

stop-redis-commander:
	@echo "Stopping redis container $(REDIS_COMMANDER_CONTAINER)"
	@docker stop $(REDIS_COMMANDER_CONTAINER)


remove-redis-commander: stop-redis-commander
	@echo "Removing redis-commander container $(REDIS_COMMANDER_CONTAINER)"
	@docker rm $(REDIS_COMMANDER_CONTAINER)
# ====================================================== Redis ====================================================== #

# ====================================================== Utils ====================================================== #
.PHONY: run-all stop-all clean test test-cover mock-broker mock-cache mock-broker mock-broker
run-all: run-mongo run-redis run-redis-commander

stop-all: stop-mongo stop-redis stop-redis-commander

clean: remove-mongo remove-redis remove-redis-commander

mock-usecase:
	  mockgen -source=internal/article/usecase.go -destination internal/article/mock/mock_usecase.go

mock-repository:
	  mockgen -source=internal/article/repository.go -destination internal/article/mock/mock_repository.go

mock-all: mock-usecase mock-repository

swagger:
	@echo "Generate swagger doc"
	@swag init -g **/**/**/*.go

lint:
	@echo "Run golangci-lint"
	@golangci-lint run -c .golangci.yml

test:
	@echo "Running tests"
	@go test ./...

test-cover:
	@echo "Running tests with cover"
	@go test -coverprofile="coverage.txt" -covermode=atomic -p 1 ./...
# ====================================================== Utils ====================================================== #