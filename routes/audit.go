package routes

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/decosblockchain/audittrail-client/config"
	"github.com/decosblockchain/audittrail-client/library"
	"github.com/decosblockchain/audittrail-client/logging"
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

	tx, err := library.CreateTransaction(dataBuffer.Bytes())
	library.SignTransaction(tx)

	txJSON, err := json.Marshal(tx)
	logging.Info.Printf("JSON:\n%s", string(txJSON))

	var auditResponse AuditResponse
	auditResponse.RecordHash = hex.EncodeToString(hash[:])
	//auditResponse.TransactionHash = tx.Hash().Hex()

	js, err := json.Marshal(auditResponse)
	if err != nil {
		logging.Error.Printf("Error encoding response to JSON: %s\n", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", config.SendUrl(), bytes.NewReader(txJSON))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.Error.Printf("Error calling server: %s\n", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated { // created
		logging.Error.Printf("Unexpected response code from server: %d %s\n", resp.StatusCode, resp.Status)

		http.Error(w, fmt.Sprintf("Received error from server: %d %s", resp.StatusCode, resp.Status), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
