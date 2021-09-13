package bridgeprog

import (
	"crypto/sha256"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

type Instruction [8]byte

var (
	InstructionCreateBridge        Instruction
	InstructionCreateMintProposal  Instruction
	InstructionApproveMintProposal Instruction
	InstructionSetOwners           Instruction
	InstructionChangeThreshold     Instruction
	InstructionSetResourceId       Instruction
)

func init() {
	createBridgeHash := sha256.Sum256([]byte("global:create_bridge"))
	copy(InstructionCreateBridge[:], createBridgeHash[:8])
	createMintProposalHash := sha256.Sum256([]byte("global:create_mint_proposal"))
	copy(InstructionCreateMintProposal[:], createMintProposalHash[:8])
	approveMintProposalHash := sha256.Sum256([]byte("global:approve_mint_proposal"))
	copy(InstructionApproveMintProposal[:], approveMintProposalHash[:8])
	setOwnersHash := sha256.Sum256([]byte("global:set_owners"))
	copy(InstructionSetOwners[:], setOwnersHash[:8])
	changeThresholdHash := sha256.Sum256([]byte("global:change_threshold"))
	copy(InstructionChangeThreshold[:], changeThresholdHash[:8])
	setResourceIdHash := sha256.Sum256([]byte("global:set_resource_id"))
	copy(InstructionSetResourceId[:], setResourceIdHash[:8])
}

func CreateBridge(
	programID,
	bridgeAccount common.PublicKey,
	owners []common.PublicKey,
	threshold uint64,
	nonce uint8,
	resourceIdToMint map[[32]byte]common.PublicKey,
	admin common.PublicKey) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction           Instruction
		Owners                []common.PublicKey
		Threshold             uint64
		Nonce                 uint8
		ResourceIdToTokenProg map[[32]byte]common.PublicKey
		Admin                 common.PublicKey
	}{
		Instruction:           InstructionCreateBridge,
		Owners:                owners,
		Threshold:             threshold,
		Nonce:                 nonce,
		ResourceIdToTokenProg: resourceIdToMint,
		Admin:                 admin,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: programID,
		Accounts: []types.AccountMeta{
			{PubKey: bridgeAccount, IsSigner: false, IsWritable: true},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

type ProposalUsedAccount struct {
	Pubkey     common.PublicKey
	IsSigner   bool
	IsWritable bool
}

func CreateMintProposal(
	programID common.PublicKey,
	bridgeAccount common.PublicKey,
	proposalAccount common.PublicKey,
	proposerAccount common.PublicKey,
	resourceId [32]byte,
	to common.PublicKey,
	amount uint64,
	tokenProgram common.PublicKey,
) types.Instruction {

	data, err := common.SerializeData(struct {
		Instruction  Instruction
		ResourceId   [32]byte
		To           common.PublicKey
		Amount       uint64
		TokenProgram common.PublicKey
	}{
		Instruction:  InstructionCreateMintProposal,
		ResourceId:   resourceId,
		To:           to,
		Amount:       amount,
		TokenProgram: tokenProgram,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: programID,
		Accounts: []types.AccountMeta{
			{PubKey: bridgeAccount, IsSigner: false, IsWritable: false},
			{PubKey: proposalAccount, IsSigner: false, IsWritable: true},
			{PubKey: proposerAccount, IsSigner: true, IsWritable: false},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func ApproveMintProposal(
	programID,
	bridgeAccount,
	multiSiner,
	proposalAccount,
	approverAccount common.PublicKey,
	mintAccount common.PublicKey,
	to common.PublicKey,
	tokenProgram common.PublicKey,
) types.Instruction {

	data, err := common.SerializeData(struct {
		Instruction Instruction
	}{
		Instruction: InstructionApproveMintProposal,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: bridgeAccount, IsSigner: false, IsWritable: false},
		{PubKey: multiSiner, IsSigner: false, IsWritable: false},
		{PubKey: proposalAccount, IsSigner: false, IsWritable: true},
		{PubKey: approverAccount, IsSigner: true, IsWritable: false},
		{PubKey: mintAccount, IsSigner: false, IsWritable: true},
		{PubKey: to, IsSigner: false, IsWritable: true},
		{PubKey: tokenProgram, IsSigner: false, IsWritable: true},
	}

	return types.Instruction{
		ProgramID: programID,
		Accounts:  accounts,
		Data:      data,
	}
}

func ChangeThreshold(
	programID,
	bridgeAccount,
	multiSiner common.PublicKey,
	threshold uint64) types.Instruction {

	data, err := common.SerializeData(struct {
		Instruction Instruction
		Threshold   uint64
	}{
		Instruction: InstructionChangeThreshold,
		Threshold:   threshold,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: bridgeAccount, IsSigner: false, IsWritable: true},
		{PubKey: multiSiner, IsSigner: true, IsWritable: false},
	}

	return types.Instruction{
		ProgramID: programID,
		Accounts:  accounts,
		Data:      data,
	}
}

func GetRemainAccounts(ins []types.Instruction) []types.AccountMeta {
	accountMetas := []types.AccountMeta{}
	accountMap := make(map[string]types.AccountMeta)
	for _, in := range ins {
		accountMetas = append(accountMetas, types.AccountMeta{
			PubKey:     in.ProgramID,
			IsSigner:   false,
			IsWritable: false,
		})
		accountMetas = append(accountMetas, in.Accounts...)
	}
	for i, _ := range accountMetas {
		addrStr := accountMetas[i].PubKey.ToBase58()
		accountMetas[i].IsWritable = accountMap[addrStr].IsWritable || accountMetas[i].IsWritable
		accountMetas[i].IsSigner = false
		accountMap[addrStr] = accountMetas[i]
	}

	ret := make([]types.AccountMeta, 0)
	for _, value := range accountMap {
		ret = append(ret, value)
	}
	return ret
}
