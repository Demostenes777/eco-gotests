package o_cloud_system_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/reportxml"
	. "github.com/rh-ecosystem-edge/eco-gotests/tests/internal/inittools"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/internal/reporter"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
	_ "github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/o-cloud/tests"
	subscriber "github.com/rh-ecosystem-edge/eco-gotests/tests/internal/oran-subscriber"
	. "github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"

)

var (
	_, currentFile, _, _ = runtime.Caller(0)
)

func TestOCloud(t *testing.T) {
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = GeneralConfig.GetJunitReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "O-Cloud SystemTests Suite", Label(ocloudparams.Labels...), reporterConfig)
}

var _ = BeforeSuite(func() {
	By("deploying the subscriber for alarm notifications")
	subscriberDomain := "oran-subscriber.apps.hub03.oran.telcoqe.eng.rdu2.dc.redhat.com"
	err := subscriber.Deploy(HubAPIClient, "oran-subscriber", subscriberDomain, "")
	Expect(err).ToNot(HaveOccurred(), "Failed to deploy subscriber")
})

var _ = AfterSuite(func() {
	err := os.RemoveAll("tmp/")
	if err != nil {
		glog.V(ocloudparams.OCloudLogLevel).Infof("removed tmp/")
	}

	By("cleaning up the subscriber deployment")
	err = subscriber.Cleanup(HubAPIClient, "oran-subscriber")
	Expect(err).ToNot(HaveOccurred(), "Failed to cleanup subscriber")
})

var _ = JustAfterEach(func() {
	reporter.ReportIfFailed(
		CurrentSpecReport(), currentFile, ocloudparams.ReporterNamespacesToDump, ocloudparams.ReporterCRDsToDump)
})

var _ = ReportAfterSuite("", func(report Report) {
	reportxml.Create(
		report, GeneralConfig.GetReportPath(), GeneralConfig.TCPrefix)
})
