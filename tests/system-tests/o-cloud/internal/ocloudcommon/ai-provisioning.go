package ocloudcommon

import (
	"fmt"
	"sync"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/openshift-kni/eco-goinfra/pkg/oran"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
)

// VerifySuccessfulSnoProvisioning verifies the successful provisioning of a SNO cluster using
// Assisted Installer
func VerifySuccessfulSnoProvisioning(ctx SpecContext) {
	ProvisionSnoCluster(
		ocloudparams.PrName1,
		ocloudparams.TemplateName,
		ocloudparams.TemplateVersion1,
		ocloudparams.NodeClusterName1,
		ocloudparams.OCloudSiteId1,
		ocloudparams.PolicyTemplateParameters,
		ocloudparams.ClusterInstanceParameters1,
		nil)
	
	pr, err := oran.PullPR(HubAPIClient, ocloudparams.PrName1)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Failed to retrieve PR %s", ocloudparams.PrName1))

	node, nodePool, ns, ci := VerifyAndRetrieveAssociatedCRs(
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1,
		ctx)
	
	nsName := ns.Object.Name
	VerifyAllPoliciesInNamespaceAreCompliant(nsName, ctx)
	glog.V(ocloudparams.OCloudLogLevel).Infof("All the policies in namespace %s are Complete", nsName)

	VerifyProvisioningRequestState(pr, ocloudparams.PrName1, "fulfilled")
	glog.V(ocloudparams.OCloudLogLevel).Infof("Provisioning request %s is fulfilled", ocloudparams.PrName1)

	DeprovisionSnoCluster(pr, ns, ci, node, nodePool, ctx)
}

// VerifyFailedSnoProvisioning verifies that the provisioning of a SNO cluster using
// Assisted Installer fails
func VerifyFailedSnoProvisioning(ctx SpecContext) {
	ProvisionSnoCluster(
		ocloudparams.PrName1,
		ocloudparams.TemplateName,
		ocloudparams.TemplateVersion2,
		ocloudparams.NodeClusterName1,
		ocloudparams.OCloudSiteId1,
		ocloudparams.PolicyTemplateParameters,
		ocloudparams.ClusterInstanceParameters1,
		nil)

	pr, err := oran.PullPR(HubAPIClient, ocloudparams.PrName1)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Failed to retrieve PR %s", ocloudparams.PrName1))
	
	node, nodePool, ns, ci := VerifyAndRetrieveAssociatedCRs(
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1,
		ctx)

	VerifyProvisioningRequestState(pr, ocloudparams.PrName1, "failed")
	glog.V(ocloudparams.OCloudLogLevel).Infof("Provisioning request %s has failed", ocloudparams.PrName1)

	DeprovisionSnoCluster(pr, ns, ci, node, nodePool, ctx)
}

// VerifySimultaneousSnoProvisioningSameClusterTemplate verifies the successful provisioning of two SNO clusters 
// simultaneously with the same cluster templates.
func VerifySimultaneousSnoProvisioningSameClusterTemplate(ctx SpecContext) {
	var wg sync.WaitGroup
	wg.Add(2)
	go ProvisionSnoCluster(
		ocloudparams.PrName1,
		ocloudparams.TemplateName,
		ocloudparams.TemplateVersion1,
		ocloudparams.NodeClusterName1,
		ocloudparams.OCloudSiteId1,
		ocloudparams.PolicyTemplateParameters,
		ocloudparams.ClusterInstanceParameters1,
		&wg)
	go ProvisionSnoCluster(
		ocloudparams.PrName2,
		ocloudparams.TemplateName,
		ocloudparams.TemplateVersion1,
		ocloudparams.NodeClusterName2,
		ocloudparams.OCloudSiteId2,
		ocloudparams.PolicyTemplateParameters,
		ocloudparams.ClusterInstanceParameters2,
		&wg)
	wg.Wait()

	_, _, ns1, _ := VerifyAndRetrieveAssociatedCRs(
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1,
		ctx)

	_, _, ns2, _ := VerifyAndRetrieveAssociatedCRs(
		ocloudparams.NodeClusterName2, 
		ocloudparams.NodeClusterName2, 
		ocloudparams.NodeClusterName2, 
		ocloudparams.NodeClusterName2,
		ctx)
	
	var wg2 sync.WaitGroup
	wg2.Add(2)
	go	VerifyAllPoliciesInNamespaceAreCompliant(ns1.Object.Name, ctx)
	go	VerifyAllPoliciesInNamespaceAreCompliant(ns2.Object.Name, ctx)
	wg2.Wait()
}

// VerifySimultaneousSnoDeprovisioningSameClusterTemplate verifies the successful deletion of 
// two SNO clusters with the same cluster template.
func VerifySimultaneousSnoDeprovisioningSameClusterTemplate(ctx SpecContext) {
	prName1 := ocloudparams.PrName1
	prName2 := ocloudparams.PrName2

	pr1, err := oran.PullPR(HubAPIClient, prName1)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to get PR %s", prName1))
	VerifyProvisioningRequestState(pr1, prName1, "fulfilled")

	pr2, err := oran.PullPR(HubAPIClient, prName2)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to get PR %s", prName2))
	VerifyProvisioningRequestState(pr2, prName2, "fulfilled")

	By(fmt.Sprintf("Verify that %s PR and %s PR are using the same template version", prName1, prName2))
	pr1TemplateVersion := pr1.Object.Spec.TemplateName + pr1.Object.Spec.TemplateVersion
	pr2TemplateVersion := pr2.Object.Spec.TemplateName + pr1.Object.Spec.TemplateVersion
	Expect(pr1TemplateVersion).To(Equal(pr2TemplateVersion),
		fmt.Sprintf("PR %s and %s are not using the same cluster template", prName1, prName2))

	node1, nodePool1, ns1, ci1 := VerifyAndRetrieveAssociatedCRs(
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1,
		ctx)

	node2, nodePool2, ns2, ci2 := VerifyAndRetrieveAssociatedCRs(
		ocloudparams.NodeClusterName2, 
		ocloudparams.NodeClusterName2, 
		ocloudparams.NodeClusterName2, 
		ocloudparams.NodeClusterName2,
		ctx)

	var waitGroup sync.WaitGroup

	waitGroup.Add(10)
	go VerifyProvisioningRequestIsDeleted(pr1, &waitGroup, ctx)
	go VerifyProvisioningRequestIsDeleted(pr2, &waitGroup, ctx)
	go VerifyNamespaceDoesNotExist(ns1, &waitGroup, ctx)
	go VerifyNamespaceDoesNotExist(ns2, &waitGroup, ctx)
	go VerifyClusterInstanceDoesNotExist(ci1, &waitGroup, ctx)
	go VerifyClusterInstanceDoesNotExist(ci2, &waitGroup, ctx)
	go VerifyOranNodeDoesNotExist(node1, &waitGroup, ctx)
	go VerifyOranNodeDoesNotExist(node2, &waitGroup, ctx)
	go VerifyOranNodePoolDoesNotExist(nodePool1, &waitGroup, ctx)
	go VerifyOranNodePoolDoesNotExist(nodePool2, &waitGroup, ctx)
	waitGroup.Wait()
}

// VerifySimultaneousSnoProvisioningDifferentClusterTemplates verifies the successful provisioning of 
// two SNO clusters simultaneously with different cluster templates.
func VerifySimultaneousSnoProvisioningDifferentClusterTemplates(ctx SpecContext) {
	var wg sync.WaitGroup
	wg.Add(2)
	go ProvisionSnoCluster(
		ocloudparams.PrName1,
		ocloudparams.TemplateName,
		ocloudparams.TemplateVersion1,
		ocloudparams.NodeClusterName1,
		ocloudparams.OCloudSiteId1,
		ocloudparams.PolicyTemplateParameters,
		ocloudparams.ClusterInstanceParameters1,
		&wg)
	go ProvisionSnoCluster(
		ocloudparams.PrName2,
		ocloudparams.TemplateName,
		ocloudparams.TemplateVersion3,
		ocloudparams.NodeClusterName2,
		ocloudparams.OCloudSiteId2,
		ocloudparams.PolicyTemplateParameters,
		ocloudparams.ClusterInstanceParameters2,
		&wg)
	wg.Wait()

	_, _, ns1, _ := VerifyAndRetrieveAssociatedCRs(
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1,
		ctx)

	_, _, ns2, _ := VerifyAndRetrieveAssociatedCRs(
		ocloudparams.NodeClusterName2, 
		ocloudparams.NodeClusterName2, 
		ocloudparams.NodeClusterName2, 
		ocloudparams.NodeClusterName2,
		ctx)
	
	var wg2 sync.WaitGroup
	wg2.Add(2)
	go	VerifyAllPoliciesInNamespaceAreCompliant(ns1.Object.Name, ctx)
	go	VerifyAllPoliciesInNamespaceAreCompliant(ns2.Object.Name, ctx)
	wg2.Wait()
}

// VerifySimultaneousSnoDeprovisioningDifferentClusterTemplates verifies the successful deletion of 
// two SNO clusters with different cluster templates.
func VerifySimultaneousSnoDeprovisioningDifferentClusterTemplates(ctx SpecContext) {
	prName1 := ocloudparams.PrName1
	prName2 := ocloudparams.PrName2

	pr1, err := oran.PullPR(HubAPIClient, prName1)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to get PR %s", prName1))
	VerifyProvisioningRequestState(pr1, prName1, "fulfilled")

	pr2, err := oran.PullPR(HubAPIClient, prName2)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to get PR %s", prName2))
	VerifyProvisioningRequestState(pr2, prName2, "fulfilled")

	By(fmt.Sprintf("Verify that %s PR and %s PR are using different cluster template versions", prName1, prName2))
	pr1TemplateVersion := pr1.Object.Spec.TemplateName + pr1.Object.Spec.TemplateVersion
	pr2TemplateVersion := pr2.Object.Spec.TemplateName + pr1.Object.Spec.TemplateVersion
	Expect(pr1TemplateVersion).NotTo(Equal(pr2TemplateVersion),
		fmt.Sprintf("PR %s and %s are using the same cluster template", prName1, prName2))

	node1, nodePool1, ns1, ci1 := VerifyAndRetrieveAssociatedCRs(
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1, 
		ocloudparams.NodeClusterName1,
		ctx)

	node2, nodePool2, ns2, ci2 := VerifyAndRetrieveAssociatedCRs(
		ocloudparams.NodeClusterName2, 
		ocloudparams.NodeClusterName2, 
		ocloudparams.NodeClusterName2, 
		ocloudparams.NodeClusterName2,
		ctx)

	var waitGroup sync.WaitGroup

	waitGroup.Add(10)
	go VerifyProvisioningRequestIsDeleted(pr1, &waitGroup, ctx)
	go VerifyProvisioningRequestIsDeleted(pr2, &waitGroup, ctx)
	go VerifyNamespaceDoesNotExist(ns1, &waitGroup, ctx)
	go VerifyNamespaceDoesNotExist(ns2, &waitGroup, ctx)
	go VerifyClusterInstanceDoesNotExist(ci1, &waitGroup, ctx)
	go VerifyClusterInstanceDoesNotExist(ci2, &waitGroup, ctx)
	go VerifyOranNodeDoesNotExist(node1, &waitGroup, ctx)
	go VerifyOranNodeDoesNotExist(node2, &waitGroup, ctx)
	go VerifyOranNodePoolDoesNotExist(nodePool1, &waitGroup, ctx)
	go VerifyOranNodePoolDoesNotExist(nodePool2, &waitGroup, ctx)
	waitGroup.Wait()
}

