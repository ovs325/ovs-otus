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
params:
  login:
    lenQueue: 10000
    capacity: 10
    leakageRate: 0.16666666666666666
  pass:
    lenQueue: 10000
    capacity: 10
    leakageRate: 1.6666666666666667
  ip:
    lenQueue: 10000
    capacity: 10
    leakageRate: 16.666666666666668
default:
  lenQueue: 10000
  capacity: 10
  leakageRate: 0.5
`
	configPath := filepath.Join(tempDir, "config.yaml")
	assert.Nil(t, os.WriteFile(configPath, []byte(configContent), 0644), "failed to write config file")

	// Загружаем конфигурацию
	config, err := LoadConfig(tempDir)
	assert.NoError(t, err)

	// Проверяем параметры
	assert.Equal(t, int64(10000), config.Params["login"].LenQueue)
	assert.Equal(t, int64(10), config.Params["login"].Capacity)
	assert.Equal(t, float64(0.16666666666666666), config.Params["login"].LeakageRate)

	assert.Equal(t, int64(10000), config.Params["pass"].LenQueue)
	assert.Equal(t, int64(10), config.Params["pass"].Capacity)
	assert.Equal(t, float64(1.6666666666666667), config.Params["pass"].LeakageRate)

	assert.Equal(t, int64(10000), config.Params["ip"].LenQueue)
	assert.Equal(t, int64(10), config.Params["ip"].Capacity)
	assert.Equal(t, float64(16.666666666666668), config.Params["ip"].LeakageRate)

	// Проверяем значения по умолчанию
	assert.Equal(t, int64(10000), config.Deault.LenQueue)
	assert.Equal(t, int64(10), config.Deault.Capacity)
	assert.Equal(t, float64(0.5), config.Deault.LeakageRate)
}
