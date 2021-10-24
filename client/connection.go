package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	DevnetRPCEndpoint  = "https://api.devnet.solana.com"
	TestnetRPCEndpoint = "https://testnet.solana.com"
	MainnetRPCEndpoint = "https://api.mainnet-beta.solana.com"
	retryLimit         = 60 * 10
	waitTime           = time.Second * 3
)

type Commitment string

const (
	CommitmentFinalized Commitment = "finalized"
	CommitmentConfirmed Commitment = "confirmed"
	CommitmentProcessed Commitment = "processed"
	CommitmentRecent    Commitment = "recent"
)

type Client struct {
	index        int
	indexMutex   sync.Mutex
	endpointList []string
}

func NewClient(endpointList []string) *Client {
	if len(endpointList) == 0 {
		panic("endpoint empty")
	}
	return &Client{endpointList: endpointList}
}

func (s *Client) Endpoint() string {
	return s.endpointList[s.index]
}

func (s *Client) ChangeEndpoint() {
	s.indexMutex.Lock()
	defer s.indexMutex.Unlock()

	next := (s.index + 1) % len(s.endpointList)
	s.index = next
}

// err will retry: 1) connection err 2) body read err 3) status code err
// err will return: 1) reach retry err 2) rpc res error
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

	retry := 0
	for {
		s.ChangeEndpoint()
		if retry > retryLimit {
			return fmt.Errorf("httpclient reach retry limit, err: %s", err)
		}
		// post request
		var req *http.Request
		req, err = http.NewRequestWithContext(ctx, "POST", s.Endpoint(), bytes.NewBuffer(j))
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/json")
		// http client and send request
		httpclient := &http.Client{}

		var res *http.Response
		res, err = httpclient.Do(req)
		if err != nil {
			time.Sleep(waitTime)
			retry++
			continue
		}

		// parse body
		var body []byte
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			time.Sleep(waitTime)
			retry++
			res.Body.Close()
			continue
		}

		if len(body) != 0 {
			if err = json.Unmarshal(body, &response); err != nil {
				time.Sleep(waitTime)
				retry++
				res.Body.Close()
				continue
			}
		}
		// check status code
		if res.StatusCode < 200 || res.StatusCode > 300 {
			err = fmt.Errorf("get status code: %d", res.StatusCode)
			time.Sleep(waitTime)
			retry++
			res.Body.Close()
			continue
		}

		res.Body.Close()
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
