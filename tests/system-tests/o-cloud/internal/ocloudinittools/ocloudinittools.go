package ocloudinittools

import (
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-gotests/tests/internal/inittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudconfig"
)

var (
	// APIClient provides API access to cluster.
	APIClient *clients.Settings
	// OCLoudTestConfig provides access to O-Cloud system tests configuration parameters.
	OCloudTestConfig *ocloudconfig.OCloudConfig
)

// init loads all variables automatically when this package is imported. Once package is imported a user has full
// access to all vars within init function. It is recommended to import this package using dot import.
func init() {
	OCloudTestConfig = ocloudconfig.NewOCloudConfig()
	APIClient = inittools.APIClient
}
