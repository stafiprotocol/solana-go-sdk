package rsolprog

import (
	"github.com/stafiprotocol/solana-go-sdk/common"
)

type StakeManager struct {
	Admin                common.PublicKey
	Balancer             common.PublicKey
	RSolMint             common.PublicKey
	FeeRecipient         common.PublicKey
	PoolSeedBump         uint8
	RentExemptForPoolAcc uint64

	MinStakeAmount        uint64
	UnstakeFeeCommission  uint64 // decimals 9
	ProtocolFeeCommission uint64 // decimals 9
	RateChangeLimit       uint64 // decimals 9
	StakeAccountsLenLimit uint64
	SplitAccountsLenLimit uint64
	UnbondingDuration     uint64

	LatestEra        uint64
	Rate             uint64 // decimals 9
	EraBond          uint64
	EraUnbond        uint64
	Active           uint64
	TotalRSolSupply  uint64
	TotalProtocolFee uint64
	Validators       []common.PublicKey
	StakeAccounts    []common.PublicKey
	SplitAccounts    []common.PublicKey
	EraProcessData   EraProcessData
}

type EraProcessData struct {
	NeedBond             uint64
	NeedUnbond           uint64
	OldActive            uint64
	NewActive            uint64
	PendingStakeAccounts []common.PublicKey
}

type UnstakeAccount struct {
	StakeManager common.PublicKey
	Recipient    common.PublicKey
	Amount       uint64
	CreatedEpoch uint64
}
