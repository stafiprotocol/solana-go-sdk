package rsolprog

import (
	"crypto/sha256"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

type Instruction [8]byte
type Event [8]byte

var StakeManagerAccountLengthDefault = uint64(2000000)
var UnstakeAccountLengthDefault = uint64(100)
var (
	InstructionInitialize          Instruction
	InstructionMigrateStakeAccount Instruction

	InstructionStake    Instruction
	InstructionUnstake  Instruction
	InstructionWithdraw Instruction

	InstructionSetActive     Instruction
	InstructionSetRsolSupply Instruction

	InstructionEraNew          Instruction
	InstructionEraBond         Instruction
	InstructionEraUnbond       Instruction
	InstructionEraUpdateActive Instruction
	InstructionEraUpdateRate   Instruction
	InstructionEraMerge        Instruction
)

func init() {
	initializeHash := sha256.Sum256([]byte("global:initialize"))
	copy(InstructionInitialize[:], initializeHash[:8])

	migrateStakeAccountHash := sha256.Sum256([]byte("global:migrate_stake_account"))
	copy(InstructionMigrateStakeAccount[:], migrateStakeAccountHash[:8])

	stakeHash := sha256.Sum256([]byte("global:stake"))
	copy(InstructionStake[:], stakeHash[:8])
	unstakeHash := sha256.Sum256([]byte("global:unstake"))
	copy(InstructionUnstake[:], unstakeHash[:8])
	withdrawHash := sha256.Sum256([]byte("global:withdraw"))
	copy(InstructionWithdraw[:], withdrawHash[:8])

	setActiveHash := sha256.Sum256([]byte("global:set_active"))
	copy(InstructionSetActive[:], setActiveHash[:8])
	setRsolSupplyHash := sha256.Sum256([]byte("global:set_rsol_supply"))
	copy(InstructionSetRsolSupply[:], setRsolSupplyHash[:8])

	eraNewHash := sha256.Sum256([]byte("global:era_new"))
	copy(InstructionEraNew[:], eraNewHash[:8])
	eraBondHash := sha256.Sum256([]byte("global:era_bond"))
	copy(InstructionEraBond[:], eraBondHash[:8])
	eraUnbondHash := sha256.Sum256([]byte("global:era_unbond"))
	copy(InstructionEraUnbond[:], eraUnbondHash[:8])
	eraUpdateActiveHash := sha256.Sum256([]byte("global:era_update_active"))
	copy(InstructionEraUpdateActive[:], eraUpdateActiveHash[:8])
	eraUpdateRateHash := sha256.Sum256([]byte("global:era_update_rate"))
	copy(InstructionEraUpdateRate[:], eraUpdateRateHash[:8])
	eraMergeHash := sha256.Sum256([]byte("global:era_merge"))
	copy(InstructionEraMerge[:], eraMergeHash[:8])

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

func MigrateStakeAccount(
	rSolProgramID,
	stakeManager,
	stakePool,
	stakeAccount,
	stakeAuthority common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
	}{
		Instruction: InstructionMigrateStakeAccount,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: rSolProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: stakePool, IsSigner: false, IsWritable: false},
			{PubKey: stakeAccount, IsSigner: false, IsWritable: true},
			{PubKey: stakeAuthority, IsSigner: true, IsWritable: false},
			{PubKey: common.StakeProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func SetActive(
	rSolProgramID,
	stakeManager,
	admin common.PublicKey,
	active uint64,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
		Active      uint64
	}{
		Instruction: InstructionSetActive,
		Active:      active,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: rSolProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: admin, IsSigner: true, IsWritable: false},
		},
		Data: data,
	}
}

func SetRsolSupply(
	rSolProgramID,
	stakeManager,
	admin common.PublicKey,
	total uint64,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
		Total       uint64
	}{
		Instruction: InstructionSetRsolSupply,
		Total:       total,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: rSolProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: admin, IsSigner: true, IsWritable: false},
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

func Unstake(
	rSolProgramID,
	stakeManager,
	rSolMint,
	burnRsolFrom,
	burnRsolAuthority,
	unstakeAccount,
	feeRecipient common.PublicKey,
	unstakeAmount uint64,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction   Instruction
		UnstakeAmount uint64
	}{
		Instruction:   InstructionUnstake,
		UnstakeAmount: unstakeAmount,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: rSolProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: rSolMint, IsSigner: false, IsWritable: true},
			{PubKey: burnRsolFrom, IsSigner: false, IsWritable: true},
			{PubKey: burnRsolAuthority, IsSigner: true, IsWritable: false},
			{PubKey: unstakeAccount, IsSigner: true, IsWritable: true},
			{PubKey: feeRecipient, IsSigner: false, IsWritable: true},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.TokenProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func Withdraw(
	rSolProgramID,
	stakeManager,
	stakePool,
	unstakeAccount,
	recipient common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
	}{
		Instruction: InstructionWithdraw,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: rSolProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: stakePool, IsSigner: false, IsWritable: true},
			{PubKey: unstakeAccount, IsSigner: false, IsWritable: true},
			{PubKey: recipient, IsSigner: false, IsWritable: true},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func EraNew(
	rSolProgramID,
	stakeManager common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
	}{
		Instruction: InstructionEraNew,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: rSolProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func EraBond(
	rSolProgramID,
	stakeManager,
	validator,
	stakePool,
	stakeAccount,
	rentPayer common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
	}{
		Instruction: InstructionEraBond,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: rSolProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: validator, IsSigner: false, IsWritable: true},
			{PubKey: stakePool, IsSigner: false, IsWritable: true},
			{PubKey: stakeAccount, IsSigner: true, IsWritable: true},
			{PubKey: rentPayer, IsSigner: true, IsWritable: true},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.StakeConfigPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarStakeHistoryPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.StakeProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func EraUnbond(
	rSolProgramID,
	stakeManager,
	stakePool,
	stakeAccount,
	splitStakeAccount,
	validator,
	rentPayer common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
	}{
		Instruction: InstructionEraUnbond,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: rSolProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: stakePool, IsSigner: false, IsWritable: false},
			{PubKey: stakeAccount, IsSigner: false, IsWritable: true},
			{PubKey: splitStakeAccount, IsSigner: true, IsWritable: true},
			{PubKey: validator, IsSigner: false, IsWritable: true},
			{PubKey: rentPayer, IsSigner: true, IsWritable: true},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.StakeProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func EraUpdateActive(
	rSolProgramID,
	stakeManager,
	stakeAccount common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
	}{
		Instruction: InstructionEraUpdateActive,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: rSolProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: stakeAccount, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func EraUpdateRate(
	rSolProgramID,
	stakeManager,
	stakePool,
	mintManager,
	rsolMint,
	feeRecipient,
	mintAuthority,
	minterProgramID common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
	}{
		Instruction: InstructionEraUpdateRate,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: rSolProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: stakePool, IsSigner: false, IsWritable: false},
			{PubKey: mintManager, IsSigner: false, IsWritable: false},
			{PubKey: rsolMint, IsSigner: false, IsWritable: true},
			{PubKey: feeRecipient, IsSigner: false, IsWritable: true},
			{PubKey: mintAuthority, IsSigner: false, IsWritable: false},
			{PubKey: minterProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.TokenProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func EraMerge(
	rSolProgramID,
	stakeManager,
	srcStakeAccount,
	dstStakeAccount,
	stakePool common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
	}{
		Instruction: InstructionEraMerge,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: rSolProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: srcStakeAccount, IsSigner: false, IsWritable: true},
			{PubKey: dstStakeAccount, IsSigner: false, IsWritable: true},
			{PubKey: stakePool, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarStakeHistoryPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.StakeProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}
