package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	errLen := ValidationErrors{}
	AppErr(&errLen, "Phones", "длина элемента 1 должна быть = ", 11)

	errAll := ValidationErrors{}
	AppErr(&errAll, "ID", "длинна должна быть = ", 36)
	AppErr(&errAll, "Age", "должно быть не менее ", 18)
	AppErr(&errAll, "Email", "не соответствует регулярному выражению ", "^\\w+@\\w+\\.\\w+$")
	AppErr(&errAll, "Role", "должен быть из множества ", "admin,stuff")
	AppErr(&errAll, "Phones", "длина элемента 0 должна быть = ", 11)
	AppErr(&errAll, "Phones", "длина элемента 1 должна быть = ", 11)

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Chikatilo",
				Age:    20,
				Email:  "tea@bk.ru",
				Role:   "admin",
				Phones: []string{"79605025555"},
				meta:   json.RawMessage([]byte("Странное сообщение")),
			},
			expectedErr: error(nil),
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Chikatilo",
				Age:    20,
				Email:  "tea@bk.ru",
				Role:   "admin",
				Phones: []string{"79605025555", "9621620101"},
				meta:   json.RawMessage([]byte("Странное сообщение")),
			},
			expectedErr: errLen,
		},
		{
			in: User{
				ID:     "1234567890123456",
				Name:   "ЗдесьНеПроверяют",
				Age:    10,
				Email:  "@bk.ru",
				Role:   "luser",
				Phones: []string{"755", "962"},
				meta:   json.RawMessage([]byte("Здесь не проверяют")),
			},
			expectedErr: errAll,
		},
		{
			in:          ValidationErrors{},
			expectedErr: fmt.Errorf("ожидалась структура, получена %s", "slice"),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			errRes := Validate(tt.in)
			assert.Equal(t, errRes, tt.expectedErr)
		})
	}
}
