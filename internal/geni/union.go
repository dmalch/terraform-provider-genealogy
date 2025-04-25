package geni

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type UnionRequest struct {
	// Marriage date and location
	Marriage *EventElement `json:"marriage,omitempty"`
	// Divorce date and location
	Divorce *EventElement `json:"divorce,omitempty"`
}

type UnionBulkResponse struct {
	Results []UnionResponse `json:"results,omitempty"`
}

type UnionResponse struct {
	// The union's id
	Id string `json:"id,omitempty"`
	// AdoptedChildren is a subset of the children array, indicating which children are adopted
	AdoptedChildren []string `json:"adopted_children,omitempty"`
	// Children is an array of children in the union (urls or ids, if requested)
	Children []string `json:"children,omitempty"`
	// FosterChildren is a subset of the children array, indicating which children are foster
	FosterChildren []string `json:"foster_children,omitempty"`
	// Partners is an array of partners in the union (urls or ids, if requested)
	Partners []string `json:"partners,omitempty"`
	// Marriage date and location
	Marriage *EventElement `json:"marriage,omitempty"`
	// Divorce date and location
	Divorce *EventElement `json:"divorce,omitempty"`
	// Status of the union (spouse|ex_spouse)
	Status string `json:"status,omitempty"`
}

func (c *Client) GetUnion(ctx context.Context, unionId string) (*UnionResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + unionId
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	//WithRequestKey(func() string {
	//	return unionId
	//}),
	//WithPrepareBulkRequest(func(req *http.Request, urlMap *sync.Map) {
	//	// Add a new ids parameter containing IDs of all unions to be fetched in
	//	// addition to the current one. First, we need to get the IDs from the map.
	//	ids := make([]string, 0)
	//
	//	ids = append(ids, unionId)
	//
	//	urlMap.Range(func(key, value interface{}) bool {
	//		if _, ok := value.(context.CancelFunc); ok {
	//			if keyString, ok := key.(string); ok && strings.Contains(keyString, "union") {
	//				ids = append(ids, keyString)
	//			}
	//		}
	//		return true
	//	})
	//
	//	if len(ids) > 1 {
	//		query := req.URL.Query()
	//		query.Add("ids", strings.Join(ids, ","))
	//		req.URL.RawQuery = query.Encode()
	//	}
	//}),
	//WithParseBulkResponse(func(req *http.Request, body []byte, urlMap *sync.Map) ([]byte, error) {
	//	// If only one union is requested, we can skip the bulk response parsing
	//	if !req.URL.Query().Has("ids") {
	//		return body, nil
	//	}
	//
	//	// Parse the response to get the union ID
	//	var response UnionBulkResponse
	//	err := json.Unmarshal(body, &response)
	//	if err != nil {
	//		slog.Error("Error unmarshaling bulk response", "error", err)
	//		return nil, err
	//	}
	//
	//	var requestedUnionRes []byte
	//
	//	// Store the response in the map using the union ID as the key
	//	for _, union := range response.Results {
	//
	//		jsonBody, err := json.Marshal(&union)
	//		if err != nil {
	//			slog.Error("Error marshaling request", "error", err)
	//			return nil, err
	//		}
	//
	//		if union.Id == unionId {
	//			requestedUnionRes = jsonBody
	//			continue
	//		}
	//
	//		previous, loaded := urlMap.Swap(union.Id, jsonBody)
	//		if loaded {
	//			// If the previous value is context cancel function, cancel it
	//			if cancelFunc, ok := previous.(context.CancelFunc); ok {
	//				cancelFunc()
	//			}
	//		}
	//	}
	//
	//	return requestedUnionRes, nil
	//}))
	if err != nil {
		return nil, err
	}

	var union UnionResponse
	err = json.Unmarshal(body, &union)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &union, nil
}

func (c *Client) GetUnions(ctx context.Context, unionIds []string) (*UnionBulkResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/union"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	query := req.URL.Query()
	query.Add("ids", strings.Join(unionIds, ","))
	req.URL.RawQuery = query.Encode()

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var union UnionBulkResponse
	err = json.Unmarshal(body, &union)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &union, nil
}

func (c *Client) GetUnionAsync(ctx context.Context, unionId string) (*UnionResponse, error) {
	response := make(chan *UnionResponse)
	errors := make(chan error)

	c.unionRequests <- UnionAsyncRequest{
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
	batch := make([]UnionAsyncRequest, 0, batchSize)
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

func (c *Client) processBatchOfUnions(ctx context.Context, batch []UnionAsyncRequest) {
	if len(batch) == 1 {
		req := batch[0]
		res, err := c.GetUnion(ctx, req.UnionId)
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

		res, err := c.GetUnions(ctx, unionIds)
		if err != nil {
			for _, req := range batch {
				req.Error <- err
			}
			return
		}

		unionIdToResponse := make(map[string]*UnionResponse)
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

func (c *Client) UpdateUnion(ctx context.Context, unionId string, request *UnionRequest) (*UnionResponse, error) {
	jsonBody, err := json.Marshal(request)
	if err != nil {
		slog.Error("Error marshaling request", "error", err)
		return nil, err
	}

	jsonStr := strings.ReplaceAll(string(jsonBody), "\\\\", "\\")
	jsonStr = escapeString(jsonStr)

	url := BaseUrl(c.useSandboxEnv) + "api/" + unionId + "/update"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(jsonStr))
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var union UnionResponse
	err = json.Unmarshal(body, &union)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &union, nil
}
