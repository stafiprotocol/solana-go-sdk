package client

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/common"
)

type GetBridgeAccountInfo struct {
	Admin            [32]uint8
	Owners           [][32]uint8
	Threshold        uint64
	Nonce            uint8
	OwnerSetSeqno    uint32
	DepositCount     map[uint8]uint64
	ResourceIdToMint map[[32]uint8]common.PublicKey
}

var BridgeAccountLengthDefault = uint64(2000)
var GetBridgeAccountInfoCfgDefault = GetAccountInfoConfig{
	Encoding: GetAccountInfoConfigEncodingBase64,
	DataSlice: GetAccountInfoConfigDataSlice{
		Offset: 0,
		Length: MultisigTxAccountLengthDefault,
	},
}

func (s *Client) GetBridgeAccountInfo(ctx context.Context, account string) (*GetBridgeAccountInfo, error) {
	accountInfo, err := s.GetAccountInfo(ctx, account, GetBridgeAccountInfoCfgDefault)
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

	bridgeAccountInfo := GetBridgeAccountInfo{}
	err = borsh.Deserialize(&bridgeAccountInfo, accountDataBts[8:])
	if err != nil {
		return nil, fmt.Errorf("Deserialize err: %s", err)
	}
	return &bridgeAccountInfo, nil
}
