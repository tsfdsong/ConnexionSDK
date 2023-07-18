package skywalking

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/pool"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"

	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

func CreateLocalSpan(ctx context.Context, opts ...go2sky.SpanOption) (s go2sky.Span, c context.Context, err error) {
	return GetSkyTrace().CreateLocalSpan(ctx, opts...)
}

func CreateExitSpan(ctx context.Context, operationName string, peer string, injector propagation.Injector) (s go2sky.Span, err error) {
	return GetSkyTrace().CreateExitSpan(ctx, operationName, peer, injector)
}

func SkyPostRequest(url string, body interface{}, v interface{}, ctx context.Context) error {
	client := pool.Pool.Get().(http.Client)
	defer pool.Pool.Put(client)

	rawBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request body failed, %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(rawBody))
	if err != nil {
		return fmt.Errorf("new post request failed, %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	//sign span
	signSpan, err := CreateExitSpan(ctx, "SignMachine", url, func(key, value string) error {
		req.Header.Set(key, value)
		return nil
	})
	if err != nil {
		return fmt.Errorf("create exit span failed, %v", err)
	}

	signSpan.SetComponent(2)
	signSpan.SetSpanLayer(agentv3.SpanLayer_RPCFramework)
	defer signSpan.End()

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

	signSpan.Tag(go2sky.TagHTTPMethod, http.MethodGet)
	signSpan.Tag(go2sky.TagURL, url)

	signSpan.Log(time.Now(), "[Sign]", fmt.Sprintf(" sign hash: %s \n", string(readBytes)))

	return nil
}
