package slcminittools

import (
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-gotests/tests/internal/inittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/slcm/internal/slcmconfig"
)

var (
	// APIClient provides API access to cluster.
	APIClient *clients.Settings
	// SLCMConfig provides access to SLCM tests configuration parameters.
	SLCMConfig *slcmconfig.SLCMConfig
)

// init loads all variables automatically when this package is imported. Once package is imported a user has full
// access to all vars within init function. It is recommended to import this package using dot import.
func init() {
	SLCMConfig = slcmconfig.NewSLCMConfig()
	APIClient = inittools.APIClient
}
