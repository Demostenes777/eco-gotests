package tests

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/slcm/internal/slcmcommon"
)

var _ = Describe(
	"Basic sanity tests for DU deployment",
	Ordered,
	ContinueOnFailure,
	Label("du"), func() {
		Context("Deployment and Basic Sanity", Label("deployment", "test_basic_sanity"), func() {
			It("Verifies the number of pods in the specified namespace matches the expected number",
				slcmcommon.TestDUPodsCount)
			It("Verifies all pods in the specified namespace are in the 'Running' state and have the "+
				"correct condition status", slcmcommon.TestDUPodsStatus)
		})
	})
