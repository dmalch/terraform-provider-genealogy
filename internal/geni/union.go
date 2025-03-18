package geni

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

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
	Marriage EventElement `json:"marriage,omitempty"`
	// Divorce date and location
	Divorce EventElement `json:"divorce,omitempty"`
	// Status of the union (spouse|ex_spouse)
	Status string `json:"status,omitempty"`
}

func GetUnion(accessToken, unionId string) (*UnionResponse, error) {
	baseUrl := geniUrl + "api/" + unionId
	req, err := http.NewRequest(http.MethodGet, baseUrl, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	addStandardHeadersAndQueryParams(req, accessToken)

	body, err := doRequest(req)
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

func AddChildToUnion(accessToken, unionId string) (*ProfileResponse, error) {
	baseUrl := geniUrl + "api/" + unionId + "/add-child"
	req, err := http.NewRequest(http.MethodPost, baseUrl, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	addStandardHeadersAndQueryParams(req, accessToken)

	body, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	var profile ProfileResponse
	err = json.Unmarshal(body, &profile)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	fixResponse(&profile)

	return &profile, nil
}
