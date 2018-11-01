// Copyright (C) 2018

package thunderdome

import (
	"context"
)

type Client interface {
	GetUploadBandwidth(ctx context.Context) (bw int64, err error)
	GetDownloadBandwidth(ctx context.Context) (bw int64, err error)
}

type Server interface {
	Run(ctx context.Context) (err error)
}
