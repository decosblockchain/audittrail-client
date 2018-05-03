package routes

import (
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<html><head><title>Decos Blockchain Audittrail Client</title></head><body><h1>It works!</h1><p>Decos Blockchain Audittrail Client is active.</p></body></html>"))
}
