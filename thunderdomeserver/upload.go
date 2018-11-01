// Copyright (C) 2018

package thunderdomeserver

import (
	"context"
	"io"
	"net"

	"github.com/spacemonkeygo/monotime"
)

func (s *Server) StartUploadServer(ctx context.Context) (err error) {
	defer mon.Task()(&ctx)(&err)

	// start server
	listener, err := net.Listen("tcp", s.upload_server_address)
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

		go s.handleUploadConnection(ctx, conn)
	}
}

func (s *Server) handleUploadConnection(ctx context.Context,
	conn net.Conn) (err error) {
	defer mon.Task()(&ctx)(&err)
	defer conn.Close()

	tcp, ok := conn.(*net.TCPConn)
	if !ok {
		return IncorrectProtocol.New("upload server")
	}

	// set keep alives for connection
	tcp.SetKeepAlive(true)
	tcp.SetKeepAlivePeriod(*serverKeepAlivePeriod)

	// Set a deadline for reading.
	// Read operation will fail if no data sent
	tcp.SetReadDeadline(monotime.Now().Add(*serverConnectionTimeout))

	// allocate buffers
	buf := make([]byte, 0, *uploadMaxBufferSize)
	tmp := make([]byte, *uploadMaxBufferSizePerSocketRead)

	for {
		// read into tmp buffer
		n, err := conn.Read(tmp)
		if err != nil {
			// record all errors except EOF (end of transfer)
			if err != io.EOF {
				logger.Errore(err)
			}
			break
		}

		// buffer overflow check
		if (len(buf) + n) > *uploadMaxBufferSize {
			break
		}
		buf = append(buf, tmp[:n]...)
	}

	return nil
}
