package o_cloud_system_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	. "github.com/openshift-kni/eco-gotests/tests/internal/inittools"
	"github.com/openshift-kni/eco-gotests/tests/internal/reporter"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
	_ "github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/tests"
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

var _ = AfterSuite(func() {
	err := os.RemoveAll("tmp/")
	if err != nil {
		glog.V(ocloudparams.OCloudLogLevel).Infof("removed tmp/")
	}

})

var _ = JustAfterEach(func() {
	reporter.ReportIfFailed(
		CurrentSpecReport(), currentFile, ocloudparams.ReporterNamespacesToDump, ocloudparams.ReporterCRDsToDump)
})

var _ = ReportAfterSuite("", func(report Report) {
	reportxml.Create(
		report, GeneralConfig.GetReportPath(), GeneralConfig.TCPrefix)
})
