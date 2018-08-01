package library

import (
	"encoding/json"
	"fmt"

	"github.com/decosblockchain/audittrail-client/library/fastsha256"
)

type CkTransaction struct {
	Inputs  []CkTransactionInput  `json:"inputs"`
	Outputs []CkTransactionOutput `json:"outputs"`
}

type CkTransactionInput struct {
	OutputId string                 `json:"outputId"`
	Data     CkTransactionInputData `json:"data,omitempty"`
}

type CkTransactionInputData struct {
	Signature []byte `json:"signature"`
}

type CkTransactionOutput struct {
	Value int64                   `json:"value"`
	Nonce uint64                  `json:"nonce"`
	Data  CkTransactionOutputData `json:"data"`
}

type CkTransactionOutputData struct {
	PublicKey []byte `json:"publicKey,omitempty"`
	AuditData []byte `json:"auditData,omitempty"`
}

func (d CkTransactionOutputData) String() (string, error) {
	js, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	return string(js), nil
}

func (o CkTransactionOutput) Id() (string, error) {
	dataJson, err := o.Data.String()
	if err != nil {
		return "", err
	}
	hashInput := fmt.Sprintf("%d%d%s", o.Value, o.Nonce, dataJson)
	hash := fastsha256.Sum256([]byte(hashInput))
	return fmt.Sprintf("%x", hash), nil
}
