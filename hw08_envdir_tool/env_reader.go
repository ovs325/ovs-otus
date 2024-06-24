package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue помогает отличить пустые файлы от файлов с первой пустой строкой.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir считывает указанный каталог и возвращает отображение переменных env.
// Переменные представлены в виде файлов, где filename - это имя переменной, первая строка файла - это значение.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			firstStr, needRemove, err := GetEnvVar(filepath.Join(dir, name))
			if err != nil {
				return nil, err
			}
			env[name] = EnvValue{Value: firstStr, NeedRemove: needRemove}
		}
	}
	return env, nil
}

// Считываем содержимое файла, определяем необходимость удаления.
func GetEnvVar(path string) (string, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false, err
	}
	text := string(bytes.ReplaceAll(data, []byte{0x00}, []byte{'\n'}))
	text = strings.TrimRight(text, " \t")
	if text == "" {
		return "", true, nil
	}
	return strings.SplitN(text, "\n", 2)[0], false, nil
}
