package ocloudparams

import (
	"time"
)

const (
	// Label represents O-Cloud system tests label that can be used for test cases selection.
	Label = "ocloud"
	// DefaultTimeout is the timeout used for test resources creation.
	DefaultTimeout = 900 * time.Second
	// OCloudLogLevel configures logging level for O-Cloud related tests.
	OCloudLogLevel = 90
)
