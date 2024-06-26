package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(
			os.Stderr,
			"Слишком мало аргументов. Используйте следующий формат: %s <directory> <command> [args...]\n",
			os.Args[0],
		)
		os.Exit(1)
	}

	dir := os.Args[1]
	command := os.Args[2]
	args := os.Args[3:]

	env, err := ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading environment directory: %v\n", err)
		os.Exit(1)
	}

	commands := []string{command}
	commands = append(commands, args...)

	returnCode := RunCmd(commands, env)
	os.Exit(returnCode)
}
