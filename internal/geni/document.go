package geni

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"
)

type DocumentRequest struct {
	// Title is the document's title
	Title string `json:"title,omitempty"`
	// Description is the document's description
	Description *string `json:"description,omitempty"`
	// ContentType is the document's content type
	ContentType *string `json:"content_type,omitempty"`
	// Date is the document's date
	Date *DateElement `json:"date,omitempty"`
	// Location is the document's location
	Location *LocationElement `json:"location,omitempty"`
	// Labels is the document's comma separated labels
	Labels *string `json:"labels,omitempty"`
	// File is the Base64 encoded file to create a document from
	File *string `json:"file,omitempty"`
	// FileName is the name of the file, required if the file is provided
	FileName *string `json:"file_name,omitempty"`
	// SourceUrl is the source URL for the document
	SourceUrl *string `json:"source_url,omitempty"`
	// Text is the text to create a document from
	Text *string `json:"text,omitempty"`
}

type DocumentBulkResponse struct {
	Results []DocumentResponse `json:"results,omitempty"`
}
type DocumentResponse struct {
	// Id is the document's id
	Id string `json:"id,omitempty"`
	// Title is the document's title
	Title string `json:"title,omitempty"`
	// Description is the document's description
	Description *string `json:"description"`
	// ContentType is the document's content type
	ContentType *string `json:"content_type"`
	// Date is the document's date
	Date *DateElement `json:"date"`
	// Location is the document's location
	Location *LocationElement `json:"location,omitempty"`
	// Profiles is the list of profiles tagged in the document
	Tags []string `json:"tags"`
	// Labels is the list of labels associated with the document
	Labels []string `json:"labels"`
	// UpdatedAt is the timestamp of when the document was last updated
	UpdatedAt string `json:"updated_at"`
	// CreatedAt is the timestamp of when the document was created
	CreatedAt string `json:"created_at"`
}

func (c *Client) CreateDocument(ctx context.Context, request *DocumentRequest) (*DocumentResponse, error) {
	jsonBody, err := json.Marshal(request)
	if err != nil {
		slog.Error("Error marshaling request", "error", err)
		return nil, err
	}

	jsonStr := strings.ReplaceAll(string(jsonBody), "\\\\", "\\")
	jsonStr = escapeString(jsonStr)

	url := BaseUrl(c.useSandboxEnv) + "api/document/add"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(jsonStr))
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var document DocumentResponse
	err = json.Unmarshal(body, &document)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &document, nil
}

func (c *Client) GetDocument(ctx context.Context, documentId string) (*DocumentResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + documentId
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	body, err := c.doRequest(ctx, req,
		WithRequestKey(func() string {
			return documentId
		}),
		WithPrepareBulkRequest(func(req *http.Request, urlMap *sync.Map) {
			// Add a new ids parameter containing IDs of all documents to be fetched in
			// addition to the current one. First, we need to get the IDs from the map.
			ids := make([]string, 0)

			ids = append(ids, documentId)

			urlMap.Range(func(key, value interface{}) bool {
				if _, ok := value.(context.CancelFunc); ok {
					if keyString, ok := key.(string); ok {
						ids = append(ids, keyString)
					}
				}
				return true
			})

			if len(ids) > 1 {
				query := req.URL.Query()
				query.Add("ids", strings.Join(ids, ","))
				req.URL.RawQuery = query.Encode()
			}
		}),
		WithParseBulkResponse(func(req *http.Request, body []byte, urlMap *sync.Map) ([]byte, error) {
			// If only one document is requested, we can skip the bulk response parsing
			if !req.URL.Query().Has("ids") {
				return body, nil
			}

			// Parse the response to get the document ID
			var response DocumentBulkResponse
			err := json.Unmarshal(body, &response)
			if err != nil {
				slog.Error("Error unmarshaling bulk response", "error", err)
				return nil, err
			}

			var requestedRes []byte

			// Store the response in the map using the document ID as the key
			for _, document := range response.Results {

				jsonBody, err := json.Marshal(&document)
				if err != nil {
					slog.Error("Error marshaling request", "error", err)
					return nil, err
				}

				if document.Id == documentId {
					requestedRes = jsonBody
					continue
				}

				previous, loaded := urlMap.Swap(document.Id, jsonBody)
				if loaded {
					// If the previous value is context cancel function, cancel it
					if cancelFunc, ok := previous.(context.CancelFunc); ok {
						cancelFunc()
					}
				}
			}

			return requestedRes, nil
		}))
	if err != nil {
		return nil, err
	}

	var document DocumentResponse
	err = json.Unmarshal(body, &document)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &document, nil
}

func (c *Client) DeleteDocument(ctx context.Context, documentId string) error {
	url := BaseUrl(c.useSandboxEnv) + "api/" + documentId + "/delete"
	req, err := http.NewRequest(http.MethodPost, url, nil)

	if err != nil {
		slog.Error("Error creating request", "error", err)
		return err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}

	var result ResultResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return err
	}

	return nil
}

func (c *Client) UpdateDocument(ctx context.Context, profileId string, request *DocumentRequest) (*DocumentResponse, error) {
	jsonBody, err := json.Marshal(request)
	if err != nil {
		slog.Error("Error marshaling request", "error", err)
		return nil, err
	}

	jsonStr := strings.ReplaceAll(string(jsonBody), "\\\\", "\\")
	jsonStr = escapeString(jsonStr)

	url := BaseUrl(c.useSandboxEnv) + "api/" + profileId + "/update"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(jsonStr))
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var document DocumentResponse
	err = json.Unmarshal(body, &document)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &document, nil
}
