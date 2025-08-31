package main

import (
	"testing"

	"github.com/sing3demons/go-common-kp/kp/pkg/kp"
)

func TestHealthzHandler(t *testing.T) {
	// Create a mock context using the helper function from the kp package
	ctx, _, recorder, mockKafka, mockAppLog, mockCustomLog := kp.CreateMockContextForTesting(t)

	// Test the healthz handler
	err := healthzHandler(ctx)

	// Verify the result
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify HTTP response
	if recorder.Code != 200 {
		t.Errorf("Expected status code 200, got %d", recorder.Code)
	}

	// Verify logging was called
	if len(mockAppLog.InfoCalls) != 1 {
		t.Errorf("Expected 1 Info call, got %d", len(mockAppLog.InfoCalls))
	}

	// Verify Kafka publish was called
	if len(mockKafka.PublishCalls) != 1 {
		t.Errorf("Expected 1 Publish call, got %d", len(mockKafka.PublishCalls))
	}

	// Verify the Kafka topic
	if mockKafka.PublishCalls[0].Topic != "Stg-BulkNormal-CorrelatorTx" {
		t.Errorf("Expected topic 'Stg-BulkNormal-CorrelatorTx', got '%s'", mockKafka.PublishCalls[0].Topic)
	}

	// Verify custom logger was called for JSON response
	if len(mockCustomLog.InfoCalls) != 1 {
		t.Errorf("Expected 1 custom log Info call, got %d", len(mockCustomLog.InfoCalls))
	}

	// Verify custom logger End was called
	if len(mockCustomLog.EndCalls) != 1 {
		t.Errorf("Expected 1 custom log End call, got %d", len(mockCustomLog.EndCalls))
	}
}
