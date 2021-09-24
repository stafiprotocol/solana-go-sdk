package client_test

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/bridgeprog"
	"github.com/stafiprotocol/solana-go-sdk/client"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/types"
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
	// c := client.NewClient("https://api.devnet.solana.com")
	c := client.NewClient("https://solana-dev-rpc.wetez.io")
	info, err := c.GetBridgeAccountInfo(context.Background(), "63ytYLeNDaaUx2u94KHJcoueaLzA7gryB26p2w8E53oh")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
	t.Log(common.PublicKeyFromBytes(info.FeeReceiver[:]).ToBase58())
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
				t.Log(strings.TrimPrefix(log, bridgeprog.ProgramLogPrefix))
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
	// msg := "7arrB4Lk4L+33mMucKYMb78cH5By6eymggY2XBfqajtrBnTVEmFjbAUAAAABAQEBAQoAAAAAAAAAAQECAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAA="
	msg := "7arrB4Lk4L8DxOYp6nBnEQYF6Kx+u2D/FSd+muH+uTMW3s/snnL2JCAAAAB0g0gRxgiA0CZ5M+McJT6TfhSFT1Ls3R8l0mvcGR4tEICWmAAAAAAAAQAAAAAAAAAAAAAAAAAAAGWbkw+FaJUst7DIt+2jBgsBAgAAAAAAAAA="
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
	// t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:9]))
	// t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:10]))
	// t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:11]))
	// t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:12]))
	// t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:13]))
	// t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:14]))
	// t.Log(base64.StdEncoding.EncodeToString(accountDataBts[:15]))

	multiTxAccountInfo := EventTransferOut{}
	err = borsh.Deserialize(&multiTxAccountInfo, accountDataBts[8:])
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(multiTxAccountInfo.Receiver))
	t.Log(multiTxAccountInfo.Amount)

	// pubkey:=common.PublicKeyFromString("9Riwnxn53S4wmy5h5nbQN1gxTCm1EvgqB4Gc5aKDAPyc")
	pubkey := common.PublicKeyFromString("2cTdCXvyeLfNvoKinFVWGYWnWYxaY45gydtnnbJpSJE3")
	t.Log(pubkey)
	t.Log(hex.EncodeToString(pubkey.Bytes()))

	bts, err := hex.DecodeString("98d9634ad58009cda11726a718073b5ba525d51483cbf8e8bef127cb6b70e900")
	t.Log(common.PublicKeyFromBytes(bts).ToBase58())
}

// FRzXkJ4p1knQkFdBCtLCt8Zuvykr7Wd5yKTrryQV3K51

func TestGetSignaturesForAddress(t *testing.T) {
	c := client.NewClient("https://api.devnet.solana.com")
	info, err := c.GetConfirmedSignaturesForAddress(context.Background(), "FRzXkJ4p1knQkFdBCtLCt8Zuvykr7Wd5yKTrryQV3K51", client.GetConfirmedSignaturesForAddressConfig{
		Limit:      1000,
		Before:     "5yhpbdfLBJvstkpv2RaE4A98xGiEanrznt2yAV22ooxedLTSThQFXmvUyRboJX38e2UKokZtBvYMcQonLxQ8j6SD",
		Until:      "2T64SSqK3X6xQsbqEgx5THTXFtKHmL14gfbMX1sZaXdqfpGbA3CcTab57p2jw9qEHnYnHbYavKtoyz1wxYZP8vDi",
		Commitment: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", info))
}

func TestGetTokenAccount(t *testing.T) {
	c := client.NewClient("https://api.mainnet-beta.solana.com")
	miniMumBalance200, err := c.GetMinimumBalanceForRentExemption(context.Background(), 200)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(miniMumBalance200)

	feePayer := types.AccountFromPrivateKeyBytes([]byte{179, 95, 213, 234, 125, 167, 246, 188, 230, 134, 181, 219, 31, 146, 239, 75, 190, 124, 112, 93, 187, 140, 178, 119, 90, 153, 207, 178, 137, 5, 53, 71, 116, 28, 190, 12, 249, 238, 110, 135, 109, 21, 196, 36, 191, 19, 236, 175, 229, 204, 68, 180, 130, 102, 71, 239, 41, 53, 152, 159, 175, 124, 180, 6})
	_, err = c.RequestAirdrop(context.Background(), feePayer.PublicKey.ToBase58(), 10e9)
	if err != nil {
		fmt.Println(err)
	}
	fromBytes, _ := hex.DecodeString("cf0b31c9a3ca108ffe22d4e9b73af6be36c87fc4cfabe52a938ca60ce28c20143429f41f8636e46a8f7a90a11c1e652787bbee64a60a04650f7f5b8e55f0a739")
	fromAccount := types.AccountFromPrivateKeyBytes(fromBytes)
	fmt.Println("fromAccount", fromAccount.PublicKey.ToBase58())
	accountInfo, err := c.GetTokenAccountInfo(context.Background(), fromAccount.PublicKey.ToBase58())
	if err != nil {
		t.Log(err)
	}
	t.Log(fmt.Sprintf("%+v", accountInfo))
	t.Log(fmt.Sprintf("%+v", accountInfo.Mint.ToBase58()))
	t.Log(hex.EncodeToString(bridgeprog.InstructionTransferOut[:]))
}
