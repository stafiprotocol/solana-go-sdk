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

var bridgeProgramIdDev = common.PublicKeyFromString("GF5hXVTvkErn2LTL5myFVbgqPXHnZYj2CkVGU6ZTEtyK")
var minterProgramIdDev = common.PublicKeyFromString("HDb577JnkPHLFpfbTg1ncX9jmVHGjzX6S9bgZvNnXjVj")
var localClient = []string{"https://api.devnet.solana.com"}
var id = types.AccountFromPrivateKeyBytes([]byte{179, 95, 213, 234, 125, 167, 246, 188, 230, 134, 181, 219, 31, 146, 239, 75, 190, 124, 112, 93, 187, 140, 178, 119, 90, 153, 207, 178, 137, 5, 53, 71, 116, 28, 190, 12, 249, 238, 110, 135, 109, 21, 196, 36, 191, 19, 236, 175, 229, 204, 68, 180, 130, 102, 71, 239, 41, 53, 152, 159, 175, 124, 180, 6})
var id2 = types.AccountFromPrivateKeyBytes([]byte{12, 118, 31, 12, 142, 132, 83, 25, 46, 59, 254, 109, 3, 206, 1, 153, 178, 123, 50, 146, 96, 83, 237, 214, 94, 147, 87, 127, 42, 39, 97, 56, 62, 33, 157, 80, 212, 54, 114, 143, 17, 90, 115, 208, 188, 27, 52, 104, 139, 106, 39, 235, 193, 194, 9, 133, 204, 227, 135, 55, 224, 76, 179, 74})

func TestCreateBridge(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	var mintAccountPubkey = common.PublicKeyFromString("F6KFk1jzBNQis7HdVdUyFLYQ6L3dVZoYL4VwwgQvnjBE") // rsol_mint.json

	feePayer := id

	bridgeAccount := types.NewAccount()
	accountA := types.NewAccount()
	accountB := types.NewAccount()
	accountC := types.NewAccount()
	feeReceiver := types.NewAccount()
	accountAdmin := id2
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
					{1, 2, 3}: mintAccountPubkey,
				},
				accountAdmin.PublicKey,
				feeReceiver.PublicKey,
				map[uint8]uint64{1: 2},
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

	fmt.Println("createBridge txHash:", txHash)
	fmt.Println("feePayer:", feePayer.PublicKey.ToBase58())
	fmt.Println("bridge account:", bridgeAccount.PublicKey.ToBase58())
	fmt.Println("bridge pda account nonce", nonce)
	fmt.Println("pda address:", multiSigner.ToBase58())
	fmt.Println("accountA", accountA.PublicKey.ToBase58(), hex.EncodeToString(accountA.PrivateKey))
	fmt.Println("accountB", accountB.PublicKey.ToBase58(), hex.EncodeToString(accountB.PrivateKey))
	fmt.Println("accountC", accountC.PublicKey.ToBase58(), hex.EncodeToString(accountC.PrivateKey))
	fmt.Println("accountAdmin", accountAdmin.PublicKey.ToBase58(), hex.EncodeToString(accountAdmin.PrivateKey))
	fmt.Println("feeReceiver", feeReceiver.PublicKey.ToBase58(), hex.EncodeToString(feeReceiver.PrivateKey))

}

func TestBridgeMint(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}
	feePayer := types.AccountFromPrivateKeyBytes([]byte{179, 95, 213, 234, 125, 167, 246, 188, 230, 134, 181, 219, 31, 146, 239, 75, 190, 124, 112, 93, 187, 140, 178, 119, 90, 153, 207, 178, 137, 5, 53, 71, 116, 28, 190, 12, 249, 238, 110, 135, 109, 21, 196, 36, 191, 19, 236, 175, 229, 204, 68, 180, 130, 102, 71, 239, 41, 53, 152, 159, 175, 124, 180, 6})

	_, err = c.RequestAirdrop(context.Background(), feePayer.PublicKey.ToBase58(), 10e9)
	if err != nil {
		fmt.Println(err)
	}

	var mintAccountPubkey = common.PublicKeyFromString("ApyYYc8URTrmvzTko5ffYuyJCZtVnULWq5qxM2tm1mYj")
	accountABytes, _ := hex.DecodeString("71bfe36bed6af18e074a703168beb568fe2032a03c6e424ab7193b392cb5a07b2fb30c0b530a85f2653c4137ba50ca993588d729676a6955a33ed43263a961e3")
	accountBBytes, _ := hex.DecodeString("d7a82e87057e2a73e573e3c4b9d6699f388c06733b787fac893050c5b5a7bdb86d6a21fd14eecabcb8b5f6d0bc3d33854f6f6dbd3be08fb9ef64ea370e6ac830")
	accountCBytes, _ := hex.DecodeString("94e5bacc1bb1afab057900f103ed6cd05567556f196305b15d732e6ae4e099aeebee5c0fcb5f30f2c86fc4736dd3b39b5682947333f852c40c491953f2b0767d")
	pdaPubkey, nonce := common.PublicKeyFromString("3WEkp3TcCBSfV78jnLFpj6xfKUfEForJ2TecSgjsw2Qt"), 255

	bridgeAccountPubkey := common.PublicKeyFromString("AvwyyCXyrUSvunnTmDAdUfqQJYRbvYzpcW5kuf6FEMF")
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

	res, err = c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
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

	res, err = c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}
	accountToPubKey := accountTo.PublicKey

	rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			bridgeprog.CreateMintProposal(
				bridgeProgramIdDev,
				bridgeAccountPubkey,
				mintProposalPubkey,
				accountToPubKey,
				accountA.PublicKey,
				[32]byte{1, 2, 3},
				1000,
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
	minterManagerAccount := common.PublicKeyFromString("55GGz9kCyU8guxJBTtGSscWbM6WS9RsZ4nDmKZU19ubF")
	mintAuthority := common.PublicKeyFromString("8fXWpVJfVyeh6RnS3p1FtNV6iEPxqddgw1Xa2BHyLxvV")
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
					minterManagerAccount,
					mintAuthority,
					minterProgramIdDev,
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

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}
	var mintAccountPubkey = common.PublicKeyFromString("ApyYYc8URTrmvzTko5ffYuyJCZtVnULWq5qxM2tm1mYj")
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

	res, err = c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
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

	res, err = c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}
}

func TestSetResourceId(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
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

	res, err = c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
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

func TestSetMintAuthorities(t *testing.T) {
	c := client.NewClient(localClient)

	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}
	feePayer := id
	bridgeAccountPubkey := common.PublicKeyFromString("51pE2sAyDAkmPxv1FLFXQsZTyuPMkx23XfZ44UdkJM3T")
	admin := id2
	bridgeSigner := common.PublicKeyFromString("9HaoxMY56uWzyrv7LXv8ToHYxXx7Kzubvc6vU4uipa9F")
	newMintAuthority := common.PublicKeyFromString("8fXWpVJfVyeh6RnS3p1FtNV6iEPxqddgw1Xa2BHyLxvV")
	var mintAccountPubkey = common.PublicKeyFromString("F6KFk1jzBNQis7HdVdUyFLYQ6L3dVZoYL4VwwgQvnjBE") // rsol_mint.json

	res, err = c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
	if err != nil {
		fmt.Printf("get recent block hash error, err: %v\n", err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			bridgeprog.SetMintAuthority(
				bridgeProgramIdDev,
				bridgeAccountPubkey,
				admin.PublicKey,
				bridgeSigner,
				mintAccountPubkey,
				newMintAuthority,
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
	fmt.Println("set mint authority tx  hash ", txHash)
}

func TestTransferOut(t *testing.T) {
	c := client.NewClient(localClient)
	var mintAccountPubkey = common.PublicKeyFromString("ApyYYc8URTrmvzTko5ffYuyJCZtVnULWq5qxM2tm1mYj")
	feePayer := types.AccountFromPrivateKeyBytes([]byte{179, 95, 213, 234, 125, 167, 246, 188, 230, 134, 181, 219, 31, 146, 239, 75, 190, 124, 112, 93, 187, 140, 178, 119, 90, 153, 207, 178, 137, 5, 53, 71, 116, 28, 190, 12, 249, 238, 110, 135, 109, 21, 196, 36, 191, 19, 236, 175, 229, 204, 68, 180, 130, 102, 71, 239, 41, 53, 152, 159, 175, 124, 180, 6})
	_, err := c.RequestAirdrop(context.Background(), feePayer.PublicKey.ToBase58(), 10e9)
	if err != nil {
		fmt.Println(err)
	}
	fromBytes, _ := hex.DecodeString("a38dd7b13e4a1a889d1289f3c5ad6413259755601c86fb8178f49dec95e629c5fd9e41c290546ab525e153fbc0e460fe2c150f492eabdcd8145e5c6d1ebe99ff")
	fromAccount := types.AccountFromPrivateKeyBytes(fromBytes)
	fmt.Println("fromAccount", fromAccount.PublicKey.ToBase58())

	bridgeAccountPubkey := common.PublicKeyFromString("AvwyyCXyrUSvunnTmDAdUfqQJYRbvYzpcW5kuf6FEMF")

	receiver, _ := hex.DecodeString("306721211d5404bd9da88e0204360a1a9ab8b87c66c1bc2fcdd37f3c2222cc20")
	feeReceiverBts, _ := hex.DecodeString("2c9ae658eb59cddeb35fe5a9f35ee1e8d10f042b93a916d24f1e4e3c1c762a80c467557cbe571da11d7b6ac731d0d9f23260347062b6a55fa2aa3a26a160dabe")
	feeReceiver := types.AccountFromPrivateKeyBytes(feeReceiverBts)
	fmt.Println("feeReceiver", feeReceiver.PublicKey.ToBase58())

	chainId := 1
	amount := 1
	res, err := c.GetLatestBlockhash(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: client.CommitmentConfirmed,
	})
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
				feeReceiver.PublicKey,
				common.TokenProgramID,
				common.SystemProgramID,
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

func TestSerilize(t *testing.T) {
	receiver, _ := hex.DecodeString("306721211d5404bd9da88e0204360a1a9ab8b87c66c1bc2fcdd37f3c2222cc20")
	data, _ := common.SerializeData(struct {
		Instruction bridgeprog.Instruction
		Amount      uint64
		Receiver    []byte
		DestChainId uint8
	}{
		Instruction: bridgeprog.InstructionTransferOut,
		Amount:      10000000,
		Receiver:    receiver,
		DestChainId: 1,
	})
	t.Log(hex.EncodeToString(data))
}
