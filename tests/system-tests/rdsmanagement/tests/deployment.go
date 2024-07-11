package rdsmanagement_system_test

import (
	. "github.com/onsi/ginkgo/v2"

	"github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementcommon"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementparams"
)

var _ = Describe(
	"Management Basic Deployment Suite",
	Ordered,
	ContinueOnFailure,
	Label(rdsmanagementparams.Label), func() {

		It("Check system reserved memory for master nodes",
			rdsmanagementcommon.SetSystemReservedMemoryForMasterNodes)

		It("Checking clusterversion",
			rdsmanagementcommon.CheckingClusterversion)

		It("Verifies that all node are ready",
			rdsmanagementcommon.WaitAllNodesAreReady)

		It("Verify that the cluster is operational",
			rdsmanagementcommon.VerifyClusterIsOperational)

	})
