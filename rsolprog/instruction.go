package rsolprog

import (
	"crypto/sha256"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

type Instruction [8]byte
type Event [8]byte

var (
	InstructionInitialize Instruction
	InstructionStake      Instruction

	StakeManagerAccountLengthDefault = uint64(2000000)
)

func init() {
	initializeHash := sha256.Sum256([]byte("global:initialize"))
	copy(InstructionInitialize[:], initializeHash[:8])

	stakeHash := sha256.Sum256([]byte("global:stake"))
	copy(InstructionStake[:], stakeHash[:8])

}

type InitializeData struct {
	RSolMint         common.PublicKey
	Validator        common.PublicKey
	Bond             uint64
	Unbond           uint64
	Active           uint64
	LatestEra        uint64
	Rate             uint64
	TotalRSolSupply  uint64
	TotalProtocolFee uint64
}

func Initialize(
	rSolProgramID,
	stakeManager,
	stakePool,
	feeRecipient,
	rSolMint,
	admin common.PublicKey,
	initData InitializeData,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction    Instruction
		InitializeData InitializeData
	}{
		Instruction:    InstructionInitialize,
		InitializeData: initData,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: rSolProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: stakePool, IsSigner: false, IsWritable: true},
			{PubKey: feeRecipient, IsSigner: false, IsWritable: false},
			{PubKey: rSolMint, IsSigner: false, IsWritable: false},
			{PubKey: admin, IsSigner: true, IsWritable: false},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func Stake(
	rSolProgramID,
	stakeManager,
	stakePool,
	from,
	mintManager,
	rSolMint,
	mintTo,
	mintAuthority,
	minterProgramId common.PublicKey,
	stakeAmount uint64,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
		StakeAmount uint64
	}{
		Instruction: InstructionStake,
		StakeAmount: stakeAmount,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: rSolProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: stakePool, IsSigner: false, IsWritable: true},
			{PubKey: from, IsSigner: true, IsWritable: true},
			{PubKey: mintManager, IsSigner: false, IsWritable: false},
			{PubKey: rSolMint, IsSigner: false, IsWritable: true},
			{PubKey: mintTo, IsSigner: false, IsWritable: true},
			{PubKey: mintAuthority, IsSigner: false, IsWritable: false},
			{PubKey: minterProgramId, IsSigner: false, IsWritable: false},
			{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.TokenProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}
