package common

import (
	"errors"
	"reflect"
	"testing"
)

func TestPublicKeyFromBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want PublicKey
	}{
		{
			args: args{b: []byte{1}},
			want: PublicKey([32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
		},
		{
			args: args{b: []byte{1, 2, 3}},
			want: PublicKey([32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3}),
		},
		{
			args: args{b: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3}},
			want: PublicKey([32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PublicKeyFromBytes(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PublicKeyFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublicKeyFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want PublicKey
	}{
		{
			args: args{s: "11111111111111111111111111111111"},
			want: PublicKey([32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		},
		{
			args: args{s: "Config1111111111111111111111111111111111111"},
			want: PublicKey([32]byte{0x03, 0x06, 0x4a, 0xa3, 0x00, 0x2f, 0x74, 0xdc, 0xc8, 0x6e, 0x43, 0x31, 0x0f, 0x0c, 0x05, 0x2a, 0xf8, 0xc5, 0xda, 0x27, 0xf6, 0x10, 0x40, 0x19, 0xa3, 0x23, 0xef, 0xa0, 0x00, 0x00, 0x00, 0x00}),
		},
		{
			args: args{s: "Stake11111111111111111111111111111111111111"},
			want: PublicKey([32]byte{0x06, 0xa1, 0xd8, 0x17, 0x91, 0x37, 0x54, 0x2a, 0x98, 0x34, 0x37, 0xbd, 0xfe, 0x2a, 0x7a, 0xb2, 0x55, 0x7f, 0x53, 0x5c, 0x8a, 0x78, 0x72, 0x2b, 0x68, 0xa4, 0x9d, 0xc0, 0x00, 0x00, 0x00, 0x00}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PublicKeyFromString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PublicKeyFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateProgramAddress(t *testing.T) {
	type args struct {
		seeds     [][]byte
		ProgramID PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    PublicKey
		wantErr error
	}{
		{
			args: args{
				seeds:     [][]byte{{0x1}},
				ProgramID: PublicKeyFromString("EmPaWGCw48Sxu9Mu9pVrxe4XL2JeXUNTfoTXLuLz31gv"),
			},
			want:    PublicKeyFromString("65JQyZBU2RzNpP9vTdW5zSzujZR5JHZyChJsDWvkbM8u"),
			wantErr: nil,
		},
		{
			args: args{
				seeds:     [][]byte{{0x2}},
				ProgramID: PublicKeyFromString("EmPaWGCw48Sxu9Mu9pVrxe4XL2JeXUNTfoTXLuLz31gv"),
			},
			want:    PublicKey{},
			wantErr: errors.New("Invalid seeds, address must fall off the curve"),
		},
		{
			args: args{
				seeds:     [][]byte{[]byte("123456789012345678901234567890123")},
				ProgramID: PublicKeyFromString("EmPaWGCw48Sxu9Mu9pVrxe4XL2JeXUNTfoTXLuLz31gv"),
			},
			want:    PublicKey{},
			wantErr: errors.New("Max seed length exceeded"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateProgramAddress(tt.args.seeds, tt.args.ProgramID)
			if tt.wantErr != nil && errors.Is(err, tt.wantErr) {
				t.Errorf("CreateProgramAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateProgramAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindAssociatedTokenAddress(t *testing.T) {
	type args struct {
		walletAddress    PublicKey
		tokenMintAddress PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    PublicKey
		want1   int
		wantErr bool
	}{
		{
			args: args{
				walletAddress:    PublicKeyFromString("EvN4kgKmCmYzdbd5kL8Q8YgkUW5RoqMTpBczrfLExtx7"),
				tokenMintAddress: PublicKeyFromString("8765cK2Vucsic6NA5nm4cfkrCzusaFVqBf6Pk31tGkXH"),
			},
			want:    PublicKeyFromString("HLzppk6ohPg9Ab99XTFhsa6FcG14Au3rTijGe9c8QHp1"),
			want1:   254,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := FindAssociatedTokenAddress(tt.args.walletAddress, tt.args.tokenMintAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAssociatedTokenAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAssociatedTokenAddress() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindAssociatedTokenAddress() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCreateWithSeed(t *testing.T) {
	type args struct {
		from      PublicKey
		seed      string
		programID PublicKey
	}
	tests := []struct {
		name string
		args args
		want PublicKey
	}{
		{
			args: args{
				from:      PublicKeyFromString("EvN4kgKmCmYzdbd5kL8Q8YgkUW5RoqMTpBczrfLExtx7"),
				seed:      "0",
				programID: SystemProgramID,
			},
			want: PublicKeyFromString("DTA7FmUNYuQs2mScj2Lx8gQV63SEL1zGtzCSvPxtijbi"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateWithSeed(tt.args.from, tt.args.seed, tt.args.programID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateWithSeed() = %v, want %v", got, tt.want)
			}
		})
	}
}
