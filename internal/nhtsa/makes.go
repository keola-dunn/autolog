package nhtsavpic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type GetAllMakesOutput struct {
	// Count is the number of makes returned from the search
	Count int `json:"Count"`

	// Message returned from the API
	Message string `json:"Message"`

	// Results are the results returned for getting all makes
	Results []GetAllMakesResult `json:"Results"`
}

type GetAllMakesResult struct {
	MakeID   int    `json:"Make_ID"`
	MakeName string `json:"Make_Name"`
}

// GetAllMakes wraps the NHTSA vPIC Get All Makes API call
func (c *Client) GetAllMakes(ctx context.Context) (GetAllMakesOutput, error) {

	url := url.URL{
		Scheme:   "https",
		Host:     baseURL,
		Path:     "/api/vehicles/getallmakes",
		RawQuery: "format=json",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return GetAllMakesOutput{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return GetAllMakesOutput{}, fmt.Errorf("failed to get all makes: %w", err)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetAllMakesOutput{}, fmt.Errorf("failed to read response data: %w", err)
	}

	var out GetAllMakesOutput
	if err := json.Unmarshal(respData, &out); err != nil {
		return GetAllMakesOutput{}, fmt.Errorf("failed to unmarhsal response as expected: %w", err)
	}

	return out, nil
}
