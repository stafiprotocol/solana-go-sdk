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

var multisigProgramIDDev = common.PublicKeyFromString("3STyarNLwxhXay9oCXV1LKHTgrVUn3VFgJ2u2mMQFDuV")
var token = common.PublicKeyFromString("97XzCoNwKUqWuyHixgNxrBKc1RZTuNQau2ho87JjteZH")
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
	multiSigner, nonce, err := common.FindProgramAddress([][]byte{bridgeAccount.PublicKey.Bytes()}, multisigProgramIDDev)
	if err != nil {
		fmt.Println(err)
	}
	owners := []common.PublicKey{accountA.PublicKey, accountB.PublicKey, accountC.PublicKey}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.CreateAccount(
				feePayer.PublicKey,
				bridgeAccount.PublicKey,
				multisigProgramIDDev,
				1000000000,
				2000,
			),
			bridgeprog.CreateBridge(
				multisigProgramIDDev,
				bridgeAccount.PublicKey,
				owners,
				2,
				uint8(nonce),
				map[[32]byte]common.PublicKey{
					[32]byte{1}: token,
				},
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
	transactionAccountPubkey := common.CreateWithSeed(feePayer.PublicKey, seed, multisigProgramIDDev)

	fmt.Println("createBridge txHash:", txHash)
	fmt.Println("feePayer:", feePayer.PublicKey.ToBase58())
	fmt.Println("bridge account:", bridgeAccount.PublicKey.ToBase58())
	fmt.Println("bridge account nonce", nonce)
	fmt.Println("proposal account:", transactionAccountPubkey.ToBase58())
	fmt.Println("multiSigner:", multiSigner.ToBase58())
	fmt.Println("accountA", accountA.PublicKey.ToBase58(), hex.EncodeToString(accountA.PrivateKey))
	fmt.Println("accountB", accountB.PublicKey.ToBase58(), hex.EncodeToString(accountB.PrivateKey))
	fmt.Println("accountC", accountC.PublicKey.ToBase58(), hex.EncodeToString(accountC.PrivateKey))

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

	accountABytes, _ := hex.DecodeString("0e6d725aebffe99ea2df8be16f222a5fcf1fef0b4d4ae9883172eec9d7c93e126899273b7195a741b994231d48d746fc13c32c8a5f040e55e36fe2a710650535")
	accountBBytes, _ := hex.DecodeString("688879a73c93cbf70a719f10d98b055508885ec4db9264d0ec553c8894806b690e284891732850987dd6bddf34a59a36833d971bac36b2ef24b26742ea68ef5b")
	accountCBytes, _ := hex.DecodeString("0d2bbf4118201453c89674f86f1e01a3286bace6d1afa7f3fda8aaa1263ae4388028f190c2780590b8e7656b259c2c40d55180131f2d70c03c8e732379b83a8f")
	multiSigner, nonce := common.PublicKeyFromString("EJceTPUiJJiT9RjqQWXRh65EgTbztqGPUbRgPcwpqv14"), 251

	bridgeAccountPubkey := common.PublicKeyFromString("EbSMYeu5NXowH3vd2rJENML9qAzzQxCRzRyKSGRKFrqE")
	accountA := types.AccountFromPrivateKeyBytes(accountABytes)
	accountB := types.AccountFromPrivateKeyBytes(accountBBytes)
	accountC := types.AccountFromPrivateKeyBytes(accountCBytes)
	accountD := types.NewAccount()
	rand.Seed(time.Now().Unix())
	seed := fmt.Sprintf("8111%d", rand.Int())
	transactionAccountPubkey := common.CreateWithSeed(feePayer.PublicKey, seed, multisigProgramIDDev)

	fmt.Println("feePayer:", feePayer.PublicKey.ToBase58())
	fmt.Println("multisig account:", bridgeAccountPubkey.ToBase58())
	fmt.Println("multisig account nonce", nonce)
	fmt.Println("transaction account:", transactionAccountPubkey.ToBase58())
	fmt.Println("multiSigner:", multiSigner.ToBase58())
	fmt.Println("accountA", accountA.PublicKey.ToBase58())
	fmt.Println("accountB", accountB.PublicKey.ToBase58())
	fmt.Println("accountC", accountC.PublicKey.ToBase58())
	fmt.Println("accountD", accountD.PublicKey.ToBase58())

	res, err = c.GetRecentBlockhash(context.Background())
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	mintInstruct := tokenprog.MintTo(
		token,
		accountD.PublicKey,
		multiSigner,
		[]common.PublicKey{},
		10000000000,
	)
	programIds := make([]common.PublicKey, 0)
	accountMetas := make([][]types.AccountMeta, 0)
	datas := make([][]byte, 0)
	instructions := make([]types.Instruction, 0)

	programIds = append(programIds, common.TokenProgramID)
	accountMetas = append(accountMetas, mintInstruct.Accounts)
	datas = append(datas, mintInstruct.Data)
	instructions = append(instructions, mintInstruct)

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.CreateAccount(
				feePayer.PublicKey,
				accountD.PublicKey,
				common.TokenProgramID,
				1000000000,
				165,
			),
			tokenprog.InitializeAccount(
				accountD.PublicKey,
				token,
				multiSigner,
			),
			sysprog.CreateAccountWithSeed(
				feePayer.PublicKey,
				transactionAccountPubkey,
				feePayer.PublicKey,
				multisigProgramIDDev,
				seed,
				1000000000,
				1000,
			),
		},
		Signers:         []types.Account{feePayer, accountD},
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
	fmt.Println("create transaction account hash ", txHash)

	res, err = c.GetRecentBlockhash(context.Background())
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			bridgeprog.CreateProposal(
				multisigProgramIDDev,
				programIds,
				accountMetas,
				datas,
				bridgeAccountPubkey,
				transactionAccountPubkey,
				accountA.PublicKey,
			),
		},
		Signers:         []types.Account{accountA, feePayer},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})

	if err != nil {
		fmt.Printf("generate createTransaction tx error, err: %v\n", err)
	}

	// t.Log("rawtx base58:", base58.Encode(rawTx))
	txHash, err = c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		fmt.Printf("send tx error, err: %v\n", err)
	}
	fmt.Println("Create Transaction txHash:", txHash)

	remainingAccounts := bridgeprog.GetRemainAccounts(instructions)
	rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			bridgeprog.Approve(
				multisigProgramIDDev,
				bridgeAccountPubkey,
				multiSigner,
				transactionAccountPubkey,
				accountB.PublicKey,
				remainingAccounts,
			),
		},
		Signers:         []types.Account{accountB, feePayer},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})

	if err != nil {
		fmt.Printf("b generate Approve tx error, err: %v\n", err)
	}

	// t.Log("rawtx base58:", base58.Encode(rawTx))
	txHash, err = c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		fmt.Printf("send tx error, err: %v\n", err)
	}
	fmt.Println("b Approve txHash:", txHash)

	rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			bridgeprog.Approve(
				multisigProgramIDDev,
				bridgeAccountPubkey,
				multiSigner,
				transactionAccountPubkey,
				accountA.PublicKey,
				remainingAccounts,
			),
		},
		Signers:         []types.Account{accountA, feePayer},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})

	if err != nil {
		fmt.Printf("a generate Approve tx error, err: %v\n", err)
	}

	// t.Log("rawtx base58:", base58.Encode(rawTx))
	txHash, err = c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		fmt.Printf("send tx error, err: %v\n", err)
	}
	fmt.Println("a Approve txHash:", txHash)

	rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			bridgeprog.Approve(
				multisigProgramIDDev,
				bridgeAccountPubkey,
				multiSigner,
				transactionAccountPubkey,
				accountC.PublicKey,
				remainingAccounts,
			),
		},
		Signers:         []types.Account{accountC, feePayer},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})

	if err != nil {
		fmt.Printf("c generate Approve tx error, err: %v\n", err)
	}

	// t.Log("rawtx base58:", base58.Encode(rawTx))
	txHash, err = c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		fmt.Printf("send tx error, err: %v\n", err)
	}
	fmt.Println("c Approve txHash:", txHash)
}
