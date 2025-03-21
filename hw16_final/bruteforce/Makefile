# Имя образа
IMAGE_NAME := bruteforce_service
CONTAINER_NAME := $(IMAGE_NAME)_container
NET := brutenetwork

# Запуск Docker-контейнера
doc: docker-build
	docker run -d --rm -p 3009:3009 --name $(CONTAINER_NAME) $(IMAGE_NAME)
# -d: Запускает контейнер в фоновом режиме.
# --rm: Удаляет контейнер автоматически после его остановки.
# -p 3005:3005: Перебрасывает порт 3005 из контейнера на хост.
# --name $(CONTAINER_NAME): Устанавливает имя контейнера для удобства управления.
# --network host - для запроса из Постмана поменять сеть!

# Сборка бинарного файла
build:
	go build -o ./bruteforce main.go

# Запуск приложения локально
run: build
	@go run bruteforce

# Удаление бинарного файла
del:
	rm bruteforce

# Сборка Docker-образа
docker-build:
	docker build -t $(IMAGE_NAME) .


# Запущенные контейнеры
ps: 
	docker ps

# образы
im:
	docker images

# логи контейнера
log:
	docker logs $(CONTAINER_NAME)

# зайти в контейнер
in:
	docker exec -it $(CONTAINER_NAME) /bin/sh

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

test:
	go test -race -count 100 ./...

# Дерево каталогов в буфер обмена
tree:
	tree | xclip -selection clipboard

#  Проверки:
# curl --location 'http://127.0.0.1:3009/params/all?add=0'
# curl -I http://127.0.0.1:3009
# nc -zv 127.0.0.1 3009
# netstat -tuln | grep 3009

# ps aux | grep bruteforce_service
# pgrep -fl bruteforce_service