package client

import (
	"context"
	"errors"
)

var ErrTxNotFound = errors.New("TxNotFound")
var DefaultMaxSupportedTransactionVersion = uint8(0)

type GetTransactionWithLimitConfig struct {
	// TODO custom encoding
	// Encoding   string     `json:"encoding"`          // either "json", "jsonParsed", "base58" (slow), "base64", default: json
	Commitment                     Commitment `json:"commitment,omitempty"`                     // "processed" is not supported. If parameter not provided, the default is "finalized".
	MaxSupportedTransactionVersion *uint8     `json:"maxSupportedTransactionVersion,omitempty"` // default: nil legacy only
}

type GetTransaction struct {
	Slot        uint64          `json:"slot"`
	Meta        TransactionMeta `json:"meta"`
	Transaction Transaction     `json:"transaction"`
}

type GetTransactionResponse struct {
	Slot        uint64          `json:"slot"`
	Meta        TransactionMeta `json:"meta"`
	Transaction Transaction     `json:"transaction"`
}

// NEW: This method is only available in solana-core v1.7 or newer. Please use getConfirmedTransaction for solana-core v1.6
// GetConfirmedTransaction returns transaction details for a confirmed transaction
func (s *Client) GetTransaction(ctx context.Context, txhash string, cfg GetTransactionWithLimitConfig) (GetTransactionResponse, error) {
	res := struct {
		GeneralResponse
		Result GetTransactionResponse `json:"result"`
	}{}
	err := s.request(ctx, "getTransaction", []interface{}{txhash, cfg}, &res)
	if err != nil {
		return GetTransactionResponse{}, err
	}
	return res.Result, nil
}

// NEW: This method is only available in solana-core v1.7 or newer. Please use getConfirmedTransaction for solana-core v1.6
// GetConfirmedTransaction returns transaction details for a confirmed transaction
func (s *Client) GetTransactionV2(ctx context.Context, txhash string) (GetTransactionResponse, error) {
	res := struct {
		GeneralResponse
		Result GetTransactionResponse `json:"result"`
	}{}
	err := s.request(ctx, "getTransaction", []interface{}{txhash, GetTransactionWithLimitConfig{Commitment: CommitmentFinalized, MaxSupportedTransactionVersion: &DefaultMaxSupportedTransactionVersion}}, &res)
	if err != nil {
		return GetTransactionResponse{}, err
	}

	if res.Result.Slot == 0 {
		return GetTransactionResponse{}, ErrTxNotFound
	}

	return res.Result, nil
}
