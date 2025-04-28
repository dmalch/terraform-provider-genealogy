package genibatch

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

type Client struct {
	client           *geni.Client
	unionRequests    chan unionAsyncRequest
	profileRequests  chan profileAsyncRequest
	documentRequests chan documentAsyncRequest
}

func NewClient(client *geni.Client) *Client {
	return &Client{
		client:           client,
		unionRequests:    make(chan unionAsyncRequest),
		profileRequests:  make(chan profileAsyncRequest),
		documentRequests: make(chan documentAsyncRequest),
	}
}

type unionAsyncRequest struct {
	Id       string
	Response chan *geni.UnionResponse
	Error    chan error
}

type profileAsyncRequest struct {
	Id       string
	Response chan *geni.ProfileResponse
	Error    chan error
}

type documentAsyncRequest struct {
	Id       string
	Response chan *geni.DocumentResponse
	Error    chan error
}

func (c *Client) GetUnion(ctx context.Context, id string) (*geni.UnionResponse, error) {
	response := make(chan *geni.UnionResponse)
	errors := make(chan error)

	c.unionRequests <- unionAsyncRequest{
		Id:       id,
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

func (c *Client) GetProfile(ctx context.Context, id string) (*geni.ProfileResponse, error) {
	response := make(chan *geni.ProfileResponse)
	errors := make(chan error)

	c.profileRequests <- profileAsyncRequest{
		Id:       id,
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

func (c *Client) GetDocument(ctx context.Context, id string) (*geni.DocumentResponse, error) {
	response := make(chan *geni.DocumentResponse)
	errors := make(chan error)

	c.documentRequests <- documentAsyncRequest{
		Id:       id,
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
				// copy the batch to a new slice
				requests := make([]unionAsyncRequest, len(batch))
				copy(requests, batch)
				batch = batch[:0] // Reset batch

				go c.processBatchOfUnions(ctx, requests)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				// copy the batch to a new slice
				requests := make([]unionAsyncRequest, len(batch))
				copy(requests, batch)
				batch = batch[:0] // Reset batch

				go c.processBatchOfUnions(ctx, requests)
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

func (c *Client) ProfileBulkProcessor(ctx context.Context) {
	batch := make([]profileAsyncRequest, 0, batchSize)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case req := <-c.profileRequests:
			batch = append(batch, req)
			if len(batch) >= batchSize {
				// copy the batch to a new slice
				requests := make([]profileAsyncRequest, len(batch))
				copy(requests, batch)
				batch = batch[:0] // Reset batch

				go c.processBatchOfProfiles(ctx, requests)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				// copy the batch to a new slice
				requests := make([]profileAsyncRequest, len(batch))
				copy(requests, batch)
				batch = batch[:0] // Reset batch

				go c.processBatchOfProfiles(ctx, requests)
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

func (c *Client) DocumentBulkProcessor(ctx context.Context) {
	batch := make([]documentAsyncRequest, 0, batchSize)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case req := <-c.documentRequests:
			batch = append(batch, req)
			if len(batch) >= batchSize {
				// copy the batch to a new slice
				requests := make([]documentAsyncRequest, len(batch))
				copy(requests, batch)
				batch = batch[:0] // Reset batch

				go c.processBatchOfDocuments(ctx, requests)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				// copy the batch to a new slice
				requests := make([]documentAsyncRequest, len(batch))
				copy(requests, batch)
				batch = batch[:0] // Reset batch

				go c.processBatchOfDocuments(ctx, requests)
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
	// Create a hashset to store unique IDs
	ids := make(map[string]struct{}, len(batch))
	for _, req := range batch {
		ids[req.Id] = struct{}{}
	}

	// Get keys from the hashset as a slice
	keys := make([]string, 0, len(ids))
	for id := range ids {
		keys = append(keys, id)
	}

	if len(keys) == 1 {
		result, err := c.client.GetUnion(ctx, keys[0])
		if err != nil {
			for _, req := range batch {
				req.Error <- err
			}
			return
		}

		for _, req := range batch {
			req.Response <- result
		}
	}

	if len(keys) > 1 {
		res, err := c.client.GetUnions(ctx, keys)
		if err != nil {
			for _, req := range batch {
				req.Error <- err
			}
			return
		}

		idToResponse := make(map[string]*geni.UnionResponse)
		for _, result := range res.Results {
			idToResponse[result.Id] = &result
		}

		for _, req := range batch {
			if result, ok := idToResponse[req.Id]; ok {
				req.Response <- result
			} else {
				req.Error <- fmt.Errorf("union %s not found in the response", req.Id)
			}
		}
	}
}

func (c *Client) processBatchOfProfiles(ctx context.Context, batch []profileAsyncRequest) {
	// Create a hashset to store unique IDs
	ids := make(map[string]struct{}, len(batch))
	for _, req := range batch {
		ids[req.Id] = struct{}{}
	}

	// Get keys from the hashset as a slice
	keys := make([]string, 0, len(ids))
	for id := range ids {
		keys = append(keys, id)
	}

	if len(keys) == 1 {
		result, err := c.client.GetProfile(ctx, keys[0])
		if err != nil {
			for _, req := range batch {
				req.Error <- err
			}
			return
		}

		for _, req := range batch {
			req.Response <- result
		}
	}

	if len(keys) > 1 {
		res, err := c.client.GetProfiles(ctx, keys)
		if err != nil {
			for _, req := range batch {
				req.Error <- err
			}
			return
		}

		idToResponse := make(map[string]*geni.ProfileResponse)
		for _, result := range res.Results {
			idToResponse[result.Id] = &result
		}

		for _, req := range batch {
			if result, ok := idToResponse[req.Id]; ok {
				req.Response <- result
			} else {
				req.Error <- fmt.Errorf("profile %s not found in the response", req.Id)
			}
		}
	}
}

func (c *Client) processBatchOfDocuments(ctx context.Context, batch []documentAsyncRequest) {
	// Create a hashset to store unique IDs
	ids := make(map[string]struct{}, len(batch))
	for _, req := range batch {
		ids[req.Id] = struct{}{}
	}

	// Get keys from the hashset as a slice
	keys := make([]string, 0, len(ids))
	for id := range ids {
		keys = append(keys, id)
	}

	if len(keys) == 1 {
		result, err := c.client.GetDocument(ctx, keys[0])
		if err != nil {
			for _, req := range batch {
				req.Error <- err
			}
			return
		}

		for _, req := range batch {
			req.Response <- result
		}
	}

	if len(keys) > 1 {
		res, err := c.client.GetDocuments(ctx, keys)
		if err != nil {
			for _, req := range batch {
				req.Error <- err
			}
			return
		}

		idToResponse := make(map[string]*geni.DocumentResponse)
		for _, result := range res.Results {
			idToResponse[result.Id] = &result
		}

		for _, req := range batch {
			if result, ok := idToResponse[req.Id]; ok {
				req.Response <- result
			} else {
				req.Error <- fmt.Errorf("document %s not found in the response", req.Id)
			}
		}
	}
}
