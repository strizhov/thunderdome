// Copyright (C) 2018

package thunderdomeserver

import (
	"context"
	crypto_rand "crypto/rand"
	"net"

	"github.com/spacemonkeygo/monotime"
)

func (s *Server) StartDownloadServer(ctx context.Context) (err error) {
	defer mon.Task()(&ctx)(&err)

	// start server
	listener, err := net.Listen("tcp", s.download_server_address)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Errore(err)
			continue
		}

		go s.handleDownloadConnection(ctx, conn)
	}
}

func (s *Server) handleDownloadConnection(ctx context.Context,
	conn net.Conn) (err error) {
	defer mon.Task()(&ctx)(&err)
	defer conn.Close()

	tcp, ok := conn.(*net.TCPConn)
	if !ok {
		return IncorrectProtocol.New("download server")
	}

	// set keep alives for connection
	tcp.SetKeepAlive(true)
	tcp.SetKeepAlivePeriod(*serverKeepAlivePeriod)

	// Set a deadline for reading.
	// Read operation will fail if no data sent
	tcp.SetWriteDeadline(monotime.Now().Add(*serverConnectionTimeout))

	// generate some random data
	buf := make([]byte, *downloadMaxTestDataSize)
	_, err = crypto_rand.Read(buf)
	if err != nil {
		logger.Errore(err)
		return err
	}

	// write data to the socket
	_, err = conn.Write(buf)
	if err != nil {
		logger.Errore(err)
		return err
	}

	return nil
}
