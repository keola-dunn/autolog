package nhtsavpic

import (
	"fmt"
	"strconv"
	"strings"
)

type ErrorCode int

const (
	// ErrorCodeSuccess is the error code returned when there isn't an issue decoding the VIN
	ErrorCodeSuccess = ErrorCode(0)

	// ErrorCodeIncompleteVIN is the error returned when a VIN is incomplete
	ErrorCodeIncompleteVIN = ErrorCode(6)

	// ErrorCodeManufacturerNotRegistered is the error code when a manufacturer cannot be found
	ErrorCodeManufacturerNotRegistered = ErrorCode(7)

	// ErrorCodeModelYearWarning is the error code when a provided model year doesn't match the VIN
	ErrorCodeModelYearWarning = ErrorCode(12)
)

func (d *DecodeVINFlatResult) ErrorCodes() ([]ErrorCode, error) {
	codes := strings.Split(d.ErrorCode, ",")

	var resp = []ErrorCode{}

	for _, code := range codes {
		errCode, err := strconv.Atoi(strings.ToLower(code))
		if err != nil {
			return resp, fmt.Errorf("failed to parse error codes: %w", err)
		}

		resp = append(resp, ErrorCode(errCode))
	}

	return resp, nil
}
