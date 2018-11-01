// Copyright (C) 2018

package thunderdomeserver

import (
	"context"
	"flag"
	"net"
	"sync"
	"time"
)

var (
	serverKeepAlivePeriod = flag.Duration(
		"thunderdomeserver.server_keep_alive_period",
		30*time.Second,
		"duration to set for tcp keep alive period on uptime sockets")
	serverConnectionTimeout = flag.Duration(
		"thundersdomeserver.server_connection_timeout",
		120*time.Second,
		"server time to wait for data from the connection")
	uploadMaxBufferSize = flag.Int(
		"thunderdomeserver.upload_max_buffer_size",
		40960000, // 40.96 Mbytes
		"size of the max upload allowed")
	uploadMaxBufferSizePerSocketRead = flag.Int(
		"thunderdomeserver.upload_max_buffer_size_per_socket_read",
		4096,
		"size of buffer to read from socket")
	downloadMaxTestDataSize = flag.Int(
		"thunderdomeserver.download_max_test_data_size",
		10240000, // 10.24 MBytes
		"size of the file used in download server")
)

type Server struct {
	upload_server_address   string
	download_server_address string
}

func New(upload_server_address, download_server_address string) (
	s *Server, err error) {

	// resolve first address
	_, err = net.ResolveTCPAddr("tcp", upload_server_address)
	if err != nil {
		return nil, err
	}

	// resolve second address
	_, err = net.ResolveTCPAddr("tcp", download_server_address)
	if err != nil {
		return nil, err
	}

	return &Server{
		upload_server_address:   upload_server_address,
		download_server_address: download_server_address,
	}, nil
}

func (s *Server) Run(ctx context.Context) (err error) {
	defer mon.Task()(&ctx)(&err)

	var wg sync.WaitGroup
	defer wg.Done()
	errors := make(chan error, 2)

	// start upload server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Infof("starting up upload server on %s",
			s.upload_server_address)
		errors <- s.StartUploadServer(ctx)
	}()

	// start download server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Infof("starting up download server on %s",
			s.download_server_address)
		errors <- s.StartDownloadServer(ctx)
	}()

	select {
	case err = <-errors:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
