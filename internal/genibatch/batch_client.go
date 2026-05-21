package genibatch

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/dmalch/go-geni"
	genidocument "github.com/dmalch/go-geni/document"
	geniprofile "github.com/dmalch/go-geni/profile"
	geniunion "github.com/dmalch/go-geni/union"
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
	Response chan *geniunion.Union
	Error    chan error
}

type profileAsyncRequest struct {
	Id       string
	Response chan *geniprofile.Profile
	Error    chan error
}

type documentAsyncRequest struct {
	Id       string
	Response chan *genidocument.Document
	Error    chan error
}

func (c *Client) GetUnion(ctx context.Context, id string) (*geniunion.Union, error) {
	response := make(chan *geniunion.Union)
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

func (c *Client) GetProfile(ctx context.Context, id string) (*geniprofile.Profile, error) {
	response := make(chan *geniprofile.Profile)
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

func (c *Client) GetDocument(ctx context.Context, id string) (*genidocument.Document, error) {
	response := make(chan *genidocument.Document)
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
		result, err := c.client.Union().Get(ctx, keys[0])
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
		res, err := c.client.Union().GetBulk(ctx, keys)
		if err != nil {
			for _, req := range batch {
				req.Error <- err
			}
			return
		}

		fulfillUnionRequests(batch, res.Results)
	}
}

// fulfillUnionRequests dispatches per-request results from a bulk union response.
// IDs absent from the bulk results are treated as not-found, because the Geni bulk
// endpoint silently omits missing IDs from its response — that absence is the domain
// signal that the union no longer exists.
func fulfillUnionRequests(batch []unionAsyncRequest, results []geniunion.Union) {
	idToResponse := make(map[string]*geniunion.Union, len(results))
	for i := range results {
		idToResponse[results[i].ID] = &results[i]
	}

	for _, req := range batch {
		if result, ok := idToResponse[req.Id]; ok {
			req.Response <- result
		} else {
			req.Error <- fmt.Errorf("union %s not found in the response: %w", req.Id, geni.ErrResourceNotFound)
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
		result, err := c.client.Profile().Get(ctx, keys[0])
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
		res, err := c.client.Profile().GetBulk(ctx, keys)
		if err != nil {
			for _, req := range batch {
				req.Error <- err
			}
			return
		}

		fulfillProfileRequests(batch, res.Results)
	}
}

// fulfillProfileRequests dispatches per-request results from a bulk profile response.
// IDs absent from the bulk results are treated as not-found, because the Geni bulk
// endpoint silently omits missing IDs from its response.
func fulfillProfileRequests(batch []profileAsyncRequest, results []geniprofile.Profile) {
	idToResponse := make(map[string]*geniprofile.Profile, len(results))
	for i := range results {
		idToResponse[results[i].ID] = &results[i]
	}

	for _, req := range batch {
		if result, ok := idToResponse[req.Id]; ok {
			req.Response <- result
		} else {
			req.Error <- fmt.Errorf("profile %s not found in the response: %w", req.Id, geni.ErrResourceNotFound)
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
		result, err := c.client.Document().Get(ctx, keys[0])
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
		res, err := c.client.Document().GetBulk(ctx, keys)
		if err != nil {
			for _, req := range batch {
				req.Error <- err
			}
			return
		}

		fulfillDocumentRequests(batch, res.Results)
	}
}

// fulfillDocumentRequests dispatches per-request results from a bulk document response.
// IDs absent from the bulk results are treated as not-found, because the Geni bulk
// endpoint silently omits missing IDs from its response.
func fulfillDocumentRequests(batch []documentAsyncRequest, results []genidocument.Document) {
	idToResponse := make(map[string]*genidocument.Document, len(results))
	for i := range results {
		idToResponse[results[i].ID] = &results[i]
	}

	for _, req := range batch {
		if result, ok := idToResponse[req.Id]; ok {
			req.Response <- result
		} else {
			req.Error <- fmt.Errorf("document %s not found in the response: %w", req.Id, geni.ErrResourceNotFound)
		}
	}
}
