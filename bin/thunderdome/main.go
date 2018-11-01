// Copyright (C) 2018

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spacemonkeygo/errors"
	"github.com/spacemonkeygo/spacelog"
	monkit "gopkg.in/spacemonkeygo/monkit.v2"

	"github.com/strizhov/thunderdome/thunderdomeserver"
)

var (
	upload_server_address = flag.String(
		"upload_server_address",
		":9090",
		"upload server address to listen on")
	download_server_address = flag.String(
		"download_server_address",
		":9091",
		"download server address to listen on")

	Error     = errors.NewClass("thunderdome")
	flagError = Error.NewClass("flag", errors.NoCaptureStack())
	logger    = spacelog.GetLogger()
	mon       = monkit.Package()
)

func main() {
	// create the parent context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// set up signals to cancel the context
	sig := make(chan os.Signal, 10)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() { <-sig; cancel() }()

	// run Main command
	if err := Main(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}

func Main(ctx context.Context) (err error) {
	defer mon.Task()(&ctx)(&err)

	err = checkFlags()
	if err != nil {
		return err
	}

	server, err := thunderdomeserver.New(*upload_server_address, *download_server_address)
	if err != nil {
		return err
	}

	return server.Run(ctx)
}

func checkFlags() error {
	if *upload_server_address == "" {
		return flagError.New("upload server address is required")
	}
	if *download_server_address == "" {
		return flagError.New("download server address is required")
	}
	return nil
}
