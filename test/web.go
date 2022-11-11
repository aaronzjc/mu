package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaronzjc/mu/internal/api/res"

	"github.com/gin-gonic/gin"
)

type Request struct {
	t *testing.T
	g *gin.Engine
	h gin.HandlerFunc

	Method string
	Url    string
}

func NewRequest(t *testing.T) *Request {
	return &Request{
		t: t,
		g: gin.Default(),
	}
}

func (r *Request) Handler(h gin.HandlerFunc) *Request {
	r.h = h
	return r
}

func (r *Request) Get(url string) *Request {
	r.Method = "GET"
	r.Url = url
	r.g.GET(url, r.h)
	return r
}

func (r *Request) Exec() *Response {
	resp := &Response{
		resp: httptest.NewRecorder(),
	}
	req, _ := http.NewRequest(r.Method, r.Url, nil)
	r.g.ServeHTTP(resp.resp, req)
	return resp
}

type Response struct {
	resp *httptest.ResponseRecorder
}

func (r *Response) Code() int {
	return r.resp.Code
}

func (r *Response) Body() string {
	return r.resp.Body.String()
}

func (r *Response) TryDecode() (code int, msg string, data any, err error) {
	var dataSt res.RespSt
	err = json.Unmarshal([]byte(r.Body()), &dataSt)
	if err != nil {
		return
	}
	code, msg, data = dataSt.Code, dataSt.Msg, dataSt.Data
	return
}
