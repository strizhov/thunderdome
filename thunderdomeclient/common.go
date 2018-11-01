// Copyright (C) 2018

package thunderdomeclient

import (
	"github.com/spacemonkeygo/errors"
	"github.com/spacemonkeygo/spacelog"
	"gopkg.in/spacemonkeygo/monkit.v2"
)

var (
	Error  = errors.NewClass("thunderdomeclient", errors.NoCaptureStack())
	logger = spacelog.GetLogger()
	mon    = monkit.Package()
)
