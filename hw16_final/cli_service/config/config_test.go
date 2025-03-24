package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Создаем временную директорию
	tempDir := t.TempDir()

	// Создаем файл config.yaml с тестовыми данными
	configContent := `
logger:
  level: info
http_server:
  scheme: http
  host: 127.0.0.1
  port: 3009
is_test: true
`
	configPath := filepath.Join(tempDir, "config.yaml")
	assert.Nil(t, os.WriteFile(configPath, []byte(configContent), 0o644), "failed to write config file")

	// Загружаем конфигурацию
	config, err := LoadConfig(tempDir)
	assert.NoError(t, err)

	// Проверяем параметры
	assert.Equal(t, "info", config.Logger.Level)
	assert.Equal(t, "http", config.HTTPServer.Scheme)
	assert.Equal(t, "127.0.0.1", config.HTTPServer.Host)
	assert.Equal(t, "3009", config.HTTPServer.Port)
	assert.True(t, config.IsTest)
}
