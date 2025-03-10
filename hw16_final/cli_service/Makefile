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

# run: 
# 	./cli