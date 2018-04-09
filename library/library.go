package library

import (
	"crypto/ecdsa"
	"encoding/binary"
	"io/ioutil"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
)

var keyMutex = &sync.Mutex{}
var nonceMutex = &sync.Mutex{}

var key *ecdsa.PrivateKey
var nonceInitialized bool
var nonce uint64

func GetAddress() (string, error) {
	key, err := GetKey()
	if err != nil {
		return "", err
	}

	addr := crypto.PubkeyToAddress(key.PublicKey)
	return addr.String(), nil
}

func GetKey() (*ecdsa.PrivateKey, error) {

	keyMutex.Lock()
	if key == nil {
		if _, err := os.Stat("data/keyfile.hex"); os.IsNotExist(err) {
			generatedKey, err := crypto.GenerateKey()
			if err != nil {
				keyMutex.Unlock()
				return nil, err
			}
			err = crypto.SaveECDSA("data/keyfile.hex", generatedKey)
			if err != nil {
				keyMutex.Unlock()
				return nil, err
			}
		}

		privateKey, err := crypto.LoadECDSA("data/keyfile.hex")
		if err != nil {
			keyMutex.Unlock()
			return nil, err
		}
		key = privateKey

	}
	keyMutex.Unlock()

	return key, nil
}

func writeNonce() error {

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, nonce)

	err := ioutil.WriteFile("data/nonce.hex", b, 0600)
	return err
}

func readNonce() error {
	if _, err := os.Stat("data/nonce.hex"); os.IsNotExist(err) {
		err := writeNonce()
		if err != nil {
			return err
		}
	}
	b, err := ioutil.ReadFile("data/nonce.hex")
	if err != nil {
		return err
	}

	nonce = binary.LittleEndian.Uint64(b)
	return nil
}

func GetNonce() (uint64, error) {

	nonceMutex.Lock()
	if !nonceInitialized {
		err := readNonce()
		if err != nil {
			nonceMutex.Unlock()
			return 0, err
		}
		nonceInitialized = true
	}
	returnValue := uint64(nonce)

	// Write next nonce
	nonce++
	err := writeNonce()
	if err != nil {
		nonceMutex.Unlock()
		return 0, err
	}
	nonceMutex.Unlock()

	return returnValue, nil
}

func CancelNonce() error {
	nonceMutex.Lock()
	if !nonceInitialized {
		err := readNonce()
		if err != nil {
			nonceMutex.Unlock()
			return err
		}
		nonceInitialized = true
	}
	nonce--
	err := writeNonce()
	if err != nil {
		nonceMutex.Unlock()
		return err
	}
	nonceMutex.Unlock()
	return nil
}
