package hdkeystore

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
)

type HDKeyStore struct {
	//文件所在路径
	keyDirPath string
	//生成加密文件的参数N
	scryptN int
	//生成加密文件的参数P
	scryptP int
	//keystore对应的key
	Key keystore.Key
}

func NewHDKeyStore(path string, privateKey *ecdsa.PrivateKey) *HDKeyStore {
	key := keystore.Key{
		Id:         uuid.New(),
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey), //地址信息
		PrivateKey: privateKey, //私钥信息
	}

	return &HDKeyStore{
		keyDirPath: path,
		scryptN:    keystore.LightScryptN, //固定参数
		scryptP:    keystore.LightScryptP, //固定参数
		Key:        key,
	}
}

//支持无私钥创建HDKeyStore
func NewHDKeyStoreWithoutKey(path string) *HDKeyStore {
	return &HDKeyStore{
		keyDirPath: path,
		scryptN:    keystore.LightScryptN,
		scryptP:    keystore.LightScryptP,
		Key:        keystore.Key{},
	}
}

//储存Key为keystore文件
func (ks HDKeyStore) StoreKey(filename string, key *keystore.Key, auth string) error {
	//编码
	keyjson, err := keystore.EncryptKey(key, auth, ks.scryptN, ks.scryptP)
	if err != nil {
		fmt.Println("Failed to EncryptKey: ", err)
		return err
	}
	//写入文件
	return WriteKeyFile(filename, keyjson)
}

//以太坊中的WriteKeyFile源码略作修改
func WriteKeyFile(file string, content []byte) error {
	//根据正确权限创建keystore文件
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return err
	}
	//原子写: 首先创建一个临时隐藏文件
	//然后把它移到位。TempFile分配模式0600。
	f, err := ioutil.TempFile(filepath.Dir(file), "." + filepath.Base(file) + ".tmp")
	if err != nil {
		return err
	}
	if _, err := f.Write(content); err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return err
	}
	f.Close()
	return os.Rename(f.Name(), file)
}

//实现路径拼接接口，用于将路径与文件名拼接
func (ks HDKeyStore) JoinPath(filename string) string {
	//如果filename是绝对路径则直接返回
	if filepath.IsAbs(filename) {
		return filename
	}
	//将路径与文件名拼接
	return filepath.Join(ks.keyDirPath, filename)
}

//实现keystore文件解析接口
func (ks *HDKeyStore) GetKey(addr common.Address, filename, auth string) (*keystore.Key, error) {
	keyjson, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	//利用以太坊DecryptKey解码json文件
	key, err := keystore.DecryptKey(keyjson, auth)
	if err != nil {
		return nil, err
	}
	//如果地址不同则代表解析失败
	if key.Address != addr {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, addr)
	}
	ks.Key = *key
	return key, nil
}

//交易签名方法
func (ks *HDKeyStore) SignTx(address common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, ks.Key.PrivateKey)
	if err != nil {
		return nil, err
	}
	//验证签名
	msg, err := signedTx.AsMessage(types.HomesteadSigner{}, nil)
	if err != nil {
		return nil, err
	}
	sender := msg.From()
	if sender != address {
		return nil, fmt.Errorf("signer mismatch: expected %s, got %s", address.Hex(), sender.Hex())
	}
	return signedTx, nil
}

//利用keystore生成token合约调用身份
func (ks HDKeyStore) NewTransactOpts() *bind.TransactOpts {
	return bind.NewKeyedTransactor(ks.Key.PrivateKey)
}