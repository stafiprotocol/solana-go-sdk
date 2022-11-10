package client

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/common"
)

type GetMintProposalINfo struct {
	Bridge        [32]uint8
	Signers       []uint8
	DidExecute    uint8
	OwnerSetSeqno uint32
	Mint          common.PublicKey
	To            common.PublicKey
	Amount        uint64
	TokenProgram  common.PublicKey
}

var MintProposalInfoLengthDefault = uint64(200)
var MintProposalInfoCfgDefault = GetAccountInfoConfig{
	Encoding: GetAccountInfoConfigEncodingBase64,
	DataSlice: GetAccountInfoConfigDataSlice{
		Offset: 0,
		Length: MintProposalInfoLengthDefault,
	},
}

func (s *Client) GetMintProposalInfo(ctx context.Context, account string) (*GetMintProposalINfo, error) {
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

	mintProposalInfo := GetMintProposalINfo{}
	err = borsh.Deserialize(&mintProposalInfo, accountDataBts[8:])
	if err != nil {
		return nil, fmt.Errorf("deserialize err: %s", err)
	}
	return &mintProposalInfo, nil
}
