package lsdprog

import (
	"github.com/stafiprotocol/solana-go-sdk/common"
)

type Stack struct {
	Admin                  common.PublicKey
	StackFeeCommission     uint64 // decimals 9
	StakeManagersLenLimit  uint64
	EntrustedStakeManagers []common.PublicKey
}

type StakeManager struct {
	Admin                common.PublicKey
	Balancer             common.PublicKey
	Stack                common.PublicKey
	LsdTokenMint         common.PublicKey
	PoolSeedBump         uint8
	RentExemptForPoolAcc uint64

	MinStakeAmount        uint64
	PlatformFeeCommission uint64 // decimals 9
	StackFeeCommission    uint64 // decimals 9
	RateChangeLimit       uint64 // decimals 9
	StakeAccountsLenLimit uint64
	SplitAccountsLenLimit uint64
	UnbondingDuration     uint64

	LatestEra        uint64
	Rate             uint64 // decimals 9
	EraBond          uint64
	EraUnbond        uint64
	Active           uint64
	TotalPlatformFee uint64
	Validators       []common.PublicKey
	StakeAccounts    []common.PublicKey
	SplitAccounts    []common.PublicKey
	EraRates         []EraRate
	EraProcessData   EraProcessData
}

type EraProcessData struct {
	NeedBond             uint64
	NeedUnbond           uint64
	OldActive            uint64
	NewActive            uint64
	PendingStakeAccounts []common.PublicKey
}

type EraRate struct {
	Era  uint64
	Rate uint64
}

type UnstakeAccount struct {
	StakeManager common.PublicKey
	Recipient    common.PublicKey
	Amount       uint64
	CreatedEpoch uint64
}

type StackFeeAccount struct {
	Bump   uint8
	Amount uint64
}
