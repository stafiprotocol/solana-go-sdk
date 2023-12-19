package minterprog

import "github.com/stafiprotocol/solana-go-sdk/common"

type MintManager struct {
	Admin                 common.PublicKey
	RSolMint              common.PublicKey
	MintAuthoritySeedBump uint8
	ExtMintAuthorities    []common.PublicKey
}
