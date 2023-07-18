package pool

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/machinebox/graphql"
)

var (
	graphClientPool = sync.Pool{
		New: func() interface{} {
			return http.DefaultClient
		},
	}
)

func InitGraphClient(PoolSize int) {
	for i := 0; i < PoolSize; i++ {
		Client := &http.Client{
			Transport: &http.Transport{
				DisableKeepAlives: true,
			},
		}
		graphClientPool.Put(Client)
	}
}

// ???  memory leak?
func GraphRequest(url, req string, response interface{}) error {
	httpClient := graphClientPool.Get().(*http.Client)
	defer graphClientPool.Put(httpClient)

	graphqlClient := graphql.NewClient(url, graphql.WithHTTPClient(httpClient))
	if graphqlClient == nil {
		return fmt.Errorf("emtpy market graph url")
	}

	graphqlReq := graphql.NewRequest(req)
	graphqlReq.Header.Set("Cache-Control", "no-cache")
	err := graphqlClient.Run(context.Background(), graphqlReq, &response)
	//logger.Logrus.WithFields(logrus.Fields{"ret": response}).Info("GraphRequest")
	return err
}
