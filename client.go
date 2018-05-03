package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/decosblockchain/audittrail-client/config"
	"github.com/decosblockchain/audittrail-client/library"
	"github.com/decosblockchain/audittrail-client/logging"
	"github.com/decosblockchain/audittrail-client/routes"

	"github.com/gorilla/mux"
	"github.com/kardianos/service"
	"github.com/regorov/logwriter"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {

	config.Init()

	if len(os.Args) > 1 && os.Args[1] == "console" {
		logging.Info.Println("Running in console mode")
	}

	address, err := library.GetAddress()
	if err != nil {
		logging.Error.Printf("Error getting address: %s\n", err.Error())
		os.Exit(4)
	}
	logging.Info.Printf("My address is %s\n", address)

	r := mux.NewRouter()
	r.HandleFunc("/audit", routes.AuditHandler)
	r.HandleFunc("/", routes.HomeHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", config.ListenPort()), r))
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	// Create data, log directories if they do not exist
	config.EnsurePathsExist()

	cfg := &logwriter.Config{
		BufferSize:       0,                 // no buffering
		FreezeInterval:   24 * time.Hour,    // freeze log file every hour
		HotMaxSize:       10 * logwriter.MB, // 10 MB max file size
		CompressColdFile: false,             // compress cold file
		HotPath:          path.Join(config.BaseDir(), "log"),
		ColdPath:         path.Join(config.BaseDir(), "log", "archive"),
		Mode:             logwriter.ProductionMode, // write to file only
	}

	if len(os.Args) > 1 && os.Args[1] == "console" {
		cfg.Mode = logwriter.DebugMode // Write to file and console
	}

	lw, err := logwriter.NewLogWriter("audit-client",
		cfg,
		true, // freeze hot file if exists
		nil)

	if err != nil {
		panic(err)
	}

	logging.Init(lw, lw, lw, lw)

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
			logging.Error.Println(err)
			os.Exit(2)
		}
		if len(os.Args) > 1 {
			err = service.Control(s, os.Args[1])
			if err != nil {
				logging.Error.Println(err)
				os.Exit(3)
			}
			log.Printf("Service control action [%s] executed succesfully\n", os.Args[1])
			os.Exit(0)
		}
		err = s.Run()
		if err != nil {
			logging.Error.Println(err)
		}
	}
}
