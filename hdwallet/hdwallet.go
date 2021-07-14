package hdwallet

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"receipt/hdkeystore"
)

type HDWallet struct {
	Address common.Address
	HDKeystore *hdkeystore.HDKeyStore
}

func NewHDWallet(keypath string) (*HDWallet, error) {
	//1.创建助记词
	mne, err := createMnemonic()
	if err != nil {
		fmt.Println("Failed to NewHDWallet: ", err)
		return nil, err
	}
	//2. 推导私钥
	privateKey, err := DerivePrivateKeyFromMnemonic(mne)
	if err != nil {
		fmt.Println("Failed to DerivePrivateKeyFromMnemonic: ", err)
		return nil, err
	}
	//3. 获取地址
	publicKey, err := DerivePublicKey(privateKey)
	if err != nil {
		fmt.Println("Failed to DerivePublicKey: ", err)
		return nil, err
	}
	//通过公钥推导地址
	address := crypto.PubkeyToAddress(*publicKey)
	//4. 创建keystore
	hdks := hdkeystore.NewHDKeyStore(keypath, privateKey)
	//5. 创建钱包
	return &HDWallet{
		Address:    address,
		HDKeystore: hdks,
	}, nil
}

////通过账户文件来构建钱包文件，以用户输入非明文方式获得私钥
//func LoadWallet(filename, keypath string) (*HDWallet, error) {
//	//在无私钥时创建钱包
//	hdks := hdkeystore.NewHDKeyStoreWithoutKey(keypath)
//	fmt.Println("Please input password for: ", filename)
//	pass, _ := gopass.GetPasswd()
//	//filename也是账户地址
//	fromaddr := common.HexToAddress(filename)
//	_, err := hdks.GetKey(fromaddr, hdks.JoinPath(filename), string(pass))
//	if err != nil {
//		log.Panic("Failed to GetKey: ", err)
//	}
//	return &HDWallet{
//		Address:    fromaddr,
//		HDKeystore: hdks,
//	}, nil
//}

func LoadWalletByPassword(filename, datadir, pw string) (*HDWallet, error) {
	hdks := hdkeystore.NewHDKeyStoreWithoutKey(datadir)
	fromaddr := common.HexToAddress(filename)
	_, err := hdks.GetKey(fromaddr, hdks.JoinPath(filename), pw)
	if err != nil {
		log.Panic("Failed to GetKey: ", err)
	}
	return &HDWallet{fromaddr, hdks}, nil
}

func (hw *HDWallet) StoreKey(pass string) error {
	//账户即文件名
	filename := hw.HDKeystore.JoinPath(hw.Address.Hex())
	return hw.HDKeystore.StoreKey(filename, &hw.HDKeystore.Key, pass)
}