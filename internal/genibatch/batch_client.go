package genibatch

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/oauth2"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

type Client struct {
	client        *geni.Client
	unionRequests chan unionAsyncRequest
}

func NewClient(tokenSource oauth2.TokenSource, useSandboxEnv bool) *Client {
	return &Client{
		client:        geni.NewClient(tokenSource, useSandboxEnv),
		unionRequests: make(chan unionAsyncRequest),
	}
}

type unionAsyncRequest struct {
	UnionId  string
	Response chan *geni.UnionResponse
	Error    chan error
}

func (c *Client) GetUnion(ctx context.Context, unionId string) (*geni.UnionResponse, error) {
	response := make(chan *geni.UnionResponse)
	errors := make(chan error)

	c.unionRequests <- unionAsyncRequest{
		UnionId:  unionId,
		Response: response,
		Error:    errors,
	}

	select {
	case res := <-response:
		return res, nil
	case err := <-errors:
		tflog.Error(ctx, "Error processing request", map[string]interface{}{"error": err})
		return nil, err
	case <-ctx.Done():
		tflog.Error(ctx, "Context done", map[string]interface{}{"error": ctx.Err()})
		return nil, ctx.Err()
	}
}

const batchSize = 50

func (c *Client) UnionBulkProcessor(ctx context.Context) {
	batch := make([]unionAsyncRequest, 0, batchSize)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case req := <-c.unionRequests:
			batch = append(batch, req)
			if len(batch) >= batchSize {
				c.processBatchOfUnions(ctx, batch)
				batch = batch[:0] // Reset batch
			}
		case <-ticker.C:
			if len(batch) > 0 {
				c.processBatchOfUnions(ctx, batch)
				batch = batch[:0] // Reset batch
			}
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil {
				tflog.Error(ctx, "Context done", map[string]interface{}{"error": err})
			}
			return
		}
	}
}

func (c *Client) processBatchOfUnions(ctx context.Context, batch []unionAsyncRequest) {
	if len(batch) == 1 {
		req := batch[0]
		res, err := c.client.GetUnion(ctx, req.UnionId)
		if err != nil {
			req.Error <- err
			return
		}

		req.Response <- res
	}

	if len(batch) > 1 {
		unionIds := make([]string, len(batch))
		for i, req := range batch {
			unionIds[i] = req.UnionId
		}

		res, err := c.client.GetUnions(ctx, unionIds)
		if err != nil {
			for _, req := range batch {
				req.Error <- err
			}
			return
		}

		unionIdToResponse := make(map[string]*geni.UnionResponse)
		for _, resUnion := range res.Results {
			unionIdToResponse[resUnion.Id] = &resUnion
		}

		for _, req := range batch {
			if resUnion, ok := unionIdToResponse[req.UnionId]; ok {
				req.Response <- resUnion
			} else {
				req.Error <- fmt.Errorf("union %s not found in the response", req.UnionId)
			}
		}
	}
}
