package ratelimitmock

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrBlocked reports if service is blocked.
var ErrBlocked = errors.New("blocked")

// Service defines external service that can process batches of items.
type Service interface {
	GetLimits() (n uint64, p time.Duration)
	Process(ctx context.Context, batch Batch) error
}

// Batch is a batch of items.
type Batch []Item

// Item is some abstract item.
type Item struct{}

type Client struct {
	service Service
	n       uint64
	p       time.Duration
	batchCh chan Batch
}

func NewClient(service Service, queuelimit int) *Client {
	n, p := service.GetLimits()
	client := &Client{
		service: service,
		n:       n,
		p:       p,
		batchCh: make(chan Batch, queuelimit),
	}
	go client.processBatches()
	return client
}

func (c *Client) SendBatch(batch Batch) {
	c.batchCh <- batch
}

func (c *Client) processBatches() {
	timer := time.NewTimer(c.p)
	for batch := range c.batchCh {
		<-timer.C
		c.processBatch(batch)
		timer.Reset(c.p)
	}
}

func (c *Client) processBatch(batch Batch) {
	var j int
	for i := 0; i < len(batch); i += int(c.n) {
		j = i + int(c.n)
		if j > len(batch) {
			j = len(batch) // beyond the limit
		}

		if err := c.service.Process(context.Background(), batch[i:j]); err != nil {
			fmt.Println("Error processing batch:", err)
		}
	}
}
