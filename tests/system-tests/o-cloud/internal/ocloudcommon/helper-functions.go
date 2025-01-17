package ocloudcommon

import (
	"context"
	"fmt"
	"time"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"

	"github.com/openshift-kni/eco-goinfra/pkg/ocm"
	"github.com/openshift-kni/eco-goinfra/pkg/olm"
	"github.com/openshift-kni/eco-goinfra/pkg/oran"
	"github.com/openshift-kni/eco-goinfra/pkg/namespace"
	"github.com/openshift-kni/eco-goinfra/pkg/pod"
	"github.com/openshift-kni/eco-goinfra/pkg/siteconfig"

	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/apiobjectshelper"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/csv"
)

func VerifyNamespaceExists(nsname string) {
	By(fmt.Sprintf("Verifying that %s namespace exists", nsname))
	err := apiobjectshelper.VerifyNamespaceExists(HubAPIClient, nsname, time.Second)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to pull namespace %q; %v", nsname, err))
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

func VerifyProvisioningRequestState(pr *oran.ProvisioningRequestBuilder, prName string, expectedState string) {
	By(fmt.Sprintf("Verifying that %s PR if %s", prName, expectedState))
	actualState := pr.Object.Status.ProvisioningStatus.ProvisioningState
	Expect(fmt.Sprintf("%v", actualState)).To(Equal(expectedState),
		fmt.Sprintf("PR %s not fulfilled (status: %s)", prName, actualState))
}

func VerifyClusterInstanceCompleted(prName string, ns string, ciName string) {
	By(fmt.Sprintf("Verifying that %s PR has a Cluster Instance CR associated in namespace %s", prName, ns))

	ci1, err := siteconfig.PullClusterInstance(HubAPIClient, ciName, ns)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to pull Cluster Instance %q; %v", ns, err))

	found := false
	for _, value := range ci1.Object.ObjectMeta.Labels {
		if value == prName {
			found = true
			break
		}
	}
	Expect(found).To(BeTrue(), "Failed to verify that Cluster Instance contains the label %s", prName)
}

func VerifyPoliciesInNamespace(namespace string) {
	By(fmt.Sprintf("Verifying that all the policies in namespace %s are Compliant", namespace))
	policies, err := ocm.ListPoliciesInAllNamespaces(HubAPIClient)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to pull policies from all namespaces: %v", err))

	compliant := true
	for _, policy := range policies {
		if policy.Object.ObjectMeta.Namespace == namespace {
			if policy.Object.Status.ComplianceState != "Compliant" {
				compliant = false
				break
			}
		}
	}

	Expect(compliant).To(BeTrue(), "Failed to verify that all the policies in namespace %s are Compliant", namespace)
}

func VerifyOranNode(node string) {
	By(fmt.Sprintf("Verifying that ORAN node %s exists in oran-hwmgr-plugin namespace", node))
	// todo -  check with Kirsten
}

func VerifyOranNodePool(nodePool string) {
	By(fmt.Sprintf("Verifying that ORAN node pool %s exists in oran-hwmgr-plugin namespace", nodePool))
	// todo -  check with Kirsten
}

func VerifyProvisioningRequestIsDeleted(pr *oran.ProvisioningRequestBuilder, wg *sync.WaitGroup, ctx SpecContext) {
	defer wg.Done()

	prName := pr.Object.Name
	//err := pr.Delete()
	//Expect(err).ToNot(HaveOccurred(),
	//	fmt.Sprintf("Failed to delete PR %s: %v", prName, err))

	// TODO - Idea, we need to verify as well that the ns, ci, nodepool, etc... are removed. Can we have goroutins inside goroutines?
	
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

func VerifyNamespaceIsDeleted(ns *namespace.Builder, wg *sync.WaitGroup, ctx SpecContext) {
	defer wg.Done()

	nsName := ns.Object.Name
	//err := ns.Delete()
	//Expect(err).ToNot(HaveOccurred(),
	//	fmt.Sprintf("Failed to delete PR %s: %v", prName, err))
	
	By(fmt.Sprintf("Verifying that namespace %s is deleted", nsName))
	Eventually(func(ctx context.Context) bool {
		exists := ns.Exists()
		if !exists {
			return true
		}
		return false
	}).WithTimeout(25*time.Minute).WithPolling(time.Second).WithContext(ctx).Should(BeTrue(),
		fmt.Sprintf("Error deleting namespace %s", nsName))
}

