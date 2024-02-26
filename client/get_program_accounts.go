package client

import "context"

type GetProgramAccountsConfig struct {
	Commitment  *Commitment                    `json:"commitment,omitempty"` // "processed" is not supported. If parameter not provided, the default is "finalized".
	Encoding    GetAccountInfoConfigEncoding   `json:"encoding"`
	DataSlice   *GetAccountInfoConfigDataSlice `json:"dataSlice,omitempty"`
	Filters     []interface{}                  `json:"filters,omitempty"`
	WithContext bool                           `json:"withContext,omitempty"`
}

type Memcmp struct {
	Offset uint64 `json:"offset"`
	Bytes  string `json:"bytes"`
}

type DataSize struct {
	DataSize uint64
}

type GetProgramAccountsResponse struct {
	Pubkey  string                 `json:"pubkey"`
	Account GetAccountInfoResponse `json:"account"`
}

func (s *Client) GetProgramAccounts(ctx context.Context, programId string, cfg GetProgramAccountsConfig) ([]GetProgramAccountsResponse, error) {
	res := struct {
		GeneralResponse
		Result []GetProgramAccountsResponse `json:"result"`
	}{}
	err := s.request(ctx, "getProgramAccounts", []interface{}{programId, cfg}, &res)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}
