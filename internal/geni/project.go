package geni

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type ProjectBulkResponse struct {
	Results []ProjectResponse `json:"results,omitempty"`
}

type ProjectResponse struct {
	// The project's id
	Id string `json:"id,omitempty"`
	// The project's name
	Name string `json:"name,omitempty"`
	// The project's description
	Description *string `json:"description,omitempty"`
	// UpdatedAt is the timestamp of when the project was last updated
	UpdatedAt string `json:"updated_at,omitempty"`
	// CreatedAt is the timestamp of when the project was created
	CreatedAt string `json:"created_at,omitempty"`
}

func (c *Client) GetProject(ctx context.Context, projectId string) (*ProjectResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + projectId
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var project ProjectResponse
	err = json.Unmarshal(body, &project)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &project, nil
}

func (c *Client) AddProfileToProject(ctx context.Context, profileId, projectId string) (*ProfileResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + projectId + "/add_profiles"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	query := req.URL.Query()
	query.Add("profile_ids", profileId)
	req.URL.RawQuery = query.Encode()

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var profileResponse ProfileResponse
	err = json.Unmarshal(body, &profileResponse)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &profileResponse, nil
}

func (c *Client) AddDocumentToProject(ctx context.Context, docimentId, projectId string) (*DocumentBulkResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + projectId + "/add_documents"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	query := req.URL.Query()
	query.Add("ids", docimentId)
	req.URL.RawQuery = query.Encode()

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var documentResponse DocumentBulkResponse
	err = json.Unmarshal(body, &documentResponse)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &documentResponse, nil
}
