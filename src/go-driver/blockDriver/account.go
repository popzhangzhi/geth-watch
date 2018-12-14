package blockDriver

import (
	"crypto/ecdsa"
	crand "crypto/rand"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
)

/**
创造离线eth秘钥对
*/
func DoCreate() (string, string, error) {
	rand := crand.Reader
	address, privateKey, err := storeNewKey(rand)
	return address, privateKey, err

}
func storeNewKey(rand io.Reader) (string, string, error) {
	addr, pri, err := newKey(rand)
	if err != nil {
		return "", "", err
	}

	return addr.String(), hexutil.Encode(pri), err

}
func newKey(rand io.Reader) (common.Address, []byte, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand)
	if err != nil {
		return common.Address{}, nil, err
	}
	address, pri := newKeyFromECDSA(privateKeyECDSA)
	return address, pri, nil
}
func newKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) (common.Address, []byte) {
	//id := uuid.NewRandom()
	//key := &keystore.Key{
	//	Id:         id,
	//	Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
	//	PrivateKey: privateKeyECDSA,
	//}
	address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
	privateKey := crypto.FromECDSA(privateKeyECDSA)

	return address, privateKey
}
