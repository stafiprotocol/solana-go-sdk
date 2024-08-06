package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/mr-tron/base58"
	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/rsolprog"
)

func (s *Client) GetUnstakeAccount(ctx context.Context, programId string, stakeManager string, recipient string) ([]rsolprog.UnstakeAccount, error) {

	stakeManagerPubkey := common.PublicKeyFromString(stakeManager)
	recipientPubkey := common.PublicKeyFromString(recipient)

	accounts, err := s.GetProgramAccounts(
		context.Background(),
		programId,
		GetProgramAccountsConfig{
			Encoding: GetAccountInfoConfigEncodingBase64,
			Filters: []interface{}{
				map[string]interface{}{"memcmp": Memcmp{
					Offset: 8,
					Bytes:  base58.Encode(stakeManagerPubkey[:]),
				}},
				map[string]interface{}{"memcmp": Memcmp{
					Offset: 40,
					Bytes:  base58.Encode(recipientPubkey[:]),
				}},
				map[string]interface{}{
					"dataSize": 88,
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
