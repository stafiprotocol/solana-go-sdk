package client

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/stafiprotocol/solana-go-sdk/tokenprog"
)

var GetTokenAccountInfoCfgDefault = GetAccountInfoConfig{
	Encoding: GetAccountInfoConfigEncodingBase64,
	DataSlice: GetAccountInfoConfigDataSlice{
		Offset: 0,
		Length: tokenprog.TokenAccountSize,
	},
}

func (s *Client) GetTokenAccountInfo(ctx context.Context, account string) (*tokenprog.TokenAccount, error) {

	accountInfo, err := s.GetAccountInfo(ctx, account, GetMultsigTxAccountInfoCfgDefault)
	if err != nil {
		return nil, err
	}

	accountDataInterface, ok := accountInfo.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("account data err")
	}
	if len(accountDataInterface) != 2 {
		return nil, fmt.Errorf("account data length err")
	}
	accountDataBase64, ok := accountDataInterface[0].(string)
	if !ok {
		return nil, fmt.Errorf("get account base64 failed")
	}

	accountDataBts, err := base64.StdEncoding.DecodeString(accountDataBase64)
	if err != nil {
		return nil, err
	}

	return tokenprog.TokenAccountFromData(accountDataBts[:])
}
