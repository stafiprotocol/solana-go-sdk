package bridgeprog_test

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/bridgeprog"
	"github.com/stafiprotocol/solana-go-sdk/client"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/sysprog"
	"github.com/stafiprotocol/solana-go-sdk/tokenprog"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

var bridgeProgramIdDev = common.PublicKeyFromString("FRzXkJ4p1knQkFdBCtLCt8Zuvykr7Wd5yKTrryQV3K51")
var mintAccountPubkey = common.PublicKeyFromString("9qab2RkbcDbkKjbSAfaN6CCLPgZ7h39npsksVmWPbW6e")
var localClient = "https://api.devnet.solana.com"

func TestCreateBridge(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}
	feePayer := types.AccountFromPrivateKeyBytes([]byte{179, 95, 213, 234, 125, 167, 246, 188, 230, 134, 181, 219, 31, 146, 239, 75, 190, 124, 112, 93, 187, 140, 178, 119, 90, 153, 207, 178, 137, 5, 53, 71, 116, 28, 190, 12, 249, 238, 110, 135, 109, 21, 196, 36, 191, 19, 236, 175, 229, 204, 68, 180, 130, 102, 71, 239, 41, 53, 152, 159, 175, 124, 180, 6})

	_, err = c.RequestAirdrop(context.Background(), feePayer.PublicKey.ToBase58(), 10e9)
	if err != nil {
		fmt.Println(err)
	}

	bridgeAccount := types.NewAccount()
	accountA := types.NewAccount()
	accountB := types.NewAccount()
	accountC := types.NewAccount()
	accountAdmin := types.NewAccount()
	multiSigner, nonce, err := common.FindProgramAddress([][]byte{bridgeAccount.PublicKey.Bytes()}, bridgeProgramIdDev)
	if err != nil {
		fmt.Println(err)
	}
	owners := []common.PublicKey{accountA.PublicKey, accountB.PublicKey, accountC.PublicKey}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.CreateAccount(
				feePayer.PublicKey,
				bridgeAccount.PublicKey,
				bridgeProgramIdDev,
				1000000000,
				client.BridgeAccountLengthDefault,
			),
			bridgeprog.CreateBridge(
				bridgeProgramIdDev,
				bridgeAccount.PublicKey,
				owners,
				2,
				uint8(nonce),
				[]uint8{1},
				map[[32]byte]common.PublicKey{
					[32]byte{1, 2, 3}: mintAccountPubkey,
				},
				accountAdmin.PublicKey,
			),
		},
		Signers:         []types.Account{feePayer, bridgeAccount},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		fmt.Printf("generate tx error, err: %v\n", err)
	}
	txHash, err := c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		fmt.Printf("send tx error, err: %v\n", err)
	}
	seed := "8111"
	mintProposalPubkey := common.CreateWithSeed(feePayer.PublicKey, seed, bridgeProgramIdDev)

	fmt.Println("createBridge txHash:", txHash)
	fmt.Println("feePayer:", feePayer.PublicKey.ToBase58())
	fmt.Println("bridge account:", bridgeAccount.PublicKey.ToBase58())
	fmt.Println("bridge pda account nonce", nonce)
	fmt.Println("pda address:", multiSigner.ToBase58())
	fmt.Println("proposal account:", mintProposalPubkey.ToBase58())
	fmt.Println("accountA", accountA.PublicKey.ToBase58(), hex.EncodeToString(accountA.PrivateKey))
	fmt.Println("accountB", accountB.PublicKey.ToBase58(), hex.EncodeToString(accountB.PrivateKey))
	fmt.Println("accountC", accountC.PublicKey.ToBase58(), hex.EncodeToString(accountC.PrivateKey))
	fmt.Println("accountAdmin", accountAdmin.PublicKey.ToBase58(), hex.EncodeToString(accountAdmin.PrivateKey))

}

func TestBridgeMint(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}
	feePayer := types.AccountFromPrivateKeyBytes([]byte{179, 95, 213, 234, 125, 167, 246, 188, 230, 134, 181, 219, 31, 146, 239, 75, 190, 124, 112, 93, 187, 140, 178, 119, 90, 153, 207, 178, 137, 5, 53, 71, 116, 28, 190, 12, 249, 238, 110, 135, 109, 21, 196, 36, 191, 19, 236, 175, 229, 204, 68, 180, 130, 102, 71, 239, 41, 53, 152, 159, 175, 124, 180, 6})

	_, err = c.RequestAirdrop(context.Background(), feePayer.PublicKey.ToBase58(), 10e9)
	if err != nil {
		fmt.Println(err)
	}

	accountABytes, _ := hex.DecodeString("5342fdb647371247b1b8ea7fa9284bc693b77d511b74ab6fb5ce4ea9e2e30cdc59ef57818a3ef5788d7364220b5d00cd532254efd83e94c7fbaea73926c4f1f4")
	accountBBytes, _ := hex.DecodeString("94ec1ddf5f8d8df9fbe646788d39e052e184f9ecde92074a184f55ced09594d94e08c08466313bb1d9a92e7ac7094493cbfd83e2393dd0b0af1aa8fc8487ac75")
	accountCBytes, _ := hex.DecodeString("0700d35fb18d46e1d890143586962b66784628c3bcdf7e62fb4c0d280a0326beea72d9b4457ce89ca1b0d12f1dbab5cc44be12ddd3d3dedd8b3c7f9bb02b640f")
	pdaPubkey, nonce := common.PublicKeyFromString("36HcYj2ep7wvTTfyjNBErTLvxFG1zFs5E2yCDQTyN2dZ"), 254

	bridgeAccountPubkey := common.PublicKeyFromString("2z4iNM45St7DL6xSPpxthCaUQZYdYZsMf7KPs9m51eoh")
	accountA := types.AccountFromPrivateKeyBytes(accountABytes)
	accountB := types.AccountFromPrivateKeyBytes(accountBBytes)
	accountC := types.AccountFromPrivateKeyBytes(accountCBytes)

	accountTo := types.NewAccount()
	rand.Seed(time.Now().Unix())
	seed := fmt.Sprintf("proposal:%d", rand.Int())
	mintProposalPubkey := common.CreateWithSeed(feePayer.PublicKey, seed, bridgeProgramIdDev)

	fmt.Println("feePayer:", feePayer.PublicKey.ToBase58())
	fmt.Println("bridge account:", bridgeAccountPubkey.ToBase58())
	fmt.Println("bridge pda account nonce", nonce)
	fmt.Println("pda address:", pdaPubkey.ToBase58())
	fmt.Println("mint proposal account:", mintProposalPubkey.ToBase58())
	fmt.Println("accountA", accountA.PublicKey.ToBase58())
	fmt.Println("accountB", accountB.PublicKey.ToBase58())
	fmt.Println("accountC", accountC.PublicKey.ToBase58())
	fmt.Println("accountTo", accountTo.PublicKey.ToBase58(), hex.EncodeToString(accountTo.PrivateKey))

	res, err = c.GetRecentBlockhash(context.Background())
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.CreateAccount(
				feePayer.PublicKey,
				accountTo.PublicKey,
				common.TokenProgramID,
				1000000000,
				165,
			),
			tokenprog.InitializeAccount(
				accountTo.PublicKey,
				mintAccountPubkey, //mint must == token mintAccount
				feePayer.PublicKey,
			),
			sysprog.CreateAccountWithSeed(
				feePayer.PublicKey,
				mintProposalPubkey,
				feePayer.PublicKey,
				bridgeProgramIdDev,
				seed,
				1000000000,
				client.MintProposalInfoLengthDefault,
			),
		},
		Signers:         []types.Account{feePayer, accountTo},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})

	if err != nil {
		fmt.Printf("generate create account tx error, err: %v\n", err)
	}
	txHash, err := c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		fmt.Printf("send tx error, err: %v\n", err)
	}
	fmt.Println("create mint proposal and to account hash ", txHash)

	res, err = c.GetRecentBlockhash(context.Background())
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}
	accountToPubKey := accountTo.PublicKey
	// accountToPubKey := common.PublicKeyFromString("9RM7zLSC521zDHRQaxFZhnExs3Giba8BtnJYS2peQBJf")

	rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			bridgeprog.CreateMintProposal(
				bridgeProgramIdDev,
				bridgeAccountPubkey,
				mintProposalPubkey,
				accountToPubKey,
				accountA.PublicKey,
				[32]byte{1, 2, 3},
				100,
				common.TokenProgramID,
			),
		},
		Signers:         []types.Account{accountA, feePayer},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})

	if err != nil {
		fmt.Printf("generate createTransaction tx error, err: %v\n", err)
	}

	txHash, err = c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		fmt.Printf("send tx error, err: %v\n", err)
	}
	fmt.Println("init mint proposal account txHash:", txHash)

	approve := func(approver types.Account) {
		rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
			Instructions: []types.Instruction{
				bridgeprog.ApproveMintProposal(
					bridgeProgramIdDev,
					bridgeAccountPubkey,
					pdaPubkey,
					mintProposalPubkey,
					approver.PublicKey,
					mintAccountPubkey,
					accountToPubKey,
					common.TokenProgramID,
				),
			},
			Signers:         []types.Account{approver, feePayer},
			FeePayer:        feePayer.PublicKey,
			RecentBlockHash: res.Blockhash,
		})

		if err != nil {
			fmt.Printf("generate Approve tx error, err: %v\n", err)
		}

		// t.Log("rawtx base58:", base58.Encode(rawTx))
		txHash, err = c.SendRawTransaction(context.Background(), rawTx)
		if err != nil {
			fmt.Printf("send tx error, err: %v\n", err)
		}
		fmt.Println("Approve txHash:", txHash, approver.PublicKey.ToBase58())
	}

	approve(accountA)
	approve(accountB)
	approve(accountC)
}

func TestCreateTokenAccount(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}
	feePayer := types.AccountFromPrivateKeyBytes([]byte{179, 95, 213, 234, 125, 167, 246, 188, 230, 134, 181, 219, 31, 146, 239, 75, 190, 124, 112, 93, 187, 140, 178, 119, 90, 153, 207, 178, 137, 5, 53, 71, 116, 28, 190, 12, 249, 238, 110, 135, 109, 21, 196, 36, 191, 19, 236, 175, 229, 204, 68, 180, 130, 102, 71, 239, 41, 53, 152, 159, 175, 124, 180, 6})

	_, err = c.RequestAirdrop(context.Background(), feePayer.PublicKey.ToBase58(), 10e9)
	if err != nil {
		fmt.Println(err)
	}

	accountTo := types.NewAccount()
	rand.Seed(time.Now().Unix())
	// seed := fmt.Sprintf("proposal:%d", rand.Int())
	// mintProposalPubkey := common.CreateWithSeed(feePayer.PublicKey, seed, bridgeProgramIdDev)

	fmt.Println("feePayer:", feePayer.PublicKey.ToBase58())
	// fmt.Println("mint proposal account:", mintProposalPubkey.ToBase58())
	fmt.Println("accountTo", accountTo.PublicKey.ToBase58(), hex.EncodeToString(accountTo.PrivateKey),
		hex.EncodeToString(accountTo.PublicKey.Bytes()))

	res, err = c.GetRecentBlockhash(context.Background())
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.CreateAccount(
				feePayer.PublicKey,
				accountTo.PublicKey,
				common.TokenProgramID,
				1000000000,
				165,
			),
			tokenprog.InitializeAccount(
				accountTo.PublicKey,
				mintAccountPubkey, //mint must == token mintAccount
				feePayer.PublicKey,
			),
		},
		Signers:         []types.Account{feePayer, accountTo},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})

	if err != nil {
		fmt.Printf("generate create account tx error, err: %v\n", err)
	}
	txHash, err := c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		fmt.Printf("send tx error, err: %v\n", err)
	}
	fmt.Println("create mint proposal and to account hash ", txHash)

	res, err = c.GetRecentBlockhash(context.Background())
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}
}

func TestSetResourceId(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}
	feePayer := types.AccountFromPrivateKeyBytes([]byte{179, 95, 213, 234, 125, 167, 246, 188, 230, 134, 181, 219, 31, 146, 239, 75, 190, 124, 112, 93, 187, 140, 178, 119, 90, 153, 207, 178, 137, 5, 53, 71, 116, 28, 190, 12, 249, 238, 110, 135, 109, 21, 196, 36, 191, 19, 236, 175, 229, 204, 68, 180, 130, 102, 71, 239, 41, 53, 152, 159, 175, 124, 180, 6})

	_, err = c.RequestAirdrop(context.Background(), feePayer.PublicKey.ToBase58(), 10e9)
	if err != nil {
		fmt.Println(err)
	}

	adminBytes, _ := hex.DecodeString("77953dc75b228e6a89ab7588163529faf8e7fd992fd1f36641f871425f75990484ec530062f49221821481fe403d34b01f9ad90c3ffb303df44132bfd5648cb7")

	bridgeAccountPubkey := common.PublicKeyFromString("GxjjC75vjSHfSjb6izyRi9JnZBVxYsdzb5dX4zhGKZuv")
	admin := types.AccountFromPrivateKeyBytes(adminBytes)

	fmt.Println("admin", admin.PublicKey.ToBase58())

	res, err = c.GetRecentBlockhash(context.Background())
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			bridgeprog.SetResourceId(
				bridgeProgramIdDev,
				bridgeAccountPubkey,
				admin.PublicKey,
				[32]byte{1},
				common.PublicKeyFromString("5FDo833zdrHtdmxMS4fG6g3A9bsgYN2UwtiArEQSTVM4"),
			),
		},
		Signers:         []types.Account{feePayer, admin},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})

	if err != nil {
		fmt.Printf("generate set resource tx error, err: %v\n", err)
	}
	txHash, err := c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		fmt.Printf("send tx error, err: %v\n", err)
	}
	fmt.Println("set resourceId tx  hash ", txHash)
}

func TestTransferOut(t *testing.T) {
	c := client.NewClient(localClient)

	feePayer := types.AccountFromPrivateKeyBytes([]byte{179, 95, 213, 234, 125, 167, 246, 188, 230, 134, 181, 219, 31, 146, 239, 75, 190, 124, 112, 93, 187, 140, 178, 119, 90, 153, 207, 178, 137, 5, 53, 71, 116, 28, 190, 12, 249, 238, 110, 135, 109, 21, 196, 36, 191, 19, 236, 175, 229, 204, 68, 180, 130, 102, 71, 239, 41, 53, 152, 159, 175, 124, 180, 6})
	_, err := c.RequestAirdrop(context.Background(), feePayer.PublicKey.ToBase58(), 10e9)
	if err != nil {
		fmt.Println(err)
	}
	fromBytes, _ := hex.DecodeString("cf0b31c9a3ca108ffe22d4e9b73af6be36c87fc4cfabe52a938ca60ce28c20143429f41f8636e46a8f7a90a11c1e652787bbee64a60a04650f7f5b8e55f0a739")
	bridgeAccountPubkey := common.PublicKeyFromString("ByorhXUES7EQHxx5epzgD7PM5fyorqE7L4XmAY2Qz6Vm")
	fromAccount := types.AccountFromPrivateKeyBytes(fromBytes)
	receiver, _ := hex.DecodeString("306721211d5404bd9da88e0204360a1a9ab8b87c66c1bc2fcdd37f3c2222cc20")
	fmt.Println("fromAccount", fromAccount.PublicKey.ToBase58())
	chainId := 1
	amount := 5
	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			bridgeprog.TransferOut(
				bridgeProgramIdDev,
				bridgeAccountPubkey,
				feePayer.PublicKey,
				mintAccountPubkey,
				fromAccount.PublicKey,
				common.TokenProgramID,
				uint64(amount),
				receiver,
				uint8(chainId),
			),
		},
		Signers:         []types.Account{feePayer},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})

	if err != nil {
		fmt.Printf("generate set resource tx error, err: %v\n", err)
	}
	txHash, err := c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		fmt.Printf("send tx error, err: %v\n", err)
	}
	fmt.Println("transfer out tx  hash ", txHash)
}


