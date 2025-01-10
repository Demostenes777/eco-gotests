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
		Context("Configured hub cluster", Label("ocloud-day2-configuration"), func() {
			It("Successful operator upgrade",
				ocloudcommon.SuccessfulOperatorUpgrade)
			
			It("Failed operator upgrade in all the SNOs",
				ocloudcommon.FailedOperatorUpgrade)

			It("Failed operator upgrade in a subset of the SNOs",
				ocloudcommon.FailedPartialOperatorUpgrade)
		})
	})