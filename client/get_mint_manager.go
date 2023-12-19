package client

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/minterprog"
)

var GetMintManagerCfgDefault = GetAccountInfoConfig{
	Encoding: GetAccountInfoConfigEncodingBase64,
	DataSlice: GetAccountInfoConfigDataSlice{
		Offset: 0,
		Length: minterprog.MinterManagerAccountLengthDefault,
	},
}

func (s *Client) GetMintManager(ctx context.Context, account string) (*minterprog.MintManager, error) {
	accountInfo, err := s.GetAccountInfo(ctx, account, GetMintManagerCfgDefault)
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
	if len(accountDataBts) <= 8 {
		return nil, fmt.Errorf("no account data bytes")
	}

	mintManager := minterprog.MintManager{}
	err = borsh.Deserialize(&mintManager, accountDataBts[8:])
	if err != nil {
		return nil, fmt.Errorf("deserialize err: %s", err.Error())
	}
	return &mintManager, nil
}
