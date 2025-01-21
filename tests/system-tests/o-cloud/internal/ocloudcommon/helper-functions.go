package ocloudcommon

import (
	"context"
	"fmt"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"

	"github.com/openshift-kni/eco-goinfra/pkg/namespace"
	"github.com/openshift-kni/eco-goinfra/pkg/ocm"
	"github.com/openshift-kni/eco-goinfra/pkg/olm"
	"github.com/openshift-kni/eco-goinfra/pkg/oran"
	"github.com/openshift-kni/eco-goinfra/pkg/pod"
	"github.com/openshift-kni/eco-goinfra/pkg/siteconfig"

	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/apiobjectshelper"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/csv"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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

func VerifyPodsRunning(nsname string) {
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

func VerifyProvisioningRequestCreation(prName string, templateName string, templateVersion string, nodeClusterName string, oCloudSiteId string, policyTemplateParameters map[string]any, clusterInstanceParameters map[string]any) *oran.ProvisioningRequestBuilder {
	By(fmt.Sprintf("Verifing the successful creation of the %s PR", prName))
	pr := oran.NewPRBuilder(HubAPIClient, prName, templateName, templateVersion)
	pr.WithTemplateParameter("nodeClusterName", nodeClusterName)
	pr.WithTemplateParameter("oCloudSiteId", oCloudSiteId)
	pr.WithTemplateParameter("policyTemplateParameters", policyTemplateParameters)
	pr.WithTemplateParameter("clusterInstanceParameters", clusterInstanceParameters)
	pr, err := pr.Create()
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to create PR %s", prName))

	VerifyProvisioningRequestState(pr, prName, "progressing")
	return pr
}

func VerifyProvisioningRequestState(pr *oran.ProvisioningRequestBuilder, prName string, expectedState string) {
	By(fmt.Sprintf("Verifying that %s PR if %s", prName, expectedState))
	actualState := pr.Object.Status.ProvisioningStatus.ProvisioningState
	Expect(fmt.Sprintf("%v", actualState)).To(Equal(expectedState),
		fmt.Sprintf("PR %s not fulfilled (status: %s)", prName, actualState))
}

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
	}).WithTimeout(25*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("ClusterInstance %s is not Completed", ciName))

	return ci
}

func VerifyAllPoliciesInNamespaceAreCompliant(namespace string, ctx SpecContext) {
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
	}).WithTimeout(25*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("Failed to verify that all the policies in namespace %s are Compliant", namespace))
}

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

	Expect(nodeFound).To(BeTrue(),
		fmt.Sprintf("Failed to pull the oran node with the HW MGR ID %s from namespace %s", nodeId, nsName))
	
	if nodeFound {
		return oranNodes[i]
	} 

	return nil
}

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
	}).WithTimeout(25*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("Error deleting PR %s", prName))
}

func VerifyNamespaceDoesNotExist(ns *namespace.Builder, wg *sync.WaitGroup, ctx SpecContext) {
	if wg != nil {
		defer wg.Done()
	}

	nsName := ns.Object.Name
	By(fmt.Sprintf("Verifying that namespace %s does not exist", nsName))
	Eventually(func(ctx context.Context) bool {
		return !ns.Exists()
	}).WithTimeout(25*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("Namespace %s still exists", nsName))
}

func VerifyClusterInstanceDoesNotExist(ci *siteconfig.CIBuilder, wg *sync.WaitGroup, ctx SpecContext) {
	if wg != nil {
		defer wg.Done()
	}

	ciName := ci.Object.Name
	By(fmt.Sprintf("Verifying that clusterinstance %s does not exist", ciName))
	Eventually(func(ctx context.Context) bool {
		return !ci.Exists()
	}).WithTimeout(25*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("ClusterInstance %s still exists", ciName))
}

func VerifyOranNodeDoesNotExist(node *oran.NodeBuilder, wg *sync.WaitGroup, ctx SpecContext) {
	if wg != nil {
		defer wg.Done()
	}

	nodeName := node.Object.Name
	By(fmt.Sprintf("Verifying that oran node %s does not exist", nodeName))
	Eventually(func(ctx context.Context) bool {
		return !node.Exists()
	}).WithTimeout(25*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("Oran node %s still exists", nodeName))
}

func VerifyOranNodePoolDoesNotExist(nodePool *oran.NodePoolBuilder, wg *sync.WaitGroup, ctx SpecContext) {
	if wg != nil {
		defer wg.Done()
	}

	nodePoolName := nodePool.Object.Name
	By(fmt.Sprintf("Verifying that oran node pool %s does not exist", nodePoolName))
	Eventually(func(ctx context.Context) bool {
		return !nodePool.Exists()
	}).WithTimeout(25*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("Oran node pool %s still exists", nodePoolName))
}
