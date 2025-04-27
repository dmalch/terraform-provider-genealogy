package geni

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	Results    []DocumentResponse `json:"results,omitempty"`
	Page       int                `json:"page,omitempty"`
	TotalCount int                `json:"total_count,omitempty"`
}
type DocumentResponse struct {
	// Id is the document's id
	Id string `json:"id,omitempty"`
	// Title is the document's title
	Title string `json:"title,omitempty"`
	// Description is the document's description
	Description *string `json:"description"`
	// SourceUrl is the document's source URL
	SourceUrl *string `json:"source_url"`
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
		tflog.Error(ctx, "Error marshaling request", map[string]interface{}{"error": err})
		return nil, err
	}

	jsonStr := strings.ReplaceAll(string(jsonBody), "\\\\", "\\")
	jsonStr = escapeString(jsonStr)

	url := BaseUrl(c.useSandboxEnv) + "api/document/add"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(jsonStr))
	if err != nil {
		tflog.Error(ctx, "Error marshaling request", map[string]interface{}{"error": err})
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var document DocumentResponse
	err = json.Unmarshal(body, &document)
	if err != nil {
		tflog.Error(ctx, "Error unmarshaling response", map[string]interface{}{"error": err})
		return nil, err
	}

	return &document, nil
}

func (c *Client) GetDocument(ctx context.Context, documentId string) (*DocumentResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + documentId
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		tflog.Error(ctx, "Error marshaling request", map[string]interface{}{"error": err})
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var document DocumentResponse
	err = json.Unmarshal(body, &document)
	if err != nil {
		tflog.Error(ctx, "Error unmarshaling response", map[string]interface{}{"error": err})
		return nil, err
	}

	return &document, nil
}

func (c *Client) GetDocuments(ctx context.Context, documentIds []string) (*DocumentBulkResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/document"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		tflog.Error(ctx, "Error marshaling request", map[string]interface{}{"error": err})
		return nil, err
	}

	query := req.URL.Query()
	query.Add("ids", strings.Join(documentIds, ","))
	req.URL.RawQuery = query.Encode()

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var document DocumentBulkResponse
	err = json.Unmarshal(body, &document)
	if err != nil {
		tflog.Error(ctx, "Error unmarshaling response", map[string]interface{}{"error": err})
		return nil, err
	}

	return &document, nil
}

func (c *Client) GetUploadedDocuments(ctx context.Context, page int) (*DocumentBulkResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/user/uploaded-documents"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		tflog.Error(ctx, "Error marshaling request", map[string]interface{}{"error": err})
		return nil, err
	}

	query := req.URL.Query()
	query.Add("page", strconv.Itoa(page))
	req.URL.RawQuery = query.Encode()

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var document DocumentBulkResponse
	err = json.Unmarshal(body, &document)
	if err != nil {
		tflog.Error(ctx, "Error unmarshaling response", map[string]interface{}{"error": err})
		return nil, err
	}

	return &document, nil
}

func (c *Client) DeleteDocument(ctx context.Context, documentId string) error {
	url := BaseUrl(c.useSandboxEnv) + "api/" + documentId + "/delete"
	req, err := http.NewRequest(http.MethodPost, url, nil)

	if err != nil {
		tflog.Error(ctx, "Error marshaling request", map[string]interface{}{"error": err})
		return err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}

	var result ResultResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		tflog.Error(ctx, "Error unmarshaling response", map[string]interface{}{"error": err})
		return err
	}

	return nil
}

func (c *Client) UpdateDocument(ctx context.Context, documentId string, request *DocumentRequest) (*DocumentResponse, error) {
	jsonBody, err := json.Marshal(request)
	if err != nil {
		tflog.Error(ctx, "Error marshaling request", map[string]interface{}{"error": err})
		return nil, err
	}

	jsonStr := strings.ReplaceAll(string(jsonBody), "\\\\", "\\")
	jsonStr = escapeString(jsonStr)

	url := BaseUrl(c.useSandboxEnv) + "api/" + documentId + "/update"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(jsonStr))
	if err != nil {
		tflog.Error(ctx, "Error marshaling request", map[string]interface{}{"error": err})
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var document DocumentResponse
	err = json.Unmarshal(body, &document)
	if err != nil {
		tflog.Error(ctx, "Error unmarshaling response", map[string]interface{}{"error": err})
		return nil, err
	}

	return &document, nil
}

func (c *Client) TagDocument(ctx context.Context, documentId, profileId string) (*ProfileBulkResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + documentId + "/tag/" + profileId

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		tflog.Error(ctx, "Error marshaling request", map[string]interface{}{"error": err})
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var profiles ProfileBulkResponse
	err = json.Unmarshal(body, &profiles)
	if err != nil {
		tflog.Error(ctx, "Error unmarshaling response", map[string]interface{}{"error": err})
		return nil, err
	}

	return &profiles, nil
}

func (c *Client) UntagDocument(ctx context.Context, documentId, profileId string) (*ProfileBulkResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + documentId + "/untag/" + profileId

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		tflog.Error(ctx, "Error marshaling request", map[string]interface{}{"error": err})
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var profiles ProfileBulkResponse
	err = json.Unmarshal(body, &profiles)
	if err != nil {
		tflog.Error(ctx, "Error unmarshaling response", map[string]interface{}{"error": err})
		return nil, err
	}

	return &profiles, nil
}
