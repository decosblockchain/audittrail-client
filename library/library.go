package library

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/decosblockchain/audittrail-client/config"
	"github.com/ethereum/go-ethereum/crypto"
)

var keyMutex = &sync.Mutex{}

var key *ecdsa.PrivateKey
var nonceInitialized bool
var nonce uint64

func GetAddress() (string, error) {
	key, err := GetPubKey()
	if err != nil {
		return "", err
	}

	addr := base64.StdEncoding.EncodeToString(key)
	return addr, nil
}

func GetUnspentOutpoint() (string, int64, error) {
	return "8271fa1effc361d9c034309537e25df3f9d858744a1d8f07efae49bcb5cbb259", 1000000000, nil
}

func CreateTransaction(auditData []byte) (*CkTransaction, error) {
	tx := new(CkTransaction)

	outputId, outputSize, err := GetUnspentOutpoint()
	if err != nil {
		return nil, err
	}

	pubKey, err := GetPubKey()
	if err != nil {
		return nil, err
	}

	fee := int64(500)

	tx.Inputs = make([]CkTransactionInput, 1)
	tx.Inputs[0] = CkTransactionInput{}
	tx.Inputs[0].OutputId = outputId
	tx.Inputs[0].Data = CkTransactionInputData{}

	tx.Outputs = make([]CkTransactionOutput, 2)
	tx.Outputs[0] = CkTransactionOutput{}
	tx.Outputs[0].Value = 0
	tx.Outputs[0].Nonce, err = GetNonce()
	if err != nil {
		return nil, err
	}
	tx.Outputs[0].Data = CkTransactionOutputData{}
	tx.Outputs[0].Data.AuditData = auditData

	tx.Outputs[1] = CkTransactionOutput{}
	tx.Outputs[1].Value = outputSize - fee
	tx.Outputs[1].Nonce, err = GetNonce()
	if err != nil {
		return nil, err
	}
	tx.Outputs[1].Data = CkTransactionOutputData{}
	tx.Outputs[1].Data.PublicKey = pubKey

	return tx, nil
}

func SignTransaction(tx *CkTransaction) error {
	return nil
}

func GetPubKey() ([]byte, error) {
	key, err := GetKey()
	if err != nil {
		return []byte{}, err
	}

	return SerializePublicKey(key.PublicKey), nil
}

func GetKey() (*ecdsa.PrivateKey, error) {

	keyMutex.Lock()
	if key == nil {
		if _, err := os.Stat(path.Join(config.BaseDir(), "data", "keyfile.hex")); os.IsNotExist(err) {
			generatedKey, err := crypto.GenerateKey()
			if err != nil {
				keyMutex.Unlock()
				return nil, err
			}
			err = crypto.SaveECDSA(path.Join(config.BaseDir(), "data", "keyfile.hex"), generatedKey)
			if err != nil {
				keyMutex.Unlock()
				return nil, err
			}
		}

		privateKey, err := crypto.LoadECDSA(path.Join(config.BaseDir(), "data", "keyfile.hex"))
		if err != nil {
			keyMutex.Unlock()
			return nil, err
		}
		key = privateKey

	}
	keyMutex.Unlock()

	return key, nil
}

// paddedAppend appends the src byte slice to dst, returning the new slice.
// If the length of the source is smaller than the passed size, leading zero
// bytes are appended to the dst slice before appending src.
func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

func SerializePublicKey(p ecdsa.PublicKey) []byte {
	b := make([]byte, 0, 65)
	b = append(b, 0x4)
	b = paddedAppend(32, b, p.X.Bytes())
	return paddedAppend(32, b, p.Y.Bytes())
}

func GetNonce() (uint64, error) {
	b := make([]byte, 8)
	n, err := rand.Read(b)
	if n != 8 {
		return 0, fmt.Errorf("Could not read 8 random bytes, read %d instead", n)
	} else if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b), nil
}
