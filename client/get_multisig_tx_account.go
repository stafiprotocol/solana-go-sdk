package client

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/near/borsh-go"
)

type GetMultisigTxAccountInfo struct {
	Multisig   [32]uint8
	ProgramID  [32]uint8
	Accounts   []TransactionAccount
	Data       []uint8
	Signers    []uint8
	DidExecute uint8
}

type TransactionAccount struct {
	Pubkey     [32]uint8
	IsSigner   uint8
	IsWritable uint8
}

var MultisigTxAccountLengthDefault = uint64(1000)
var GetMultsigTxAccountInfoCfgDefault = GetAccountInfoConfig{
	Encoding: GetAccountInfoConfigEncodingBase64,
	DataSlice: GetAccountInfoConfigDataSlice{
		Offset: 0,
		Length: MultisigTxAccountLengthDefault,
	},
}

func (s *Client) GetMultisigTxAccountInfo(ctx context.Context, account string) (*GetMultisigTxAccountInfo, error) {

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
	if len(accountDataBts) <= 8 {
		return nil, fmt.Errorf("no account data bytes")
	}

	multiTxAccountInfo := GetMultisigTxAccountInfo{}
	err = borsh.Deserialize(&multiTxAccountInfo, accountDataBts[8:])
	if err != nil {
		return nil, err
	}
	return &multiTxAccountInfo, nil
}
