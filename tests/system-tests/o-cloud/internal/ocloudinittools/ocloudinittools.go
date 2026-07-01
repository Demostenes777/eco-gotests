package ocloudinittools

import (
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/bmc"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/clients"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/internal/inittools"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/internal/platform"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/o-cloud/internal/ocloudconfig"
	"k8s.io/klog/v2"
)

var (
	// HubAPIClient provides API access to hub cluster.
	HubAPIClient *clients.Settings
	// OCloudConfig provides access to O-Cloud system tests configuration parameters.
	OCloudConfig *ocloudconfig.OCloudConfig
	// BMCClient provides API access to BMC.
	BMCClient *bmc.BMC
)

// init loads all variables automatically when this package is imported. Once package is imported a user has full
// access to all vars within init function. It is recommended to import this package using dot import.
func init() {
	OCloudConfig = ocloudconfig.NewOCloudConfig()
	HubAPIClient = inittools.APIClient

	if HubAPIClient != nil {
		hubOCPVersion, err := platform.GetOCPVersion(HubAPIClient)
		if err != nil {
			klog.Warningf("Failed to retrieve hub OCP version: %v", err)
		} else {
			OCloudConfig.HubOCPVersion = hubOCPVersion
			klog.V(90).Infof("Detected hub OCP version: %s", hubOCPVersion)
		}
	}
}
