package kp

import (
	"context"
	"net/http/httptest"
	"net/url"
	"testing"

	config "github.com/sing3demons/go-common-kp/kp/configs"
	"github.com/sing3demons/go-common-kp/kp/pkg/kafka"
	"github.com/sing3demons/go-common-kp/kp/pkg/logger"
)

// MockKafkaClient implements kafka.Client interface for testing
type MockKafkaClient struct {
	PublishCalls     []PublishCall
	SubscribeCalls   []SubscribeCall
	CreateTopicCalls []CreateTopicCall
	DeleteTopicCalls []DeleteTopicCall
	CloseCalls       int
	PublishError     error
	SubscribeError   error
	CreateTopicError error
	DeleteTopicError error
	CloseError       error
}

type PublishCall struct {
	Topic   string
	Message []byte
}

type SubscribeCall struct {
	Topic string
}

type CreateTopicCall struct {
	Name string
}

type DeleteTopicCall struct {
	Name string
}

func (m *MockKafkaClient) Publish(ctx context.Context, topic string, message []byte) error {
	m.PublishCalls = append(m.PublishCalls, PublishCall{Topic: topic, Message: message})
	return m.PublishError
}

func (m *MockKafkaClient) Subscribe(ctx context.Context, topic string) (*kafka.Message, error) {
	m.SubscribeCalls = append(m.SubscribeCalls, SubscribeCall{Topic: topic})
	if m.SubscribeError != nil {
		return nil, m.SubscribeError
	}
	return &kafka.Message{}, nil
}

func (m *MockKafkaClient) CreateTopic(ctx context.Context, name string) error {
	m.CreateTopicCalls = append(m.CreateTopicCalls, CreateTopicCall{Name: name})
	return m.CreateTopicError
}

func (m *MockKafkaClient) DeleteTopic(ctx context.Context, name string) error {
	m.DeleteTopicCalls = append(m.DeleteTopicCalls, DeleteTopicCall{Name: name})
	return m.DeleteTopicError
}

func (m *MockKafkaClient) Close() error {
	m.CloseCalls++
	return m.CloseError
}

// MockLoggerService implements logger.LoggerService interface for testing
type MockLoggerService struct {
	InfoCalls   []string
	DebugCalls  []string
	ErrorCalls  []string
	DebugfCalls []DebugfCall
	LogfCalls   []LogfCall
	LogCalls    []string
	ErrorfCalls []ErrorfCall
	SyncCalls   int
}

type DebugfCall struct {
	Format string
	Args   []any
}

type LogfCall struct {
	Format string
	Args   []any
}

type ErrorfCall struct {
	Format string
	Args   []any
}

func (m *MockLoggerService) Debugf(format string, args ...any) {
	m.DebugfCalls = append(m.DebugfCalls, DebugfCall{Format: format, Args: args})
}

func (m *MockLoggerService) Debug(msg string) {
	m.DebugCalls = append(m.DebugCalls, msg)
}

func (m *MockLoggerService) Logf(format string, args ...any) {
	m.LogfCalls = append(m.LogfCalls, LogfCall{Format: format, Args: args})
}

func (m *MockLoggerService) Log(msg string) {
	m.LogCalls = append(m.LogCalls, msg)
}

func (m *MockLoggerService) Info(msg string) {
	m.InfoCalls = append(m.InfoCalls, msg)
}

func (m *MockLoggerService) Errorf(format string, args ...any) {
	m.ErrorfCalls = append(m.ErrorfCalls, ErrorfCall{Format: format, Args: args})
}

func (m *MockLoggerService) Error(msg string) {
	m.ErrorCalls = append(m.ErrorCalls, msg)
}

func (m *MockLoggerService) Sync() error {
	m.SyncCalls++
	return nil
}

// MockCustomLoggerService implements logger.CustomLoggerService interface for testing
type MockCustomLoggerService struct {
	InitCalls       []logger.LogDto
	InfoCalls       []CustomLogCall
	DebugCalls      []CustomLogCall
	ErrorCalls      []CustomLogCall
	UpdateCalls     []UpdateCall
	GetLogDtoCalls  int
	FlushCalls      int
	SetSummaryCalls []logger.LogEventTag
	EndCalls        []EndCall
	AddFieldCalls   []AddFieldCall
}

type CustomLogCall struct {
	Action logger.LoggerAction
	Data   any
	Masks  []logger.MaskingOptionDto
}

type UpdateCall struct {
	Key   string
	Value any
}

type EndCall struct {
	Code        int
	Description string
}

type AddFieldCall struct {
	Key   string
	Value any
}

func (m *MockCustomLoggerService) Init(dto logger.LogDto) {
	m.InitCalls = append(m.InitCalls, dto)
}

func (m *MockCustomLoggerService) GetLogDto() logger.LogDto {
	m.GetLogDtoCalls++
	return logger.LogDto{}
}

func (m *MockCustomLoggerService) Update(key string, value any) {
	m.UpdateCalls = append(m.UpdateCalls, UpdateCall{Key: key, Value: value})
}

func (m *MockCustomLoggerService) Info(action logger.LoggerAction, data any, masks ...logger.MaskingOptionDto) {
	m.InfoCalls = append(m.InfoCalls, CustomLogCall{Action: action, Data: data, Masks: masks})
}

func (m *MockCustomLoggerService) Debug(action logger.LoggerAction, data any, masks ...logger.MaskingOptionDto) {
	m.DebugCalls = append(m.DebugCalls, CustomLogCall{Action: action, Data: data, Masks: masks})
}

func (m *MockCustomLoggerService) Error(action logger.LoggerAction, data any, masks ...logger.MaskingOptionDto) {
	m.ErrorCalls = append(m.ErrorCalls, CustomLogCall{Action: action, Data: data, Masks: masks})
}

func (m *MockCustomLoggerService) Flush() {
	m.FlushCalls++
}

func (m *MockCustomLoggerService) SetSummary(summary logger.LogEventTag) logger.CustomLoggerService {
	m.SetSummaryCalls = append(m.SetSummaryCalls, summary)
	return m
}

func (m *MockCustomLoggerService) End(code int, description string) {
	m.EndCalls = append(m.EndCalls, EndCall{Code: code, Description: description})
}

func (m *MockCustomLoggerService) AddField(key string, value any) {
	m.AddFieldCalls = append(m.AddFieldCalls, AddFieldCall{Key: key, Value: value})
}

// MockMaskingService implements logger.MaskingServiceInterface for testing
type MockMaskingService struct {
	MaskingCalls []MaskingCall
	MaskResult   string
}

type MaskingCall struct {
	Value       string
	MaskingType logger.MaskingType
}

func (m *MockMaskingService) Masking(value string, maskingType logger.MaskingType) string {
	m.MaskingCalls = append(m.MaskingCalls, MaskingCall{Value: value, MaskingType: maskingType})
	if m.MaskResult != "" {
		return m.MaskResult
	}
	return value
}

// Test helper functions for creating mock contexts
func NewMockRequestForTesting() *MockRequest {
	return &MockRequest{
		ParamsMap:         make(map[string]string),
		PathParamsMap:     make(map[string]string),
		QueryVals:         make(url.Values),
		BodyStr:           "",
		ErrBody:           nil,
		methodsToCallTime: make(map[string]uint),
		AddDataStr:        make(map[string]string),
	}
}

func CreateMockContextForTesting(t *testing.T) (*Context, *MockRequest, *httptest.ResponseRecorder, *MockKafkaClient, *MockLoggerService, *MockCustomLoggerService) {
	mockReq := NewMockRequestForTesting()
	mockReq.ParamsMap["topic"] = "test-topic"
	mockReq.PathParamsMap["id"] = "123"
	mockReq.QueryVals.Set("foo", "bar")
	mockReq.BodyStr = `{"key":"value"}`
	mockReq.AddDataStr["ClientIP"] = "192.168.1.1"
	mockReq.AddDataStr["UserAgent"] = "TestAgent"
	mockReq.AddDataStr["Method"] = "POST"
	mockReq.AddDataStr["URL"] = "http://localhost/test"
	mockReq.AddDataStr["SessionId"] = "test-session"
	mockReq.AddDataStr["RequestId"] = "test-request"

	recorder := httptest.NewRecorder()
	mockKafka := &MockKafkaClient{}
	mockAppLog := &MockLoggerService{}
	mockDetailLog := &MockLoggerService{}
	mockSummaryLog := &MockLoggerService{}
	mockMaskingService := &MockMaskingService{}
	mockCustomLog := &MockCustomLoggerService{}

	logService := LogService{
		appLog:         mockAppLog,
		detailLog:      mockDetailLog,
		summaryLog:     mockSummaryLog,
		maskingService: mockMaskingService,
	}

	// Create a basic config.Config for the newContext function
	cfg := &config.Config{}
	cfg.App.Name = "test-service"
	cfg.App.Version = "1.0.0"

	ctx := newContext(recorder, mockReq, mockKafka, logService, cfg)
	ctx.detail = mockCustomLog // Override with our mock

	return ctx, mockReq, recorder, mockKafka, mockAppLog, mockCustomLog
}

// Example test functions for Context methods
func TestContextInfo(t *testing.T) {
	ctx, _, _, _, mockAppLog, _ := CreateMockContextForTesting(t)

	// Test
	ctx.Info(map[string]any{"test": "data"})

	// Verify
	if len(mockAppLog.InfoCalls) != 1 {
		t.Errorf("Expected 1 Info call, got %d", len(mockAppLog.InfoCalls))
	}
}

func TestContextDebug(t *testing.T) {
	ctx, _, _, _, mockAppLog, _ := CreateMockContextForTesting(t)

	// Test
	ctx.Debug(map[string]any{"debug": "info"})

	// Verify
	if len(mockAppLog.DebugCalls) != 1 {
		t.Errorf("Expected 1 Debug call, got %d", len(mockAppLog.DebugCalls))
	}
}

func TestContextError(t *testing.T) {
	ctx, _, _, _, mockAppLog, _ := CreateMockContextForTesting(t)

	// Test
	ctx.Error(map[string]any{"error": "message"})

	// Verify
	if len(mockAppLog.ErrorCalls) != 1 {
		t.Errorf("Expected 1 Error call, got %d", len(mockAppLog.ErrorCalls))
	}
}

func TestContextJSON(t *testing.T) {
	ctx, _, recorder, _, _, mockCustomLog := CreateMockContextForTesting(t)

	// Test
	testData := map[string]any{"message": "success"}
	err := ctx.JSON(200, testData)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if recorder.Code != 200 {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}
	if recorder.Header().Get("Content-Type") != "application/json; charset=UTF8" {
		t.Errorf("Expected Content-Type 'application/json; charset=UTF8', got '%s'", recorder.Header().Get("Content-Type"))
	}
	if len(mockCustomLog.InfoCalls) != 1 {
		t.Errorf("Expected 1 Info call, got %d", len(mockCustomLog.InfoCalls))
	}
	if len(mockCustomLog.EndCalls) != 1 {
		t.Errorf("Expected 1 End call, got %d", len(mockCustomLog.EndCalls))
	}
}

func TestContextGetIncoming(t *testing.T) {
	ctx, _, _, _, _, _ := CreateMockContextForTesting(t)

	// Test
	incoming := ctx.GetIncoming()

	// Verify
	if incoming.URL != "http://localhost/test" {
		t.Errorf("Expected URL 'http://localhost/test', got '%s'", incoming.URL)
	}
	if incoming.IP != "192.168.1.1" {
		t.Errorf("Expected IP '192.168.1.1', got '%s'", incoming.IP)
	}
	if incoming.Method != "POST" {
		t.Errorf("Expected Method 'POST', got '%s'", incoming.Method)
	}
}

func TestContextLogAuto(t *testing.T) {
	ctx, _, _, _, _, mockCustomLog := CreateMockContextForTesting(t)

	// Test
	logService := ctx.LogAuto()

	// Verify
	if logService == nil {
		t.Error("Expected non-nil log service")
	}
	if len(mockCustomLog.InfoCalls) != 1 {
		t.Errorf("Expected 1 Info call, got %d", len(mockCustomLog.InfoCalls))
	}
}

func TestContextLog(t *testing.T) {
	ctx, _, _, _, _, mockCustomLog := CreateMockContextForTesting(t)

	// Test
	logService := ctx.Log()

	// Verify
	if logService != mockCustomLog {
		t.Error("Expected log service to match mock")
	}
}
