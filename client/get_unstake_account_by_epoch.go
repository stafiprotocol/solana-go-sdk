package client

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/mr-tron/base58"
	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/rsolprog"
)

func (s *Client) GetUnstackAccountByEpoch(ctx context.Context, programId string, epoch uint64) ([]rsolprog.UnstakeAccount, error) {
	hexStr := strconv.FormatUint(epoch, 16)
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	hexBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	newHexBytes := make([]byte, len(hexBytes))
	for i, bt := range hexBytes {
		newHexBytes[len(hexBytes)-1-i] = bt
	}

	accounts, err := s.GetProgramAccounts(
		context.Background(),
		programId,
		GetProgramAccountsConfig{
			Encoding: GetAccountInfoConfigEncodingBase64,
			Filters: []interface{}{
				map[string]interface{}{"memcmp": Memcmp{
					Offset: 80,
					Bytes:  base58.Encode(newHexBytes),
				}},
				map[string]interface{}{
					"dataSize": 100,
				},
			}})
	if err != nil {
		return nil, err
	}
	ret := make([]rsolprog.UnstakeAccount, 0)
	for _, accountInfo := range accounts {
		accountDataInterface, ok := accountInfo.Account.Data.([]interface{})
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

		unstakeAccount := rsolprog.UnstakeAccount{}
		err = borsh.Deserialize(&unstakeAccount, accountDataBts[8:])
		if err != nil {
			return nil, fmt.Errorf("deserialize err: %s", err.Error())
		}
		ret = append(ret, unstakeAccount)
	}
	return ret, nil
}
