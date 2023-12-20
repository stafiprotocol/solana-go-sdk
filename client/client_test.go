package client_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
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

// var c = client.NewClient([]string{"https://api.devnet.solana.com"})

var c = client.NewClient([]string{"https://solana-dev-rpc.stafi.io"})

// var c = client.NewClient([]string{client.MainnetRPCEndpoint})
// var c = client.NewClient([]string{"https://solana-rpc1.stafi.io"})

// var c = client.NewClient([]string{"https://mainnet-rpc.wetez.io/solana/v1/6e0a86ceca790361d95a588efcd1af0b"})
// var c = client.NewClient([]string{"https://mainnet-rpc.wetez.io/solana/v1/308aa4d20d1624a5a35e2d7fca8624f9"})

// var c = client.NewClient([]string{"https://try-rpc.mainnet.solana.blockdaemon.tech"})

// var c = client.NewClient([]string{"https://rpc.ankr.com/solana_mainnet"})
// var c = client.NewClient([]string{"https://solana-mainnet.g.alchemy.com/v2/jfqvfqIeeKDImPdksQEH-SL62h-fExgv"})
// var c = client.NewClient([]string{"https://try.blockdaemon.com/rpc/solana"})
// var c = client.NewClient([]string{"https://try-rpc.mainnet.solana.blockdaemon.tech"})

// var c = client.NewClient([]string{"https://solana.public-rpc.com"})

// var c = client.NewClient([]string{"https://solana-mainnet.phantom.tech"})

// var c = client.NewClient([]string{"https://free.rpcpool.com"})

// var c = client.NewClient([]string{"https://solana.public-rpc.com"})
// var c = client.NewClient([]string{"https://free.rpcpool.com"})

// var c = client.NewClient([]string{"https://solana-mainnet.phantom.tech"})
// 4天前
// era=314 active=1044170088955   https://solana.public-rpc.com

// era=314 active=1058670098955   https://rpc.ankr.com/solana

// era=314 active=1058670098955   https://solana.public-rpc.com

// 今天
// era=314 active=1058831750514   https://rpc.ankr.com/solana

func GetStakeAccountPubkey(baseAccount common.PublicKey, era uint32) (common.PublicKey, string) {
	seed := fmt.Sprintf("stake:%d", era)
	return common.CreateWithSeed(baseAccount, seed, common.StakeProgramID), seed
}

func TestGetSubAccount(t *testing.T) {
	// pubkey := common.PublicKeyFromString("D6tm58oqeMz1VSLNFXNnpyJi8S2A9JHJEp24sDpBo3Dm")
	// subPubKey, _ := GetStakeAccountPubkey(pubkey, 316)
	info, err := c.GetStakeAccountInfo(context.Background(), "D6tm58oqeMz1VSLNFXNnpyJi8S2A9JHJEp24sDpBo3Dm")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info.StakeAccount.Info.Stake.Delegation.Stake)
	t.Log(info.Lamports)

	t.Log(info.StakeAccount.IsStakeAndNoDeactive())
}

func TestAccountInfo(t *testing.T) {
	//user CWVd9HtYD2txbiiSwV3Ss33TGMqUVrS2F5sTs7XZQKWN
	//tx 3mQXBo3FSJ3bvXj9moJx7mW3424mz8DnQjBjrCrzRp3T4bPT5xTtMnzib5Q7NCJf6fLyRSgpWaa5EBfL8EijLi2D
	//block CEjkgbUm169E1bRaeUT7kWg2imJf3j2XZ2qjJCW2CHcU

	wg := sync.WaitGroup{}
	wg.Add(30)

	for i := 0; i < 300; i++ {
		time.Sleep(1 * time.Second)
		accountInfo, err := c.GetAccountInfo(context.Background(), "a1exwPymWZ9Z3ouEsYTrjLt3g7Fsf7DyfSF9BfmGser", client.GetAccountInfoConfig{
			Encoding: client.GetAccountInfoConfigEncodingBase58,
		})
		if err != nil {
			t.Log("err", i, err)
		} else {
			t.Log("success", i, fmt.Sprintf("%+v", accountInfo))
		}
	}

	wg.Wait()
}

func TestGetVersion(t *testing.T) {
	accountActivateInfo, err := c.GetVersion(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", accountActivateInfo)

	// sigs, err := c.GetSignaturesForAddress(context.Background(), "7hUdUTkJLwdcmt3jSEeqx4ep91sm1XwBxMDaJae6bD5D", client.GetSignaturesForAddressConfig{})
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// for _, sig := range sigs {
	// 	t.Log(sig.Signature)
	// }
	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentFinalized,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
func TestGetStakeActivation(t *testing.T) {
	accountActivateInfo, err := c.GetStakeActivation(context.Background(), "G7x84EPhC635pFoBqtWYiHPs5Dc7FsNwxJ6rsdXGeTL6", client.GetStakeActivationConfig{})
	if err != nil {
		if strings.Contains(err.Error(), "account not found") {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	}

	t.Logf("%+v", accountActivateInfo)
}

func TestGetStakeAccountInfo(t *testing.T) {
	accountActivateInfo, err := c.GetStakeAccountInfo(context.Background(), "B6gVbwSfRxonjx6VWHaF5bZycRKRnW9aMw7MtxuzAGTg")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", accountActivateInfo)
	return
	accountActivateInfoBase1, err := c.GetStakeAccountInfo(context.Background(), "4ackc4eexr1DN5eNwzQ5DnNNCAVJiCU84Ev4abUMRKau")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", accountActivateInfoBase1)

	accountActivateInfo2, err := c.GetStakeActivation(context.Background(), "4ackc4eexr1DN5eNwzQ5DnNNCAVJiCU84Ev4abUMRKau", client.GetStakeActivationConfig{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", accountActivateInfo2)
	return
	accountActivateInfoBase, err := c.GetStakeActivation(context.Background(), "J6L2EyHooCuRLKR17ABFmLmCD9Uq9xwDuboJUpZ5wdH7", client.GetStakeActivationConfig{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", accountActivateInfoBase)

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
	t.Logf("%+v", account)

}

func TestGetMultisigTxInfo(t *testing.T) {
	info, err := c.GetMultisigTxAccountInfo(context.Background(), "Gn3Wzs1rbeJcTefiEwZ8c8vJZjNZeSm5WUxbYC5ji74F")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
}

func TestGetStakeManager(t *testing.T) {
	info, err := c.GetStakeManager(context.Background(), "FccgufF6s9WivdfZYKsR52DWyN9fFMyELvKjyJNCeDkj")
	if err != nil {
		t.Fatal(err)
	}
	bts, _ := json.MarshalIndent(info, "", "  ")
	t.Log(fmt.Printf("%s", string(bts)))
}

func TestGetMintManager(t *testing.T) {
	info, err := c.GetMintManager(context.Background(), "55GGz9kCyU8guxJBTtGSscWbM6WS9RsZ4nDmKZU19ubF")
	if err != nil {
		t.Fatal(err)
	}
	bts, _ := json.MarshalIndent(info, "", "  ")
	t.Log(fmt.Printf("%s", string(bts)))
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
	height, err := c.GetBlockHeight(context.Background(), client.GetBlockHeightConfig{client.CommitmentFinalized})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(height)
	info, err := c.GetBlock(context.Background(), 169469698, client.GetBlockConfig{
		Commitment:                     client.CommitmentFinalized,
		MaxSupportedTransactionVersion: &client.DefaultMaxSupportedTransactionVersion,
	})
	if err != nil {
		t.Fatal(err)
	}
	// t.Log(fmt.Printf("%+v", info))

	for _, tx := range info.Transactions {
		for _, log := range tx.Meta.LogMessages {
			t.Log(log)
			if strings.HasPrefix(log, bridgeprog.EventTransferOutPrefix) {
				t.Log(strings.TrimPrefix(log, bridgeprog.ProgramLogPrefix))
			}
		}
		// t.Log(tx.Meta.LogMessages)
	}
}

func TestGetTransaction(t *testing.T) {

	// sigs, _ := c.GetSignaturesForAddress(context.Background(), "EPfxck35M3NJwsjreExLLyQAgAL3y5uWfzddY6cHBrGy", client.GetSignaturesForAddressConfig{})
	// for _, sig := range sigs {
	// 	t.Log(sig.Signature)

	// }
	info3, err := c.GetTransaction(context.Background(), "4WXPE52ce1erEiE6HEnDJijCqjdCKsHDBFin1hcR8A49spDbr8ceyWPeZ9K4GyAf9T4s25kqqArKTsDkM6QizbPq", client.GetTransactionWithLimitConfig{
		Commitment:                     client.CommitmentFinalized,
		MaxSupportedTransactionVersion: &client.DefaultMaxSupportedTransactionVersion})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", info3.Meta)
	t.Log(info3.Meta.Err)
	// blockHeight, err := c.GetBlockHeight(context.Background(), client.GetBlockHeightConfig{client.CommitmentFinalized})
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(blockHeight)
	// slot, err := c.GetSlot(context.Background(), client.GetSlotConfig{client.CommitmentFinalized})
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(slot)
	// time, err := c.GetBlockTime(context.Background(), slot)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(time)
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

	bts, _ := hex.DecodeString("98d9634ad58009cda11726a718073b5ba525d51483cbf8e8bef127cb6b70e900")
	t.Log(common.PublicKeyFromBytes(bts).ToBase58())
}

// FRzXkJ4p1knQkFdBCtLCt8Zuvykr7Wd5yKTrryQV3K51

func TestGetSignaturesForAddress(t *testing.T) {
	info, err := c.GetSignaturesForAddress(context.Background(), "H3mPx8i41Zn4dLC6ZQRBzNRe1cqYdbcDP1WpojnaiAVo", client.GetSignaturesForAddressConfig{
		Until: "2xMo6H3wAerJgBKNPw2c1Mo4n6bbrgLMVozDY9S2mVus2L1ZAks4ebPGHAvxh6oTX9e86TBGmCVNmpthAkT69KLU",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)

	for _, sig := range info {
		usesig := sig.Signature
		t.Log("sig", sig)
		tx, err := c.GetTransaction(context.Background(), usesig, client.GetTransactionWithLimitConfig{
			Commitment:                     client.CommitmentFinalized,
			MaxSupportedTransactionVersion: &client.DefaultMaxSupportedTransactionVersion,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("rpcClient.GetConfirmedTransaction err: %s", err.Error()))
		}
		//skip failed tx
		if tx.Meta.Err != nil {
			if err != nil {
				t.Fatal(err)
			}
			continue
		}
		t.Log("fffff")
		//skip zero instruction
		if len(tx.Transaction.Message.Instructions) == 0 {
			t.Fatal("11111")
			continue
		}
		for _, instruct := range tx.Transaction.Message.Instructions {

			accountKeys := tx.Transaction.Message.AccountKeys
			programIdIndex := instruct.ProgramIDIndex
			if len(accountKeys) <= int(programIdIndex) {
				t.Fatal(fmt.Errorf("accounts or programIdIndex err, %v", tx))
			}
			//skip if it doesn't call  bridge program
			if !strings.EqualFold(accountKeys[programIdIndex], "H3mPx8i41Zn4dLC6ZQRBzNRe1cqYdbcDP1WpojnaiAVo") {
				t.Log("222")
				continue
			}

			// check instruction data
			if len(instruct.Data) == 0 {
				t.Log("3333")
				continue
			}
			dataBts, err := base58.Decode(instruct.Data)
			if err != nil {
				t.Fatal(err)
			}
			if len(dataBts) < 8 {
				t.Log("ttttt")
				continue
			}
			// skip if it doesn't call transferOut func
			if !bytes.Equal(dataBts[:8], bridgeprog.InstructionTransferOut[:]) {
				t.Fatal("call func is not transferOut", "tx", tx)

				continue
			}
			// check bridge account
			if len(instruct.Accounts) == 0 {
				t.Fatal("444")
				continue
			}
			if !strings.EqualFold(accountKeys[instruct.Accounts[0]], "Ev64NXXeKdtBgJbXyuJKEw77pxaw5q4BkUb2eKeV5xDy") {
				t.Fatal("bridge account not equal", "tx", tx)
				continue
			}
			t.Log(tx.Meta.LogMessages)

			for _, logMessage := range tx.Meta.LogMessages {
				if strings.HasPrefix(logMessage, bridgeprog.EventTransferOutPrefix) {
					t.Log("find log", "log", logMessage, "signature", usesig)
					use_log := strings.TrimPrefix(logMessage, bridgeprog.ProgramLogPrefix)
					logBts, err := base64.StdEncoding.DecodeString(use_log)
					if err != nil {
						t.Fatal(err)
					}
					if len(logBts) <= 8 {
						t.Fatal(fmt.Errorf("event pase length err"))
					}

					eventTransferOut := EventTransferOut{}
					err = borsh.Deserialize(&eventTransferOut, logBts[8:])
					if err != nil {
						t.Fatal(err)
					}
					t.Logf("555 %+v", eventTransferOut)

				}

			}
		}

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
	t.Logf("%+v", accountInfo)
	t.Logf("%+v", accountInfo.Mint.ToBase58())
	t.Log(hex.EncodeToString(bridgeprog.InstructionTransferOut[:]))
}

func TestDecodeAccount(t *testing.T) {

	bhbts, err := hex.DecodeString("36ca6a5226f2ae7a258a77e364723d3efe1c873b2db327fff3a243baa681719f")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(common.PublicKeyFromBytes(bhbts).ToBase58())

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

	blockHash, _ := base58.Decode("CNXkUVPfhfmpjtHB5XbJgZ5unkopeeRiEGzZGu6eN2Uq")
	t.Log(hex.EncodeToString(blockHash))

	bts, _ := hex.DecodeString("a9b8dfb4676247ed4f770ef5055f95d324b31e5d99273fec8150a4f4e83e7dc5")
	t.Log(common.PublicKeyFromBytes(bts).ToBase58())

}

func TestGetProgramAccounts(t *testing.T) {
	accounts, err := c.GetProgramAccounts(
		context.Background(),
		common.TokenProgramID.ToBase58(),
		client.GetProgramAccountsConfig{
			WithContext: true,
			Encoding:    client.GetAccountInfoConfigEncodingBase64})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(accounts)
}
