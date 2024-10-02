package rdsmanagement_system_test

import (
	. "github.com/onsi/ginkgo/v2"

	"github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementcommon"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementparams"
)

var _ = Describe(
	"IDM Suite",
	Ordered,
	ContinueOnFailure,
	Label(rdsmanagementparams.Label), func() {
		Context("Installed IDM server", Label("idm"), func() {
			It("Verifies that SSH login to IDM VM is working",
				rdsmanagementcommon.VerifyIDMInstallation)

			It("Verifies that SSH login to IDM web interface is successful",
				rdsmanagementcommon.VerifyIDMReplication)

			It("Verifies that new user accounts can be created",
				rdsmanagementcommon.VerifyOCPIntegrationWithIDM)
		})
	})
