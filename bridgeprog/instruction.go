package bridgeprog

import (
	"crypto/sha256"

	"github.com/near/borsh-go"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/types"
)

type Instruction [8]byte
type Event [8]byte

var (
	InstructionCreateBridge        Instruction
	InstructionCreateMintProposal  Instruction
	InstructionApproveMintProposal Instruction
	InstructionSetOwners           Instruction
	InstructionChangeThreshold     Instruction
	InstructionSetResourceId       Instruction
	InstructionRemoveResourceId    Instruction
	InstructionTransferOut         Instruction
	InstructionSetFeeReceiver      Instruction
	InstructionSetFeeAmount        Instruction
	InstructionSetSupportChainIds  Instruction
	InstructionSetMintAuthority    Instruction

	EventTransferOut       Event
	ProgramLogPrefix       = "Program data: "
	EventTransferOutPrefix = ProgramLogPrefix + "7arrB4Lk4L"
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
	removeResourceIdHash := sha256.Sum256([]byte("global:remove_resource_id"))
	copy(InstructionRemoveResourceId[:], removeResourceIdHash[:8])
	transferOutHash := sha256.Sum256([]byte("global:transfer_out"))
	copy(InstructionTransferOut[:], transferOutHash[:8])
	eventTransferOutHash := sha256.Sum256([]byte("event:EventTransferOut"))
	copy(EventTransferOut[:], eventTransferOutHash[:8])
	setFeeReceiverHash := sha256.Sum256([]byte("global:set_fee_receiver"))
	copy(InstructionSetFeeReceiver[:], setFeeReceiverHash[:8])
	setFeeAmountHash := sha256.Sum256([]byte("global:set_fee_amount"))
	copy(InstructionSetFeeAmount[:], setFeeAmountHash[:8])
	setSupportChainIdsHash := sha256.Sum256([]byte("global:set_support_chain_ids"))
	copy(InstructionSetSupportChainIds[:], setSupportChainIdsHash[:8])
	setMintAuthorityHash := sha256.Sum256([]byte("global:set_mint_authority"))
	copy(InstructionSetMintAuthority[:], setMintAuthorityHash[:8])
}

func CreateBridge(
	bridgeProgramID,
	bridgeAccount common.PublicKey,
	owners []common.PublicKey,
	threshold uint64,
	nonce uint8,
	supportChainIds []uint8,
	resourceIdToMint map[[32]byte]common.PublicKey,
	admin,
	feeReceiver common.PublicKey,
	feeAmounts map[uint8]uint64) types.Instruction {

	data, err := borsh.Serialize(struct {
		Instruction           Instruction
		Owners                []common.PublicKey
		Threshold             uint64
		Nonce                 uint8
		SupportChainIds       []uint8
		ResourceIdToTokenProg map[[32]byte]common.PublicKey
		Admin                 common.PublicKey
		FeeReceiver           common.PublicKey
		FeeAmounts            map[uint8]uint64
	}{
		Instruction:           InstructionCreateBridge,
		Owners:                owners,
		Threshold:             threshold,
		Nonce:                 nonce,
		SupportChainIds:       supportChainIds,
		ResourceIdToTokenProg: resourceIdToMint,
		Admin:                 admin,
		FeeReceiver:           feeReceiver,
		FeeAmounts:            feeAmounts,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: bridgeProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: bridgeAccount, IsSigner: false, IsWritable: true},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func ChangeThreshold(
	bridgeProgramID,
	bridgeAccount,
	adminAccount common.PublicKey,
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
		{PubKey: adminAccount, IsSigner: true, IsWritable: false},
	}

	return types.Instruction{
		ProgramID: bridgeProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}

func SetResourceId(
	bridgeProgramID,
	bridgeAccount,
	adminAccount common.PublicKey,
	resourceId [32]byte,
	mint common.PublicKey,
) types.Instruction {

	data, err := common.SerializeData(struct {
		Instruction Instruction
		ResourceId  [32]byte
		Mint        common.PublicKey
	}{
		Instruction: InstructionSetResourceId,
		ResourceId:  resourceId,
		Mint:        mint,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: bridgeAccount, IsSigner: false, IsWritable: true},
		{PubKey: adminAccount, IsSigner: true, IsWritable: false},
	}

	return types.Instruction{
		ProgramID: bridgeProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}

func RemoveResourceId(
	bridgeProgramID,
	bridgeAccount,
	adminAccount common.PublicKey,
	resourceId [32]byte,
) types.Instruction {

	data, err := common.SerializeData(struct {
		Instruction Instruction
		ResourceId  [32]byte
	}{
		Instruction: InstructionRemoveResourceId,
		ResourceId:  resourceId,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: bridgeAccount, IsSigner: false, IsWritable: true},
		{PubKey: adminAccount, IsSigner: true, IsWritable: false},
	}

	return types.Instruction{
		ProgramID: bridgeProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}

func SetSupportChainIds(
	bridgeProgramID,
	bridgeAccount,
	adminAccount common.PublicKey,
	chainIds []uint8,
) types.Instruction {

	data, err := common.SerializeData(struct {
		Instruction Instruction
		ChainIds    []uint8
	}{
		Instruction: InstructionSetSupportChainIds,
		ChainIds:    chainIds,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: bridgeAccount, IsSigner: false, IsWritable: true},
		{PubKey: adminAccount, IsSigner: true, IsWritable: false},
	}

	return types.Instruction{
		ProgramID: bridgeProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}

func SetOwners(
	bridgeProgramID,
	bridgeAccount,
	adminAccount common.PublicKey,
	owners []common.PublicKey,
) types.Instruction {

	data, err := common.SerializeData(struct {
		Instruction Instruction
		Owners      []common.PublicKey
	}{
		Instruction: InstructionSetOwners,
		Owners:      owners,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: bridgeAccount, IsSigner: false, IsWritable: true},
		{PubKey: adminAccount, IsSigner: true, IsWritable: false},
	}

	return types.Instruction{
		ProgramID: bridgeProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}

func SetFeeReceiver(
	bridgeProgramID,
	bridgeAccount,
	adminAccount common.PublicKey,
	feeReceiver common.PublicKey,
) types.Instruction {
	data, err := common.SerializeData(struct {
		Instruction Instruction
		FeeReceiver common.PublicKey
	}{
		Instruction: InstructionSetFeeReceiver,
		FeeReceiver: feeReceiver,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: bridgeAccount, IsSigner: false, IsWritable: true},
		{PubKey: adminAccount, IsSigner: true, IsWritable: false},
	}

	return types.Instruction{
		ProgramID: bridgeProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}

func SetFeeAmount(
	bridgeProgramID,
	bridgeAccount,
	adminAccount common.PublicKey,
	destChainId uint8,
	amount uint64,
) types.Instruction {
	data, err := common.SerializeData(struct {
		Instruction Instruction
		DestChainId uint8
		Amount      uint64
	}{
		Instruction: InstructionSetFeeAmount,
		DestChainId: destChainId,
		Amount:      amount,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: bridgeAccount, IsSigner: false, IsWritable: true},
		{PubKey: adminAccount, IsSigner: true, IsWritable: false},
	}

	return types.Instruction{
		ProgramID: bridgeProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}

func SetMintAuthority(
	bridgeProgramID,
	bridgeAccount,
	adminAccount,
	bridgeSigner,
	mint,
	newMintAuthority common.PublicKey,
) types.Instruction {
	data, err := common.SerializeData(struct {
		Instruction   Instruction
		MintAuthority common.PublicKey
	}{
		Instruction:   InstructionSetMintAuthority,
		MintAuthority: newMintAuthority,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: bridgeAccount, IsSigner: false, IsWritable: false},
		{PubKey: adminAccount, IsSigner: true, IsWritable: false},
		{PubKey: bridgeSigner, IsSigner: false, IsWritable: false},
		{PubKey: mint, IsSigner: false, IsWritable: true},
		{PubKey: common.TokenProgramID, IsSigner: false, IsWritable: true},
	}

	return types.Instruction{
		ProgramID: bridgeProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}

type ProposalUsedAccount struct {
	Pubkey     common.PublicKey
	IsSigner   bool
	IsWritable bool
}

func CreateMintProposal(
	bridgeProgramID,
	bridgeAccount,
	proposalAccount,
	toAccount,
	proposerAccount common.PublicKey,
	resourceId [32]byte,
	amount uint64,
	tokenProgram common.PublicKey,
) types.Instruction {

	data, err := common.SerializeData(struct {
		Instruction  Instruction
		ResourceId   [32]byte
		Amount       uint64
		TokenProgram common.PublicKey
	}{
		Instruction:  InstructionCreateMintProposal,
		ResourceId:   resourceId,
		Amount:       amount,
		TokenProgram: tokenProgram,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: bridgeProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: bridgeAccount, IsSigner: false, IsWritable: false},
			{PubKey: proposalAccount, IsSigner: false, IsWritable: true},
			{PubKey: toAccount, IsSigner: false, IsWritable: false},
			{PubKey: proposerAccount, IsSigner: true, IsWritable: false},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}

func ApproveMintProposal(
	bridgeProgramID,
	bridgeAccount,
	multiSiner,
	proposalAccount,
	approverAccount,
	mintAccount,
	toAccount,
	mintManager,
	mintAuthority,
	minterProgramId common.PublicKey,
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
		{PubKey: toAccount, IsSigner: false, IsWritable: true},
		{PubKey: mintManager, IsSigner: false, IsWritable: false},
		{PubKey: mintAuthority, IsSigner: false, IsWritable: false},
		{PubKey: minterProgramId, IsSigner: false, IsWritable: false},
		{PubKey: common.TokenProgramID, IsSigner: false, IsWritable: false},
	}

	return types.Instruction{
		ProgramID: bridgeProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}

func TransferOut(
	bridgeProgramID,
	bridgeAccount,
	authorityAccount,
	mintAccount,
	fromAccount,
	feeReciever,
	tokenProgram,
	systemProgram common.PublicKey,
	amount uint64,
	receiver []byte,
	destChainId uint8,
) types.Instruction {

	data, err := common.SerializeData(struct {
		Instruction Instruction
		Amount      uint64
		Receiver    []byte
		DestChainId uint8
	}{
		Instruction: InstructionTransferOut,
		Amount:      amount,
		Receiver:    receiver,
		DestChainId: destChainId,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: bridgeAccount, IsSigner: false, IsWritable: true},
		{PubKey: authorityAccount, IsSigner: true, IsWritable: false},
		{PubKey: mintAccount, IsSigner: false, IsWritable: true},
		{PubKey: fromAccount, IsSigner: false, IsWritable: true},
		{PubKey: feeReciever, IsSigner: false, IsWritable: true},
		{PubKey: tokenProgram, IsSigner: false, IsWritable: false},
		{PubKey: systemProgram, IsSigner: false, IsWritable: false},
	}

	return types.Instruction{
		ProgramID: bridgeProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}
