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
				rdsmanagementcommon.VerifySSHAccess)

			It("Verifies that SSH login to IDM web interface is successful",
				rdsmanagementcommon.VerifyWebAccess)

			It("Verifies that new user accounts can be created",
				rdsmanagementcommon.VerifyNewUserAccountCreation)

			It("Verifies that new groups can be created",
				rdsmanagementcommon.VerifyNewGroupCreation)
		})
	})
