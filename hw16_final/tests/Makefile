# Имя образа
IMAGE_NAME := test_service
CONTAINER_NAME := $(IMAGE_NAME)_container
NET := brutenetwork

# Запуск Docker-контейнера
doc: docker-build
	docker run -d --rm --name $(CONTAINER_NAME) $(IMAGE_NAME)
# -d: фон.
# --rm: Удаляет контейнер после останова.
# --name $(CONTAINER_NAME): Устанавливает имя контейнера.

docker-build:
	docker build -t $(IMAGE_NAME) .

# логи контейнера
log:
	docker logs $(CONTAINER_NAME)

# остановить контейнер 
stop:
	docker stop $(CONTAINER_NAME)

net:
	@if ! docker network ls | grep -q "$(NET)"; then \
		echo "Network '$(NET)' does not exist. Creating..."; \
		docker network create $(NET); \
	else \
		echo "Network '$(NET)' already exists."; \
	fi
	
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.8

lint: install-lint-deps
	golangci-lint run ./...
