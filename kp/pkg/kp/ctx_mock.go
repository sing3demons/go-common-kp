package kp

import (
	"context"
	"net/url"
)

type MockContext struct {
	context.Context
	MockRequest
}

type MockRequest struct {
	ParamsMap     map[string]string
	PathParamsMap map[string]string
	QueryVals     url.Values
	BodyStr       string
	ErrBody       error

	methodsToCallTime map[string]uint
	AddDataStr        map[string]string
}

func (m *MockRequest) Context() context.Context {
	m.methodsToCallTime["Context"]++
	return context.Background()
}
func (m *MockRequest) Param(key string) string {
	m.methodsToCallTime["Param"]++
	return m.ParamsMap[key]
}
func (m *MockRequest) PathParam(key string) string {
	m.methodsToCallTime["PathParam"]++
	return m.PathParamsMap[key]
}
func (m *MockRequest) Bind(any) error {
	m.methodsToCallTime["Bind"]++
	return nil
}
func (m *MockRequest) HostName() string {
	m.methodsToCallTime["HostName"]++
	return "localhost"
}
func (m *MockRequest) Params(key string) []string {
	m.methodsToCallTime["Params"]++
	return []string{m.ParamsMap[key]}
}
func (m *MockRequest) ClientIP() string {
	m.methodsToCallTime["ClientIP"]++
	if m.AddDataStr["ClientIP"] != "" {
		return m.AddDataStr["ClientIP"]
	}
	return "127.0.0.1"
}
func (m *MockRequest) UserAgent() string {
	m.methodsToCallTime["UserAgent"]++
	if m.AddDataStr["UserAgent"] != "" {
		return m.AddDataStr["UserAgent"]
	}
	return "MockAgent"
}

func (m *MockRequest) Referer() string {
	m.methodsToCallTime["Referer"]++
	if m.AddDataStr["Referer"] != "" {
		return m.AddDataStr["Referer"]
	}
	return ""
}

func (m *MockRequest) Method() string {
	m.methodsToCallTime["Method"]++
	if m.AddDataStr["Method"] != "" {
		return m.AddDataStr["Method"]
	}
	return "GET"
}

func (m *MockRequest) URL() string {
	m.methodsToCallTime["URL"]++
	if m.AddDataStr["URL"] != "" {
		return m.AddDataStr["URL"]
	}
	return "http://localhost/mock"
}
func (m *MockRequest) TransactionId() string {
	m.methodsToCallTime["TransactionId"]++
	if m.AddDataStr["TransactionId"] != "" {
		return m.AddDataStr["TransactionId"]
	}
	return "txid"
}

func (m *MockRequest) SessionId() string {
	m.methodsToCallTime["SessionId"]++
	if m.AddDataStr["SessionId"] != "" {
		return m.AddDataStr["SessionId"]
	}
	return "sessionid"
}

func (m *MockRequest) RequestId() string {
	m.methodsToCallTime["RequestId"]++
	if m.AddDataStr["RequestId"] != "" {
		return m.AddDataStr["RequestId"]
	}
	return "requestid"
}
func (m *MockRequest) Body() (string, error) {
	m.methodsToCallTime["Body"]++
	if m.AddDataStr["Body"] != "" {
		return m.AddDataStr["Body"], nil
	}
	return m.BodyStr, m.ErrBody
}
func (m *MockRequest) Query() url.Values {
	m.methodsToCallTime["Query"]++
	return m.QueryVals
}

func (m *MockRequest) PathParams() map[string]string {
	m.methodsToCallTime["PathParams"]++
	return m.PathParamsMap
}
