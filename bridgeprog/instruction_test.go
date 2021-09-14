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

var bridgeProgramIdDev = common.PublicKeyFromString("ExSFgFAnMLSGSY9MJaeBhhCgiwKJ3G2hDv5UsyphMgqi")
var mintAccountPubkey = common.PublicKeyFromString("EuEw3HUYJ8A2HMbCeZKN9CHykAitSdNPZz3HbgorCgvp")
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
	accountE := types.NewAccount()
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
				2000,
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
				accountE.PublicKey,
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
	fmt.Println("accountE", accountE.PublicKey.ToBase58(), hex.EncodeToString(accountE.PrivateKey))

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

	accountABytes, _ := hex.DecodeString("79a1a79cc114c0d7e39085c53991d9458c8a721163d03580d395bf4e386f6232e2d1187f8c7cc750c2675e0fbf34f3355aac23b4cbf56de3aa505e744acad975")
	accountBBytes, _ := hex.DecodeString("9da195a5edcf6eebd8096ee1e72ce1089fae78562499db339d60995eb91e050eac8c6545265b7e20e2d5ece1e6e541fb75fe508a8341fb9778ab5f00bce04a74")
	accountCBytes, _ := hex.DecodeString("29fb8d928f919961c0208e6c834110cfbb4b905092665b502ce3d3ec9bdaf4c1c315673697e073a6f48433b3f9212cc0456a44415d1a20614d1449a6811a858e")
	multiSigner, nonce := common.PublicKeyFromString("E1cbabe5sfiLHYn7ANVoRwWMYFz7p25Aq9apJuhpivzR"), 255

	bridgeAccountPubkey := common.PublicKeyFromString("8B29iREQvQgmiyZSdzRrgEsh56W3m2Mpna5xvGG6jAEf")
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
	fmt.Println("pda address:", multiSigner.ToBase58())
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
				mintAccountPubkey,
				multiSigner,
			),
			sysprog.CreateAccountWithSeed(
				feePayer.PublicKey,
				mintProposalPubkey,
				feePayer.PublicKey,
				bridgeProgramIdDev,
				seed,
				1000000000,
				1000,
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

	rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			bridgeprog.CreateMintProposal(
				bridgeProgramIdDev,
				bridgeAccountPubkey,
				mintProposalPubkey,
				accountA.PublicKey,
				[32]byte{1},
				accountTo.PublicKey,
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
	fmt.Println("Create mint proposal account txHash:", txHash)

	approve := func(approver types.Account) {
		rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
			Instructions: []types.Instruction{
				bridgeprog.ApproveMintProposal(
					bridgeProgramIdDev,
					bridgeAccountPubkey,
					multiSigner,
					mintProposalPubkey,
					approver.PublicKey,
					mintAccountPubkey,
					accountTo.PublicKey,
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
}
