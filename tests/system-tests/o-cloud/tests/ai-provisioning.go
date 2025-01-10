package o_cloud_system_test

import (
	. "github.com/onsi/ginkgo/v2"

	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudcommon"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
)

var _ = Describe(
	"Assisted Installer based SNO provisioning Test Suite",
	Ordered,
	ContinueOnFailure,
	Label(ocloudparams.Label), func() {
		Context("Configured hub cluster", Label("ocloud-ai-provisioning"), func() {
			It("Verifies the successful provisioning of a single SNO cluster",
				ocloudcommon.VerifySuccessfulSingleAiSnoProvisioning)
			
			It("Verifies the successful deprovisioning of a single SNO cluster",
				ocloudcommon.VerifySuccessfulSingleAiSnoDeprovisioning)

			It("Verifies the failed provisioning of a single SNO cluster",
				ocloudcommon.VerifyFailedSingleAiSnoProvisioning)

			It("Verifies the succesful E2E simultaneous provisioning of SNO clusters with the same cluster template",
				ocloudcommon.VerifySimultaneousSnoProvisioningSameClusterTemplate)

			It("Verifies the succesful E2E simultaneous deprovisioning of SNO clusters with the same cluster template",
				ocloudcommon.VerifySimultaneousSnoDeprovisioningSameClusterTemplate)

			It("Verifies the succesful E2E simultaneous provisioning of SNO clusters with different cluster templates",
				ocloudcommon.VerifySimultaneousSnoProvisioningDifferentClusterTemplate)

			It("Verifies the succesful E2E simultaneous deprovisioning of SNO clusters with different cluster templates",
				ocloudcommon.VerifySimultaneousSnoDeprovisioningDifferentClusterTemplate)
		})
	})