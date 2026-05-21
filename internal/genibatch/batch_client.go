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
	unionRequests    chan asyncRequest[geniunion.Union]
	profileRequests  chan asyncRequest[geniprofile.Profile]
	documentRequests chan asyncRequest[genidocument.Document]
}

func NewClient(client *geni.Client) *Client {
	return &Client{
		client:           client,
		unionRequests:    make(chan asyncRequest[geniunion.Union]),
		profileRequests:  make(chan asyncRequest[geniprofile.Profile]),
		documentRequests: make(chan asyncRequest[genidocument.Document]),
	}
}

// asyncRequest is a single batched read awaiting fulfilment by a batch worker.
// Exactly one value is ever delivered, on Response or on Error.
type asyncRequest[T any] struct {
	Id       string
	Response chan *T
	Error    chan error
}

func (c *Client) GetUnion(ctx context.Context, id string) (*geniunion.Union, error) {
	response := make(chan *geniunion.Union)
	errors := make(chan error)

	c.unionRequests <- asyncRequest[geniunion.Union]{
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

	c.profileRequests <- asyncRequest[geniprofile.Profile]{
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

	c.documentRequests <- asyncRequest[genidocument.Document]{
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
	batch := make([]asyncRequest[geniunion.Union], 0, batchSize)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case req := <-c.unionRequests:
			batch = append(batch, req)
			if len(batch) >= batchSize {
				// copy the batch to a new slice
				requests := make([]asyncRequest[geniunion.Union], len(batch))
				copy(requests, batch)
				batch = batch[:0] // Reset batch

				go c.processBatchOfUnions(ctx, requests)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				// copy the batch to a new slice
				requests := make([]asyncRequest[geniunion.Union], len(batch))
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
	batch := make([]asyncRequest[geniprofile.Profile], 0, batchSize)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case req := <-c.profileRequests:
			batch = append(batch, req)
			if len(batch) >= batchSize {
				// copy the batch to a new slice
				requests := make([]asyncRequest[geniprofile.Profile], len(batch))
				copy(requests, batch)
				batch = batch[:0] // Reset batch

				go c.processBatchOfProfiles(ctx, requests)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				// copy the batch to a new slice
				requests := make([]asyncRequest[geniprofile.Profile], len(batch))
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
	batch := make([]asyncRequest[genidocument.Document], 0, batchSize)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case req := <-c.documentRequests:
			batch = append(batch, req)
			if len(batch) >= batchSize {
				// copy the batch to a new slice
				requests := make([]asyncRequest[genidocument.Document], len(batch))
				copy(requests, batch)
				batch = batch[:0] // Reset batch

				go c.processBatchOfDocuments(ctx, requests)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				// copy the batch to a new slice
				requests := make([]asyncRequest[genidocument.Document], len(batch))
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

// recoverBatch converts a panic in a batch worker into an error delivered to
// every request in the batch. Batch workers run in their own goroutines, so an
// unrecovered panic would crash the entire provider process. Because a worker
// only panics before it has answered any request (the panic-prone work — the
// API call — happens up front, before any channel send), broadcasting the error
// also unblocks every caller that would otherwise wait forever on its response
// channel. The send is abandoned if ctx is cancelled so the worker cannot leak.
func recoverBatch[T any](ctx context.Context, kind string, batch []asyncRequest[T]) {
	r := recover()
	if r == nil {
		return
	}

	err := fmt.Errorf("recovered from panic while processing %s batch: %v", kind, r)
	tflog.Error(ctx, "Panic in batch processor", map[string]interface{}{"error": err})

	for _, req := range batch {
		select {
		case req.Error <- err:
		case <-ctx.Done():
		}
	}
}

func (c *Client) processBatchOfUnions(ctx context.Context, batch []asyncRequest[geniunion.Union]) {
	defer recoverBatch(ctx, "union", batch)

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
func fulfillUnionRequests(batch []asyncRequest[geniunion.Union], results []geniunion.Union) {
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

func (c *Client) processBatchOfProfiles(ctx context.Context, batch []asyncRequest[geniprofile.Profile]) {
	defer recoverBatch(ctx, "profile", batch)

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
func fulfillProfileRequests(batch []asyncRequest[geniprofile.Profile], results []geniprofile.Profile) {
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

func (c *Client) processBatchOfDocuments(ctx context.Context, batch []asyncRequest[genidocument.Document]) {
	defer recoverBatch(ctx, "document", batch)

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
func fulfillDocumentRequests(batch []asyncRequest[genidocument.Document], results []genidocument.Document) {
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
