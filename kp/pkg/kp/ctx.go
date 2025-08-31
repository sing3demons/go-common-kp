package kp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	config "github.com/sing3demons/go-common-kp/kp/configs"
	"github.com/sing3demons/go-common-kp/kp/pkg/kafka"
	"github.com/sing3demons/go-common-kp/kp/pkg/logger"
	"go.opentelemetry.io/otel/trace"
)

type IncomingReq struct {
	URL     string            `json:"url"`
	IP      string            `json:"ip"`
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
	Params  map[string]string `json:"params"`
	Query   map[string]any    `json:"query"`
	Body    map[string]any    `json:"body"`
}

type Context struct {
	context.Context
	Request
	http.ResponseWriter
	kafka.Client
	detail   logger.CustomLoggerService
	incoming IncomingReq
	metaData logger.Metadata
	conf     *config.Config
	appLog   logger.LoggerService
}
type SubscribeFunc func(c *Context) error

type Request interface {
	Context() context.Context
	Param(string) string
	PathParam(string) string
	Bind(any) error
	HostName() string
	Params(string) []string
	ClientIP() string
	UserAgent() string
	Referer() string
	Method() string
	URL() string
	TransactionId() string
	SessionId() string
	RequestId() string
	Body() (string, error)
	Query() url.Values
	PathParams() map[string]string
}

type LogService struct {
	appLog         logger.LoggerService
	detailLog      logger.LoggerService
	summaryLog     logger.LoggerService
	maskingService logger.MaskingServiceInterface
}

func newContext(w http.ResponseWriter, r Request, k kafka.Client, log LogService, conf *config.Config) *Context {
	c := r.Context()
	if c == nil {
		c = context.Background()
	}
	traceID := trace.SpanFromContext(c).SpanContext().TraceID().String()
	spanId := trace.SpanFromContext(c).SpanContext().SpanID().String()

	t := logger.NewTimer()
	kpLog := logger.NewCustomLogger(log.detailLog, log.summaryLog, t, log.maskingService)
	ctx := &Context{
		Context:        c,
		Request:        r,
		ResponseWriter: w,
		Client:         k,
		conf:           conf,
		appLog:         log.appLog,
	}

	isHTTP := true

	broker := "none"
	source := "api"
	if w == nil {
		broker = r.HostName()
		source = "event-source"
		isHTTP = false
	}

	meta := logger.Metadata{
		ClientIP:  r.ClientIP(),
		UserAgent: r.UserAgent(),
		Referer:   r.Referer(),
		Method:    r.Method(),
		URL:       r.URL(),
		Source:    source,
		Broker:    broker,
		TraceId:   traceID,
		SpanId:    spanId,
	}
	hostName, _ := os.Hostname()
	ctx.metaData = meta

	customLog := logger.LogDto{
		ServiceName:      conf.App.Name,
		LogType:          "detail",
		ComponentVersion: conf.App.Version,
		Instance:         hostName,
		Metadata:         meta,
		SessionId:        ctx.SessionId(),
		RequestId:        ctx.RequestId(),
	}
	kpLog.Init(customLog)
	if !isHTTP {
		topic := r.Param("topic")
		summary := logger.LogEventTag{
			Node:        "consumer",
			Command:     topic,
			Code:        "200",
			Description: "",
		}
		body, err := ctx.Body()
		if err != nil {
			summary.Code = "500"
			summary.Description = err.Error()
			kpLog.SetSummary(summary).Error(logger.NewConsuming(topic, "kafka"+"_consumer"), map[string]any{
				"topic":  topic,
				"broker": broker,
				"error":  err.Error(),
			})
		} else {
			kpLog.SetSummary(summary).Info(logger.NewConsuming(topic, "kafka"+"_consumer"), map[string]any{
				"topic":  topic,
				"broker": broker,
				"body":   body,
			})
		}
	} else {
		body := map[string]any{}
		rawBody, err := r.Body()
		if err == nil {
			// convert to map
			json.Unmarshal([]byte(rawBody), &body)
		}
		query := map[string]any{}
		if r.Query() != nil {
			for k, v := range r.Query() {
				query[k] = v
			}
		}
		headers := map[string]string{}
		for k, v := range ctx.Header() {
			headers[k] = ""
			if len(v) > 0 {
				headers[k] = v[0]
				if len(v) > 1 {
					// If multiple values, join them with commas
					headers[k] = ""
					for i, val := range v {
						if i > 0 {
							headers[k] += ","
						}
						headers[k] += val
					}
				}
			}
		}

		ctx.incoming = IncomingReq{
			URL:     ctx.URL(),
			IP:      ctx.ClientIP(),
			Headers: headers,
			Method:  ctx.Method(),
			Params:  ctx.PathParams(),
			Query:   query,
			Body:    body,
		}
	}

	ctx.detail = kpLog
	return ctx
}

type AppLogStruct struct {
	LogType     string `json:"logType"`
	LogLevel    string `json:"logLevel"`
	Message     string `json:"message"`
	ServiceName string `json:"serviceName"`
	RequestId   string `json:"requestId"`
	SessionId   string `json:"sessionId"`
}

func (c *Context) Info(msg any) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		jsonMsg = []byte("failed to marshal message")
	}

	appLog := AppLogStruct{
		LogType:     "app",
		LogLevel:    "info",
		Message:     string(jsonMsg),
		ServiceName: c.conf.App.Name,
		RequestId:   c.RequestId(),
		SessionId:   c.SessionId(),
	}
	strMsg, _ := json.Marshal(appLog)
	c.appLog.Info(string(strMsg))
}

func (c *Context) Debug(msg any) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		jsonMsg = []byte("failed to marshal message")
	}

	appLog := AppLogStruct{
		LogType:     "app",
		LogLevel:    "debug",
		Message:     string(jsonMsg),
		ServiceName: c.conf.App.Name,
		RequestId:   c.RequestId(),
		SessionId:   c.SessionId(),
	}
	strMsg, _ := json.Marshal(appLog)
	c.appLog.Debug(string(strMsg))
}

func (c *Context) Error(msg any) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		jsonMsg = []byte("failed to marshal message")
	}

	appLog := AppLogStruct{
		LogType:     "app",
		LogLevel:    "error",
		Message:     string(jsonMsg),
		ServiceName: c.conf.App.Name,
		RequestId:   c.RequestId(),
		SessionId:   c.SessionId(),
	}
	strMsg, _ := json.Marshal(appLog)
	c.appLog.Error(string(strMsg))
}

func (c *Context) GetConfig(key string) string {
	return c.conf.Get(key)
}

func (c *Context) GetConfigOrDefault(key, defaultValue string) string {
	return c.conf.GetOrDefault(key, defaultValue)
}

func (c *Context) GetIncoming() IncomingReq {
	return c.incoming
}

func (c *Context) JSON(code int, v any) error {
	if c.ResponseWriter != nil {
		c.ResponseWriter.Header().Del("Content-Length")
		c.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8") // set first
		c.ResponseWriter.WriteHeader(code)                                               // then send status

		if err := json.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
			if c.detail != nil {
				c.detail.AddField("Error", err.Error())
			}
		}
		if c.detail != nil {
			c.detail.Info(logger.NewOutbound("client", ""), v)
			c.detail.End(code, "")
		}
	}

	return nil
}

func (c *Context) LogAuto(masks ...logger.MaskingOptionDto) logger.CustomLoggerService {
	if c.incoming.URL != "" && c.incoming.Method != "" {
		c.detail.Info(logger.NewInbound("client", ""), c.incoming, masks...)
		c.incoming = IncomingReq{}
	}
	return c.detail
}

func (c *Context) Log() logger.CustomLoggerService {
	return c.detail
}
