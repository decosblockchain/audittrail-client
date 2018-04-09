package routes

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/decosblockchain/audittrail-client/config"
	"github.com/decosblockchain/audittrail-client/library"
	"github.com/decosblockchain/audittrail-client/logging"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type AuditRecord struct {
	Header  AuditRecordHeader   `json:"header"`
	Details []AuditRecordDetail `json:"details"`
}

type AuditRecordHeader struct {
	Actor  string `json:"actor"`
	Intent string `json:"intent"`
	Object string `json:"object"`
}

type AuditRecordDetail struct {
	Key   string `json:"k"`
	Value string `json:"v"`
}

type AuditResponse struct {
	RecordHash      string `json:"recordHash"`
	TransactionHash string `json:"transactionHash"`
}

func AuditHandler(w http.ResponseWriter, r *http.Request) {
	logging.Info.Printf("Received request to /audit\n")
	if r.Method != "POST" {
		logging.Error.Printf("Received illegal HTTP Method in request to /audit: %s\n", r.Method)
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Decode the audit record
	decoder := json.NewDecoder(r.Body)
	var auditRecord AuditRecord
	err := decoder.Decode(&auditRecord)
	if err != nil {
		logging.Error.Printf("Error decoding JSON: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Re-encode the JSON to prevent formatting differences
	inputJson, err := json.Marshal(auditRecord)
	if err != nil {
		logging.Error.Printf("Error re-encoding JSON: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hash := sha256.Sum256(inputJson)

	dataBuffer := new(bytes.Buffer)
	actorHash := sha256.Sum256([]byte(auditRecord.Header.Actor))
	intentHash := sha256.Sum256([]byte(auditRecord.Header.Intent))
	objectHash := sha256.Sum256([]byte(auditRecord.Header.Object))

	dataBuffer.Write(actorHash[:])
	dataBuffer.Write(intentHash[:])
	dataBuffer.Write(objectHash[:])
	dataBuffer.Write(hash[:])

	nonce, err := library.GetNonce()
	if err != nil {
		logging.Error.Printf("Error getting nonce: %s\n", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key, err := library.GetKey()
	if err != nil {
		_ = library.CancelNonce()
		logging.Error.Printf("Error getting signing key: %s\n", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addr := crypto.PubkeyToAddress(key.PublicKey)
	signer := types.NewEIP155Signer(big.NewInt(192001))
	tx, err := types.SignTx(types.NewTransaction(nonce, addr, big.NewInt(0), 1000000, big.NewInt(0), dataBuffer.Bytes()), signer, key)
	if err != nil {
		_ = library.CancelNonce()
		logging.Error.Printf("Error signing TX: %s\n", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawTx := bytes.NewBuffer(nil)
	err = tx.EncodeRLP(rawTx)
	if err != nil {
		_ = library.CancelNonce()
		logging.Error.Printf("Error encoding transaction: %s\n", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var auditResponse AuditResponse
	auditResponse.RecordHash = hex.EncodeToString(hash[:])
	auditResponse.TransactionHash = tx.Hash().Hex()

	js, err := json.Marshal(auditResponse)
	if err != nil {
		_ = library.CancelNonce()
		logging.Error.Printf("Error encoding response to JSON: %s\n", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", config.SendUrl(), rawTx)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		_ = library.CancelNonce()
		logging.Error.Printf("Error calling server: %s\n", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated { // created
		_ = library.CancelNonce()
		logging.Error.Printf("Unexpected response code from server: %d %s\n", resp.StatusCode, resp.Status)

		http.Error(w, fmt.Sprintf("Received error from server: %d %s", resp.StatusCode, resp.Status), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
