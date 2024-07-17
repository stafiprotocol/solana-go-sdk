package lsdprog

import (
	"crypto/sha256"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

type Instruction [8]byte
type Event [8]byte

var StakeManagerAccountLengthDefault = uint64(100000)
var StackAccountLengthDefault = uint64(1000)
var StackFeeAccountLengthDefault = uint64(17)

var (
	InstructionInitializeStack        Instruction
	InstructionInitializeStakeManager Instruction

	InstructionAddEntrustedStakeManager Instruction

	InstructionAddValidator             Instruction
	InstructionRemoveValidator          Instruction
	InstructionRedelegate               Instruction
	InstructionSetRateChangeLimit       Instruction
	InstructionSetPlatformFeeCommission Instruction
	InstructionSetUnbondingDuration     Instruction
	InstructionReallocStakeManager      Instruction

	InstructionStake    Instruction
	InstructionUnstake  Instruction
	InstructionWithdraw Instruction

	InstructionEraNew          Instruction
	InstructionEraBond         Instruction
	InstructionEraUnbond       Instruction
	InstructionEraUpdateActive Instruction
	InstructionEraUpdateRate   Instruction
	InstructionEraMerge        Instruction
	InstructionEraWithdraw     Instruction
)

func init() {
	initializeStackHash := sha256.Sum256([]byte("global:initialize_stack"))
	copy(InstructionInitializeStack[:], initializeStackHash[:8])
	initializeStakeManagerHash := sha256.Sum256([]byte("global:initialize_stake_manager"))
	copy(InstructionInitializeStakeManager[:], initializeStakeManagerHash[:8])

	addEntrustedStakeManagerHash := sha256.Sum256([]byte("global:add_entrusted_stake_manager"))
	copy(InstructionAddEntrustedStakeManager[:], addEntrustedStakeManagerHash[:8])

	addValidatorHash := sha256.Sum256([]byte("global:add_validator"))
	copy(InstructionAddValidator[:], addValidatorHash[:8])
	removeValidatorHash := sha256.Sum256([]byte("global:remove_validator"))
	copy(InstructionRemoveValidator[:], removeValidatorHash[:8])
	redelegateHash := sha256.Sum256([]byte("global:redelegate"))
	copy(InstructionRedelegate[:], redelegateHash[:8])
	setRateChangeLimitHash := sha256.Sum256([]byte("global:set_rate_change_limit"))
	copy(InstructionSetRateChangeLimit[:], setRateChangeLimitHash[:8])
	setPlatformFeeCommissionHash := sha256.Sum256([]byte("global:set_platform_fee_commission"))
	copy(InstructionSetPlatformFeeCommission[:], setPlatformFeeCommissionHash[:8])
	setUnbondingDurationHash := sha256.Sum256([]byte("global:set_unbonding_duration"))
	copy(InstructionSetUnbondingDuration[:], setUnbondingDurationHash[:8])
	reallocStakeManagerHash := sha256.Sum256([]byte("global:realloc_stake_manager"))
	copy(InstructionReallocStakeManager[:], reallocStakeManagerHash[:8])

	stakeHash := sha256.Sum256([]byte("global:stake"))
	copy(InstructionStake[:], stakeHash[:8])
	unstakeHash := sha256.Sum256([]byte("global:unstake"))
	copy(InstructionUnstake[:], unstakeHash[:8])
	withdrawHash := sha256.Sum256([]byte("global:withdraw"))
	copy(InstructionWithdraw[:], withdrawHash[:8])

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
	eraWithdrawHash := sha256.Sum256([]byte("global:era_withdraw"))
	copy(InstructionEraWithdraw[:], eraWithdrawHash[:8])

}

func InitializeStack(
	lsdProgramID,
	stack,
	rentPayer,
	admin common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
	}{
		Instruction: InstructionInitializeStack,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stack, IsSigner: true, IsWritable: true},
			{PubKey: rentPayer, IsSigner: true, IsWritable: true},
			{PubKey: admin, IsSigner: true, IsWritable: false},
			{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func InitializeStakeManager(
	lsdProgramID,
	stakeManager,
	stack,
	stakePool,
	stackFeeAccount,
	lsdTokenMint,
	validator,
	rentPayer,
	admin common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
	}{
		Instruction: InstructionInitializeStakeManager,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: true, IsWritable: true},
			{PubKey: stack, IsSigner: false, IsWritable: false},
			{PubKey: stakePool, IsSigner: false, IsWritable: true},
			{PubKey: stackFeeAccount, IsSigner: false, IsWritable: true},
			{PubKey: lsdTokenMint, IsSigner: false, IsWritable: false},
			{PubKey: validator, IsSigner: false, IsWritable: false},
			{PubKey: rentPayer, IsSigner: true, IsWritable: true},
			{PubKey: admin, IsSigner: true, IsWritable: false},
			{PubKey: common.SPLAssociatedTokenAccountProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func Redelegate(
	lsdProgramID,
	stakeManager,
	admin,
	to_validator,
	stakePool,
	fromStakeAccount,
	splitStakeAccount,
	toStakeAccount,
	rentPayer common.PublicKey,
	amount uint64,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
		Amount      uint64
	}{
		Instruction: InstructionRedelegate,
		Amount:      amount,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: admin, IsSigner: true, IsWritable: false},
			{PubKey: to_validator, IsSigner: false, IsWritable: true},
			{PubKey: stakePool, IsSigner: false, IsWritable: false},
			{PubKey: fromStakeAccount, IsSigner: false, IsWritable: true},
			{PubKey: splitStakeAccount, IsSigner: true, IsWritable: true},
			{PubKey: toStakeAccount, IsSigner: true, IsWritable: true},
			{PubKey: rentPayer, IsSigner: false, IsWritable: true},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.StakeConfigPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarStakeHistoryPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.StakeProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func AddEntrustedStakeManager(
	lsdProgramID,
	stack,
	admin,
	entrustedStakeManager common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction  Instruction
		StakeManager common.PublicKey
	}{
		Instruction:  InstructionAddEntrustedStakeManager,
		StakeManager: entrustedStakeManager,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stack, IsSigner: false, IsWritable: true},
			{PubKey: admin, IsSigner: true, IsWritable: false},
		},
		Data: data,
	}
}

func AddValidator(
	lsdProgramID,
	stakeManager,
	admin,
	newValidator common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction  Instruction
		NewValidator common.PublicKey
	}{
		Instruction:  InstructionAddValidator,
		NewValidator: newValidator,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: admin, IsSigner: true, IsWritable: false},
		},
		Data: data,
	}
}

func RemoveValidator(
	lsdProgramID,
	stakeManager,
	admin,
	removeValidator common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction     Instruction
		RemoveValidator common.PublicKey
	}{
		Instruction:     InstructionRemoveValidator,
		RemoveValidator: removeValidator,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: admin, IsSigner: true, IsWritable: false},
		},
		Data: data,
	}
}

func SetRateChangeLimit(
	lsdProgramID,
	stakeManager,
	admin common.PublicKey,
	rateChangeLimit uint64,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction     Instruction
		RateChangeLimit uint64
	}{
		Instruction:     InstructionSetRateChangeLimit,
		RateChangeLimit: rateChangeLimit,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: admin, IsSigner: true, IsWritable: false},
		},
		Data: data,
	}
}

func SetPlatformFeeCommission(
	lsdProgramID,
	stakeManager,
	admin common.PublicKey,
	platformFeeCommission uint64,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction           Instruction
		platformFeeCommission uint64
	}{
		Instruction:           InstructionSetPlatformFeeCommission,
		platformFeeCommission: platformFeeCommission,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: admin, IsSigner: true, IsWritable: false},
		},
		Data: data,
	}
}

func SetUnbondingDuration(
	lsdProgramID,
	stakeManager,
	admin common.PublicKey,
	unbondingDuration uint64,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction       Instruction
		UnbondingDuration uint64
	}{
		Instruction:       InstructionSetUnbondingDuration,
		UnbondingDuration: unbondingDuration,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: admin, IsSigner: true, IsWritable: false},
		},
		Data: data,
	}
}

func ReallocStakeManager(
	lsdProgramID,
	stakeManager,
	admin,
	rentPayer common.PublicKey,
	newSize uint32,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
		NewSize     uint32
	}{
		Instruction: InstructionReallocStakeManager,
		NewSize:     newSize,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: admin, IsSigner: true, IsWritable: false},
			{PubKey: rentPayer, IsSigner: true, IsWritable: true},
			{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func Stake(
	lsdProgramID,
	stakeManager,
	stakePool,
	from,
	lsdTokenMint,
	mintTo common.PublicKey,
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
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: stakePool, IsSigner: false, IsWritable: true},
			{PubKey: from, IsSigner: true, IsWritable: true},
			{PubKey: lsdTokenMint, IsSigner: false, IsWritable: true},
			{PubKey: mintTo, IsSigner: false, IsWritable: true},
			{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.TokenProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func Unstake(
	lsdProgramID,
	stakeManager,
	lsdTokenMint,
	burnLsdTokenFrom,
	burnLsdTokenAuthority,
	unstakeAccount,
	rentPayer common.PublicKey,
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
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: lsdTokenMint, IsSigner: false, IsWritable: true},
			{PubKey: burnLsdTokenFrom, IsSigner: false, IsWritable: true},
			{PubKey: burnLsdTokenAuthority, IsSigner: true, IsWritable: false},
			{PubKey: unstakeAccount, IsSigner: true, IsWritable: true},
			{PubKey: rentPayer, IsSigner: true, IsWritable: true},
			{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.TokenProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func Withdraw(
	lsdProgramID,
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
		ProgramID: lsdProgramID,
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
	lsdProgramID,
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
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func EraBond(
	lsdProgramID,
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
		ProgramID: lsdProgramID,
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
	lsdProgramID,
	stakeManager,
	stakePool,
	fromStakeAccount,
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
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: stakePool, IsSigner: false, IsWritable: false},
			{PubKey: fromStakeAccount, IsSigner: false, IsWritable: true},
			{PubKey: splitStakeAccount, IsSigner: true, IsWritable: true},
			{PubKey: validator, IsSigner: false, IsWritable: true},
			{PubKey: rentPayer, IsSigner: true, IsWritable: true},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarStakeHistoryPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.StakeProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func EraUpdateActive(
	lsdProgramID,
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
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: stakeAccount, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func EraUpdateRate(
	lsdProgramID,
	stakeManager,
	stack,
	stakePool,
	lsdTokenMint,
	platformFeeRecipient,
	stackFeeRecipient,
	stackFeeAccount common.PublicKey,
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
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: stack, IsSigner: false, IsWritable: false},
			{PubKey: stakePool, IsSigner: false, IsWritable: false},
			{PubKey: lsdTokenMint, IsSigner: false, IsWritable: true},
			{PubKey: platformFeeRecipient, IsSigner: false, IsWritable: true},
			{PubKey: stackFeeRecipient, IsSigner: false, IsWritable: true},
			{PubKey: stackFeeAccount, IsSigner: false, IsWritable: true},
			{PubKey: common.SPLAssociatedTokenAccountProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.TokenProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func EraMerge(
	lsdProgramID,
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
		ProgramID: lsdProgramID,
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

func EraWithdraw(
	lsdProgramID,
	stakeManager,
	stakePool,
	stakeAccount common.PublicKey,
) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
	}{
		Instruction: InstructionEraWithdraw,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: lsdProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: stakeManager, IsSigner: false, IsWritable: true},
			{PubKey: stakePool, IsSigner: false, IsWritable: true},
			{PubKey: stakeAccount, IsSigner: false, IsWritable: true},
			{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarStakeHistoryPubkey, IsSigner: false, IsWritable: false},
			{PubKey: common.StakeProgramID, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}
