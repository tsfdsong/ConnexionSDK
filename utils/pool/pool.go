package pool

import (
	"net/http"
	"sync"
	"time"
)

var (
	Pool = sync.Pool{
		New: func() interface{} {
			Client := http.Client{
				Transport: &http.Transport{
					DisableKeepAlives: true,
				},
				Timeout: time.Duration(time.Second * 10),
			}
			return Client
		},
	}
)

func InitClient(PoolSize int) {
	for i := 0; i < PoolSize; i++ {
		Client := http.Client{
			Transport: &http.Transport{
				DisableKeepAlives: true,
			},
		}

		Pool.Put(Client)
	}
}
