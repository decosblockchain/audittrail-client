package library

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/decosblockchain/audittrail-client/library/fastsha256"
)

type CkTransaction struct {
	Inputs    []CkTransactionInput  `json:"inputs"`
	Outputs   []CkTransactionOutput `json:"outputs"`
	Timestamp int64                 `json:"timestamp"`
}

type CkTransactionInput struct {
	OutputId string                 `json:"outputId"`
	Data     CkTransactionInputData `json:"data,omitempty"`
}

type CkTransactionInputData struct {
	SignatureBase64 string `json:"signature"`
}

type CkTransactionOutput struct {
	Value int64                   `json:"value"`
	Nonce uint64                  `json:"nonce"`
	Data  CkTransactionOutputData `json:"data"`
}

type CkTransactionOutputData struct {
	Contract  []byte `json:"contract,omitempty"`
	PublicKey []byte `json:"publicKey,omitempty"`
	AuditData []byte `json:"auditData,omitempty"`
}

func (d CkTransactionOutputData) String() (string, error) {
	js, err := json.Marshal(d)
	if err != nil {
		log.Printf("Error marshalling outputdata to json: %s", err.Error())
		return "", err
	}

	return string(js), nil
}

func (d CkTransactionInputData) String() (string, error) {
	js, err := json.Marshal(d)
	if err != nil {
		log.Printf("Error marshalling inputdata to json: %s", err.Error())
		return "", err
	}

	return string(js), nil
}

func (tx *CkTransaction) GetFee() int64 {
	fee := int64(0)
	for _, in := range tx.Inputs {
		js, _ := json.Marshal(in.Data)
		fee += int64(len(js)) * int64(100)
	}
	for _, out := range tx.Outputs {
		js, _ := json.Marshal(out.Data)
		fee += int64(len(js)) * int64(100)
	}
	return fee
}

func (tx CkTransaction) GetOutputSetId() string {
	outputIds := make([]*big.Int, len(tx.Outputs))
	for i, o := range tx.Outputs {
		outputIds[i], _ = new(big.Int).SetString(o.Id(), 16)
	}
	node := MakeMerkleTree(outputIds)
	return node.MerkleRoot().Text(16)
}

func (tx CkTransaction) GetInputSetId() string {
	if len(tx.Inputs) == 0 {
		return ""
	}

	inputIds := make([]*big.Int, len(tx.Inputs))
	for i, inp := range tx.Inputs {
		inputIds[i], _ = new(big.Int).SetString(inp.Id(), 16)
	}
	node := MakeMerkleTree(inputIds)
	return node.MerkleRoot().Text(16)
}

func (tx *CkTransaction) Id() string {
	hashInput := fmt.Sprintf("%s%s%d", tx.GetInputSetId(), tx.GetOutputSetId(), tx.Timestamp)
	hash := fastsha256.Sum256([]byte(hashInput))
	return fmt.Sprintf("%x", hash)
}

func (i CkTransactionInput) Id() string {
	dataJson, _ := i.Data.String()
	hashInput := fmt.Sprintf("%s%s", i.OutputId, dataJson)
	hash := fastsha256.Sum256([]byte(hashInput))
	return fmt.Sprintf("%x", hash)
}

func (o CkTransactionOutput) Id() string {
	dataJson, _ := o.Data.String()
	hashInput := fmt.Sprintf("%d%d%s", o.Value, o.Nonce, dataJson)
	hash := fastsha256.Sum256([]byte(hashInput))
	return fmt.Sprintf("%x", hash)
}
