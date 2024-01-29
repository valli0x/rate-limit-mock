package ratelimitmock

import (
	"context"
	"testing"
	"time"
)

var (
	queuelimit = 10
)

type MockService struct{}

func (ms *MockService) GetLimits() (n uint64, p time.Duration) {
	return 10, time.Second
}

func (ms *MockService) Process(ctx context.Context, batch Batch) error {
	return nil
}

func TestNewClient(t *testing.T) {
	service := &MockService{}
	client := NewClient(service, 10)
	if client.n != 10 || client.p != time.Second {
		t.Errorf("Expected n=10 and p=1s, got n=%d and p=%s", client.n, client.p)
	}
}

func TestClient_SendBatch(t *testing.T) {
	service := &MockService{}
	client := NewClient(service, 10)

	batch := make(Batch, 20)
	go client.SendBatch(batch)
	batchFromChannel := <-client.batchCh

	if len(batchFromChannel) != 20 {
		t.Errorf("Expected batch size 20, got: %d", len(batchFromChannel))
	}
}
