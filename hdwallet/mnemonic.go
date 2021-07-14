package hdwallet

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"log"
)

func createMnemonic() (string, error){
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Panic("Failed to NewEntropy: ", err, entropy)
	}
	//生成助记词
	mne, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Panic("Failed to NewMnemonic: ", err)
	}
	fmt.Println(mne)
	return mne, nil
}

const defaultPath = "m/44'/60'/0'/0/1"

func DeriveAddressFromMnemonic(mne string) string {
	//1. 推导目录
	path, err := accounts.ParseDerivationPath(defaultPath)
	if err != nil {
		log.Panic("Failed to ParseDerivationPath: ", err)
	}

	//2. 通过助记词生成种子
	seed, err := bip39.NewSeedWithErrorChecking(mne, "")
	if err != nil {
		log.Panic("Failed to NewSeedWithErrorChecking: ", err)
	}

	//3. 通过seed获取master key
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		log.Panic("Failed to NewMaster: ", err)
	}

	//4. 推导私钥
	privateKey, err := DerivePrivateKey(path, masterKey)
	if err != nil {
		log.Panic("Failed to DerivePrivateKey: ", err)
	}
	//5. 推导公钥
	publicKey, err := DerivePublicKey(privateKey)
	if err != nil {
		log.Panic("Failed to DerivePublicKey: ", err)
	}
	//6. 利用公钥推导地址
	address := crypto.PubkeyToAddress(*publicKey)
	return address.Hex()
}

//通过助记词推导私钥
func DerivePrivateKeyFromMnemonic(mne string) (*ecdsa.PrivateKey, error) {
	//1. 推导目录
	path, err := accounts.ParseDerivationPath(defaultPath)
	if err != nil {
		log.Panic("Failed to ParseDerivationPath: ", err)
	}

	//2. 通过助记词生成种子
	seed, err := bip39.NewSeedWithErrorChecking(mne, "")
	if err != nil {
		log.Panic("Failed to NewSeedWithErrorChecking: ", err)
	}

	//3. 通过seed获取master key
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		log.Panic("Failed to NewMaster: ", err)
	}

	//4. 推导私钥
	return DerivePrivateKey(path, masterKey)
}

//推导私钥
func DerivePrivateKey(path accounts.DerivationPath, masterKey *hdkeychain.ExtendedKey) (*ecdsa.PrivateKey, error) {
	var err error
	key := masterKey
	for _, n := range path {
		//按照路径迭代获取最终的Key
		key, err = key.Derive(n)
		if err != nil {
			return nil, err
		}
	}
	//将key转化为ecdsa私钥
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	privateKeyECDSA := privateKey.ToECDSA()
	return privateKeyECDSA, nil
}

//推导公钥
func DerivePublicKey(privateKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}
	return publicKeyECDSA, nil
}