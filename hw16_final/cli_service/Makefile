tree:
	tree | xclip -selection clipboard

# --------- CLI -----------
# Переменные
CLI_BINARY = ./main.go
CLI_OUTPUT = ../cli

# Команда для сборки бинарника cli
cli:
	go build -o $(CLI_OUTPUT) $(CLI_BINARY)

# Команда для удаления бинарника
del:
	rm -f $(CLI_OUTPUT)

	
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.8

lint: install-lint-deps
	golangci-lint run ./...

# run: 
# 	./cli