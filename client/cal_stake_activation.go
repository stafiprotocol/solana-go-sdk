package client

import (
	"context"
	"fmt"
)

func (s *Client) CalStakeActivation(ctx context.Context, address string) (*GetStakeActivationResponse, error) {
	stakeAccount, err := s.GetStakeAccountInfo(ctx, address)
	if err != nil {
		return nil, err
	}
	if stakeAccount.StakeAccount.Type == 0 {
		return nil, fmt.Errorf("stake account not init")
	}

	rentExemptReserve := stakeAccount.StakeAccount.Info.Meta.RentExemptReserve
	delegation := stakeAccount.StakeAccount.Info.Stake.Delegation
	if delegation == (Delegation{}) {
		return &GetStakeActivationResponse{
			State:    StakeActivationStateInactive,
			Active:   0,
			Inactive: stakeAccount.Lamports - uint64(rentExemptReserve),
		}, nil
	}

	stakeHistories, err := s.GetStakeHistory(ctx)
	if err != nil {
		return nil, err
	}
	epochInfo, err := s.GetEpochInfo(ctx, CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	stakeActivationStatus := delegation.StakeActivatingAndDeactivating(uint64(epochInfo.Epoch), stakeHistories.StakeHistories)

	stake_activation_state := StakeActivationStateInactive
	if stakeActivationStatus.Deactivating > 0 {
		stake_activation_state = StakeActivationStateDeactivating
	} else if stakeActivationStatus.Activating > 0 {
		stake_activation_state = StakeActivationStateActivating
	} else if stakeActivationStatus.Effective > 0 {
		stake_activation_state = StakeActivationStateActive
	} else {
		stake_activation_state = StakeActivationStateInactive
	}

	inactive_stake := uint64(0)
	if stakeAccount.Lamports > stakeActivationStatus.Effective+uint64(rentExemptReserve) {
		inactive_stake = stakeAccount.Lamports - (stakeActivationStatus.Effective + uint64(rentExemptReserve))
	}

	return &GetStakeActivationResponse{
		State:    stake_activation_state,
		Active:   stakeActivationStatus.Effective,
		Inactive: inactive_stake,
	}, nil
}
