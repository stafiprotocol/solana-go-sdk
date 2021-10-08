package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	DevnetRPCEndpoint  = "https://api.devnet.solana.com"
	TestnetRPCEndpoint = "https://testnet.solana.com"
	MainnetRPCEndpoint = "https://api.mainnet-beta.solana.com"
	retryLimit         = 60 * 4
	waitTime           = time.Second * 5
)

type Commitment string

const (
	CommitmentFinalized Commitment = "finalized"
	CommitmentConfirmed Commitment = "confirmed"
	CommitmentProcessed Commitment = "processed"
	CommitmentRecent    Commitment = "recent"
)

type Client struct {
	endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{endpoint: endpoint}
}

func (s *Client) request(ctx context.Context, method string, params []interface{}, response interface{}) error {
	// post data
	j, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      0,
		"method":  method,
		"params":  params,
	})
	if err != nil {
		return err
	}

	// post request
	req, err := http.NewRequestWithContext(ctx, "POST", s.endpoint, bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	// http client and send request
	httpclient := &http.Client{}

	retry := 0
	var res *http.Response
	for {
		if retry > retryLimit {
			return fmt.Errorf("httpclient reach retry limit, err: %s", err)
		}

		res, err = httpclient.Do(req)
		if err != nil {
			time.Sleep(waitTime)
			retry++
			continue
		}

		defer res.Body.Close()
		// parse body
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			time.Sleep(waitTime)
			retry++
			continue
		}

		if len(body) != 0 {
			if err = json.Unmarshal(body, &response); err != nil {
				time.Sleep(waitTime)
				retry++
				continue
			}

			//check err object
			ge := GeneralResponse{}
			err = json.Unmarshal(body, &ge)
			if err != nil {
				time.Sleep(waitTime)
				retry++
				continue
			} else if ge.Error != (ErrorResponse{}) {
				err = errors.New(ge.Error.Message)
				time.Sleep(waitTime)
				retry++
				continue
			}

		}
		// return result
		if res.StatusCode < 200 || res.StatusCode > 300 {
			err = fmt.Errorf("get status code: %d", res.StatusCode)
			time.Sleep(waitTime)
			retry++
			continue
		}
		break
	}
	return nil
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Context struct {
	Slot uint64 `json:"slot"`
}

type GeneralResponse struct {
	JsonRPC string        `json:"jsonrpc"`
	ID      uint64        `json:"id"`
	Error   ErrorResponse `json:"error"`
}
