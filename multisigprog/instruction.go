package multisigprog

import (
	"crypto/sha256"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

type Instruction [8]byte

var (
	InstructionCreateMultisig    Instruction
	InstructionCreateTransaction Instruction
	InstructionApprove           Instruction
	InstructionSetOwners         Instruction
	InstructionChangeThreshold   Instruction
)

func init() {
	createMultisigHash := sha256.Sum256([]byte("global:create_multisig"))
	copy(InstructionCreateMultisig[:], createMultisigHash[:8])
	createTransactionHash := sha256.Sum256([]byte("global:create_transaction"))
	copy(InstructionCreateTransaction[:], createTransactionHash[:8])
	approveHash := sha256.Sum256([]byte("global:approve"))
	copy(InstructionApprove[:], approveHash[:8])
	setOwnersHash := sha256.Sum256([]byte("global:set_owners"))
	copy(InstructionSetOwners[:], setOwnersHash[:8])
	changeThresholdHash := sha256.Sum256([]byte("global:change_threshold"))
	copy(InstructionChangeThreshold[:], changeThresholdHash[:8])
}

func CreateMultisig(
	programID,
	multisigAccount common.PublicKey,
	owners []common.PublicKey,
	threshold uint64,
	nonce uint8) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction Instruction
		Owners      []common.PublicKey
		Threshold   uint64
		Nonce       uint8
	}{
		Instruction: InstructionCreateMultisig,
		Owners:      owners,
		Threshold:   threshold,
		Nonce:       nonce,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: programID,
		Accounts: []types.AccountMeta{
			{PubKey: multisigAccount, IsSigner: false, IsWritable: true},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

type TransactionUsedAccount struct {
	Pubkey     common.PublicKey
	IsSigner   bool
	IsWritable bool
}

func CreateTransaction(
	programID common.PublicKey,
	txUsedProgramID []common.PublicKey,
	txUsedAccounts [][]types.AccountMeta,
	txInstructionData [][]byte,
	multisigAccount common.PublicKey,
	txAccount common.PublicKey,
	proposalAccount common.PublicKey) types.Instruction {

	data, err := common.SerializeData(struct {
		Instruction       Instruction
		TxUsedProgramID   []common.PublicKey
		TxUsedAccounts    [][]types.AccountMeta
		TxInstructionData [][]byte
	}{
		Instruction:       InstructionCreateTransaction,
		TxUsedProgramID:   txUsedProgramID,
		TxUsedAccounts:    txUsedAccounts,
		TxInstructionData: txInstructionData,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: programID,
		Accounts: []types.AccountMeta{
			{PubKey: multisigAccount, IsSigner: false, IsWritable: false},
			{PubKey: txAccount, IsSigner: false, IsWritable: true},
			{PubKey: proposalAccount, IsSigner: true, IsWritable: false},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func Approve(
	programID,
	multisigAccount,
	multiSiner,
	txAccount,
	approverAccount common.PublicKey,
	remainingAccounts []types.AccountMeta) types.Instruction {

	data, err := common.SerializeData(struct {
		Instruction Instruction
	}{
		Instruction: InstructionApprove,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: multisigAccount, IsSigner: false, IsWritable: false},
		{PubKey: multiSiner, IsSigner: false, IsWritable: false},
		{PubKey: txAccount, IsSigner: false, IsWritable: true},
		{PubKey: approverAccount, IsSigner: true, IsWritable: false},
	}

	return types.Instruction{
		ProgramID: programID,
		Accounts:  append(accounts, remainingAccounts...),
		Data:      data,
	}
}

func ChangeThreshold(
	programID,
	multisigAccount,
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
		{PubKey: multisigAccount, IsSigner: false, IsWritable: true},
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
