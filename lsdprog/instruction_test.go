package lsdprog_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mr-tron/base58"
	"github.com/stafiprotocol/solana-go-sdk/client"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/lsdprog"
	"github.com/stafiprotocol/solana-go-sdk/sysprog"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

var lsdprogramIdDev = common.PublicKeyFromString("795MBfkwwtAX4fWiFqZcJK8D91P9tqqtiSRrSNhBvGzq")

var feeRecipient = common.PublicKeyFromString("344uJfqqsMji7jkcoGY6vcHpExsupcygpex6bJvq2ywG") //random
// var localClient = []string{"https://api.devnet.solana.com"}
var localClient = []string{"https://solana-dev-rpc.stafi.io"}

var id = types.AccountFromPrivateKeyBytes([]byte{179, 95, 213, 234, 125, 167, 246, 188, 230, 134, 181, 219, 31, 146, 239, 75, 190, 124, 112, 93, 187, 140, 178, 119, 90, 153, 207, 178, 137, 5, 53, 71, 116, 28, 190, 12, 249, 238, 110, 135, 109, 21, 196, 36, 191, 19, 236, 175, 229, 204, 68, 180, 130, 102, 71, 239, 41, 53, 152, 159, 175, 124, 180, 6})
var admin = types.AccountFromPrivateKeyBytes([]byte{142, 61, 202, 203, 179, 165, 19, 161, 233, 247, 36, 152, 120, 184, 62, 139, 88, 69, 120, 227, 94, 87, 244, 241, 207, 94, 29, 115, 12, 177, 134, 33, 252, 93, 7, 42, 197, 184, 34, 111, 171, 84, 21, 195, 106, 93, 249, 214, 173, 78, 212, 191, 16, 138, 230, 43, 25, 124, 41, 12, 133, 211, 37, 242})

func TestInitialize(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	rSolMint := common.PublicKeyFromString("F6KFk1jzBNQis7HdVdUyFLYQ6L3dVZoYL4VwwgQvnjBE") // rsol_mint.json
	feePayer := id
	admin := admin

	stakeManager := types.NewAccount()

	stakePool, _, err := common.FindProgramAddress([][]byte{stakeManager.PublicKey.Bytes(), []byte("pool_seed")}, lsdprogramIdDev)
	if err != nil {
		t.Fatal(err)
	}

	stakePoolRent, err := c.GetMinimumBalanceForRentExemption(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}

	stakeManagerRent, err := c.GetMinimumBalanceForRentExemption(context.Background(), lsdprog.StakeManagerAccountLengthDefault)
	if err != nil {
		t.Fatal(err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.Transfer(
				feePayer.PublicKey,
				stakePool,
				stakePoolRent,
			),
			sysprog.CreateAccount(
				feePayer.PublicKey,
				stakeManager.PublicKey,
				lsdprogramIdDev,
				stakeManagerRent,
				lsdprog.StakeManagerAccountLengthDefault,
			),
			lsdprog.InitializeStakeManager(
				lsdprogramIdDev,
				stakeManager.PublicKey,
				stakePool,
				feeRecipient,
				rSolMint,
				admin.PublicKey,
				admin.PublicKey,
				admin.PublicKey,
				admin.PublicKey,
			),
		},
		Signers:         []types.Account{feePayer, stakeManager, admin},
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

	fmt.Println("createStakeManager txHash:", txHash)
	fmt.Println("stakeManager account:", stakeManager.PublicKey.ToBase58())
	fmt.Println("stakePool account:", stakePool.ToBase58())
	fmt.Println("admin", admin.PublicKey.ToBase58())
	fmt.Println("feePayer:", feePayer.PublicKey.ToBase58())
	fmt.Println("stake pool rent:", stakePoolRent)
	fmt.Println("stake manager rent:", stakeManagerRent)

	//	createStakeManager txHash: 5DekF87gaqf1EN16199WrCrKsLqDfqxHCU8X8poVC7KmAJ1T9aCFPva5xRtvjbAF5gMVUv25cVsnkdSh539QqPeP
	//
	// stakeManager account: CThKc2gVW9fZUaz9g5UEZikMRusPjThKaFGohR1tkQhk
	// stakePool account: 33aoSpaFKDuKqh35a1N5eGopFH4nr51DENxh9bkzvnKe
	// admin Hz81pzkXTqhaZ6v4M6ERCZU4x3aaXrqq2rCafLDwNE1w
	// feePayer: 8pFiM2vyEzyYL7oJqaK2CgHPnARFdziM753rDHWsnhU1
	// stake pool rent: 890880
	// stake manager rent: 13920890880
}

func TestAddValidator(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	adminBts, _ := base58.Decode("2u6qDjEobBnbQuCsW18ELizXx8AUn1SF3JF42c88BbDrw97ADrKg1zw7tokJ1F5fRort8Tzjb9iPfVcDJ4FRXhrd")
	admin := types.AccountFromPrivateKeyBytes(adminBts)
	feePayer := id
	stakeManager := common.PublicKeyFromString("FccgufF6s9WivdfZYKsR52DWyN9fFMyELvKjyJNCeDkj")
	newValidator := common.PublicKeyFromString("5ZWgXcyqrrNpQHCme5SdC5hCeYb2o3fEJhF7Gok3bTVN")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			lsdprog.AddValidator(
				lsdprogramIdDev,
				stakeManager,
				admin.PublicKey,
				newValidator,
			),
		},
		Signers:         []types.Account{feePayer, admin},
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

	fmt.Println("migrate stake account txHash:", txHash)

}

func TestReallocStakeManager(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	adminBts, _ := base58.Decode("2u6qDjEobBnbQuCsW18ELizXx8AUn1SF3JF42c88BbDrw97ADrKg1zw7tokJ1F5fRort8Tzjb9iPfVcDJ4FRXhrd")
	admin := types.AccountFromPrivateKeyBytes(adminBts)
	feePayer := id
	rentPayer := id
	stakeManager := common.PublicKeyFromString("FccgufF6s9WivdfZYKsR52DWyN9fFMyELvKjyJNCeDkj")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			lsdprog.ReallocStakeManager(
				lsdprogramIdDev,
				stakeManager,
				admin.PublicKey,
				rentPayer.PublicKey,
				100,
			),
		},
		Signers:         []types.Account{feePayer, admin, rentPayer},
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

	fmt.Println("ReallocStakeManager stake account txHash:", txHash)

}

func TestRedelegate(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	adminBts, _ := base58.Decode("2u6qDjEobBnbQuCsW18ELizXx8AUn1SF3JF42c88BbDrw97ADrKg1zw7tokJ1F5fRort8Tzjb9iPfVcDJ4FRXhrd")
	admin := types.AccountFromPrivateKeyBytes(adminBts)
	feePayer := id
	rentPayer := id
	stakeManager := common.PublicKeyFromString("FccgufF6s9WivdfZYKsR52DWyN9fFMyELvKjyJNCeDkj")
	newValidator := common.PublicKeyFromString("5ZWgXcyqrrNpQHCme5SdC5hCeYb2o3fEJhF7Gok3bTVN")
	stakePool := common.PublicKeyFromString("GYoZ5kSumbV2zqCbRYp9jex1AFaCWjbFYQS9URDmswFG")
	fromStakeAccount := common.PublicKeyFromString("FGnk3JMdmGQDeYCVCtR6DuUPVUUpuRyBN2qAWnf2Zi2z")

	splitStakeAccount := types.NewAccount()
	toStakeAccount := types.NewAccount()

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			lsdprog.Redelegate(
				lsdprogramIdDev,
				stakeManager,
				admin.PublicKey,
				newValidator,
				stakePool,
				fromStakeAccount,
				splitStakeAccount.PublicKey,
				toStakeAccount.PublicKey,
				rentPayer.PublicKey,
				200000000,
			),
		},
		Signers:         []types.Account{feePayer, admin, splitStakeAccount, toStakeAccount},
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

	fmt.Println("migrate stake account txHash:", txHash)

}

func TestStake(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	rSolMint := common.PublicKeyFromString("6zX4Gn6NXeDF4bqWZCcvoUnScF1xuPNJ8d2g6nftW6k1") // rsol_mint.json
	feePayer := id
	from := id

	stakeManager := common.PublicKeyFromString("JAAGMA3nXSFq3QhSMC9Trkf5hneMGoGRLaGtEkmL1Nmj")
	stakePool := common.PublicKeyFromString("3NCv41v4MTbNPw38yu2ZJFPcujmBDz1FnDA7AaiPPrwB")

	mintTo := common.PublicKeyFromString("AzjXwoUUEzyjTqEykYnH2sRnHA8mAUfHTKh9VStVc8qM")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			lsdprog.Stake(
				lsdprogramIdDev,
				stakeManager,
				stakePool,
				from.PublicKey,
				rSolMint,
				mintTo,
				1e5,
			),
		},
		Signers:         []types.Account{feePayer, from},
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

	fmt.Println("stake txHash:", txHash)

}

func TestUnstake(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}
	rSolMint := common.PublicKeyFromString("99Wg1Vb9vA3S1GRGDWorHwDbGNhUGUCrbE5VeEqmg1p6") // rsol_mint.json
	feePayer := id

	stakeManager := common.PublicKeyFromString("HPaeDVBXtN2xdx3A56MHf4xx9jxqF97QmNA9w8b5zmTz")

	burnRsolAuthority := id
	unstakeAccount := types.NewAccount()

	burnRsolFrom := common.PublicKeyFromString("6m5F4LMeGeHvVD46N4oWorxGftFbNTYb4dUNdDFK5wFG")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			lsdprog.Unstake(
				lsdprogramIdDev,
				stakeManager,
				rSolMint,
				burnRsolFrom,
				burnRsolAuthority.PublicKey,
				unstakeAccount.PublicKey,
				feePayer.PublicKey,
				500000000,
			),
		},
		Signers:         []types.Account{feePayer, burnRsolAuthority, unstakeAccount},
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

	fmt.Println("unstake txHash:", txHash)

}

func TestWithdraw(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id
	recipient := id

	stakeManager := common.PublicKeyFromString("HPaeDVBXtN2xdx3A56MHf4xx9jxqF97QmNA9w8b5zmTz")
	stakePool := common.PublicKeyFromString("7jZyhr2HCfc9FUBfFjKrw9NZr9BToDhRANYFkSJsrs3b")
	unstakeAccount := common.PublicKeyFromString("GkBsi7ia8k2XGDyCnMKZrunyehRYe9sFK1BX7FSr3TGb")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			lsdprog.Withdraw(
				lsdprogramIdDev,
				stakeManager,
				stakePool,
				unstakeAccount,
				recipient.PublicKey,
			),
		},
		Signers:         []types.Account{feePayer},
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

	fmt.Println("withdraw txHash:", txHash)

}

func TestFindProgramAddress(t *testing.T) {
	minterManagerAccount := common.PublicKeyFromString("7ZSPwtsvFHcMvSGXtRjtHSR2AkQaix1g82gBm5Y5R3VQ")
	a, _, err := common.FindProgramAddress([][]byte{minterManagerAccount.Bytes(), []byte("mint")}, lsdprogramIdDev)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(a.ToBase58())
}
