// Copyright (C) 2018

package thunderdomeclient

import (
	"context"
	"flag"

	"sm/tools/thunderdome"
)

var (
	concurrentConnectionNumber = flag.Int(
		"thunderdomeclient.concurrent_connection_number",
		3,
		"number of concurrent connections to the server")
)

type Client struct {
	upload_server_address   string
	download_server_address string
}

var _ thunderdome.Client = (*Client)(nil)

func New(upload_server_address, download_server_address string) *Client {
	return &Client{
		upload_server_address:   upload_server_address,
		download_server_address: download_server_address,
	}
}

func (c *Client) GetDownloadBandwidth(ctx context.Context) (bw int64, err error) {
	defer mon.Task()(&ctx)(&err)

	results := make(chan int64, *concurrentConnectionNumber)
	for i := 0; i < *concurrentConnectionNumber; i++ {
		go c.getDownloadBandwidth(ctx, results)
	}

	var download_bw int64
	for i := 0; i < *concurrentConnectionNumber; i++ {
		select {
		case res := <-results:
			download_bw += res
		case <-ctx.Done():
			return download_bw, ctx.Err()
		}
	}

	return download_bw, nil
}

func (c *Client) GetUploadBandwidth(ctx context.Context) (bw int64, err error) {
	defer mon.Task()(&ctx)(&err)

	results := make(chan int64, *concurrentConnectionNumber)
	for i := 0; i < *concurrentConnectionNumber; i++ {
		go c.getUploadBandwidth(ctx, results)
	}

	var upload_bw int64
	for i := 0; i < *concurrentConnectionNumber; i++ {
		select {
		case res := <-results:
			upload_bw += res
		case <-ctx.Done():
			return upload_bw, ctx.Err()
		}
	}

	return upload_bw, nil
}
