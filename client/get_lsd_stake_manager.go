package client

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/lsdprog"
)

var GetLsdStakeManagerCfgDefault = GetAccountInfoConfig{
	Encoding: GetAccountInfoConfigEncodingBase64,
	DataSlice: GetAccountInfoConfigDataSlice{
		Offset: 0,
		Length: lsdprog.StakeManagerAccountLengthDefault,
	},
}

func (s *Client) GetLsdStakeManager(ctx context.Context, account string) (*lsdprog.StakeManager, error) {
	accountInfo, err := s.GetAccountInfo(ctx, account, GetLsdStakeManagerCfgDefault)
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

	stakeManager := lsdprog.StakeManager{}
	err = borsh.Deserialize(&stakeManager, accountDataBts[8:])
	if err != nil {
		return nil, fmt.Errorf("deserialize err: %s", err.Error())
	}
	return &stakeManager, nil
}
