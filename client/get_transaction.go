package client

import (
	"context"
	"errors"
)

var ErrTxNotFound = errors.New("TxNotFound")

type GetTransactionWithLimitConfig struct {
	// TODO custom encoding
	// Encoding   string     `json:"encoding"`          // either "json", "jsonParsed", "base58" (slow), "base64", default: json
	Commitment Commitment `json:"commitment,omitempty"` // "processed" is not supported. If parameter not provided, the default is "finalized".
}

type GetTransaction struct {
	Slot        uint64          `json:"slot"`
	Meta        TransactionMeta `json:"meta"`
	Transaction Transaction     `json:"transaction"`
}

// NEW: This method is only available in solana-core v1.7 or newer. Please use getConfirmedTransaction for solana-core v1.6
// GetConfirmedTransaction returns transaction details for a confirmed transaction
func (s *Client) GetTransaction(ctx context.Context, txhash string, cfg GetTransactionWithLimitConfig) (GetConfirmedTransactionResponse, error) {
	res := struct {
		GeneralResponse
		Result GetConfirmedTransactionResponse `json:"result"`
	}{}
	err := s.request(ctx, "getTransaction", []interface{}{txhash, cfg}, &res)
	if err != nil {
		return GetConfirmedTransactionResponse{}, err
	}
	return res.Result, nil
}

// NEW: This method is only available in solana-core v1.7 or newer. Please use getConfirmedTransaction for solana-core v1.6
// GetConfirmedTransaction returns transaction details for a confirmed transaction
func (s *Client) GetTransactionV2(ctx context.Context, txhash string) (GetConfirmedTransactionResponse, error) {
	res := struct {
		GeneralResponse
		Result GetConfirmedTransactionResponse `json:"result"`
	}{}
	err := s.request(ctx, "getTransaction", []interface{}{txhash, GetTransactionWithLimitConfig{CommitmentFinalized}}, &res)
	if err != nil {
		return GetConfirmedTransactionResponse{}, err
	}

	if res.Result.Slot == 0 {
		return GetConfirmedTransactionResponse{}, ErrTxNotFound
	}

	return res.Result, nil
}
