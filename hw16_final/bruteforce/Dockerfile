# Environment
FROM golang:1.21.7 AS build-env

RUN mkdir -p /opt/bruteforce_service
WORKDIR /opt/bruteforce_service
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /opt/service/bruteforce_service

# Release
FROM alpine:latest

RUN apk add --no-cache curl

ARG PTH="/bin/"
# Устанавливаем переменную окружения DOCKER_CONF_PATH
ENV DOCKER_CONF_PATH=${PTH}

# Копируем бинарный файл и config.yaml в образ
COPY --from=build-env /opt/bruteforce_service/config.yaml /${PTH}/config.yaml
COPY --from=build-env /opt/service/bruteforce_service /${PTH}/bruteforce_service

EXPOSE 3009
ENTRYPOINT ["/bin/bruteforce_service"]



