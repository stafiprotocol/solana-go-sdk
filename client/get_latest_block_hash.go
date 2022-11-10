package client

import (
	"context"
	"errors"
)

type GetLatestBlockHashResponse struct {
	Blockhash              string `json:"blockhash"`
	LatestValidBlockHeight uint64 `json:"lastValidBlockHeight"`
}

type GetLatestBlockhashConfig struct {
	Commitment Commitment `json:"commitment,omitempty"`
}

func (s *Client) GetLatestBlockhash(ctx context.Context, cfg GetLatestBlockhashConfig) (GetLatestBlockHashResponse, error) {
	res := struct {
		GeneralResponse
		Result struct {
			Context Context                    `json:"context"`
			Value   GetLatestBlockHashResponse `json:"value"`
		} `json:"result"`
	}{}
	err := s.request(ctx, "getLatestBlockhash", []interface{}{}, &res)
	if err != nil {
		return GetLatestBlockHashResponse{}, err
	}
	if res.Error != (ErrorResponse{}) {
		return GetLatestBlockHashResponse{}, errors.New(res.Error.Message)
	}
	return res.Result.Value, nil
}
