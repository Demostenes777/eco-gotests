package o_cloud_system_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	subscriber "github.com/rh-ecosystem-edge/eco-gotests/tests/internal/oran-subscriber"
	. "github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/o-cloud/internal/ocloudcommon"
	. "github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
)

var _ = Describe(
	"ORAN Alarms Tests", Ordered, ContinueOnFailure, Label(ocloudparams.Label, "ocloud-alarms"), func() {

		BeforeEach(func() {
			By("deploying the subscriber for alarm notifications")
			subscriberDomain := "oran-subscriber.apps.hub03.oran.telcoqe.eng.rdu2.dc.redhat.com"
			err := subscriber.Deploy(HubAPIClient, "oran-subscriber", subscriberDomain, "")
			Expect(err).ToNot(HaveOccurred(), "Failed to deploy subscriber")
		})

		AfterEach(func() {
			By("cleaning up the subscriber deployment")
			err := subscriber.Cleanup(HubAPIClient, "oran-subscriber")
			Expect(err).ToNot(HaveOccurred(), "Failed to cleanup subscriber")
		})

		It("Successful alarm retrieval from the API",
			VerifySuccessfulAlarmRetrieval)

		It("Successful alarms cleanup from the database after the retention period",
			VerifySuccessfulAlarmsCleanup)

	})
