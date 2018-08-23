package main

import (
	"crypto/ecdsa"
	crand "crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pborman/uuid"
	"io"
)

func DoCreate() {
	rand := crand.Reader
	storeNewKey(rand)

}
func storeNewKey(rand io.Reader) (*keystore.Key, accounts.Account, error) {
	key, err := newKey(rand)
	if err != nil {
		return nil, accounts.Account{}, err
	}
	a := accounts.Account{Address: key.Address}

	return key, a, err
}
func newKey(rand io.Reader) (*keystore.Key, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand)
	if err != nil {
		return nil, err
	}
	return newKeyFromECDSA(privateKeyECDSA), nil
}
func newKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) *keystore.Key {
	id := uuid.NewRandom()
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	fmt.Println("address:", key.Address)
	fmt.Println("address:", key.Address.String())
	fmt.Println("privateKey:", privateKeyECDSA.D)
	a := crypto.FromECDSA(privateKeyECDSA)
	fmt.Println(hexutil.Encode(a))
	return key
}
