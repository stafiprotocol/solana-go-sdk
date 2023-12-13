package rsolprog_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stafiprotocol/solana-go-sdk/client"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/rsolprog"
	"github.com/stafiprotocol/solana-go-sdk/sysprog"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

var rSolProgramIdDev = common.PublicKeyFromString("5N1PkgbPx5Qs3eGaJre16AHsNMRPYM9JSwxXDG83tWX9")
var minterProgramIdDev = common.PublicKeyFromString("HDb577JnkPHLFpfbTg1ncX9jmVHGjzX6S9bgZvNnXjVj")

var rSolMint = common.PublicKeyFromString("Fa8Xy1hHUQejskxk4XEbbnPfAg2igs53tayBVdN3nXXo")
var validator = common.PublicKeyFromString("FwR3PbjS5iyqzLiLugrBqKSa5EKZ4vK9SKs7eQXtT59f")
var feeRecipient = common.PublicKeyFromString("7U6YTbX2NZb1nTyTUZtqvq5c3EgGpvoAEa956yhz8m6w")
var localClient = []string{"https://api.devnet.solana.com"}

var id = types.AccountFromPrivateKeyBytes([]byte{179, 95, 213, 234, 125, 167, 246, 188, 230, 134, 181, 219, 31, 146, 239, 75, 190, 124, 112, 93, 187, 140, 178, 119, 90, 153, 207, 178, 137, 5, 53, 71, 116, 28, 190, 12, 249, 238, 110, 135, 109, 21, 196, 36, 191, 19, 236, 175, 229, 204, 68, 180, 130, 102, 71, 239, 41, 53, 152, 159, 175, 124, 180, 6})
var id2 = types.AccountFromPrivateKeyBytes([]byte{12, 118, 31, 12, 142, 132, 83, 25, 46, 59, 254, 109, 3, 206, 1, 153, 178, 123, 50, 146, 96, 83, 237, 214, 94, 147, 87, 127, 42, 39, 97, 56, 62, 33, 157, 80, 212, 54, 114, 143, 17, 90, 115, 208, 188, 27, 52, 104, 139, 106, 39, 235, 193, 194, 9, 133, 204, 227, 135, 55, 224, 76, 179, 74})
var staker = types.AccountFromPrivateKeyBytes([]byte{90, 111, 119, 62, 149, 35, 16, 87, 135, 90, 47, 202, 31, 47, 85, 140, 65, 17, 88, 226, 229, 193, 38, 9, 103, 255, 72, 136, 150, 213, 224, 50, 47, 183, 28, 18, 35, 161, 125, 133, 219, 9, 124, 130, 85, 200, 82, 75, 251, 232, 246, 67, 137, 238, 173, 105, 146, 126, 153, 90, 190, 88, 30, 81})

func TestInitialize(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id
	admin := id

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
					TotalRSolSupply:  0,
					TotalProtocolFee: 0,
				},
			),
		},
		Signers:         []types.Account{feePayer, stakeManager},
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

	// stakemanager 3pumTykbNdCnSHPfMchX7T2qFmhDeNsTU4nVcFxJxC8K
	// sakepool 55Z1PVDQuC9zXVLN6wyWBRGZ1qggwyXaKYMge6xNZBvt
}

func TestStake(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id
	from := staker

	stakeManager := common.PublicKeyFromString("3pumTykbNdCnSHPfMchX7T2qFmhDeNsTU4nVcFxJxC8K")
	stakePool := common.PublicKeyFromString("55Z1PVDQuC9zXVLN6wyWBRGZ1qggwyXaKYMge6xNZBvt")

	mintTo := common.PublicKeyFromString("AN22h55iQBwiivXiKNZuGEA28PzHAA1JdgpnD3rrquxo")

	minterManagerAccount := common.PublicKeyFromString("7ZSPwtsvFHcMvSGXtRjtHSR2AkQaix1g82gBm5Y5R3VQ")
	mintAuthority := common.PublicKeyFromString("GBm6iLyc85BA7RTguvv21chvBk1svN1BCqMfWB57fARe")

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
				1100000,
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

func TestFindProgramAddress(t *testing.T) {
	minterManagerAccount := common.PublicKeyFromString("7ZSPwtsvFHcMvSGXtRjtHSR2AkQaix1g82gBm5Y5R3VQ")
	a, _, err := common.FindProgramAddress([][]byte{minterManagerAccount.Bytes(), []byte("mint")}, rSolProgramIdDev)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(a.ToBase58())
}
