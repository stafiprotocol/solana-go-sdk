package bridgeprog_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stafiprotocol/solana-go-sdk/bridgeprog"
	"github.com/stafiprotocol/solana-go-sdk/client"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/sysprog"
	"github.com/stafiprotocol/solana-go-sdk/tokenprog"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

var bridgeProgramIdDev = common.PublicKeyFromString("21Ayg6sP9h9jTwHAu51dafyUa8peZiK32CSNNCp9avR8")
var mintAccountPubkey = common.PublicKeyFromString("ET5vByZ5QyMKH9RRc9EHzhvQgzZMk9W23nYM1MEo77DM")
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
				map[[32]byte]common.PublicKey{
					[32]byte{1}: mintAccountPubkey,
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

	accountABytes, _ := hex.DecodeString("0413c48d32073e9ccd19ad98142ed5da7f2267a1f908675c83b92bec6114e70e054aeed660cbad9eed127d0651caa5f46298f78981e0adad63c8092d947347e4")
	accountBBytes, _ := hex.DecodeString("78af782d0f1e5ae77599072c6ade3898f65d6ca3abbd42f52b6c9fc780c2e5d35218002da19dc6dd6eb058bd776e400a744641aa91ae6d87cffa4a7d087d4ffe")
	accountCBytes, _ := hex.DecodeString("e667de57fd5ec04dc327646b9082a23e766a91d7944b09a6790674c70ebb5fa1298be3455eff7b372e3df7ab3001ecf575269fb13dd882ba7bf45e1d031c6955")
	pdaPubkey, nonce := common.PublicKeyFromString("CBYzJn9qHFzh5Q1KJCjW9nUt2dAnvR5m5voGXBgGePRo"), 254

	bridgeAccountPubkey := common.PublicKeyFromString("592485XJ5MJiwz59JUJusNVPymtcXTt8mnrHYiTH7mCX")
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
	fmt.Println("accountTo", accountTo.PublicKey.ToBase58())

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
				mintAccountPubkey,//mint must == token mintAccount
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
	accountToPubKey :=accountTo.PublicKey
	// accountToPubKey := common.PublicKeyFromString("9RM7zLSC521zDHRQaxFZhnExs3Giba8BtnJYS2peQBJf")

	rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			bridgeprog.CreateMintProposal(
				bridgeProgramIdDev,
				bridgeAccountPubkey,
				mintProposalPubkey,
				accountToPubKey,
				accountA.PublicKey,
				[32]byte{1},
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
