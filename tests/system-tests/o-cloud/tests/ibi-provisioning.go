package o_cloud_system_test

import (
	. "github.com/onsi/ginkgo/v2"

	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudcommon"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
)

var _ = Describe(
	"Image Based Install based SNO provisioning Test Suite",
	Ordered,
	ContinueOnFailure,
	Label(ocloudparams.Label), func() {
		Context("Configured hub cluster", Label("ocloud-ibi-provisioning"), func() {
			It("Verifies the successful provisioning of a single SNO cluster",
				ocloudcommon.VerifySuccessfulIbiSnoProvisioning)

			It("Verifies the failed provisioning of a single SNO cluster",
				ocloudcommon.VerifyFailedIbiSnoProvisioning)
		})
	})