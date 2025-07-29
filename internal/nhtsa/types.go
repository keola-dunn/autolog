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

type GetVehicleTypesForMakeByNameInput struct {
	MakeName string `json:"MakeName"`
}

type GetVehicleTypesForMakeByNameOutput struct {
	Count          int                                  `json:"Count"`
	Message        string                               `json:"Message"`
	SearchCriteria string                               `json:"SearchCriteria"`
	Results        []GetVehicleTypesForMakeByNameResult `json:"Results"`
}

type GetVehicleTypesForMakeByNameResult struct {
	MakeID          int    `json:"MakeId"`
	MakeName        string `json:"MakeName"`
	VehicleTypeID   int    `json:"VehicleTypeId"`
	VehicleTypeName string `json:"VehicleTypeName"`
}

func (c *Client) GetVehicleTypesForMakeByName(ctx context.Context, in GetVehicleTypesForMakeByNameInput) (GetVehicleTypesForMakeByNameOutput, error) {
	if strings.TrimSpace(in.MakeName) == "" {
		return GetVehicleTypesForMakeByNameOutput{}, ErrInvalidArgument
	}

	url := url.URL{
		Scheme:   "https",
		Host:     baseURL,
		Path:     fmt.Sprintf("api/vehicles/GetVehicleTypesForMake/%s", strings.TrimSpace(in.MakeName)),
		RawQuery: "format=json",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return GetVehicleTypesForMakeByNameOutput{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return GetVehicleTypesForMakeByNameOutput{}, fmt.Errorf("failed to do request: %w", err)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetVehicleTypesForMakeByNameOutput{}, fmt.Errorf("failed to read response: %w", err)
	}

	var out GetVehicleTypesForMakeByNameOutput
	if err := json.Unmarshal(respData, &out); err != nil {
		return GetVehicleTypesForMakeByNameOutput{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return out, nil
}

type GetVehicleTypesForMakeByIDInput struct {
	MakeID int `json:"MakeId"`
}

type GetVehicleTypesForMakeByIDOutput struct {
	Count          int                                `json:"Count"`
	Message        string                             `json:"Message"`
	SearchCriteria string                             `json:"SearchCriteria"`
	Results        []GetVehicleTypesForMakeByIDResult `json:"Results"`
}

type GetVehicleTypesForMakeByIDResult struct {
	VehicleTypeID   int    `json:"VehicleTypeId"`
	VehicleTypeName string `json:"VehicleTypeName"`
}

func (c *Client) GetVehicleTypesForMakeByID(ctx context.Context, in GetVehicleTypesForMakeByIDInput) (GetVehicleTypesForMakeByIDOutput, error) {
	if in.MakeID < 0 {
		return GetVehicleTypesForMakeByIDOutput{}, ErrInvalidArgument
	}

	url := url.URL{
		Scheme:   "https",
		Host:     baseURL,
		Path:     fmt.Sprintf("api/vehicles/GetVehicleTypesForMakeId/%d", in.MakeID),
		RawQuery: "format=json",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return GetVehicleTypesForMakeByIDOutput{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return GetVehicleTypesForMakeByIDOutput{}, fmt.Errorf("failed to do request: %w", err)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetVehicleTypesForMakeByIDOutput{}, fmt.Errorf("failed to read response: %w", err)
	}

	var out GetVehicleTypesForMakeByIDOutput
	if err := json.Unmarshal(respData, &out); err != nil {
		return GetVehicleTypesForMakeByIDOutput{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return out, nil
}
