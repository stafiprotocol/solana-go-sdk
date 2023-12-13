package minterprog

import (
	"crypto/sha256"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

type Instruction [8]byte
type Event [8]byte

var (
	InstructionInitialize            Instruction
	InstructionMintToken             Instruction
	InstructionSetExtMintAuthorities Instruction

	MinterManagerAccountLengthDefault = uint64(2000)
)

func init() {
	initializeHash := sha256.Sum256([]byte("global:initialize"))
	copy(InstructionInitialize[:], initializeHash[:8])

	mintTokenHash := sha256.Sum256([]byte("global:mint_token"))
	copy(InstructionMintToken[:], mintTokenHash[:8])

	setExtMintAuthoritiesHash := sha256.Sum256([]byte("global:set_ext_mint_authorities"))
	copy(InstructionSetExtMintAuthorities[:], setExtMintAuthoritiesHash[:8])

}

func Initialize(
	minterProgramID,
	mintManager,
	rsolMint,
	admin common.PublicKey,
	extMintAthorities []common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction       Instruction
		ExtMintAthorities []common.PublicKey
	}{
		Instruction:       InstructionInitialize,
		ExtMintAthorities: extMintAthorities,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: minterProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: mintManager, IsSigner: false, IsWritable: true},
			{PubKey: rsolMint, IsSigner: false, IsWritable: false},
			{PubKey: admin, IsSigner: true, IsWritable: false},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func MintToken(
	minterProgramID,
	mintManager,
	rsolMint,
	mintTo,
	mintAuthority,
	extMintAuthority,
	tokenProgram common.PublicKey,
	mintAmount uint64) types.Instruction {

	data, err := common.SerializeData(struct {
		Instruction Instruction
		MintAmount  uint64
	}{
		Instruction: InstructionMintToken,
		MintAmount:  mintAmount,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: mintManager, IsSigner: false, IsWritable: false},
		{PubKey: rsolMint, IsSigner: false, IsWritable: true},
		{PubKey: mintTo, IsSigner: false, IsWritable: true},
		{PubKey: mintAuthority, IsSigner: false, IsWritable: false},
		{PubKey: extMintAuthority, IsSigner: true, IsWritable: false},
		{PubKey: tokenProgram, IsSigner: false, IsWritable: false},
	}

	return types.Instruction{
		ProgramID: minterProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}

func SetExtMintAuthorities(
	minterProgramID,
	mintManager,
	admin common.PublicKey,
	extMintAuthorities []common.PublicKey) types.Instruction {

	data, err := common.SerializeData(struct {
		Instruction        Instruction
		ExtMintAuthorities []common.PublicKey
	}{
		Instruction:        InstructionSetExtMintAuthorities,
		ExtMintAuthorities: extMintAuthorities,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: mintManager, IsSigner: false, IsWritable: true},
		{PubKey: admin, IsSigner: true, IsWritable: false},
	}

	return types.Instruction{
		ProgramID: minterProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}
