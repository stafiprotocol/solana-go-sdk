package client

import (
	"context"
	"errors"
)

type GetSlotConfig struct {
	Commitment Commitment `json:"commitment,omitempty"`
}

// GetSlob returns the current slot  of the node
func (s *Client) GetSlot(ctx context.Context, cfg GetSlotConfig) (uint64, error) {
	res := struct {
		GeneralResponse
		Result uint64 `json:"result"`
	}{}
	err := s.request(ctx, "getSlot", []interface{}{cfg}, &res)
	if err != nil {
		return 0, err
	}
	if res.Error != (ErrorResponse{}) {
		return 0, errors.New(res.Error.Message)
	}
	return res.Result, nil
}
