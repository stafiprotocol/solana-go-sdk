package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"

	"github.com/stafiprotocol/solana-go-sdk/binary"
	"github.com/stafiprotocol/solana-go-sdk/common"
)

var StakeAccountInfoLengthDefault = uint64(200)
var GetStakeAccountInfoConfigDefault = GetAccountInfoConfig{
	Encoding: GetAccountInfoConfigEncodingBase64,
	DataSlice: GetAccountInfoConfigDataSlice{
		Offset: 0,
		Length: StakeAccountInfoLengthDefault,
	},
}

type StakeAccount struct {
	Type uint32 //0 uninitialized 1 initialized 2 delegated 3 rewardspool
	Info struct {
		Meta struct {
			RentExemptReserve int64
			Authorized        struct {
				Staker     common.PublicKey
				Withdrawer common.PublicKey
				Lockup     struct {
					UnixTimeStamp int64
					Epoch         uint64
					Custodian     common.PublicKey
				}
			}
		}
		Stake struct {
			Delegation      Delegation
			CreditsObserved uint64
		}
	}
}

type Delegation struct {
	Voter              common.PublicKey
	Stake              uint64
	ActivationEpoch    uint64 //epoch when delegate
	DeactivationEpoch  uint64 //epoch when deactive
	WarmupCooldownRate float64
}

func getHistory(h []StakeHistory, epoch uint64) *StakeHistoryEntry {
	for _, his := range h {
		if his.Epoch == epoch {
			return &his.Entry
		}
	}
	return nil
}

const NEW_WARMUP_COOLDOWN_RATE = float64(0.09)

// returned tuple is (effective, activating) stake
func (d *Delegation) StakeAndActivating(targetEpoch uint64, histories []StakeHistory) (uint64, uint64) {
	if d.ActivationEpoch == math.MaxUint64 {
		return d.Stake, 0
	}
	if d.ActivationEpoch == d.DeactivationEpoch {
		return 0, 0
	}
	if targetEpoch == d.ActivationEpoch {
		return 0, d.Stake
	}
	if targetEpoch < d.ActivationEpoch {
		return 0, 0
	}

	targetEntry := getHistory(histories, targetEpoch)
	if targetEntry != nil {

		prev_epoch := d.ActivationEpoch
		prev_cluster_stake := targetEntry

		current_epoch := uint64(0)
		current_effective_stake := uint64(0)
		for {
			current_epoch = prev_epoch + 1
			// if there is no activating stake at prev epoch, we should have been
			// fully effective at this moment
			if prev_cluster_stake.Activating == 0 {
				break
			}

			// how much of the growth in stake this account is
			//  entitled to take
			remaining_activating_stake := d.Stake - current_effective_stake
			weight := float64(remaining_activating_stake) / float64(prev_cluster_stake.Activating)
			warmup_cooldown_rate := NEW_WARMUP_COOLDOWN_RATE

			// // portion of newly effective cluster stake I'm entitled to at current epoch
			newly_effective_cluster_stake := float64(prev_cluster_stake.Effective) * warmup_cooldown_rate
			newly_effective_stake := uint64((weight * newly_effective_cluster_stake))
			if newly_effective_stake < 1 {
				newly_effective_stake = 1
			}

			current_effective_stake += newly_effective_stake
			if current_effective_stake >= d.Stake {
				current_effective_stake = d.Stake
				break
			}

			if current_epoch >= targetEpoch || current_epoch >= d.DeactivationEpoch {
				break
			}

			current_cluster_stake := getHistory(histories, current_epoch)
			if current_cluster_stake != nil {
				prev_epoch = current_epoch
				prev_cluster_stake = current_cluster_stake

			} else {
				break
			}
		}
		return current_effective_stake, d.Stake - current_effective_stake

	} else {
		return d.Stake, 0
	}
}

func (d *Delegation) StakeActivatingAndDeactivating(target_epoch uint64, histories []StakeHistory) StakeHistoryEntry {
	effective_stake, activating_stake := d.StakeAndActivating(target_epoch, histories)
	if target_epoch < d.DeactivationEpoch {
		// not deactivated
		if activating_stake == 0 {
			// StakeActivationStatus::with_effective(effective_stake)
			return StakeHistoryEntry{
				Effective: effective_stake,
			}
		} else {
			return StakeHistoryEntry{
				Effective:  effective_stake,
				Activating: activating_stake,
			}
		}
	}
	if target_epoch == d.DeactivationEpoch {
		return StakeHistoryEntry{
			Deactivating: effective_stake,
		}
	}

	targetEntry := getHistory(histories, d.DeactivationEpoch)
	if targetEntry != nil {
		prev_epoch := d.DeactivationEpoch
		prev_cluster_stake := targetEntry

		current_epoch := uint64(0)
		current_effective_stake := effective_stake
		for {
			current_epoch = prev_epoch + 1
			// if there is no deactivating stake at prev epoch, we should have been
			// fully undelegated at this moment
			if prev_cluster_stake.Deactivating == 0 {
				break
			}

			// I'm trying to get to zero, how much of the deactivation in stake
			//   this account is entitled to take
			weight := float64(current_effective_stake) / float64(prev_cluster_stake.Deactivating)
			warmup_cooldown_rate := NEW_WARMUP_COOLDOWN_RATE

			// portion of newly not-effective cluster stake I'm entitled to at current epoch
			newly_not_effective_cluster_stake := float64(prev_cluster_stake.Effective) * warmup_cooldown_rate
			newly_not_effective_stake := uint64(weight * newly_not_effective_cluster_stake)
			if newly_not_effective_stake < 1 {
				newly_not_effective_stake = 1
			}

			if current_effective_stake > newly_not_effective_stake {
				current_effective_stake = current_effective_stake - newly_not_effective_stake
			} else {
				current_effective_stake = 0
			}

			if current_effective_stake == 0 {
				break
			}

			if current_epoch >= target_epoch {
				break
			}
			current_cluster_stake := getHistory(histories, current_epoch)
			if current_cluster_stake != nil {
				prev_epoch = current_epoch
				prev_cluster_stake = current_cluster_stake
			} else {
				break
			}
		}

		return StakeHistoryEntry{
			Deactivating: current_effective_stake,
		}
	} else {
		return StakeHistoryEntry{}
	}
}

type StakeAccountRsp struct {
	Lamports     uint64
	Owner        string
	Excutable    bool
	RentEpoch    uint64
	StakeAccount StakeAccount
}

func (s *Client) GetStakeAccountInfo(ctx context.Context, account string) (*StakeAccountRsp, error) {
	accountInfo, err := s.GetAccountInfo(ctx, account, GetStakeAccountInfoConfigDefault)
	if err != nil {
		return nil, err
	}

	accountDataInterface, ok := accountInfo.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("account data err")
	}
	if len(accountDataInterface) != 2 {
		return nil, fmt.Errorf("account data length err")
	}
	accountDataBase64, ok := accountDataInterface[0].(string)
	if !ok {
		return nil, fmt.Errorf("get account base64 failed")
	}

	accountDataBts, err := base64.StdEncoding.DecodeString(accountDataBase64)
	if err != nil {
		return nil, err
	}
	if len(accountDataBts) <= 8 {
		return nil, fmt.Errorf("no account data bytes")
	}

	stakeAccountInfo := StakeAccount{}
	err = bin.NewDecoder(accountDataBts).Decode(&stakeAccountInfo)
	if err != nil {
		return nil, err
	}
	rsp := StakeAccountRsp{
		Lamports:     accountInfo.Lamports,
		Owner:        accountInfo.Owner,
		Excutable:    accountInfo.Excutable,
		RentEpoch:    accountInfo.RentEpoch,
		StakeAccount: stakeAccountInfo,
	}
	return &rsp, nil
}
