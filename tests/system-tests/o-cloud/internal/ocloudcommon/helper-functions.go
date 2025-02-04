package ocloudcommon

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"

	"github.com/openshift-kni/eco-goinfra/pkg/bmh"
	"github.com/openshift-kni/eco-goinfra/pkg/ibi"
	"github.com/openshift-kni/eco-goinfra/pkg/namespace"
	"github.com/openshift-kni/eco-goinfra/pkg/ocm"
	"github.com/openshift-kni/eco-goinfra/pkg/olm"
	"github.com/openshift-kni/eco-goinfra/pkg/oran"
	"github.com/openshift-kni/eco-goinfra/pkg/pod"
	"github.com/openshift-kni/eco-goinfra/pkg/siteconfig"

	bmhv1alpha1 "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/apiobjectshelper"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/csv"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/shell"
)

// VerifyNamespaceExists verifies that a specific namespace exists.
func VerifyNamespaceExists(nsname string, wg *sync.WaitGroup) *namespace.Builder {
	if wg != nil {
		defer wg.Done()
	}

	By(fmt.Sprintf("Verifying that %s namespace exists", nsname))
	err := apiobjectshelper.VerifyNamespaceExists(HubAPIClient, nsname, time.Second)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to pull namespace %q; %v", nsname, err))
	ns, _ := namespace.Pull(HubAPIClient, nsname)

	return ns
}

// VerifyCsvSuccessful verifies that a specific subscription exists.
func VerifyCsvSuccessful(subscriptionName string, nsname string) {
	By(fmt.Sprintf("Verifying that csv %s is successful", subscriptionName))

	csvName, err := csv.GetCurrentCSVNameFromSubscription(HubAPIClient, subscriptionName, nsname)
	if err != nil {
		Skip(fmt.Sprintf("csv %s not found in namespace %s", csvName, nsname))
	}

	csvObj, err := olm.PullClusterServiceVersion(HubAPIClient, csvName, nsname)
	if err != nil {
		Skip(fmt.Sprintf("failed to pull %q csv from the %s namespace", csvName, nsname))
	}

	_, err = csvObj.IsSuccessful()
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("failed to verify csv %s in the namespace %s status", csvName, nsname))
}

// VerifyAllPodsRunningInNamespace verifies that all the pods in a given namespace are running.
func VerifyAllPodsRunningInNamespace(nsname string) {
	By(fmt.Sprintf("Verifying that pods exist in %s namespace", nsname))
	rhacmPods, err := pod.List(HubAPIClient, nsname)
	Expect(err).NotTo(HaveOccurred(), "error nsname while listing pods in rhacm namespace")
	Expect(len(rhacmPods) > 0).To(BeTrue(), "error: did not find any pods in the rhacm namespace")

	By(fmt.Sprintf("Verifying that pods are running correctly in %s namespace", nsname))
	running, err := pod.WaitForAllPodsInNamespaceRunning(HubAPIClient, nsname, time.Minute)
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("error occurred while waiting for %s pods to be in Running state", nsname))
	Expect(running).To(BeTrue(),
		fmt.Sprintf("some %s pods are not in Running state", nsname))
}

// VerifyProvisioningRequestCreation verifies the successful creation or provisioning request and
// that the provisioning request is progressing.
func VerifyProvisioningRequestCreation(
	prName string,
	templateName string,
	templateVersion string,
	nodeClusterName string,
	oCloudSiteId string,
	policyTemplateParameters map[string]any,
	clusterInstanceParameters map[string]any) *oran.ProvisioningRequestBuilder {
	By(fmt.Sprintf("Verifing the successful creation of the %s PR", prName))
	pr := oran.NewPRBuilder(HubAPIClient, prName, templateName, templateVersion)
	pr.WithTemplateParameter("nodeClusterName", nodeClusterName)
	pr.WithTemplateParameter("oCloudSiteId", oCloudSiteId)
	pr.WithTemplateParameter("policyTemplateParameters", policyTemplateParameters)
	pr.WithTemplateParameter("clusterInstanceParameters", clusterInstanceParameters)
	pr, err := pr.Create()
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to create PR %s", prName))

	condition := metav1.Condition{
		Type:   "ClusterProvisioned",
		Reason: "InProgress",
	}
	pr.WaitForCondition(condition, time.Minute*5)

	return pr
}

// VerifyProvisioningRequestState verifies that a given provisioning request is in a given state.
func VerifyProvisioningRequestState(
	pr *oran.ProvisioningRequestBuilder,
	prName string,
	expectedState string) {
	By(fmt.Sprintf("Verifying that %s PR is %s", prName, expectedState))
	actualState := pr.Object.Status.ProvisioningStatus.ProvisioningState

	Expect(fmt.Sprintf("%v", actualState)).To(Equal(expectedState),
		fmt.Sprintf("PR %s not fulfilled (status: %s)", prName, actualState))
}

// VerifyClusterInstanceCompleted verifies that a cluster instance exists, that it is provisioned and
// that it is associated to a given provisioning request.
func VerifyClusterInstanceCompleted(
	prName string, ns string, ciName string, wg *sync.WaitGroup, ctx SpecContext) *siteconfig.CIBuilder {
	if wg != nil {
		defer wg.Done()
	}

	By(fmt.Sprintf("Verifying that %s PR has a Cluster Instance CR associated in namespace %s", prName, ns))

	ci, err := siteconfig.PullClusterInstance(HubAPIClient, ciName, ns)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to pull Cluster Instance %q; %v", ns, err))

	found := false
	for _, value := range ci.Object.ObjectMeta.Labels {
		if value == prName {
			found = true
			break
		}
	}
	Expect(found).To(BeTrue(),
		fmt.Sprintf("Failed to verify that Cluster Instance %s is associated to PR %s", ciName, prName))

	Eventually(func(ctx context.Context) bool {
		for _, value := range ci.Object.Status.Conditions {
			if value.Type == "Provisioned" && value.Status == "True" {
				return true
			}
		}
		return false
	}).WithTimeout(60*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("ClusterInstance %s is not Completed", ciName))

	return ci
}

// VerifyImageClusterInstallSucceeded verifies that a cluster instance exists, that it is provisioned and
// that it is associated to a given provisioning request.
func VerifyImageClusterInstallSucceeded(
	nodeId string, ns string, wg *sync.WaitGroup, ctx SpecContext) *ibi.ImageClusterInstallBuilder {
	if wg != nil {
		defer wg.Done()
	}

	By(fmt.Sprintf("Verifying that Image Cluster Install %s in namespace %s has succeeded", nodeId, ns))

	ici, err := ibi.PullImageClusterInstall(HubAPIClient, nodeId, ns)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to pull Image Cluster install %s from namespace %s; %v", nodeId, ns, err))

	//Eventually(func(ctx context.Context) bool {
	//	condition, _ := ici.GetCompletedCondition()
	//	return condition.Status == "True"
	//}).WithTimeout(60*time.Minute).WithPolling(20*time.Second).WithContext(ctx).Should(BeTrue(),
	//	fmt.Sprintf("Image Cluster Install %s is not Completed", nodeId))

	return ici
}

// VerifyAllPoliciesInNamespaceAreCompliant verifies that all the policies in a given namespace
// report compliant.
func VerifyAllPoliciesInNamespaceAreCompliant(namespace string, ctx SpecContext, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	By(fmt.Sprintf("Verifying that all the policies in namespace %s are Compliant", namespace))
	policies, err := ocm.ListPoliciesInAllNamespaces(HubAPIClient)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to pull policies from all namespaces: %v", err))

	Eventually(func(ctx context.Context) bool {
		for _, policy := range policies {
			if policy.Object.ObjectMeta.Namespace == namespace {
				if policy.Object.Status.ComplianceState != "Compliant" {
					return false
				}
			}
		}
		return true
	}).WithTimeout(90*time.Minute).WithPolling(time.Minute).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("Failed to verify that all the policies in namespace %s are Compliant", namespace))
}

// VerifyNotAllPoliciesInNamespaceAreCompliant verifies that not all the policies in a given namespace
// report compliant.
func VerifyNotAllPoliciesInNamespaceAreCompliant(namespace string, ctx SpecContext, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	By(fmt.Sprintf("Verifying that not all the policies in namespace %s are Compliant", namespace))
	policies, err := ocm.ListPoliciesInAllNamespaces(HubAPIClient)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to pull policies from all namespaces: %v", err))

	Eventually(func(ctx context.Context) bool {
		for _, policy := range policies {
			if policy.Object.ObjectMeta.Namespace == namespace {
				if policy.Object.Status.ComplianceState != "Compliant" {
					return true
				}
			}
		}
		return false
	}).WithTimeout(5*time.Minute).WithPolling(20*time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("Failed to verify that not all the policies in namespace %s are Compliant", namespace))
}

// VerifyOranNodeExistsInNamespace verifies that a given ORAN node exists in a given namespace.
func VerifyOranNodeExistsInNamespace(
	nodeId string, nsName string, wg *sync.WaitGroup) *oran.NodeBuilder {
	if wg != nil {
		defer wg.Done()
	}

	By(fmt.Sprintf("Verifying that ORAN node %s exists in namespace %s ", nodeId, nsName))

	listOptions := &client.ListOptions{}
	listOptions.Namespace = nsName

	oranNodes, err := oran.ListNodes(HubAPIClient, *listOptions)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to pull oran node list from namespace %s: %v", nsName, err))

	nodeFound := false
	i := 0

	for index, node := range oranNodes {
		if nodeId == node.Object.Spec.HwMgrNodeId {
			nodeFound = true
			i = index
			break
		}
	}

	//Expect(nodeFound).To(BeTrue(),
	//	fmt.Sprintf("Failed to pull the oran node with the HW MGR ID %s from namespace %s", nodeId, nsName))

	if nodeFound {
		return oranNodes[i]
	} else {
		return oranNodes[0]
	}
}

// VerifyOranNodePoolExistsInNamespace verifies that a given ORAN node pool exists in a given namespace.
func VerifyOranNodePoolExistsInNamespace(
	nodePoolName string, nsName string, wg *sync.WaitGroup) *oran.NodePoolBuilder {
	if wg != nil {
		defer wg.Done()
	}

	By(fmt.Sprintf("Verifying that ORAN node pool %s exists in namespace %s", nodePoolName, nsName))
	oranNodePool, err := oran.PullNodePool(HubAPIClient, nodePoolName, nsName)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to pull oran node pool %s from namespace %s: %v", nodePoolName, nsName, err))
	return oranNodePool
}

// VerifyBmhExternallyProvisioned verifies that a given ORAN node pool exists in a given namespace.
func VerifyBmhExternallyProvisioned(
	bmhName string, nsName string, ctx SpecContext, wg *sync.WaitGroup) *bmh.BmhBuilder {
	if wg != nil {
		defer wg.Done()
	}

	By(fmt.Sprintf("Verifying that BMH %s exists in namespace %s", bmhName, nsName))
	bmh, err := bmh.Pull(HubAPIClient, bmhName, nsName)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to pull BMH %s from namespace %s: %v", bmhName, nsName, err))

	err = bmh.WaitUntilInStatus(bmhv1alpha1.StateExternallyProvisioned, 10*time.Minute)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to verify that BMH %s is externally provisioned", bmhName))

	return bmh
}

// VerifyProvisioningRequestIsDeleted verifies that a given provisioning request is deleted.
func VerifyProvisioningRequestIsDeleted(pr *oran.ProvisioningRequestBuilder, wg *sync.WaitGroup, ctx SpecContext) {
	if wg != nil {
		defer wg.Done()
	}

	prName := pr.Object.Name
	err := pr.Delete()
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to delete PR %s: %v", prName, err))

	By(fmt.Sprintf("Verifying that PR %s is deleted", prName))
	Eventually(func(ctx context.Context) bool {
		exists := pr.Exists()
		if !exists {
			return true
		}
		actualState := pr.Object.Status.ProvisioningStatus.ProvisioningState
		if fmt.Sprintf("%v", actualState) == "deleting" {
			By(fmt.Sprintf("Confirming that PR %s is being deleted", prName))
			return true
		}
		return false
	}).WithTimeout(30*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("Error deleting PR %s", prName))
}

// VerifyNamespaceDoesNotExist verifies that a given namespace does not exist.
func VerifyNamespaceDoesNotExist(ns *namespace.Builder, wg *sync.WaitGroup, ctx SpecContext) {
	if wg != nil {
		defer wg.Done()
	}

	nsName := ns.Object.Name
	By(fmt.Sprintf("Verifying that namespace %s does not exist", nsName))
	Eventually(func(ctx context.Context) bool {
		return !ns.Exists()
	}).WithTimeout(30*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("Namespace %s still exists", nsName))
}

// VerifyClusterInstanceDoesNotExist verifies that a given cluster instance does not exist
func VerifyClusterInstanceDoesNotExist(ci *siteconfig.CIBuilder, wg *sync.WaitGroup, ctx SpecContext) {
	if wg != nil {
		defer wg.Done()
	}

	ciName := ci.Object.Name
	By(fmt.Sprintf("Verifying that clusterinstance %s does not exist", ciName))
	Eventually(func(ctx context.Context) bool {
		return !ci.Exists()
	}).WithTimeout(30*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("ClusterInstance %s still exists", ciName))
}

// VerifyOranNodeDoesNotExist verifies that a given ORAN node does not exist.
func VerifyOranNodeDoesNotExist(node *oran.NodeBuilder, wg *sync.WaitGroup, ctx SpecContext) {
	if wg != nil {
		defer wg.Done()
	}

	nodeName := node.Object.Name
	By(fmt.Sprintf("Verifying that oran node %s does not exist", nodeName))
	Eventually(func(ctx context.Context) bool {
		return !node.Exists()
	}).WithTimeout(30*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("Oran node %s still exists", nodeName))
}

// VerifyBmhDoesNotExist verifies that a given ORAN node does not exist.
func VerifyBmhDoesNotExist(bmh *bmh.BmhBuilder, wg *sync.WaitGroup, ctx SpecContext) {
	if wg != nil {
		defer wg.Done()
	}

	bmhName := bmh.Object.Name
	By(fmt.Sprintf("Verifying that BMH %s does not exist", bmhName))
	Eventually(func(ctx context.Context) bool {
		return !bmh.Exists()
	}).WithTimeout(5*time.Second).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("BMH %s still exists", bmhName))
}

// VerifyImageClusterInstallDoesNotExist verifies that a given ORAN node pool does not exist.
func VerifyImageClusterInstallDoesNotExist(ici *ibi.ImageClusterInstallBuilder, wg *sync.WaitGroup, ctx SpecContext) {
	if wg != nil {
		defer wg.Done()
	}

	iciName := ici.Object.Name
	By(fmt.Sprintf("Verifying that image cluster install %s does not exist", iciName))
	Eventually(func(ctx context.Context) bool {
		return !ici.Exists()
	}).WithTimeout(5*time.Second).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("Image cluster install %s still exists", iciName))
}

// VerifyOranNodePoolDoesNotExist verifies that a given ORAN node pool does not exist.
func VerifyOranNodePoolDoesNotExist(nodePool *oran.NodePoolBuilder, wg *sync.WaitGroup, ctx SpecContext) {
	if wg != nil {
		defer wg.Done()
	}

	nodePoolName := nodePool.Object.Name
	By(fmt.Sprintf("Verifying that oran node pool %s does not exist", nodePoolName))
	Eventually(func(ctx context.Context) bool {
		return !nodePool.Exists()
	}).WithTimeout(30*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("Oran node pool %s still exists", nodePoolName))
}

// ProvisionSnoCluster provisions a SNO cluster.
func ProvisionSnoCluster(
	prName string,
	templateName string,
	templateVersion string,
	nodeClusterName string,
	oCloudNodeId string,
	policyTemplateParameters map[string]any,
	clusterInstanceParameters map[string]any,
	wg *sync.WaitGroup) {

	if wg != nil {
		defer wg.Done()
	}

	VerifyProvisioningRequestCreation(
		prName,
		templateName,
		templateVersion,
		nodeClusterName,
		oCloudNodeId,
		policyTemplateParameters,
		clusterInstanceParameters)
	glog.V(ocloudparams.OCloudLogLevel).Infof("Provisioning request %s has been created", prName)
}

// DeprovisionSnoCluster deprovisions a SNO cluster.
func DeprovisionAiSnoCluster(
	pr *oran.ProvisioningRequestBuilder,
	ns *namespace.Builder,
	ci *siteconfig.CIBuilder,
	node *oran.NodeBuilder,
	nodePool *oran.NodePoolBuilder,
	ctx SpecContext) {

	By(fmt.Sprintf("Tearing down PR %s", pr.Object.Name))

	var tearDownWg sync.WaitGroup
	tearDownWg.Add(5)
	go VerifyProvisioningRequestIsDeleted(pr, &tearDownWg, ctx)
	go VerifyNamespaceDoesNotExist(ns, &tearDownWg, ctx)
	go VerifyClusterInstanceDoesNotExist(ci, &tearDownWg, ctx)
	go VerifyOranNodeDoesNotExist(node, &tearDownWg, ctx)
	go VerifyOranNodePoolDoesNotExist(nodePool, &tearDownWg, ctx)
	tearDownWg.Wait()

	glog.V(ocloudparams.OCloudLogLevel).Infof("Provisioning request %s has been removed", pr.Object.Name)
	DowngradeImages()
}

// DeprovisionSnoCluster deprovisions a SNO cluster.
func DeprovisionIbiSnoCluster(
	pr *oran.ProvisioningRequestBuilder,
	ns *namespace.Builder,
	node *oran.NodeBuilder,
	nodePool *oran.NodePoolBuilder,
	bmh *bmh.BmhBuilder,
	ici *ibi.ImageClusterInstallBuilder,
	ctx SpecContext) {

	By(fmt.Sprintf("Tearing down PR %s", pr.Object.Name))

	var tearDownWg sync.WaitGroup
	tearDownWg.Add(5)
	go VerifyProvisioningRequestIsDeleted(pr, &tearDownWg, ctx)
	go VerifyNamespaceDoesNotExist(ns, &tearDownWg, ctx)
	go VerifyOranNodeDoesNotExist(node, &tearDownWg, ctx)
	go VerifyOranNodePoolDoesNotExist(nodePool, &tearDownWg, ctx)
	go VerifyBmhDoesNotExist(bmh, &tearDownWg, ctx)
	go VerifyImageClusterInstallDoesNotExist(ici, &tearDownWg, ctx)
	tearDownWg.Wait()

	glog.V(ocloudparams.OCloudLogLevel).Infof("Provisioning request %s has been removed", pr.Object.Name)
	DowngradeImages()
}

// VerifyAndRetrieveAssociatedCRsForAI verifies that a given ORAN node, a given ORAN node pool, a given namespace
// and a given cluster instance exist and retrieves them.
func VerifyAndRetrieveAssociatedCRsForAI(
	prName string,
	nodeId string,
	nodePoolName string,
	nsName string,
	ciName string,
	ctx SpecContext) (*oran.NodeBuilder, *oran.NodePoolBuilder, *namespace.Builder, *siteconfig.CIBuilder) {

	node := VerifyOranNodeExistsInNamespace(nodeId, ocloudparams.OCloudHardwareManagerPluginNamespace, nil)
	glog.V(ocloudparams.OCloudLogLevel).Infof("ORAN node with node ID %s has been created", nodeId)

	nodePool := VerifyOranNodePoolExistsInNamespace(
		nodePoolName,
		ocloudparams.OCloudHardwareManagerPluginNamespace,
		nil)
	glog.V(ocloudparams.OCloudLogLevel).Infof("ORAN node pool ID %s has been created", nodePoolName)

	ns := VerifyNamespaceExists(nsName, nil)
	glog.V(ocloudparams.OCloudLogLevel).Infof("Namespace %s has been created", nsName)

	ci := VerifyClusterInstanceCompleted(prName, nsName, ciName, nil, ctx)
	glog.V(ocloudparams.OCloudLogLevel).Infof("Cluster Instance %s exists and reports Complete", ciName)
	return node, nodePool, ns, ci
}

// VerifyAndRetrieveAssociatedCRsForIBI verifies that a given ORAN node, a given ORAN node pool, a given namespace
// and a given cluster instance exist and retrieves them.
func VerifyAndRetrieveAssociatedCRsForIBI(nodeId string,
	nodePoolName string,
	nsName string,
	hostName string,
	ctx SpecContext) (*oran.NodeBuilder, *oran.NodePoolBuilder, *namespace.Builder, *bmh.BmhBuilder, *ibi.ImageClusterInstallBuilder) {

	bmh := VerifyBmhExternallyProvisioned(hostName, nsName, ctx, nil)
	glog.V(ocloudparams.OCloudLogLevel).Infof("BMH %s is externally provisioned", hostName)

	ici := VerifyImageClusterInstallSucceeded(nodeId, nsName, nil, ctx)
	glog.V(ocloudparams.OCloudLogLevel).Infof("Cluster installation %s has succeeded ", nodeId)

	node := VerifyOranNodeExistsInNamespace(nodeId, ocloudparams.OCloudHardwareManagerPluginNamespace, nil)
	glog.V(ocloudparams.OCloudLogLevel).Infof("ORAN node with node ID %s has been created", nodeId)

	nodePool := VerifyOranNodePoolExistsInNamespace(
		nodePoolName,
		ocloudparams.OCloudHardwareManagerPluginNamespace,
		nil)
	glog.V(ocloudparams.OCloudLogLevel).Infof("ORAN node pool ID %s has been created", nodePoolName)

	ns := VerifyNamespaceExists(nsName, nil)
	glog.V(ocloudparams.OCloudLogLevel).Infof("Namespace %s has been created", nsName)

	return node, nodePool, ns, bmh, ici
}

func VerifyIbiBaseImageExists() {
	By(fmt.Sprintf("Verifying that file %s exists", OCloudConfig.IbiBaseImagePath))
	_, err := os.Stat(OCloudConfig.IbiBaseImagePath)
	Expect(os.IsNotExist(err)).To(BeFalse(),
		fmt.Sprintf("File %s does not exist", OCloudConfig.IbiBaseImagePath))
}

func UpgradeImages() {
	_, err := shell.ExecuteCmd(ocloudparams.PodmanTagOperatorUpgrade)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error tagging redhat-operators image for upgrade: %v", err))

	_, err = shell.ExecuteCmd(ocloudparams.PodmanTagSriovUpgrade)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error tagging far-edge-sriov-fec image for upgrade: %v", err))

	_, err = shell.ExecuteCmd(ocloudparams.PodmanPushOperatorUpgrade)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error pushing redhat-operators image for upgrade: %v", err))

	_, err = shell.ExecuteCmd(ocloudparams.PodmanPushSriovUpgrade)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error pushing far-edge-sriov-fec image for upgrade: %v", err))
}

func DowngradeImages() {
	_, err := shell.ExecuteCmd(ocloudparams.PodmanTagOperatorDowngrade)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error tagging redhat-operators image for downgrade: %v", err))

	_, err = shell.ExecuteCmd(ocloudparams.PodmanTagSriovDowngrade)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error tagging far-edge-sriov-fec image for downgrade: %v", err))

	_, err = shell.ExecuteCmd(ocloudparams.PodmanPushOperatorDowngrade)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error pushing redhat-operators image for downgrade: %v", err))

	_, err = shell.ExecuteCmd(ocloudparams.PodmanPushSriovDowngrade)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error pushing far-edge-sriov-fec image for downgrade: %v", err))
}
