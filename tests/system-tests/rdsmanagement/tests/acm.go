package rdsmanagement_system_test

import (
	. "github.com/onsi/ginkgo/v2"

	"github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementcommon"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementparams"
)

var _ = Describe(
	"ACM Suite",
	Ordered,
	ContinueOnFailure,
	Label(rdsmanagementparams.Label), func() {

		It("Check if the namespace is exist",
			rdsmanagementcommon.VerifyACMNamespace)

		It("Verify that the ACM deployment is operational",
			rdsmanagementcommon.VerifyACMDeployment)

	})
