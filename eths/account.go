package eths

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"receipt/hdwallet"
)

func NewAccount(pw string) (string, error) {
	w, err := hdwallet.NewHDWallet("./keystore")
	if err != nil {
		fmt.Println("Failed to NewWallet: ", err)
	}
	err = w.StoreKey(pw)
	if err != nil {
		fmt.Println("Failed to StoreKey: ", err)
	}
	return w.Address.Hex(), nil
}

func KeccakHash(data []byte) []byte {
	return crypto.Keccak256(data)
}

