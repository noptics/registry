package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/noptics/registry/data"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/noptics/golog"
)

func main() {
	var l golog.Logger
	if len(os.Getenv("DEBUG")) != 0 {
		l = golog.StdOut(golog.LEVEL_DEBUG)
	} else {
		l = golog.StdOut(golog.LEVEL_ERROR)
	}

	l.Init()
	defer l.Finish()

	c := newContext()

	// Setup the Database Connection
	dbendpoint := os.Getenv("DB_ENDPOINT")
	sess, err := session.NewSession()
	if err != nil {
		l.Info("Could not create AWS session...exiting")
		os.Exit(1)
	}

	tablePrefix := os.Getenv("DB_TABLE_PREFIX")
	if len(tablePrefix) == 0 {
		l.Info("must provide DB_TABLE_PREFIX")
		os.Exit(1)
	}

	dataStore := data.NewDynamodb(sess, dbendpoint, tablePrefix, l)

	// Start the GRPC Server
	c.GRPCPort = os.Getenv("GRPC_PORT")
	if c.GRPCPort == "" {
		c.GRPCPort = "7775"
	}

	errChan := make(chan error)

	gs, err := NewGRPCServer(dataStore, c.GRPCPort, errChan, l)
	if err != nil {
		l.Infow("unable to start grpc server", "error", err.Error())
		os.Exit(1)
	}

	// start the rest server
	c.RESTPort = os.Getenv("REST_PORT")
	if c.RESTPort == "" {
		c.RESTPort = "7776"
	}

	// We don't support specific address binding right now...
	c.Host = "0.0.0.0"

	rs := NewRestServer(dataStore, c.RESTPort, errChan, l, c)

	l.Info("started")

	// go until told to stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-sigs:
	case <-errChan:
		l.Infow("error", "error", err.Error())
	}

	l.Info("shutting down")
	gs.Stop()

	err = rs.Stop()
	if err != nil {
		l.Infow("error shutting down rest server", "error", err.Error())
	}

	l.Info("finished")
}
