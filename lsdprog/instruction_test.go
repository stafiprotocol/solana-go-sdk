package lsdprog_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mr-tron/base58"
	"github.com/stafiprotocol/solana-go-sdk/client"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/rsolprog"
	"github.com/stafiprotocol/solana-go-sdk/sysprog"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

var rSolProgramIdDev = common.PublicKeyFromString("5N1PkgbPx5Qs3eGaJre16AHsNMRPYM9JSwxXDG83tWX9")
var minterProgramIdDev = common.PublicKeyFromString("HDb577JnkPHLFpfbTg1ncX9jmVHGjzX6S9bgZvNnXjVj")

var validator = common.PublicKeyFromString("vgcDar2pryHvMgPkKaZfh8pQy4BJxv7SpwUG7zinWjG")
var feeRecipient = common.PublicKeyFromString("344uJfqqsMji7jkcoGY6vcHpExsupcygpex6bJvq2ywG") //random
var localClient = []string{"https://api.devnet.solana.com"}

var id = types.AccountFromPrivateKeyBytes([]byte{179, 95, 213, 234, 125, 167, 246, 188, 230, 134, 181, 219, 31, 146, 239, 75, 190, 124, 112, 93, 187, 140, 178, 119, 90, 153, 207, 178, 137, 5, 53, 71, 116, 28, 190, 12, 249, 238, 110, 135, 109, 21, 196, 36, 191, 19, 236, 175, 229, 204, 68, 180, 130, 102, 71, 239, 41, 53, 152, 159, 175, 124, 180, 6})
var admin = types.AccountFromPrivateKeyBytes([]byte{142, 61, 202, 203, 179, 165, 19, 161, 233, 247, 36, 152, 120, 184, 62, 139, 88, 69, 120, 227, 94, 87, 244, 241, 207, 94, 29, 115, 12, 177, 134, 33, 252, 93, 7, 42, 197, 184, 34, 111, 171, 84, 21, 195, 106, 93, 249, 214, 173, 78, 212, 191, 16, 138, 230, 43, 25, 124, 41, 12, 133, 211, 37, 242})
var staker = types.AccountFromPrivateKeyBytes([]byte{90, 111, 119, 62, 149, 35, 16, 87, 135, 90, 47, 202, 31, 47, 85, 140, 65, 17, 88, 226, 229, 193, 38, 9, 103, 255, 72, 136, 150, 213, 224, 50, 47, 183, 28, 18, 35, 161, 125, 133, 219, 9, 124, 130, 85, 200, 82, 75, 251, 232, 246, 67, 137, 238, 173, 105, 146, 126, 153, 90, 190, 88, 30, 81})

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

	stakePool, _, err := common.FindProgramAddress([][]byte{stakeManager.PublicKey.Bytes(), []byte("pool_seed")}, rSolProgramIdDev)
	if err != nil {
		t.Fatal(err)
	}

	stakePoolRent, err := c.GetMinimumBalanceForRentExemption(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}

	stakeManagerRent, err := c.GetMinimumBalanceForRentExemption(context.Background(), rsolprog.StakeManagerAccountLengthDefault)
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
				rSolProgramIdDev,
				stakeManagerRent,
				rsolprog.StakeManagerAccountLengthDefault,
			),
			rsolprog.Initialize(
				rSolProgramIdDev,
				stakeManager.PublicKey,
				stakePool,
				feeRecipient,
				rSolMint,
				admin.PublicKey,
				rsolprog.InitializeData{
					RSolMint:         rSolMint,
					Validator:        validator,
					Bond:             0,
					Unbond:           0,
					Active:           0,
					LatestEra:        611,
					Rate:             1000000000,
					TotalRSolSupply:  107717120,
					TotalProtocolFee: 0,
				},
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

func TestMigrateStakeAccount(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id
	stakeAuthority := id

	stakeManager := common.PublicKeyFromString("CThKc2gVW9fZUaz9g5UEZikMRusPjThKaFGohR1tkQhk")
	stakePool := common.PublicKeyFromString("33aoSpaFKDuKqh35a1N5eGopFH4nr51DENxh9bkzvnKe")
	stakeAccount := common.PublicKeyFromString("5jTc9Q44AF9avDtKGcQKNYNUZbNYtiigBygoj4bLwmdh")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			rsolprog.MigrateStakeAccount(
				rSolProgramIdDev,
				stakeManager,
				stakePool,
				stakeAccount,
				stakeAuthority.PublicKey,
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

	fmt.Println("migrate stake account txHash:", txHash)

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
			rsolprog.AddValidator(
				rSolProgramIdDev,
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
			rsolprog.ReallocStakeManager(
				rSolProgramIdDev,
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

func TestUpgradeStakeManager(t *testing.T) {
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
			rsolprog.UpgradeStakeManager(
				rSolProgramIdDev,
				stakeManager,
				admin.PublicKey,
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

	fmt.Println("UpgradeStakeManager txHash:", txHash)

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
			rsolprog.Redelegate(
				rSolProgramIdDev,
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

	rSolMint := common.PublicKeyFromString("6jnyhgA2dPWDpw1WqgTaCyjp8otXkP5655DQ6RSnwbv5") // rsol_mint.json
	feePayer := id
	from := id

	stakeManager := common.PublicKeyFromString("FccgufF6s9WivdfZYKsR52DWyN9fFMyELvKjyJNCeDkj")
	stakePool := common.PublicKeyFromString("GYoZ5kSumbV2zqCbRYp9jex1AFaCWjbFYQS9URDmswFG")

	mintTo := common.PublicKeyFromString("61wKZgejN7CEL8QYQoxUmjWgQmgx3aNqUX6ivGKKxPkF")

	minterManagerAccount := common.PublicKeyFromString("5ou6pU6ByghiA148DokoVQLpPqGcnww9qS8TQzwcmQcx")
	mintAuthority := common.PublicKeyFromString("66DBm2GT5ELRvXfZ8GVbvurMrwcsxK59rNvGULnsnXvW")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			rsolprog.Stake(
				rSolProgramIdDev,
				stakeManager,
				stakePool,
				from.PublicKey,
				minterManagerAccount,
				rSolMint,
				mintTo,
				mintAuthority,
				minterProgramIdDev,
				1e9,
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
	rSolMint := common.PublicKeyFromString("6jnyhgA2dPWDpw1WqgTaCyjp8otXkP5655DQ6RSnwbv5") // rsol_mint.json
	feePayer := id

	stakeManager := common.PublicKeyFromString("FccgufF6s9WivdfZYKsR52DWyN9fFMyELvKjyJNCeDkj")

	burnRsolAuthority := id
	unstakeAccount := types.NewAccount()

	burnRsolFrom := common.PublicKeyFromString("61wKZgejN7CEL8QYQoxUmjWgQmgx3aNqUX6ivGKKxPkF")
	unstakeAccountRent, err := c.GetMinimumBalanceForRentExemption(context.Background(), rsolprog.UnstakeAccountLengthDefault)
	if err != nil {
		t.Fatal(err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.CreateAccount(
				feePayer.PublicKey,
				unstakeAccount.PublicKey,
				rSolProgramIdDev,
				unstakeAccountRent,
				rsolprog.UnstakeAccountLengthDefault,
			),
			rsolprog.Unstake(
				rSolProgramIdDev,
				stakeManager,
				rSolMint,
				burnRsolFrom,
				burnRsolAuthority.PublicKey,
				unstakeAccount.PublicKey,
				feeRecipient,
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
	recipient := staker

	stakeManager := common.PublicKeyFromString("CThKc2gVW9fZUaz9g5UEZikMRusPjThKaFGohR1tkQhk")
	stakePool := common.PublicKeyFromString("33aoSpaFKDuKqh35a1N5eGopFH4nr51DENxh9bkzvnKe")
	unstakeAccount := common.PublicKeyFromString("Cjxm5bHvrxTcnwgwL2uLSpJDTRzaPkQkSnvSjvyfw71i")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			rsolprog.Withdraw(
				rSolProgramIdDev,
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

func TestEraNew(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id
	stakeManager := common.PublicKeyFromString("FccgufF6s9WivdfZYKsR52DWyN9fFMyELvKjyJNCeDkj")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			rsolprog.EraNew(
				rSolProgramIdDev,
				stakeManager,
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

	fmt.Println("era new txHash:", txHash)

}

func TestEraBond(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id
	rentPayer := id
	stakeAccount := types.NewAccount()
	stakeManager := common.PublicKeyFromString("FccgufF6s9WivdfZYKsR52DWyN9fFMyELvKjyJNCeDkj")
	stakePool := common.PublicKeyFromString("GYoZ5kSumbV2zqCbRYp9jex1AFaCWjbFYQS9URDmswFG")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			rsolprog.EraBond(
				rSolProgramIdDev,
				stakeManager,
				validator,
				stakePool,
				stakeAccount.PublicKey,
				rentPayer.PublicKey,
			),
		},
		Signers:         []types.Account{feePayer, stakeAccount},
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

	fmt.Println("era bond txHash:", txHash)

}

func TestEraUnbond(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id
	rentPayer := id
	splitStakeAccount := types.NewAccount()
	stakeManager := common.PublicKeyFromString("FccgufF6s9WivdfZYKsR52DWyN9fFMyELvKjyJNCeDkj")
	stakePool := common.PublicKeyFromString("GYoZ5kSumbV2zqCbRYp9jex1AFaCWjbFYQS9URDmswFG")
	// stakeAccount := common.PublicKeyFromString("Gawre8qmHnyKs5zpaDFPXSMpZq9D9YBCxmvQ4A18wue3")
	stakeAccount := common.PublicKeyFromString("Db8kTcMbMRrHN1jkXBEAsyDHzPtsHh6Rcm1ae7HHRGSy")
	// stakeAccount := common.PublicKeyFromString("APZuLDgxQNh2zgidnrnhPKAE1HsQmUMSSURQDkM6s7ps")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			rsolprog.EraUnbond(
				rSolProgramIdDev,
				stakeManager,
				stakePool,
				stakeAccount,
				splitStakeAccount.PublicKey,
				validator,
				rentPayer.PublicKey,
			),
		},
		Signers:         []types.Account{feePayer, splitStakeAccount, rentPayer},
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

	fmt.Println("era unbond txHash:", txHash)

}

func TestEraUpdateActive(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id

	stakeManager := common.PublicKeyFromString("FccgufF6s9WivdfZYKsR52DWyN9fFMyELvKjyJNCeDkj")
	// stakeAccount := common.PublicKeyFromString("APZuLDgxQNh2zgidnrnhPKAE1HsQmUMSSURQDkM6s7ps")
	// stakeAccount := common.PublicKeyFromString("Gawre8qmHnyKs5zpaDFPXSMpZq9D9YBCxmvQ4A18wue3")
	stakeAccount := common.PublicKeyFromString("FGnk3JMdmGQDeYCVCtR6DuUPVUUpuRyBN2qAWnf2Zi2z")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			rsolprog.EraUpdateActive(
				rSolProgramIdDev,
				stakeManager,
				stakeAccount,
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

	fmt.Println("era update active txHash:", txHash)

}
func TestEraUpdateRate(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id

	stakeManager := common.PublicKeyFromString("FccgufF6s9WivdfZYKsR52DWyN9fFMyELvKjyJNCeDkj")

	rSolMint := common.PublicKeyFromString("6jnyhgA2dPWDpw1WqgTaCyjp8otXkP5655DQ6RSnwbv5") // rsol_mint.json

	stakePool, _, err := common.FindProgramAddress([][]byte{stakeManager.Bytes(), []byte("pool_seed")}, rSolProgramIdDev)
	if err != nil {
		t.Fatal(err)
	}

	minterManagerAccount := common.PublicKeyFromString("5ou6pU6ByghiA148DokoVQLpPqGcnww9qS8TQzwcmQcx")
	mintAuthority := common.PublicKeyFromString("66DBm2GT5ELRvXfZ8GVbvurMrwcsxK59rNvGULnsnXvW")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			rsolprog.EraUpdateRate(
				rSolProgramIdDev,
				stakeManager,
				stakePool,
				minterManagerAccount,
				rSolMint,
				feeRecipient,
				mintAuthority,
				minterProgramIdDev,
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

	fmt.Println("era update rate txHash:", txHash)

}

func TestEraMerge(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id

	stakeManager := common.PublicKeyFromString("CThKc2gVW9fZUaz9g5UEZikMRusPjThKaFGohR1tkQhk")
	srcStakeAccount := common.PublicKeyFromString("BbHMFJozZ8SDRgMTTHDdbDNsKuBSNLaBV4o16T4mAUKz")
	dstStakeAccount := common.PublicKeyFromString("5jTc9Q44AF9avDtKGcQKNYNUZbNYtiigBygoj4bLwmdh")
	stakePool, _, err := common.FindProgramAddress([][]byte{stakeManager.Bytes(), []byte("pool_seed")}, rSolProgramIdDev)
	if err != nil {
		t.Fatal(err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			rsolprog.EraMerge(
				rSolProgramIdDev,
				stakeManager,
				srcStakeAccount,
				dstStakeAccount,
				stakePool,
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

	fmt.Println("era merge txHash:", txHash)
}

func TestEraWithdraw(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id

	stakeManager := common.PublicKeyFromString("FccgufF6s9WivdfZYKsR52DWyN9fFMyELvKjyJNCeDkj")
	stakeAccount := common.PublicKeyFromString("HATNwBQQsCBxd3G4RNMK7ScgX5CsEhm8e4EK1TT8jcrB")
	stakePool, _, err := common.FindProgramAddress([][]byte{stakeManager.Bytes(), []byte("pool_seed")}, rSolProgramIdDev)
	if err != nil {
		t.Fatal(err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			rsolprog.EraWithdraw(
				rSolProgramIdDev,
				stakeManager,
				stakePool,
				stakeAccount,
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

	fmt.Println("era withdraw txHash:", txHash)
}

func TestFindProgramAddress(t *testing.T) {
	minterManagerAccount := common.PublicKeyFromString("7ZSPwtsvFHcMvSGXtRjtHSR2AkQaix1g82gBm5Y5R3VQ")
	a, _, err := common.FindProgramAddress([][]byte{minterManagerAccount.Bytes(), []byte("mint")}, rSolProgramIdDev)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(a.ToBase58())
}
