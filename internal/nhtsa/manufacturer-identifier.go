package nhtsavpic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type DecodeWorldManufacturerIdentifierInput struct {
	WMI string `json:"WMI"`
}

type DecodeWorldManufacturerIdentifierOutput struct {
	Count          int                                        `json:"Count"`
	Message        string                                     `json:"Message"`
	SearchCriteria string                                     `json:"SearchCriteria"`
	Results        []DecodeWorldManufacturerIdentifierResults `json:"Results"`
}

type DecodeWorldManufacturerIdentifierResults struct {
	CommonName            string    `json:"CommonName"`
	CreatedOn             time.Time `json:"CreatedOn"`
	DateAvailableToPublic time.Time `json:"DateAvailableToPublic"`
	Make                  string    `json:"Make"`
	ManufacturerName      string    `json:"ManufacturerName"`
	ParentCompanyName     string    `json:"ParentCompanyName"`
	URL                   string    `json:"URL"`
	UpdatedOn             time.Time `json:"UpdatedOn,omitempty"`
	VehicleType           string    `json:"VehicleType"`
}

func (c *Client) DecodeWorldManufacturerIdentifier(ctx context.Context, in DecodeWorldManufacturerIdentifierInput) (DecodeWorldManufacturerIdentifierOutput, error) {
	if strings.TrimSpace(in.WMI) == "" || len(in.WMI) > 6 || len(in.WMI) < 3 {
		return DecodeWorldManufacturerIdentifierOutput{}, ErrInvalidArgument
	}

	url := url.URL{
		Scheme:   "https",
		Host:     baseURL,
		Path:     fmt.Sprintf("api/vehicles/decodewmi/%s", strings.TrimSpace(in.WMI)),
		RawQuery: "format=json",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return DecodeWorldManufacturerIdentifierOutput{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return DecodeWorldManufacturerIdentifierOutput{}, fmt.Errorf("failed to do request: %w", err)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return DecodeWorldManufacturerIdentifierOutput{}, fmt.Errorf("failed to read response data: %w", err)
	}

	var out DecodeWorldManufacturerIdentifierOutput
	if err := json.Unmarshal(respData, &out); err != nil {
		return DecodeWorldManufacturerIdentifierOutput{}, fmt.Errorf("failed to unmarshal response data: %w", err)
	}

	return out, nil
}

type GetWorldManufacturerIdentifiersForManufacturerInput struct {
	Manufacturer string `json:"Manufacturer"`
}

type GetWorldManufacturerIdentifiersForManufacturerOutput struct {
	Count          int                                                    `json:"Count"`
	Message        string                                                 `json:"Message"`
	SearchCriteria string                                                 `json:"SearchCriteria"`
	Results        []GetWorldManufacturerIdentifiersForManufacturerResult `json:"Results"`
}

type GetWorldManufacturerIdentifiersForManufacturerResult struct {
	Country               string    `json:"Country"`
	CreatedOn             time.Time `json:"CreatedOn"`
	DateAvailableToPublic time.Time `json:"DateAvailableToPublic"`
	ID                    int       `json:"Id"`
	Name                  string    `json:"Name"`
	UpdatedOn             time.Time `json:"UpdatedOn"`
	VehicleType           string    `json:"VehicleType"`
	WMI                   string    `json:"WMI"`
}

func (c *Client) GetWorldManufacturerIdentifiersForManufacturer(ctx context.Context, in GetWorldManufacturerIdentifiersForManufacturerInput) (GetWorldManufacturerIdentifiersForManufacturerOutput, error) {
	if strings.TrimSpace(in.Manufacturer) == "" {
		return GetWorldManufacturerIdentifiersForManufacturerOutput{}, ErrInvalidArgument
	}

	url := url.URL{
		Scheme:   "https",
		Host:     baseURL,
		Path:     fmt.Sprintf("api/vehicles/GetWMIsForManufacturer/%s", strings.TrimSpace(in.Manufacturer)),
		RawQuery: "format=json",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return GetWorldManufacturerIdentifiersForManufacturerOutput{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return GetWorldManufacturerIdentifiersForManufacturerOutput{}, fmt.Errorf("failed to do request: %w", err)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetWorldManufacturerIdentifiersForManufacturerOutput{}, fmt.Errorf("failed to read response: %w", err)
	}

	var out GetWorldManufacturerIdentifiersForManufacturerOutput
	if err := json.Unmarshal(respData, &out); err != nil {
		return GetWorldManufacturerIdentifiersForManufacturerOutput{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return out, nil
}
