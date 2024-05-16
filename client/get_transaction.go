package client

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"golang.org/x/sync/errgroup"
)

var ErrTxNotFound = errors.New("TxNotFound")
var DefaultMaxSupportedTransactionVersion = uint8(0)

type GetTransactionWithLimitConfig struct {
	// TODO custom encoding
	// Encoding   string     `json:"encoding"`          // either "json", "jsonParsed", "base58" (slow), "base64", default: json
	Commitment                     Commitment `json:"commitment,omitempty"`                     // "processed" is not supported. If parameter not provided, the default is "finalized".
	MaxSupportedTransactionVersion *uint8     `json:"maxSupportedTransactionVersion,omitempty"` // default: nil legacy only
}

type GetTransaction struct {
	Slot        uint64          `json:"slot"`
	Meta        TransactionMeta `json:"meta"`
	Transaction Transaction     `json:"transaction"`
}

type GetTransactionResponse struct {
	Slot        uint64          `json:"slot"`
	Meta        TransactionMeta `json:"meta"`
	Transaction Transaction     `json:"transaction"`
}

// NEW: This method is only available in solana-core v1.7 or newer. Please use getConfirmedTransaction for solana-core v1.6
// GetConfirmedTransaction returns transaction details for a confirmed transaction
func (s *Client) GetTransaction(ctx context.Context, txhash string, cfg GetTransactionWithLimitConfig) (GetTransactionResponse, error) {
	res := struct {
		GeneralResponse
		Result GetTransactionResponse `json:"result"`
	}{}
	err := s.request(ctx, "getTransaction", []interface{}{txhash, cfg}, &res)
	if err != nil {
		return GetTransactionResponse{}, err
	}
	return res.Result, nil
}

// NEW: This method is only available in solana-core v1.7 or newer. Please use getConfirmedTransaction for solana-core v1.6
// GetConfirmedTransaction returns transaction details for a confirmed transaction
func (s *Client) GetTransactionV2(ctx context.Context, txhash string) (GetTransactionResponse, error) {
	res := struct {
		GeneralResponse
		Result GetTransactionResponse `json:"result"`
	}{}
	err := s.request(ctx, "getTransaction", []interface{}{txhash, GetTransactionWithLimitConfig{Commitment: CommitmentFinalized, MaxSupportedTransactionVersion: &DefaultMaxSupportedTransactionVersion}}, &res)
	if err != nil {
		return GetTransactionResponse{}, err
	}

	if res.Result.Slot == 0 {
		return GetTransactionResponse{}, ErrTxNotFound
	}

	return res.Result, nil
}

var fetchLimit = int(500)
var gNumberLimit = 5

type SolTx struct {
	Res       *GetTransactionResponse
	Signature string
	BlockTime uint64
}

func (client *Client) GetAddrRelateTxAfterSlot(addresses []string, dealtSlot uint64) ([]*SolTx, error) {
	sigsMap := make(map[string]bool)
	sigs := make([]GetSignaturesForAddress, 0)

	for _, addr := range addresses {
		sigRes, err := client.GetSignaturesForAddress(context.Background(), addr, GetSignaturesForAddressConfig{
			Limit: fetchLimit,
		})
		if err != nil {
			return nil, nil
		}

		for _, sig := range sigRes {
			if sig.Slot <= dealtSlot {
				continue
			}
			if !sigsMap[sig.Signature] {
				sigsMap[sig.Signature] = true
				sigs = append(sigs, sig)
			}
		}
		for {
			// break if no more
			if len(sigRes) < fetchLimit {
				break
			}
			// break if alrady less than dealtSlot
			minSlotSig := sigRes[len(sigRes)-1]
			if minSlotSig.Slot <= dealtSlot {
				break
			}

			sigResSub, err := client.GetSignaturesForAddress(context.Background(), addr, GetSignaturesForAddressConfig{
				Limit:  fetchLimit,
				Before: minSlotSig.Signature,
			})
			if err != nil {
				return nil, nil
			}
			sigRes = sigResSub

			for _, sig := range sigRes {
				if sig.Slot <= dealtSlot {
					continue
				}
				if !sigsMap[sig.Signature] {
					sigsMap[sig.Signature] = true
					sigs = append(sigs, sig)
				}
			}

			time.Sleep(time.Millisecond * 500)

		}
	}

	// return if no sigs
	if len(sigs) == 0 {
		return nil, nil
	}
	txs := make([]*SolTx, 0)

	dealLimitPerGoRoutine := int(math.Ceil(float64(len(sigs)) / float64(gNumberLimit)))
	retry := 0
	for {
		txChan := make(chan *SolTx, len(sigs))

		if retry > 200 {
			return nil, fmt.Errorf("GetConfirmedTransactionResponse reach retry limit")
		}

		g := new(errgroup.Group)
		g.SetLimit(int(gNumberLimit))

		for i := 0; i < len(sigs); i += dealLimitPerGoRoutine {
			start := i
			end := i + dealLimitPerGoRoutine
			if end > len(sigs) {
				end = len(sigs)
			}

			g.Go(func() error {
				for j := start; j < end; j++ {
					sig := sigs[j]

					retry := 0
					for {
						if retry > 60 {
							return fmt.Errorf("GetTransactionV2 reach retry limit")
						}
						tx, err := client.GetTransactionV2(context.Background(), sig.Signature)
						if err != nil {
							time.Sleep(time.Second)
							retry++
							continue
						}

						txChan <- &SolTx{
							Res:       &tx,
							Signature: sig.Signature,
							BlockTime: uint64(*sig.BlockTime),
						}
						break
					}

				}
				return nil
			})
		}

		err := g.Wait()
		if err != nil {
			retry++
			continue
		}

		close(txChan)

		if len(txChan) > 0 {
			for tx := range txChan {
				txs = append(txs, tx)
			}
		}
		break
	}

	sort.SliceStable(txs, func(i, j int) bool {
		return txs[i].Res.Slot < txs[j].Res.Slot
	})

	return txs, nil
}
