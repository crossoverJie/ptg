package client

import (
	"github.com/crossoverJie/ptg/meta"
	"net/http"
	"time"
)

type (
	Client interface {
		Request() (*meta.Response, error)
	}
)

const (
	Http = "http"
	Grpc = "grpc"
)

func NewClient(method, url, requestBody string, meta *meta.Meta) Client {
	if meta.Protocol() == Http {
		return &httpClient{
			Method:      method,
			Url:         url,
			RequestBody: requestBody,
			httpClient: &http.Client{
				Transport: &http.Transport{
					DisableCompression:    false,
					ResponseHeaderTimeout: time.Millisecond * time.Duration(10000),
					DisableKeepAlives:     false,
				},
			},
			meta: meta,
		}
	} else {
		return NewGrpcClient(meta)
	}

}
