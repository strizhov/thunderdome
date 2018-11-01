// Copyright (C) 2018

package thunderdomeserver

import (
	"github.com/spacemonkeygo/errors"
	"github.com/spacemonkeygo/spacelog"
	"gopkg.in/spacemonkeygo/monkit.v2"
)

var (
	Error             = errors.NewClass("thunderdomeserver", errors.NoCaptureStack())
	IncorrectProtocol = Error.NewClass("invalid type of socket")

	logger = spacelog.GetLogger()
	mon    = monkit.Package()
)

