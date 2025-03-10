// tests/integration_test.go
package tests

import (
	"bytes"
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
	IMAGE_NAME           = "bruteforce_service"
	CONTAINER_NAME       = IMAGE_NAME + "_container"
	HOST_BRUTEFORCE_PORT = "3009"
	BASE_URL             = `http://bruteforce_service_tst:` + HOST_BRUTEFORCE_PORT
	RUN_BRUTEFORCE       = false
)

type ResponseErr struct {
	Error string `json:"error"`
}

type IntegrationTest struct {
	suite.Suite
}

func (suite *IntegrationTest) SetupSuite() {
	if RUN_BRUTEFORCE {
		cmd := exec.Command("docker", "run", "-d", "--rm", "-p", HOST_BRUTEFORCE_PORT+":3009", "--name", CONTAINER_NAME, IMAGE_NAME)
		err := cmd.Run()
		suite.NoError(err)
		time.Sleep(5 * time.Second) // Ожидание запуска
	}
}

func (suite *IntegrationTest) TearDownSuite() {
	if RUN_BRUTEFORCE {
		cmd := exec.Command("docker", "stop", CONTAINER_NAME)
		err := cmd.Run()
		suite.NoError(err)
	}
}

func (suite *IntegrationTest) TestIsAllowedHandler_Ok() {
	url, err := url.JoinPath(BASE_URL, "is_allowed")
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
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonBody))
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
	url, err := url.JoinPath(BASE_URL, "is_allowed")
	suite.NoError(err)

	// Некорректное тело запроса
	bodyErr := []string{"Abrec", "pass", "1234QWER+"}

	jsonBody, err := json.Marshal(bodyErr)
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonBody))
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
	url_, err := url.JoinPath(BASE_URL, "reset")
	suite.NoError(err)

	reqURL, err := url.Parse(url_)
	suite.NoError(err)

	// Add the query parameter
	query := reqURL.Query()
	query.Add("params", "ip^137.23.45.111")
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequest("PATCH", reqURL.String(), nil)
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *IntegrationTest) TestResetBucketHandler_InvalidParams() {
	url_, err := url.JoinPath(BASE_URL, "reset")
	suite.NoError(err)

	reqURL, err := url.Parse(url_)
	suite.NoError(err)

	query := reqURL.Query()
	query.Add("params", "invalid_param")
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequest("PATCH", reqURL.String(), nil)
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
	url, err := url.JoinPath(BASE_URL, "add/black")
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", url, nil)
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
	url, err := url.JoinPath(BASE_URL, "add/black")
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", url, nil)
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
	url, err := url.JoinPath(BASE_URL, "add/white")
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", url, nil)
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
	url, err := url.JoinPath(BASE_URL, "add/white")
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", url, nil)
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
	url, err := url.JoinPath(BASE_URL, "del/black")
	suite.NoError(err)

	req, err := http.NewRequest("DELETE", url, nil)
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
	url, err := url.JoinPath(BASE_URL, "del/black")
	suite.NoError(err)

	req, err := http.NewRequest("DELETE", url, nil)
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
	url, err := url.JoinPath(BASE_URL, "del/white")
	suite.NoError(err)

	req, err := http.NewRequest("DELETE", url, nil)
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
	url, err := url.JoinPath(BASE_URL, "del/white")
	suite.NoError(err)

	req, err := http.NewRequest("DELETE", url, nil)
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
		// Create a new PATCH request
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Println("Соединение установлено")
			break // Exit loop if successful response is received
		}
		if err != nil {
			fmt.Printf("Ожидание соединения: %s \t%v\n", url, err)
		} else {
			fmt.Printf("Ожидание соединения: %s \n", url)
		}
		time.Sleep(2 * time.Second) // Wait before retrying
	}
}
func TestIntegrationTestSuite(t *testing.T) {
	waitForService(BASE_URL + "/ok")
	suite.Run(t, new(IntegrationTest))
}
