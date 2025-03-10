package common

import (
	"fmt"
	"io"
	"net/http"
)

func NewRequest(method string, url string, body io.Reader) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body) // data теперь nil
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неудачный статус ответа: %s", resp.Status)
	}

	return nil
}
