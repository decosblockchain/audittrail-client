package library

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"os"
	"path"
	"sync"
	"time"

	"github.com/decosblockchain/audittrail-client/config"
	"github.com/decosblockchain/audittrail-client/library/fastsha256"
	"github.com/ethereum/go-ethereum/crypto"
)

var keyMutex = &sync.Mutex{}
var unspentOutputId = "9bb37358a4560fd1d63607ac2025ce85f08e73c5c5d9a3cd85706e48e75ff99f"
var unspentOutputValue = int64(100000000)
var key *ecdsa.PrivateKey

type Signature struct {
	R *big.Int
	S *big.Int
}

func GetAddress() (string, error) {
	key, err := GetPubKey()
	if err != nil {
		return "", err
	}

	addr := base64.StdEncoding.EncodeToString(key)
	return addr, nil
}

func GetPrivateKeyBase64() (string, error) {
	key, err := GetKey()
	if err != nil {
		return "", err
	}

	keyBytes := SerializePrivateKey(key)
	base64Key := base64.StdEncoding.EncodeToString(keyBytes)
	return base64Key, nil
}

func GetUnspentOutpoint() (string, int64, error) {
	return unspentOutputId, unspentOutputValue, nil
}

func CreateTransaction(auditData []byte) (*CkTransaction, error) {
	tx := new(CkTransaction)

	tx.Timestamp = time.Now().Unix()

	outputId, outputSize, err := GetUnspentOutpoint()
	if err != nil {
		return nil, err
	}

	pubKey, err := GetPubKey()
	if err != nil {
		return nil, err
	}

	tx.Inputs = make([]CkTransactionInput, 1)
	tx.Inputs[0] = CkTransactionInput{}
	tx.Inputs[0].OutputId = outputId
	tx.Inputs[0].Data = CkTransactionInputData{}

	tx.Outputs = make([]CkTransactionOutput, 2)
	tx.Outputs[0] = CkTransactionOutput{}
	tx.Outputs[0].Value = 1
	tx.Outputs[0].Nonce, err = GetNonce()
	if err != nil {
		return nil, err
	}
	tx.Outputs[0].Data = CkTransactionOutputData{}
	tx.Outputs[0].Data.AuditData = auditData

	tx.Outputs[1] = CkTransactionOutput{}
	tx.Outputs[1].Value = outputSize
	tx.Outputs[1].Nonce, err = GetNonce()
	if err != nil {
		return nil, err
	}
	tx.Outputs[1].Data = CkTransactionOutputData{}
	tx.Outputs[1].Data.PublicKey = pubKey

	tx.Outputs[1].Value = outputSize - tx.GetFee() - 1

	unspentOutputId = tx.Outputs[1].Id()
	unspentOutputValue = tx.Outputs[1].Value

	data, _ := tx.Outputs[1].Data.String()

	log.Printf("New unspentOutput will be id [%s] - based on Nonce [%d], Data [%s], Value [%d]\n", unspentOutputId, tx.Outputs[1].Nonce, data, unspentOutputValue)
	log.Printf("Transaction ID is: %s\n", tx.Id())

	return tx, nil
}

func SignTransaction(tx *CkTransaction) error {
	log.Printf("Signing %d inputs\n", len(tx.Inputs))

	outputHash := tx.GetOutputSetId()
	for idx, txi := range tx.Inputs {
		signatureInput := txi.OutputId + outputHash
		log.Printf("Signature input for input %d: %s", idx, signatureInput)
		sig, err := Sign(signatureInput)
		if err != nil {
			return err
		}
		tx.Inputs[idx].Data.SignatureBase64 = sig
	}
	return nil
}

func Sign(message string) (string, error) {
	digest := fastsha256.Sum256([]byte(message))
	key, err := GetKey()
	if err != nil {
		return "", err
	}

	r, s, err := ecdsa.Sign(rand.Reader, key, digest[:])
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(PointsToDER(r, s))
	return signature, nil
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

func SerializePrivateKey(k *ecdsa.PrivateKey) []byte {
	privateKeyBytes := k.D.Bytes()
	paddedPrivateKey := make([]byte, (key.Curve.Params().N.BitLen()+7)/8)
	copy(paddedPrivateKey[len(paddedPrivateKey)-len(privateKeyBytes):], privateKeyBytes)
	return paddedPrivateKey
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

// Convert an ECDSA signature (points R and S) to a byte array using ASN.1 DER encoding.
// This is a port of Bitcore's Key.rs2DER method.
func PointsToDER(r, s *big.Int) []byte {
	// Ensure MSB doesn't break big endian encoding in DER sigs
	prefixPoint := func(b []byte) []byte {
		if len(b) == 0 {
			b = []byte{0x00}
		}
		if b[0]&0x80 != 0 {
			paddedBytes := make([]byte, len(b)+1)
			copy(paddedBytes[1:], b)
			b = paddedBytes
		}
		return b
	}

	rb := prefixPoint(r.Bytes())
	sb := prefixPoint(s.Bytes())

	// DER encoding:
	// 0x30 + z + 0x02 + len(rb) + rb + 0x02 + len(sb) + sb
	length := 2 + len(rb) + 2 + len(sb)

	der := append([]byte{0x30, byte(length), 0x02, byte(len(rb))}, rb...)
	der = append(der, 0x02, byte(len(sb)))
	der = append(der, sb...)

	return der
}
