package ocloudcommon

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/olm"
	"github.com/openshift-kni/eco-goinfra/pkg/pod"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"

	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/apiobjectshelper"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/csv"
)

func VerifyACM() {
	verifyNamespaceExists(HubAPIClient, ocloudparams.AcmNamespace)
	verifyCsvSuccessful(HubAPIClient, ocloudparams.AcmSubscriptionName, ocloudparams.AcmNamespace)
	verifyPodsRunning(HubAPIClient, ocloudparams.AcmNamespace)
}

func VerifyGitOpsOperator() {
	verifyNamespaceExists(HubAPIClient, ocloudparams.OpenshiftGitOpsNamespace)
	verifyCsvSuccessful(HubAPIClient, 
		ocloudparams.OpenshiftGitOpsSubscriptionName, 
		ocloudparams.OpenshiftGitOpsNamespace)
}

func VerifySiteConfigOperator() {
	// todo - need to find how to verify it
}

func VerifyOranO2ImsOperator() {
	verifyNamespaceExists(HubAPIClient, ocloudparams.OCloudO2ImsNamespace)
	verifyCsvSuccessful(HubAPIClient, ocloudparams.OCloudO2ImsSubscriptionName,	ocloudparams.OCloudO2ImsNamespace)
}

func VerifyOranHardwareManagerPluginOperator() {
	verifyNamespaceExists(HubAPIClient, ocloudparams.OCloudHardwareManagerPluginNamespace)
	verifyCsvSuccessful(HubAPIClient,
		ocloudparams.OCloudHardwareManagerPluginSubscriptionName,
		ocloudparams.OCloudHardwareManagerPluginNamespace)
}

func verifyNamespaceExists(apiClient *clients.Settings, nsname string) {
	By(fmt.Sprintf("Check that %s namespace exists", nsname))
	err := apiobjectshelper.VerifyNamespaceExists(apiClient, nsname, time.Second)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to pull namespace %q; %v", nsname, err))
}

func verifyCsvSuccessful(apiClient *clients.Settings, subscriptionName string, nsname string) {
	By(fmt.Sprintf("Check that csv %s is successful", subscriptionName))
	
	csvName, err := csv.GetCurrentCSVNameFromSubscription(apiClient, subscriptionName, nsname)
	if err != nil {
		Skip(fmt.Sprintf("csv %s not found in namespace %s", csvName, nsname))
	}
	
	csvObj, err := olm.PullClusterServiceVersion(apiClient, csvName, nsname)
	if err != nil {
		Skip(fmt.Sprintf("failed to pull %q csv from the %s namespace", csvName, nsname))
	}

	_, err = csvObj.IsSuccessful()
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("failed to verify csv %s in the namespace %s status", csvName, nsname))		
}

func verifyPodsRunning(apiClient *clients.Settings, nsname string) {
	By(fmt.Sprintf("Check that pods exist in %s namespace", nsname))
	rhacmPods, err := pod.List(apiClient, nsname)
	Expect(err).NotTo(HaveOccurred(), "error nsname while listing pods in rhacm namespace")
	Expect(len(rhacmPods) > 0).To(BeTrue(), "error: did not find any pods in the rhacm namespace")

	By(fmt.Sprintf("Check that pods are running correctly in %s namespace", nsname))
	running, err := pod.WaitForAllPodsInNamespaceRunning(apiClient, nsname, time.Minute)
	Expect(err).NotTo(HaveOccurred(), 
		fmt.Sprintf("error occurred while waiting for %s pods to be in Running state", nsname))
	Expect(running).To(BeTrue(), 
		fmt.Sprintf("some %s pods are not in Running state", nsname))
}