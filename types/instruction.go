package types

import "github.com/stafiprotocol/solana-go-sdk/common"

type CompiledInstruction struct {
	ProgramIDIndex int
	Accounts       []int
	Data           []byte
}

type Instruction struct {
	ProgramID common.PublicKey
	Accounts  []AccountMeta //accounts programs will use
	Data      []byte        //instruct + params
}
