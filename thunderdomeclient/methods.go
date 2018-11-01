// Copyright (C) 2018

package thunderdomeclient

import (
	"context"
	crypto_rand "crypto/rand"
	"flag"
	"io"
	"net"
	"time"

	"github.com/spacemonkeygo/monotime"
)

var (
	serverDialTimeout = flag.Duration(
		"thunderdomeclient.server_dial_timeout",
		60*time.Second,
		"time to dial to the server")
	serverConnectionTimeout = flag.Duration(
		"thunderdomeclient.server_connection_timeout",
		120*time.Second,
		"client time out for write to the server socket")
	uploadTestDataSize = flag.Int(
		"thunderdomeclient.upload_test_data_size",
		1024*1024, // 1Mib
		"size of test data we will send to server")
	downloadMaxBufferSize = flag.Int(
		"thunderdomeclient.download_max_buffer_size",
		1024*1024, // 1Mib
		"size of the max download allowed")
	downloadMaxBufferSizePerSocketRead = flag.Int(
		"thunderdomeclient."+
			"download_max_buffer_size_per_socket_read",
		4096,
		"size of buffer to read from socket")
)

func (c *Client) getDownloadBandwidth(ctx context.Context,
	results chan int64) (err error) {
	defer mon.Task()(&ctx)(&err)

	// connect to the server
	conn, err := net.DialTimeout("tcp", c.download_server_address, *serverDialTimeout)
	if err != nil {
		logger.Errore(err)
		results <- 0
		return err
	}
	defer conn.Close()

	// Set a deadline for writing to socket
	conn.SetReadDeadline(time.Now().Add(*serverConnectionTimeout))

	// allocate buffers
	buf := make([]byte, *downloadMaxBufferSizePerSocketRead)

	var bytes_read int
	start_time := monotime.Now()
	for {
		n, err := conn.Read(buf)
		bytes_read += n
		if err != nil {
			// record all errors except EOF (end of transfer)
			if err != io.EOF {
				logger.Errore(err)
			}
			break
		}

		// buffer overflow check
		if bytes_read > *downloadMaxBufferSize {
			break
		}
	}

	// get delta in float64 seconds and measure bandwidth
	delta := monotime.Now().Sub(start_time).Seconds()
	bw := (float64)(bytes_read) / delta

	// convert bw in kbps
	bw_kbps := (int64)(8 * bw / 1000)

	// put result in the channel
	results <- bw_kbps
	return nil
}

func (c *Client) getUploadBandwidth(ctx context.Context,
	results chan int64) (err error) {
	defer mon.Task()(&ctx)(&err)

	// connect to the server
	conn, err := net.DialTimeout("tcp", c.upload_server_address, *serverDialTimeout)
	if err != nil {
		logger.Errore(err)
		results <- 0
		return err
	}
	defer conn.Close()

	// Set a deadline for writing to socket
	conn.SetWriteDeadline(time.Now().Add(*serverConnectionTimeout))

	// generate some random data
	buf := make([]byte, *uploadTestDataSize)
	_, err = crypto_rand.Read(buf)
	if err != nil {
		logger.Errore(err)
		results <- 0
		return err
	}

	// write data to the socket
	start_time := monotime.Now()
	_, err = conn.Write(buf)
	if err != nil {
		logger.Errore(err)
		results <- 0
		return err
	}

	// get delta in float64 seconds and measure bandwidth
	delta := monotime.Now().Sub(start_time).Seconds()
	bw := (float64)(*uploadTestDataSize) / delta

	// convert bw in kbps
	bw_kbps := (int64)(8 * bw / 1000)

	// put result in the channel
	results <- bw_kbps
	return nil
}
