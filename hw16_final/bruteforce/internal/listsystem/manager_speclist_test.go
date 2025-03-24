package listsystem

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecList(t *testing.T) {
	// Создаем экземпляр SpecList для тестирования
	seeker := SearchEngineIP
	verifier := SubnetNameVerifier
	specList := NewSpecList("id", 1000, seeker, verifier)

	// Проверяем добавление ID
	assert.True(t, specList.AddID("192.1.1.0/25", true), "Expected AddID to succeed, but it failed")

	// Проверяем наличие ID
	isFound, isBlack, err := specList.ToFind("192.1.1.0")
	assert.NoErrorf(t, err, "Expected ToFind to find the ID, but it didn't: %v", err)
	assert.True(t, isFound)
	assert.True(t, bool(isBlack))

	isFound, isBlack, err = specList.ToFind("192.1.1.0/err")
	assert.Equal(t, fmt.Errorf("некорректный IP-адрес: 192.1.1.0/err"), err)
	assert.False(t, isFound)
	assert.False(t, bool(isBlack))

	// Проверяем удаление ID
	assert.True(t, specList.DelID("192.1.1.0/25"), "Expected DelID to succeed, but it failed")

	// Проверяем отсутствие ID после удаления
	isFound, _, _ = specList.ToFind("192.1.1.0")
	assert.False(t, isFound, "Expected ToFind to not find the ID after deletion")

	// Проверяем сброс списка
	specList.Reset()
	isFound, _, _ = specList.ToFind("192.1.1.0")
	assert.False(t, isFound, "Expected ToFind to not find the ID after reset")

	// Проверяем имя списка
	assert.Equal(t, "id", specList.GetName(), "Expected GetName to return 'id', got '%s'", specList.GetName())
}
