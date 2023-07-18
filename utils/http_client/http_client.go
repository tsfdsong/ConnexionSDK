package http_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/pool"
	"io/ioutil"
	"net/http"

	e "github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

const (
	retryTime = 5
)

var StatusCodeError = errors.New("status code not 200")

func HttpClientReq(url string, body interface{}, v interface{}) error {
	client := pool.Pool.Get().(http.Client)
	defer pool.Pool.Put(client)

	rawBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request body failed, %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rawBody))
	if err != nil {
		return fmt.Errorf("new post request failed, %v", err)
	}
	if req != nil && req.Body != nil {
		defer req.Body.Close()
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request failed, %v", err)
	}

	if resp == nil {
		return fmt.Errorf("do request and response is empty")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post request code is {%d} and status is {%s}", resp.StatusCode, resp.Status)
	}

	readBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read reponse body failed, %v", err)
	}

	err = json.Unmarshal([]byte(readBytes), v)
	if err != nil {
		return fmt.Errorf("unmarshal reponse body failed, %v", err)
	}

	return nil
}

func HttpClientReqWithGet(url string, v interface{}) error {
	client := pool.Pool.Get().(http.Client)
	defer pool.Pool.Put(client)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("new request %v", err)
	}
	if req != nil && req.Body != nil {
		defer req.Body.Close()
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return fmt.Errorf("do request %v", err)
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request status: %d {%s}", resp.StatusCode, resp.Status)
	}

	readBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response %v", err)
	}

	err = json.Unmarshal([]byte(readBytes), v)
	if err != nil {
		return fmt.Errorf("unmarshal response %v", err)
	}
	return nil
}

func HttpClientReqWithExactError(url string, body interface{}, v interface{}) (error, error) {
	client := pool.Pool.Get().(http.Client)
	defer pool.Pool.Put(client)

	rawBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rawBody))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("new http request failed")
		return errors.New("new http request failed"), nil
	}
	if req != nil && req.Body != nil {
		defer req.Body.Close()
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil || resp == nil {
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("http request failed")
		}
		if resp == nil {
			logger.Logrus.Error("http reponse with empty resp")
		}
		return nil, errors.New("http request failed")
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, StatusCodeError
	}

	readBytes, _ := ioutil.ReadAll(resp.Body)
	//logger.Logrus.Info(fmt.Sprintf("%+v", string(readBytes)))
	json.Unmarshal([]byte(readBytes), v)
	return nil, nil
}

func PostUrlRetry(url string, params map[string]string, body interface{}, headers map[string]string, v interface{}) error {
	i := 0
	data, err := PostUrl(url, params, body, headers)
	for err != nil && i < retryTime {
		data, err = PostUrl(url, params, body, headers)
		i++
	}
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return e.Wrap(err, "json unmarshal res data error")
	}

	return nil
}

func PostUrl(url string, params map[string]string, body interface{}, headers map[string]string) ([]byte, error) {
	var (
		bodyJson []byte
		req      *http.Request
		err      error
	)

	client := pool.Pool.Get().(http.Client)
	defer pool.Pool.Put(client)

	if body != nil {
		bodyJson, err = json.Marshal(body)
		if err != nil {
			return nil, e.Wrap(err, "json marshal request body error")
		}
	}

	req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, e.Wrap(err, "NewRequest error")
	}

	contentType := "Content-type"
	req.Header.Set(contentType, headers[contentType])
	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, e.Wrap(err, "client.Do error")
	}
	defer response.Body.Close()
	d, err_ := ioutil.ReadAll(response.Body)
	if err_ != nil {
		return nil, e.Wrap(err, "ioutil.ReadAll error")
	}

	return d, nil
}
