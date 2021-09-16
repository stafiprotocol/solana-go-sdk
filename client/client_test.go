package client_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/bridgeprog"
	"github.com/stafiprotocol/solana-go-sdk/client"
	"github.com/stafiprotocol/solana-go-sdk/common"
)

func TestAccountInfo(t *testing.T) {
	c := client.NewClient(client.DevnetRPCEndpoint)

	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			accountInfo, err := c.GetMultisigTxAccountInfo(context.Background(), "D6nA6QHpYQDMeudHLwZqgwyCJfRSKWfzW4kyaKqmnsr4")
			if err != nil {
				t.Fatal(err)
			}
			t.Log(fmt.Printf("%+v", accountInfo))
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestGetAccountInfo(t *testing.T) {
	c := client.NewClient(client.DevnetRPCEndpoint)
	account, err := c.GetAccountInfo(context.Background(), "DiPx1Vyo5khyG8XKTc8Tu4fL9qc57VSqfr7qh3xLxqjX",
		client.GetAccountInfoConfig{
			Encoding: client.GetAccountInfoConfigEncodingBase64,
			DataSlice: client.GetAccountInfoConfigDataSlice{
				Offset: 0,
				Length: 200,
			},
		})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(account)

	accountInfo, err := c.GetStakeAccountInfo(context.Background(), "DiPx1Vyo5khyG8XKTc8Tu4fL9qc57VSqfr7qh3xLxqjX")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountInfo))

	accountInfo, err = c.GetStakeAccountInfo(context.Background(), "mNzHTv7KtARcYyJiaXH3SU2oSnLoXKVJxcxDknC2kae")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountInfo))

	accountActivateInfo, err := c.GetStakeActivation(context.Background(), "DiPx1Vyo5khyG8XKTc8Tu4fL9qc57VSqfr7qh3xLxqjX", client.GetStakeActivationConfig{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountActivateInfo))

}

func TestGetStakeActivation(t *testing.T) {
	c := client.NewClient("http://127.0.0.1:8899")
	accountActivateInfo, err := c.GetStakeActivation(context.Background(), "D17ya9gd9xRSMwqjpzixXx341gPF2sNzKft5CMnToF8h", client.GetStakeActivationConfig{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountActivateInfo))
}

func TestGetStakeAccountInfo(t *testing.T) {
	c := client.NewClient("http://127.0.0.1:8899")
	// accountActivateInfo, err := c.GetStakeAccountInfo(context.Background(), "D17ya9gd9xRSMwqjpzixXx341gPF2sNzKft5CMnToF8h")
	accountActivateInfo, err := c.GetStakeAccountInfo(context.Background(), "Gnr9LuHUh85Dt7Qr3tayXrxFAEn32jRDfsgTAyywFhyh")
	// accountActivateInfo, err := c.GetStakeAccountInfo(context.Background(), "ATY5PSBVExLoFb2CnRj1e9nUVghcvLcrvbhcYMud1d4F")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountActivateInfo))
	accountActivateInfo2, err := c.GetStakeActivation(context.Background(), "2usyY4HuMCfZW6CGjzdUuheMxh71HPQFkyZZW61qAxAq", client.GetStakeActivationConfig{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountActivateInfo2))

	account, err := c.GetAccountInfo(context.Background(), "2usyY4HuMCfZW6CGjzdUuheMxh71HPQFkyZZW61qAxAq",
		client.GetAccountInfoConfig{
			Encoding: client.GetAccountInfoConfigEncodingBase64,
			DataSlice: client.GetAccountInfoConfigDataSlice{
				Offset: 0,
				Length: 200,
			},
		})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", account))

}

func TestGetMultisigTxInfo(t *testing.T) {
	c := client.NewClient("https://api.devnet.solana.com")
	info, err := c.GetMultisigTxAccountInfo(context.Background(), "Gn3Wzs1rbeJcTefiEwZ8c8vJZjNZeSm5WUxbYC5ji74F")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
}

func TestGetMultisigInfoAccount(t *testing.T) {
	c := client.NewClient("https://api.devnet.solana.com")
	info, err := c.GetMultisigInfoAccountInfo(context.Background(), "8TNEsKSzFsi6b56JwhpHWLZf9mR81LGDcQQka5EtVux7")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
}

func TestGetBridgeAccountInfo(t *testing.T) {
	c := client.NewClient("https://api.devnet.solana.com")
	info, err := c.GetBridgeAccountInfo(context.Background(), "8B29iREQvQgmiyZSdzRrgEsh56W3m2Mpna5xvGG6jAEf")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
}

func TestGetMintProposalInfo(t *testing.T) {
	c := client.NewClient("https://api.devnet.solana.com")
	info, err := c.GetMintProposalInfo(context.Background(), "BtgxF9MgpB9JtxsgeyUKVos6E5N5NbB8BEZLq2RbgUyo")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
}

func TestGetBlock(t *testing.T) {
	c := client.NewClient("https://api.devnet.solana.com")
	info, err := c.GetBlock(context.Background(), 80837538, client.GetBlockConfig{})
	if err != nil {
		t.Fatal(err)
	}
	// t.Log(fmt.Printf("%+v", info))

	for _, tx := range info.Transactions {
		t.Log(tx.Meta.LogMessages)
	}
}

func TestGetConfirmedBlock(t *testing.T) {
	c := client.NewClient("https://api.devnet.solana.com")
	info, err := c.GetConfirmedBlock(context.Background(), 81048933)
	if err != nil {
		t.Fatal(err)
	}
	// t.Log(fmt.Printf("%+v", info))

	for _, tx := range info.Transactions {
		for _, log := range tx.Meta.LogMessages {
			// t.Log(log)
			if strings.HasPrefix(log, bridgeprog.EventTransferOutPrefix) {
				t.Log(strings.TrimPrefix(log,bridgeprog.ProgramLogPrefix))
			}
		}
		// t.Log(tx.Meta.LogMessages)
	}
}

func TestGetTransaction(t *testing.T) {
	c := client.NewClient("https://explorer-api.devnet.solana.com")
	info, err := c.GetTransaction(context.Background(), "2hF4qEu4xYX51Pu2ErcXcGKXojzbzfahhSNeMXXmFAkCW1Rom4uk51Tur7uuWfJmpMzcqFQkRFYEabdNqsz8m7fa", client.GetTransactionWithLimitConfig{})
	if err != nil {
		t.Fatal(err)
	}

	for _, tx := range info.Meta.LogMessages {
		if strings.HasPrefix(tx, bridgeprog.EventTransferOutPrefix) {
			t.Log(tx)
		}
	}
}

func TestGetConfirmedTransaction(t *testing.T) {
	c := client.NewClient("https://explorer-api.devnet.solana.com")
	info, err := c.GetConfirmedTransaction(context.Background(), "2hF4qEu4xYX51Pu2ErcXcGKXojzbzfahhSNeMXXmFAkCW1Rom4uk51Tur7uuWfJmpMzcqFQkRFYEabdNqsz8m7fa")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", info))
	for _, tx := range info.Meta.LogMessages {
		t.Log(tx)
	}
}

func TestGetBlockHeight(t *testing.T) {
	c := client.NewClient("https://api.devnet.solana.com")
	info, err := c.GetBlockHeight(context.Background(), client.GetBlockHeightConfig{client.CommitmentFinalized})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", info))
}

type EventTransferOut struct {
	Transfer     common.PublicKey
	Receiver     []byte
	Amount       uint64
	DestChainId  uint8
	ResourceId   [32]byte
	DepositNonce uint64
}

func TestParseLog(t *testing.T) {
	msg := "7arrB4Lk4L+33mMucKYMb78cH5By6eymggY2XBfqajtrBnTVEmFjbAUAAAABAQEBAQoAAAAAAAAAAQECAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAA="
	accountDataBts, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		t.Fatal(err)
	}
	if len(accountDataBts) <= 8 {
		t.Fatal("ee")
	}
	t.Log(accountDataBts[:8])
	t.Log(bridgeprog.EventTransferOut)
	t.Log(base64.StdEncoding.EncodeToString(bridgeprog.EventTransferOut[:]))
	t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:9]))
	t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:10]))
	t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:11]))
	t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:12]))
	t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:13]))
	t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:14]))
	t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:15]))

	multiTxAccountInfo := EventTransferOut{}
	err = borsh.Deserialize(&multiTxAccountInfo, accountDataBts[8:])
	if err != nil {
		t.Fatal(err)
	}
	t.Log(multiTxAccountInfo)
}




// FRzXkJ4p1knQkFdBCtLCt8Zuvykr7Wd5yKTrryQV3K51


func TestGetSignaturesForAddress(t *testing.T) {
	c := client.NewClient("https://api.devnet.solana.com")
	info, err := c.GetSignaturesForAddress(context.Background(),"FRzXkJ4p1knQkFdBCtLCt8Zuvykr7Wd5yKTrryQV3K51",client.GetConfirmedSignaturesForAddressConfig{
		Limit:      0,
		Before:     "",
		Until:      "",
		Commitment: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", info))
}