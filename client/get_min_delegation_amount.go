package client

import (
	"context"
	"errors"
)

func (s *Client) GetMinDelegationAmount(ctx context.Context) (uint64, error) {
	res := struct {
		GeneralResponse
		Result struct {
			Value uint64 `json:"value"`
		} `json:"result"`
	}{}

	err := s.request(ctx, "getStakeMinimumDelegation", []interface{}{}, &res)
	if err != nil {
		return 0, err
	}
	if res.Error != (ErrorResponse{}) {
		return 0, errors.New(res.Error.Message)
	}
	return res.Result.Value, nil
}
