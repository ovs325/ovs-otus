// tests/integration_test.go
package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	ImageName          = "bruteforce_service"
	ContainerName      = ImageName + "_container"
	HostBruteforcePort = "3009"
	BaseURL            = `http://bruteforce_service_tst:` + HostBruteforcePort
	RunBruteforce      = false
)

type ResponseErr struct {
	Error string `json:"error"`
}

type IntegrationTest struct {
	suite.Suite
}

func (suite *IntegrationTest) SetupSuite() {
	if RunBruteforce {
		cmd := exec.Command(
			"docker",
			"run",
			"-d",
			"--rm",
			"-p",
			HostBruteforcePort+":3009",
			"--name",
			ContainerName,
			ImageName,
		)
		err := cmd.Run()
		suite.NoError(err)
		time.Sleep(5 * time.Second) // Ожидание запуска
	}
}

func (suite *IntegrationTest) TearDownSuite() {
	if RunBruteforce {
		cmd := exec.Command("docker", "stop", ContainerName)
		err := cmd.Run()
		suite.NoError(err)
	}
}

func (suite *IntegrationTest) TestIsAllowedHandler_Ok() {
	url, err := url.JoinPath(BaseURL, "is_allowed")
	suite.NoError(err)

	body := map[string]interface{}{
		"content": map[string]string{
			"ip":    "137.23.45.167",
			"login": "Abrec",
			"pass":  "1234QWER+",
		},
	}

	jsonBody, err := json.Marshal(body)
	suite.NoError(err)
	req, err := http.NewRequestWithContext(context.Background(), "PATCH", url, bytes.NewBuffer(jsonBody))
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var responseBody any
	rBody := resp.Body
	err = json.NewDecoder(rBody).Decode(&responseBody)
	suite.NoError(err)

	_, ok := (responseBody).(bool)

	suite.True(ok, "I expected the response to be of the 'Bool' type, but a different type was received.")
}

func (suite *IntegrationTest) TestIsAllowedHandler_InvalidParams() {
	url, err := url.JoinPath(BaseURL, "is_allowed")
	suite.NoError(err)

	// Некорректное тело запроса
	bodyErr := []string{"Abrec", "pass", "1234QWER+"}

	jsonBody, err := json.Marshal(bodyErr)
	suite.NoError(err)

	req, err := http.NewRequestWithContext(context.Background(), "PATCH", url, bytes.NewBuffer(jsonBody))
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	var responseBody ResponseErr
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	suite.NoError(err)

	suite.Equal(
		"ошибка при получении тела запроса: json: cannot unmarshal array into Go value of type handlers.CommonContentRequest",
		responseBody.Error,
	)
}

func (suite *IntegrationTest) TestResetBucketHandler_Ok() {
	urlRaw, err := url.JoinPath(BaseURL, "reset")
	suite.NoError(err)

	reqURL, err := url.Parse(urlRaw)
	suite.NoError(err)

	// Add the query parameter
	query := reqURL.Query()
	query.Add("params", "ip^137.23.45.111")
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(context.Background(), "PATCH", reqURL.String(), nil)
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *IntegrationTest) TestResetBucketHandler_InvalidParams() {
	urlRaw, err := url.JoinPath(BaseURL, "reset")
	suite.NoError(err)

	reqURL, err := url.Parse(urlRaw)
	suite.NoError(err)

	query := reqURL.Query()
	query.Add("params", "invalid_param")
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(context.Background(), "PATCH", reqURL.String(), nil)
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	var responseBody ResponseErr
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	suite.NoError(err)

	suite.Equal("ошибка при получении Query-параметров", responseBody.Error)
}

func (suite *IntegrationTest) TestAddToBlacklistHandler_Ok() {
	url, err := url.JoinPath(BaseURL, "add/black")
	suite.NoError(err)

	req, err := http.NewRequestWithContext(context.Background(), "PATCH", url, nil)
	q := req.URL.Query()
	q.Add("network", "137.23.45.167/25")
	req.URL.RawQuery = q.Encode()
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *IntegrationTest) TestAddToBlacklistHandler_InvalidParams() {
	url, err := url.JoinPath(BaseURL, "add/black")
	suite.NoError(err)

	req, err := http.NewRequestWithContext(context.Background(), "PATCH", url, nil)
	q := req.URL.Query()
	q.Add("errNetwork", "137.23.45.167/25")
	req.URL.RawQuery = q.Encode()
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	var responseBody ResponseErr
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	suite.NoError(err)

	suite.Equal("ошибка при получении параметра 'network'", responseBody.Error)
}

func (suite *IntegrationTest) TestAddToWhitelistHandler_Ok() {
	url, err := url.JoinPath(BaseURL, "add/white")
	suite.NoError(err)

	req, err := http.NewRequestWithContext(context.Background(), "PATCH", url, nil)
	q := req.URL.Query()
	q.Add("network", "137.23.45.168/25")
	req.URL.RawQuery = q.Encode()
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *IntegrationTest) TestAddToWhitelistHandler_InvalidParams() {
	url, err := url.JoinPath(BaseURL, "add/white")
	suite.NoError(err)

	req, err := http.NewRequestWithContext(context.Background(), "PATCH", url, nil)
	q := req.URL.Query()
	q.Add("errNetwork", "137.23.45.168/25")
	req.URL.RawQuery = q.Encode()
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	var responseBody ResponseErr
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	suite.NoError(err)

	suite.Equal("ошибка при получении параметра 'network'", responseBody.Error)
}

func (suite *IntegrationTest) TestDelFromBlacklistHandler_Ok() {
	url, err := url.JoinPath(BaseURL, "del/black")
	suite.NoError(err)

	req, err := http.NewRequestWithContext(context.Background(), "DELETE", url, nil)
	q := req.URL.Query()
	q.Add("network", "137.23.45.167/25")
	req.URL.RawQuery = q.Encode()
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *IntegrationTest) TestDelFromBlacklistHandler_InvalidParams() {
	url, err := url.JoinPath(BaseURL, "del/black")
	suite.NoError(err)

	req, err := http.NewRequestWithContext(context.Background(), "DELETE", url, nil)
	q := req.URL.Query()
	q.Add("errNetwork", "137.23.45.167/25")
	req.URL.RawQuery = q.Encode()
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	var responseBody ResponseErr
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	suite.NoError(err)

	suite.Equal("ошибка при получении параметра 'network'", responseBody.Error)
}

func (suite *IntegrationTest) TestDelFromWhitelistHandler_Ok() {
	url, err := url.JoinPath(BaseURL, "del/white")
	suite.NoError(err)

	req, err := http.NewRequestWithContext(context.Background(), "DELETE", url, nil)
	q := req.URL.Query()
	q.Add("network", "137.23.45.168/25")
	req.URL.RawQuery = q.Encode()
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *IntegrationTest) TestDelFromWhitelistHandler_InvalidParams() {
	url, err := url.JoinPath(BaseURL, "del/white")
	suite.NoError(err)

	req, err := http.NewRequestWithContext(context.Background(), "DELETE", url, nil)
	q := req.URL.Query()
	q.Add("errNetwork", "137.23.45.168/25")
	req.URL.RawQuery = q.Encode()
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	var responseBody ResponseErr
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	suite.NoError(err)

	suite.Equal("ошибка при получении параметра 'network'", responseBody.Error)
}

func waitForService(url string) {
	for {
		req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
		if err != nil {
			fmt.Println("Соединение не установлено: %w", err)
			break // Exit loop if error response is received
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Println("Соединение установлено")
			break // Exit loop if successful response is received
		}
		defer resp.Body.Close()
		if err != nil {
			fmt.Printf("Ожидание соединения: %s \t%v\n", url, err)
		} else {
			fmt.Printf("Ожидание соединения: %s \n", url)
		}
		time.Sleep(2 * time.Second) // Wait before retrying
	}
}
func TestIntegrationTestSuite(t *testing.T) {
	waitForService(BaseURL + "/ok")
	suite.Run(t, new(IntegrationTest))
}
