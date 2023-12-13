package minterprog_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stafiprotocol/solana-go-sdk/client"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/minterprog"
	"github.com/stafiprotocol/solana-go-sdk/sysprog"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

var minterProgramIdDev = common.PublicKeyFromString("HDb577JnkPHLFpfbTg1ncX9jmVHGjzX6S9bgZvNnXjVj")
var rsolMint = common.PublicKeyFromString("Fa8Xy1hHUQejskxk4XEbbnPfAg2igs53tayBVdN3nXXo")
var localClient = []string{"https://api.devnet.solana.com"}

var id = types.AccountFromPrivateKeyBytes([]byte{179, 95, 213, 234, 125, 167, 246, 188, 230, 134, 181, 219, 31, 146, 239, 75, 190, 124, 112, 93, 187, 140, 178, 119, 90, 153, 207, 178, 137, 5, 53, 71, 116, 28, 190, 12, 249, 238, 110, 135, 109, 21, 196, 36, 191, 19, 236, 175, 229, 204, 68, 180, 130, 102, 71, 239, 41, 53, 152, 159, 175, 124, 180, 6})
var id2 = types.AccountFromPrivateKeyBytes([]byte{12, 118, 31, 12, 142, 132, 83, 25, 46, 59, 254, 109, 3, 206, 1, 153, 178, 123, 50, 146, 96, 83, 237, 214, 94, 147, 87, 127, 42, 39, 97, 56, 62, 33, 157, 80, 212, 54, 114, 143, 17, 90, 115, 208, 188, 27, 52, 104, 139, 106, 39, 235, 193, 194, 9, 133, 204, 227, 135, 55, 224, 76, 179, 74})
var admin = types.AccountFromPrivateKeyBytes([]byte{142, 61, 202, 203, 179, 165, 19, 161, 233, 247, 36, 152, 120, 184, 62, 139, 88, 69, 120, 227, 94, 87, 244, 241, 207, 94, 29, 115, 12, 177, 134, 33, 252, 93, 7, 42, 197, 184, 34, 111, 171, 84, 21, 195, 106, 93, 249, 214, 173, 78, 212, 191, 16, 138, 230, 43, 25, 124, 41, 12, 133, 211, 37, 242})

func TestInitialize(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id
	admin := admin
	extMintAuthority := id2

	minterManagerAccount := types.NewAccount()

	extMintAthorities := []common.PublicKey{extMintAuthority.PublicKey}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.CreateAccount(
				feePayer.PublicKey,
				minterManagerAccount.PublicKey,
				minterProgramIdDev,
				1000000000,
				minterprog.MinterManagerAccountLengthDefault,
			),
			minterprog.Initialize(
				minterProgramIdDev,
				minterManagerAccount.PublicKey,
				rsolMint,
				admin.PublicKey,
				extMintAthorities,
			),
		},
		Signers:         []types.Account{feePayer, minterManagerAccount},
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

	fmt.Println("createMinterManager txHash:", txHash)
	fmt.Println("minterManager account:", minterManagerAccount.PublicKey.ToBase58())
	fmt.Println("admin", admin.PublicKey.ToBase58())
	fmt.Println("feePayer:", feePayer.PublicKey.ToBase58())

}

func TestMintToken(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id
	extMintAuthority := id2

	minterManagerAccount := common.PublicKeyFromString("7ZSPwtsvFHcMvSGXtRjtHSR2AkQaix1g82gBm5Y5R3VQ")
	mintToAccount := common.PublicKeyFromString("AN22h55iQBwiivXiKNZuGEA28PzHAA1JdgpnD3rrquxo")

	mintAuthority := common.PublicKeyFromString("GBm6iLyc85BA7RTguvv21chvBk1svN1BCqMfWB57fARe")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{

			minterprog.MintToken(
				minterProgramIdDev,
				minterManagerAccount,
				rsolMint,
				mintToAccount,
				mintAuthority,
				extMintAuthority.PublicKey,
				common.TokenProgramID,
				1111),
		},
		Signers:         []types.Account{feePayer, extMintAuthority},
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

	fmt.Println("mintToken txHash:", txHash)

}

func TestSetExtMintAuthorities(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	feePayer := id
	admin := admin

	minterManagerAccount := common.PublicKeyFromString("7ZSPwtsvFHcMvSGXtRjtHSR2AkQaix1g82gBm5Y5R3VQ")

	extMintAuthority := id2
	extMintAuthorityStakePool := common.PublicKeyFromString("55Z1PVDQuC9zXVLN6wyWBRGZ1qggwyXaKYMge6xNZBvt")

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			minterprog.SetExtMintAuthorities(
				minterProgramIdDev,
				minterManagerAccount,
				admin.PublicKey,
				[]common.PublicKey{extMintAuthority.PublicKey, extMintAuthorityStakePool}),
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

	fmt.Println("SetExtMintAuthorities txHash:", txHash)

}
