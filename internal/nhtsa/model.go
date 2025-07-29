package nhtsavpic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GetModelsForMakeInput struct {
	Make string `json:"Make"`
}

type GetModelsForMakeOutput struct {
	Count          int                      `json:"Count"`
	Message        string                   `json:"Message"`
	SearchCriteria string                   `json:"SearchCriteria"`
	Results        []GetModelsForMakeResult `json:"Results"`
}

type GetModelsForMakeResult struct {
	MakeID    int    `json:"Make_ID"`
	MakeName  string `json:"Make_Name"`
	ModelID   int    `json:"Model_ID"`
	ModelName string `json:"Model_Name"`
}

func (c *Client) GetModelsForMake(ctx context.Context, in GetModelsForMakeInput) (GetModelsForMakeOutput, error) {
	if strings.TrimSpace(in.Make) == "" {
		return GetModelsForMakeOutput{}, ErrInvalidArgument
	}

	url := url.URL{
		Scheme:   "https",
		Host:     baseURL,
		Path:     fmt.Sprintf("api/vehicles/GetModelsForMake/%s", strings.TrimSpace(in.Make)),
		RawQuery: "format=json",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return GetModelsForMakeOutput{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return GetModelsForMakeOutput{}, fmt.Errorf("failed to do request: %w", err)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetModelsForMakeOutput{}, fmt.Errorf("failed to read response: %w", err)
	}

	var out GetModelsForMakeOutput
	if err := json.Unmarshal(respData, &out); err != nil {
		return GetModelsForMakeOutput{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return out, nil
}

type GetModelsForMakeIDInput struct {
	MakeID int `json:"Make_ID"`
}

type GetModelsForMakeIDOutput struct {
	Count          int                      `json:"Count"`
	Message        string                   `json:"Message"`
	SearchCriteria string                   `json:"SearchCriteria"`
	Results        []GetModelsForMakeResult `json:"Results"`
}

func (c *Client) GetModelsForMakeID(ctx context.Context, in GetModelsForMakeIDInput) (GetModelsForMakeIDOutput, error) {
	if in.MakeID < 0 {
		return GetModelsForMakeIDOutput{}, ErrInvalidArgument
	}

	url := url.URL{
		Scheme:   "https",
		Host:     baseURL,
		Path:     fmt.Sprintf("api/vehicles/GetModelsForMakeId/%d", in.MakeID),
		RawQuery: "format=json",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return GetModelsForMakeIDOutput{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return GetModelsForMakeIDOutput{}, fmt.Errorf("failed to get models for make id: %w", err)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetModelsForMakeIDOutput{}, fmt.Errorf("failed to read response: %w", err)
	}

	var out GetModelsForMakeIDOutput
	if err := json.Unmarshal(respData, &out); err != nil {
		return GetModelsForMakeIDOutput{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return out, nil
}
