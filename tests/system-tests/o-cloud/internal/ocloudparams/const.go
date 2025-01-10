package ocloudparams

import (
	"time"
)

const (
	// Label represents O-Cloud system tests label that can be used for test cases selection.
	Label = "ocloud"
	// LabelAiSnoProvisioning a label to select tests for Assisted Installer provisioning validation.
	LabelAiSnoProvisioning = "ocloud-ai-provisioning"
	// LabelIbiSnoProvisioning a label to select tests for Image Based Installer provisioning validation.
	LabelIbiSnoProvisioning = "ocloud-ibi-provisioning"
	// LabelDay2Configuration a label to select tests for Day 2 configuration validation.
	LabelDay2Configuration = "ocloud-day2-configuration"
	// DefaultTimeout is the timeout used for test resources creation.
	DefaultTimeout = 900 * time.Second
	// OCloudLogLevel configures logging level for O-Cloud related tests.
	OCloudLogLevel = 90
)
