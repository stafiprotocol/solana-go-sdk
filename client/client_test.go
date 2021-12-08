package client_test

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mr-tron/base58"
	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/bridgeprog"
	"github.com/stafiprotocol/solana-go-sdk/client"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

// var c = client.NewClient([]string{"https://solana-dev-rpc.wetez.io"})
var c = client.NewClient([]string{client.MainnetRPCEndpoint})

func TestAccountInfo(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(30)

	for i := 0; i < 300; i++ {
		time.Sleep(1 * time.Second)
		accountInfo, err := c.GetAccountInfo(context.Background(), "5STUJCFCFPbsagDNk6yBcpiHSPYCwgjjzbrJdWHopC9Q", client.GetAccountInfoConfig{})
		if err != nil {
			t.Log("err", i, err)
		} else {
			t.Log("success", i, fmt.Sprintf("%+v", accountInfo))
		}
	}

	wg.Wait()
}

func TestGetStakeActivation(t *testing.T) {
	accountActivateInfo, err := c.GetStakeActivation(context.Background(), "BfFFmn4iJE5Cmy6opWx26kEHTzrphnxiKpctdeUCNHep", client.GetStakeActivationConfig{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountActivateInfo))
}

func TestGetStakeAccountInfo(t *testing.T) {
	accountActivateInfo, err := c.GetStakeAccountInfo(context.Background(), "Eq2T5683L891HMeGcQHsFbva5fE8795SrXYDJMAQ4Cnq")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountActivateInfo))
	accountActivateInfoBase1, err := c.GetStakeAccountInfo(context.Background(), "AgFCNmujMooFHY378Hb2cvMieXdQS5nP7xXdwWPVytig")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountActivateInfoBase1))

	accountActivateInfo2, err := c.GetStakeActivation(context.Background(), "BfFFmn4iJE5Cmy6opWx26kEHTzrphnxiKpctdeUCNHep", client.GetStakeActivationConfig{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountActivateInfo2))
	accountActivateInfoBase, err := c.GetStakeActivation(context.Background(), "J6L2EyHooCuRLKR17ABFmLmCD9Uq9xwDuboJUpZ5wdH7", client.GetStakeActivationConfig{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountActivateInfoBase))

	account, err := c.GetAccountInfo(context.Background(), "BfFFmn4iJE5Cmy6opWx26kEHTzrphnxiKpctdeUCNHep",
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
	info, err := c.GetMultisigTxAccountInfo(context.Background(), "Gn3Wzs1rbeJcTefiEwZ8c8vJZjNZeSm5WUxbYC5ji74F")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
}

func TestGetMultisigInfoAccount(t *testing.T) {
	info, err := c.GetMultisigInfoAccountInfo(context.Background(), "8TNEsKSzFsi6b56JwhpHWLZf9mR81LGDcQQka5EtVux7")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
}

func TestGetBridgeAccountInfo(t *testing.T) {
	info, err := c.GetBridgeAccountInfo(context.Background(), "63ytYLeNDaaUx2u94KHJcoueaLzA7gryB26p2w8E53oh")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
	t.Log(common.PublicKeyFromBytes(info.FeeReceiver[:]).ToBase58())
}

func TestGetMintProposalInfo(t *testing.T) {
	info, err := c.GetMintProposalInfo(context.Background(), "BtgxF9MgpB9JtxsgeyUKVos6E5N5NbB8BEZLq2RbgUyo")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
}

func TestGetBlock(t *testing.T) {
	info, err := c.GetBlockHeight(context.Background(), client.GetBlockHeightConfig{client.CommitmentFinalized})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("info%+v", info))
}

func TestGetConfirmedBlock(t *testing.T) {
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
	info, err := c.GetTransactionV2(context.Background(), "3r61fK261bxV7uiJKr3jEiR2ysKarF23vMRexYaP5ZDYG5TtNiNw36v6GpuRVntokG7WgJzENS3iYYW9uzZSvbAU")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", info))
	info2, err := c.GetConfirmedTransaction(context.Background(), "3r61fK261bxV7uiJKr3jEiR2ysKarF23vMRexYaP5ZDYG5TtNiNw36v6GpuRVntokG7WgJzENS3iYYW9uzZSvbAU")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", info2))
	blockHeight, err := c.GetBlockHeight(context.Background(), client.GetBlockHeightConfig{client.CommitmentFinalized})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(blockHeight)
	slot, err := c.GetSlot(context.Background(), client.GetSlotConfig{client.CommitmentFinalized})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(slot)
	time, err := c.GetBlockTime(context.Background(), slot)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(time)
}

func TestGetConfirmedTransaction(t *testing.T) {
	info, err := c.GetConfirmedTransaction(context.Background(), "xSSTW1CZoFn3hxBWHXzd6dAh6duiczZmV1H5KgTQGEL2bxY8gZwKak8N3nsmi6NX2X21pgiwnrZQHe9sHa6dwys")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", info))
	t.Log(info.Meta.PreBalances[0] - info.Meta.PostBalances[0])
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
	mintAccount, _ := base64.StdEncoding.DecodeString("AQAAAIJ1WvlDiMw3kmHeTwTkJCzhDg/le+J3e7lDcwGaMPpIAAAAAAAAAAAJAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==")
	t.Log(hex.EncodeToString(mintAccount))
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
	info, err := c.GetConfirmedSignaturesForAddress(context.Background(), "H3mPx8i41Zn4dLC6ZQRBzNRe1cqYdbcDP1WpojnaiAVo", client.GetConfirmedSignaturesForAddressConfig{
		Limit:      1000,
		Until:      "",
		Commitment: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	// t.Log(fmt.Sprintf("%+v", info))
	for _, sig := range info {
		tx, err := c.GetConfirmedTransaction(context.Background(), sig.Signature)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(tx)

	}
	// }
}

func TestGetTokenAccount(t *testing.T) {
	miniMumBalance200, err := c.GetMinimumBalanceForRentExemption(context.Background(), 300000)
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

func TestDecodeAccount(t *testing.T) {
	pool := common.PublicKeyFromString("AycgB5EyyTmuQCrKTkymFQnn6F3PPNRyKuzv6dkuwBhc")
	t.Log(hex.EncodeToString(pool.Bytes()))
	sub1 := common.PublicKeyFromString("D2Qux8umtJ6VAaBuLfDPi9VyHBHhiEB1dKhPKFocKg6q")
	t.Log(hex.EncodeToString(sub1.Bytes()))
	sub2 := common.PublicKeyFromString("9t2Lcij5eGjKN6xPnJkvvM87tyT7QXQ2P5EJyQF7t4jP")
	t.Log(hex.EncodeToString(sub2.Bytes()))
	sub3 := common.PublicKeyFromString("H92d4fR7Jdcxag7JCUAhBKAhnxiS6sWcKhxADLk4dERU")
	t.Log(hex.EncodeToString(sub3.Bytes()))
	receiver := common.PublicKeyFromString("EeTKji2jWLrBeyAzxuonVX3s3DMZip9kBdvH1s5VunET")
	t.Log(hex.EncodeToString(receiver.Bytes()))

	user := common.PublicKeyFromString("8pFiM2vyEzyYL7oJqaK2CgHPnARFdziM753rDHWsnhU1")
	t.Log(hex.EncodeToString(user.Bytes()))

	txHash, _ := base58.Decode("5KZtV2942PxsbQqVircQtitEbe9CHqPMbAHswoKpKJZWNfeR6az9mUTSvgcvAE2rQu8cYjpb1uBVtFxnxk244dny")
	t.Log(hex.EncodeToString(txHash))

	blockHash, _ := base58.Decode("EfR8YmcTSXr3QDHoehtzSDT6FMaAWMhys1i6kaiCZt85")
	t.Log(hex.EncodeToString(blockHash))

}
