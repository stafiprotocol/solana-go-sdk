package client_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stafiprotocol/solana-go-sdk/client"
)

func TestAccountInfo(t *testing.T) {
	c := client.NewClient(client.DevnetRPCEndpoint)

	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			accountInfo, err := c.GetMultisigTxAccountInfo(context.Background(), "D6nA6QHpYQDMeudHLwZqgwyCJfRSKWfzW4kyaKqmnsr4")
			if err != nil {
				t.Fatal(err)
			}
			t.Log(fmt.Printf("%+v", accountInfo))
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestGetAccountInfo(t *testing.T) {
	c := client.NewClient(client.DevnetRPCEndpoint)
	account, err := c.GetAccountInfo(context.Background(), "DiPx1Vyo5khyG8XKTc8Tu4fL9qc57VSqfr7qh3xLxqjX",
		client.GetAccountInfoConfig{
			Encoding: client.GetAccountInfoConfigEncodingBase64,
			DataSlice: client.GetAccountInfoConfigDataSlice{
				Offset: 0,
				Length: 200,
			},
		})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(account)

	accountInfo, err := c.GetStakeAccountInfo(context.Background(), "DiPx1Vyo5khyG8XKTc8Tu4fL9qc57VSqfr7qh3xLxqjX")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountInfo))

	accountInfo, err = c.GetStakeAccountInfo(context.Background(), "mNzHTv7KtARcYyJiaXH3SU2oSnLoXKVJxcxDknC2kae")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountInfo))

	accountActivateInfo, err := c.GetStakeActivation(context.Background(), "DiPx1Vyo5khyG8XKTc8Tu4fL9qc57VSqfr7qh3xLxqjX", client.GetStakeActivationConfig{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountActivateInfo))

}

func TestGetStakeActivation(t *testing.T) {
	c := client.NewClient("http://127.0.0.1:8899")
	accountActivateInfo, err := c.GetStakeActivation(context.Background(), "D17ya9gd9xRSMwqjpzixXx341gPF2sNzKft5CMnToF8h", client.GetStakeActivationConfig{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountActivateInfo))
}

func TestGetStakeAccountInfo(t *testing.T) {
	c := client.NewClient("http://127.0.0.1:8899")
	// accountActivateInfo, err := c.GetStakeAccountInfo(context.Background(), "D17ya9gd9xRSMwqjpzixXx341gPF2sNzKft5CMnToF8h")
	accountActivateInfo, err := c.GetStakeAccountInfo(context.Background(), "Gnr9LuHUh85Dt7Qr3tayXrxFAEn32jRDfsgTAyywFhyh")
	// accountActivateInfo, err := c.GetStakeAccountInfo(context.Background(), "ATY5PSBVExLoFb2CnRj1e9nUVghcvLcrvbhcYMud1d4F")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountActivateInfo))
	accountActivateInfo2, err := c.GetStakeActivation(context.Background(), "2usyY4HuMCfZW6CGjzdUuheMxh71HPQFkyZZW61qAxAq", client.GetStakeActivationConfig{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", accountActivateInfo2))

	account, err := c.GetAccountInfo(context.Background(), "2usyY4HuMCfZW6CGjzdUuheMxh71HPQFkyZZW61qAxAq",
		client.GetAccountInfoConfig{
			Encoding: client.GetAccountInfoConfigEncodingBase64,
			DataSlice: client.GetAccountInfoConfigDataSlice{
				Offset: 0,
				Length: 200,
			},
		})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", account))

}

func TestGetMultisigTxInfo(t *testing.T) {
	c := client.NewClient("https://api.devnet.solana.com")
	info, err := c.GetMultisigTxAccountInfo(context.Background(), "Gn3Wzs1rbeJcTefiEwZ8c8vJZjNZeSm5WUxbYC5ji74F")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
}


func TestGetMultisigInfoAccount(t *testing.T){
	c := client.NewClient("https://api.devnet.solana.com")
	info, err := c.GetMultisigInfoAccountInfo(context.Background(), "8TNEsKSzFsi6b56JwhpHWLZf9mR81LGDcQQka5EtVux7")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
}

func TestGetBridgeAccountInfo(t *testing.T){
	c := client.NewClient("https://api.devnet.solana.com")
	info, err := c.GetBridgeAccountInfo(context.Background(), "8B29iREQvQgmiyZSdzRrgEsh56W3m2Mpna5xvGG6jAEf")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
}


func TestGetMintProposalInfo(t *testing.T){
	c := client.NewClient("https://api.devnet.solana.com")
	info, err := c.GetMintProposalInfo(context.Background(), "HkYNY6TdY6HpfYR1bqGyJjMexz6Tzja7kyjavzQhm3rq")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", info))
}