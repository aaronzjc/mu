package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/aaronzjc/mu/pkg/logger"
)

type HttpClient struct {
	client  *http.Client
	Timeout time.Duration
}

var httpclient *HttpClient

func init() {
	httpclient = &HttpClient{
		client:  new(http.Client),
		Timeout: time.Second * 3,
	}
	httpclient.client.Timeout = httpclient.Timeout
}

func (c *HttpClient) Do(ctx context.Context, method string, url string, params map[string]interface{}, body []byte, headers map[string]string) (*http.Response, error) {
	// 处理POST的参数
	var buf io.Reader
	if len(body) > 0 {
		buf = bytes.NewBuffer(body)
	}

	req, _ := http.NewRequestWithContext(ctx, method, url, buf)

	// 补充headers
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	// 格式化GET的参数
	pasrseUrlParams(req, params)

	// 发送请求&记录请求耗时
	var err error
	var resp *http.Response
	start := time.Now().UnixMilli()
	resp, err = c.client.Do(req)
	ts, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", (float64)(time.Now().UnixMilli()-start)/1000), 64)
	logger.Request(req, resp, ts, err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return req.Response, nil
}

func (c *HttpClient) Get(ctx context.Context, url string, params map[string]interface{}) (string, error) {
	resp, err := c.Do(ctx, "GET", url, params, nil, nil)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func pasrseUrlParams(req *http.Request, params map[string]interface{}) {
	if len(params) == 0 {
		return
	}
	query := req.URL.Query()
	for k, v := range params {
		switch v := v.(type) {
		case string:
			query.Add(k, v)
		case int:
			query.Add(k, strconv.Itoa(v))
		case float64:
			query.Add(k, strconv.FormatFloat(v, 'f', -1, 64))
		}
	}
	req.URL.RawQuery = query.Encode()
}
