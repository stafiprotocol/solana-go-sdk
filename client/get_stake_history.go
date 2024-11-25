package client

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/stafiprotocol/solana-go-sdk/binary"
	"github.com/stafiprotocol/solana-go-sdk/common"
)

var GetStakeHistoryDefault = GetAccountInfoConfig{
	Encoding: GetAccountInfoConfigEncodingBase64,
	DataSlice: GetAccountInfoConfigDataSlice{
		Offset: 0,
		Length: 16392,
	},
}

type StakeHistoryRsp struct {
	Lamports       uint64
	Owner          string
	Excutable      bool
	RentEpoch      uint64
	StakeHistories []StakeHistory
}
type StakeHistory struct {
	Epoch uint64
	Entry StakeHistoryEntry
}

type StakeHistoryEntry struct {
	Effective    uint64
	Activating   uint64
	Deactivating uint64
}

func (s *Client) GetStakeHistory(ctx context.Context) (*StakeHistoryRsp, error) {
	accountInfo, err := s.GetAccountInfo(ctx, common.SysVarStakeHistoryPubkey.ToBase58(), GetStakeHistoryDefault)
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

	shy := []StakeHistory{}
	err = bin.NewDecoderWithFixedSize(accountDataBts).Decode(&shy)
	if err != nil {
		return nil, err
	}
	rsp := StakeHistoryRsp{
		Lamports:       accountInfo.Lamports,
		Owner:          accountInfo.Owner,
		Excutable:      accountInfo.Excutable,
		RentEpoch:      accountInfo.RentEpoch,
		StakeHistories: shy,
	}
	return &rsp, nil
}
