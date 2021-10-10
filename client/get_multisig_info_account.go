package client

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/near/borsh-go"
)

type GetMultisigInfoAccountInfo struct {
	Owners        [][32]uint8
	Threshold     uint64
	Nonce         uint8
	OwnerSetSeqno uint32
}

var MultisigInfoAccountLengthDefault = uint64(1000)
var GetMultsigInfoAccountInfoCfgDefault = GetAccountInfoConfig{
	Encoding: GetAccountInfoConfigEncodingBase64,
	DataSlice: GetAccountInfoConfigDataSlice{
		Offset: 0,
		Length: MultisigInfoAccountLengthDefault,
	},
}

func (s *Client) GetMultisigInfoAccountInfo(ctx context.Context, account string) (*GetMultisigInfoAccountInfo, error) {
	accountInfo, err := s.GetAccountInfo(ctx, account, GetMultsigInfoAccountInfoCfgDefault)
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

	multiInfoAccountInfo := GetMultisigInfoAccountInfo{}
	err = borsh.Deserialize(&multiInfoAccountInfo, accountDataBts[8:])
	if err != nil {
		return nil, fmt.Errorf("Deserialize err: %s", err)
	}
	return &multiInfoAccountInfo, nil
}
