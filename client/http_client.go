package client

import (
	"fmt"
	"github.com/crossoverJie/ptg/meta"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type httpClient struct {
	Method      string
	Url         string
	RequestBody string
	httpClient  *http.Client
	meta        *meta.Meta
}

func (c *httpClient) Request() (*meta.Response, error) {
	var payload io.Reader
	if len(c.RequestBody) > 0 {
		payload = strings.NewReader(`{"name":"abc"}`)
	}
	req, err := http.NewRequest(c.Method, c.Url, payload)
	//req.Close = true
	if err != nil {
		fmt.Println("An error occured doing request", err)
		return nil, err
	}
	req.Header.Add("User-Agent", "ptg")
	for k, v := range c.meta.HeaderMap() {
		req.Header.Add(k, v)
	}

	//httpClient := &http.Client{
	//	Transport: &http.Transport{
	//		ResponseHeaderTimeout: time.Millisecond * time.Duration(1000),
	//		DisableKeepAlives: true,
	//	},
	//}

	start := time.Now()
	response, err := c.httpClient.Do(req)
	r := &meta.Response{
		RequestTime: time.Since(start),
	}
	//SlowRequestTime = r.slowRequest()
	//FastRequestTime = r.fastRequest()
	meta.GetResult().SetSlowRequestTime(r.SlowRequest()).
		SetFastRequestTime(r.FastRequest())
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("response is nil")
	}
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("http code not OK: %v", response.StatusCode))
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("response bodyPath read err:%v", err))
	}
	r.ResponseSize = len(body)
	return r, nil
}
