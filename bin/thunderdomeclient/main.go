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

	"github.com/strizhov/thunderdome/thunderdomeclient"
)

var (
	uploadServerAddress = flag.String(
		"upload_server_address",
		"localhost:9090",
		"upload server address:port")
	downloadServerAddress = flag.String(
		"download_server_address",
		"localhost:9091",
		"download server address:port")

	Error     = errors.NewClass("bwclient")
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

	// make client
	client := thunderdomeclient.New(*uploadServerAddress, *downloadServerAddress)

	// get upload bw
	upload_bw, err := client.GetUploadBandwidth(ctx)
	if err != nil {
		logger.Errore(err)
		return err
	}

	// get download bw
	download_bw, err := client.GetDownloadBandwidth(ctx)
	if err != nil {
		logger.Errore(err)
		return err
	}

	fmt.Println("Upload BW (kbps):", upload_bw, "Download BW (kbps): ", download_bw)
	return nil
}
