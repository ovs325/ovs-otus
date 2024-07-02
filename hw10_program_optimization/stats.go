package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// Результат рефакторинга
// Было:              "423072614"
// Стало:              "27142424"
// Требуется не более: "31457280"

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Bytes()
		var user User
		if err := jsoniter.Unmarshal(line, &user); err != nil {
			return nil, err
		}
		if strings.Contains(user.Email, "."+domain) {
			res := strings.SplitN(user.Email, "@", 2)
			if len(res) >= 2 {
				result[strings.ToLower(res[1])]++
			}
		}
	}
	return result, scanner.Err()
}
