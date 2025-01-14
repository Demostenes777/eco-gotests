package ocloudcommon

import (
	"time"
	//"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	//. "github.com/onsi/gomega"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"
	//"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/namespace"
	"github.com/openshift-kni/eco-goinfra/pkg/olm"
	"github.com/openshift-kni/eco-goinfra/pkg/pod"
	// "github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	// "github.com/openshift-kni/eco-gotests/tests/assisted/ztp/internal/meets"
	// . "github.com/openshift-kni/eco-gotests/tests/assisted/ztp/internal/ztpinittools"
	// "github.com/openshift-kni/eco-gotests/tests/assisted/ztp/operator/internal/tsparams"
)

const (
	acmNamespace  = "rhacm"
	acmCSVPattern = "advanced-cluster-management"

	gitOpsNamespace = "openshift-gitops"
	gitOpsCSVPattern = "openshift-gitops-operator"
)

func VerifyACM() {
	By("Checking that rhacm namespace exists")
	_, err := namespace.Pull(HubAPIClient, acmNamespace)
	if err != nil {
		Skip("Advanced Cluster Management is not installed")
	}

	By("Getting clusterserviceversion")
	clusterServiceVersions, err := olm.ListClusterServiceVersionWithNamePattern(
		HubAPIClient, acmCSVPattern, acmNamespace)
	Expect(err).NotTo(HaveOccurred(), "error listing clusterserviceversions by name pattern")
	Expect(len(clusterServiceVersions)).To(Equal(1), "error did not receieve expected list of clusterserviceversions")

	success, err := clusterServiceVersions[0].IsSuccessful()
	Expect(err).NotTo(HaveOccurred(), "error checking clusterserviceversions phase")
	Expect(success).To(BeTrue(), "error advanced-cluster-management clusterserviceversion is not Succeeded")

	By("Check that pods exist in rhacm namespace")
	rhacmPods, err := pod.List(HubAPIClient, acmNamespace)
	Expect(err).NotTo(HaveOccurred(), "error occurred while listing pods in rhacm namespace")
	Expect(len(rhacmPods) > 0).To(BeTrue(), "error: did not find any pods in the hive namespace")

	By("Check that rhacm pods are running correctly")
	running, err := pod.WaitForAllPodsInNamespaceRunning(HubAPIClient, acmNamespace, time.Minute)
	Expect(err).NotTo(HaveOccurred(), "error occurred while waiting for rhacm pods to be in Running state")
	Expect(running).To(BeTrue(), "some rhacm pods are not in Running state")
}

func VerifyGitOps() {
	By("Checking that openshift-gitops namespace exists")
	_, err := namespace.Pull(HubAPIClient, gitOpsNamespace)
	if err != nil {
		Skip("The OpenShift GitOps operator is not installed")
	}

	By("Getting clusterserviceversion")
	// todo - ListClusterServiceVersionWithNamePattern is not filtering by namespace
	clusterServiceVersions, err := olm.ListClusterServiceVersionWithNamePattern(
		HubAPIClient, gitOpsCSVPattern, gitOpsNamespace)
	Expect(err).NotTo(HaveOccurred(), "error listing clusterserviceversions by name pattern")
	Expect(len(clusterServiceVersions)).To(BeNumerically(">", 1), "error did not receieve expected list of clusterserviceversions")

	success, err := clusterServiceVersions[0].IsSuccessful()
	Expect(err).NotTo(HaveOccurred(), "error checking clusterserviceversions phase")
	Expect(success).To(BeTrue(), "error openshift-gitops-operator clusterserviceversion is not Succeeded")

	By("Check that pods exist in openshift-gitops namespace")
	gitOpsPods, err := pod.List(HubAPIClient, gitOpsNamespace)
	Expect(err).NotTo(HaveOccurred(), "error occurred while listing pods in openshift-gitops namespace")
	Expect(len(gitOpsPods) > 0).To(BeTrue(), "error: did not find any pods in the openshift-gitops namespace")

	By("Check that openshift-gitops pods are running correctly")
	running, err := pod.WaitForAllPodsInNamespaceRunning(HubAPIClient, gitOpsNamespace, time.Minute)
	Expect(err).NotTo(HaveOccurred(), "error occurred while waiting for openshift-gitops pods to be in Running state")
	Expect(running).To(BeTrue(), "some openshift-gitops pods are not in Running state")
}

func VerifySiteConfigOperator() {

}
	
func VerifyOCloudManagerOperator() {

}
	
func VerifyHardwareManagerPluginOperator() {

}
	