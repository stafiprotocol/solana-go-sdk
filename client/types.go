package client

type Instruction struct {
	ProgramIDIndex uint64   `json:"programIdIndex"`
	Accounts       []uint64 `json:"accounts"`
	Data           string   `json:"data"`
}

type TransactionMeta struct {
	Fee               uint64         `json:"fee"`
	PreBalances       []int64        `json:"preBalances"`
	PostBalances      []int64        `json:"postBalances"`
	PreTokenBalances  []TokenBalance `json:"preTokenBalances"`
	PostTokenBalances []TokenBalance `json:"postTokenBalances"`
	LogMessages       []string       `json:"logMessages"`
	InnerInstructions []struct {
		Index        uint64        `json:"index"`
		Instructions []Instruction `json:"instructions"`
	} `json:"innerInstructions"`
	Err    interface{}            `json:"err"`
	Status map[string]interface{} `json:"status"`
}

type MessageHeader struct {
	NumRequiredSignatures       uint8 `json:"numRequiredSignatures"`
	NumReadonlySignedAccounts   uint8 `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts uint8 `json:"numReadonlyUnsignedAccounts"`
}

type Message struct {
	Header          MessageHeader `json:"header"`
	AccountKeys     []string      `json:"accountKeys"`
	RecentBlockhash string        `json:"recentBlockhash"`
	Instructions    []Instruction `json:"instructions"`
}

type Transaction struct {
	Signatures []string `json:"signatures"`
	Message    Message  `json:"message"`
}

type TokenBalance struct {
	AccountIndex  uint64 `json:"accountIndex"`
	Mint          string `json:"mint"`
	UiTokenAmount struct {
		Amount         string  `json:"amount"`
		Decimals       uint64  `json:"decimals"`
		UiAmount       float64 `json:"uiAmount"`
		UiAmountString string  `json:"uiAmountString"`
	} `json:"uiTokenAmount"`
}
