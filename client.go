package main

import (
	"log"
	"net/http"
	"os"

	"github.com/decosblockchain/audittrail-client/library"
	"github.com/decosblockchain/audittrail-client/logging"
	"github.com/decosblockchain/audittrail-client/routes"

	"github.com/gorilla/mux"
	"github.com/kardianos/service"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	r := mux.NewRouter()
	r.HandleFunc("/audit", routes.AuditHandler)

	log.Fatal(http.ListenAndServe(":8000", r))
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	logging.Init(os.Stdout, os.Stdout, os.Stdout, os.Stdout)

	address, _ := library.GetAddress()
	logging.Info.Printf("My address is %s", address)

	svcConfig := &service.Config{
		Name:        "DecosBlockchainAuditConnector",
		DisplayName: "Decos Blockchain Audit Connector",
		Description: "This service acts as a signing & sending proxy for the Decos Blockchain Audit Service",
	}

	prg := &program{}

	if len(os.Args) > 1 && os.Args[1] == "console" {
		prg.run()
	} else {
		s, err := service.New(prg, svcConfig)
		if err != nil {
			log.Fatal(err)
		}
		logger, err = s.Logger(nil)
		if err != nil {
			log.Fatal(err)
		}
		err = s.Run()
		if err != nil {
			logger.Error(err)
		}
	}
}
