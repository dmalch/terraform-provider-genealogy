package genibatch

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

type Client struct {
	client          *geni.Client
	unionRequests   chan unionAsyncRequest
	profileRequests chan profileAsyncRequest
}

func NewClient(client *geni.Client) *Client {
	return &Client{
		client:          client,
		unionRequests:   make(chan unionAsyncRequest),
		profileRequests: make(chan profileAsyncRequest),
	}
}

type unionAsyncRequest struct {
	UnionId  string
	Response chan *geni.UnionResponse
	Error    chan error
}

type profileAsyncRequest struct {
	ProfileId string
	Response  chan *geni.ProfileResponse
	Error     chan error
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

func (c *Client) GetProfile(ctx context.Context, profileId string) (*geni.ProfileResponse, error) {
	response := make(chan *geni.ProfileResponse)
	errors := make(chan error)

	c.profileRequests <- profileAsyncRequest{
		ProfileId: profileId,
		Response:  response,
		Error:     errors,
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

func (c *Client) ProfileBulkProcessor(ctx context.Context) {
	batch := make([]profileAsyncRequest, 0, batchSize)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case req := <-c.profileRequests:
			batch = append(batch, req)
			if len(batch) >= batchSize {
				c.processBatchOfProfiles(ctx, batch)
				batch = batch[:0] // Reset batch
			}
		case <-ticker.C:
			if len(batch) > 0 {
				c.processBatchOfProfiles(ctx, batch)
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

func (c *Client) processBatchOfProfiles(ctx context.Context, batch []profileAsyncRequest) {
	if len(batch) == 1 {
		req := batch[0]
		res, err := c.client.GetProfile(ctx, req.ProfileId)
		if err != nil {
			req.Error <- err
			return
		}

		req.Response <- res
	}

	if len(batch) > 1 {
		profileIds := make([]string, len(batch))
		for i, req := range batch {
			profileIds[i] = req.ProfileId
		}

		res, err := c.client.GetProfiles(ctx, profileIds)
		if err != nil {
			for _, req := range batch {
				req.Error <- err
			}
			return
		}

		profileIdToResponse := make(map[string]*geni.ProfileResponse)
		for _, resProfile := range res.Results {
			profileIdToResponse[resProfile.Id] = &resProfile
		}

		for _, req := range batch {
			if resUnion, ok := profileIdToResponse[req.ProfileId]; ok {
				req.Response <- resUnion
			} else {
				req.Error <- fmt.Errorf("profile %s not found in the response", req.ProfileId)
			}
		}
	}
}
