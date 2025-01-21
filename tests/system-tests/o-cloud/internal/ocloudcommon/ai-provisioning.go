package ocloudcommon

import (
	"fmt"
	"sync"

	//"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"

	"github.com/openshift-kni/eco-goinfra/pkg/oran"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
)

func VerifySuccessfulSnoProvisioning(ctx SpecContext) {
	prName := "single-sno-ai-success"
	templateName := "sno-ran-du"
	templateVersion := "v4-18-0-ec3-1"
	nodeClusterName := "sno02"
	oCloudSiteId := "sno02"
	policyTemplateParameters := map[string]any{}
	clusterInstanceParameters := map[string]any{
		"clusterName": "sno02",
		"nodes": []map[string]any{
			{
				"hostName": "sno02.oran.telcoqe.eng.rdu2.dc.redhat.comf",
			},
		},
	}

	pr := VerifyProvisioningRequestCreation(
		prName, 
		templateName, 
		templateVersion, 
		nodeClusterName, 
		oCloudSiteId, 
		policyTemplateParameters, 
		clusterInstanceParameters)

	nsName, ciName, nodePoolName, nodeId := nodeClusterName, nodeClusterName, nodeClusterName, nodeClusterName
	node := VerifyOranNodeExistsInNamespace(nodeId, ocloudparams.OCloudHardwareManagerPluginNamespace, nil) 
	nodePool := VerifyOranNodePoolExistsInNamespace(nodePoolName, ocloudparams.OCloudHardwareManagerPluginNamespace, nil)  
	ns := VerifyNamespaceExists(nsName, nil)
	ci := VerifyClusterInstanceCompleted(prName, nsName, ciName, nil, ctx)
	VerifyAllPoliciesInNamespaceAreCompliant(nsName, ctx)
	VerifyProvisioningRequestState(pr, prName, "fulfilled")

	By(fmt.Sprintf("Tearing down PR %s", prName))

	var tearDownWg sync.WaitGroup
	tearDownWg.Add(5)
	go VerifyProvisioningRequestIsDeleted(pr, &tearDownWg, ctx)
	go VerifyNamespaceDoesNotExist(ns, &tearDownWg, ctx)
	go VerifyClusterInstanceDoesNotExist(ci, &tearDownWg, ctx)
	go VerifyOranNodeDoesNotExist(node, &tearDownWg, ctx)
	go VerifyOranNodePoolDoesNotExist(nodePool, &tearDownWg, ctx)
	tearDownWg.Wait()
}

func VerifyFailedSnoProvisioning() {
	Fail("Intentional failure for demonstration purposes C")
}

func VerifySimultaneousSnoProvisioningSameClusterTemplate() {
	Fail("Intentional failure for demonstration purposes D")
}

func VerifySimultaneousSnoDeprovisioningSameClusterTemplate(ctx SpecContext) {
	prName1 := "multiple-sno-same-cluster-template-sno02"
	pr1, err := oran.PullPR(HubAPIClient, prName1)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to get PR %s", prName1))
	VerifyProvisioningRequestState(pr1, prName1, "fulfilled")

	prName2 := "multiple-sno-same-cluster-template-sno03"
	pr2, err := oran.PullPR(HubAPIClient, prName2)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to get PR %s", prName2))
	VerifyProvisioningRequestState(pr2, prName2, "fulfilled")

	By(fmt.Sprintf("Verify that %s PR and %s PR are using the same template version", prName1, prName2))
	pr1TemplateVersion := pr1.Object.Spec.TemplateName + pr1.Object.Spec.TemplateVersion
	pr2TemplateVersion := pr2.Object.Spec.TemplateName + pr1.Object.Spec.TemplateVersion
	Expect(pr1TemplateVersion).To(Equal(pr2TemplateVersion),
		fmt.Sprintf("PR %s and %s are not using the same cluster template", prName1, prName2))

	nsName1, ciName1, oranNodePoolName1, nodeId1 := "sno02", "sno02", "sno02", "sno02" // This should be part of the config
	nodePool1 := VerifyOranNodePoolExistsInNamespace(oranNodePoolName1, ocloudparams.OCloudHardwareManagerPluginNamespace, nil)
	node1 := VerifyOranNodeExistsInNamespace(nodeId1, ocloudparams.OCloudHardwareManagerPluginNamespace, nil)
	ciNs1 := VerifyNamespaceExists(nsName1, nil)
	ci1 := VerifyClusterInstanceCompleted(prName1, nsName1, ciName1, nil, ctx)
	VerifyAllPoliciesInNamespaceAreCompliant(nsName1, ctx)

	nsName2, ciName2, oranNodePoolName2, nodeId2 := "sno03", "sno03", "sno03", "sno03" // This should be part of the config
	node2 := VerifyOranNodeExistsInNamespace(nodeId2, ocloudparams.OCloudHardwareManagerPluginNamespace, nil) 
	nodePool2 := VerifyOranNodePoolExistsInNamespace(oranNodePoolName2, ocloudparams.OCloudHardwareManagerPluginNamespace, nil)  
	ciNs2 := VerifyNamespaceExists(nsName2, nil)
	ci2 := VerifyClusterInstanceCompleted(prName2, nsName2, ciName2, nil, ctx)
	VerifyAllPoliciesInNamespaceAreCompliant(nsName2, ctx)

	var waitGroup sync.WaitGroup

	waitGroup.Add(10)
	go VerifyProvisioningRequestIsDeleted(pr1, &waitGroup, ctx)
	go VerifyProvisioningRequestIsDeleted(pr2, &waitGroup, ctx)
	go VerifyNamespaceDoesNotExist(ciNs1, &waitGroup, ctx)
	go VerifyNamespaceDoesNotExist(ciNs2, &waitGroup, ctx)
	go VerifyClusterInstanceDoesNotExist(ci1, &waitGroup, ctx)
	go VerifyClusterInstanceDoesNotExist(ci2, &waitGroup, ctx)
	go VerifyOranNodeDoesNotExist(node1, &waitGroup, ctx)
	go VerifyOranNodeDoesNotExist(node2, &waitGroup, ctx)
	go VerifyOranNodePoolDoesNotExist(nodePool1, &waitGroup, ctx)
	go VerifyOranNodePoolDoesNotExist(nodePool2, &waitGroup, ctx)
	waitGroup.Wait()
}

func VerifySimultaneousSnoProvisioningDifferentClusterTemplates() {
	Fail("Intentional failure for demonstration purposes F")
}

func VerifySimultaneousSnoDeprovisioningDifferentClusterTemplates(ctx SpecContext) {
	Fail("Intentional failure for demonstration purposes G")
	prName1 := "multiple-sno-different-cluster-templates-sno02"
	pr1, err := oran.PullPR(HubAPIClient, prName1)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to get PR %s", prName1))
	VerifyProvisioningRequestState(pr1, prName1, "fulfilled")

	prName2 := "multiple-sno-different-cluster-templates-sno03"
	pr2, err := oran.PullPR(HubAPIClient, prName2)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to get PR %s", prName2))
	VerifyProvisioningRequestState(pr2, prName2, "fulfilled")

	By(fmt.Sprintf("Verify that %s PR and %s PR are using different cluster template versions", prName1, prName2))
	pr1TemplateVersion := pr1.Object.Spec.TemplateName + pr1.Object.Spec.TemplateVersion
	pr2TemplateVersion := pr2.Object.Spec.TemplateName + pr1.Object.Spec.TemplateVersion
	Expect(pr1TemplateVersion).NotTo(Equal(pr2TemplateVersion),
		fmt.Sprintf("PR %s and %s are using the same cluster template", prName1, prName2))

	nsName1, ciName1, oranNodePoolName1, nodeId1 := "sno02", "sno02", "sno02", "sno02" // This should be part of the config
	nodePool1 := VerifyOranNodePoolExistsInNamespace(oranNodePoolName1, ocloudparams.OCloudHardwareManagerPluginNamespace, nil)
	node1 := VerifyOranNodeExistsInNamespace(nodeId1, ocloudparams.OCloudHardwareManagerPluginNamespace, nil)
	ciNs1 := VerifyNamespaceExists(nsName1, nil)
	ci1 := VerifyClusterInstanceCompleted(prName1, nsName1, ciName1, nil, ctx)
	VerifyAllPoliciesInNamespaceAreCompliant(nsName1, ctx)

	nsName2, ciName2, oranNodePoolName2, nodeId2 := "sno03", "sno03", "sno03", "sno03" // This should be part of the config
	node2 := VerifyOranNodeExistsInNamespace(nodeId2, ocloudparams.OCloudHardwareManagerPluginNamespace, nil) 
	nodePool2 := VerifyOranNodePoolExistsInNamespace(oranNodePoolName2, ocloudparams.OCloudHardwareManagerPluginNamespace, nil)  
	ciNs2 := VerifyNamespaceExists(nsName2, nil)
	ci2 := VerifyClusterInstanceCompleted(prName2, nsName2, ciName2, nil, ctx)
	VerifyAllPoliciesInNamespaceAreCompliant(nsName2, ctx)

	var waitGroup sync.WaitGroup

	waitGroup.Add(10)
	go VerifyProvisioningRequestIsDeleted(pr1, &waitGroup, ctx)
	go VerifyProvisioningRequestIsDeleted(pr2, &waitGroup, ctx)
	go VerifyNamespaceDoesNotExist(ciNs1, &waitGroup, ctx)
	go VerifyNamespaceDoesNotExist(ciNs2, &waitGroup, ctx)
	go VerifyClusterInstanceDoesNotExist(ci1, &waitGroup, ctx)
	go VerifyClusterInstanceDoesNotExist(ci2, &waitGroup, ctx)
	go VerifyOranNodeDoesNotExist(node1, &waitGroup, ctx)
	go VerifyOranNodeDoesNotExist(node2, &waitGroup, ctx)
	go VerifyOranNodePoolDoesNotExist(nodePool1, &waitGroup, ctx)
	go VerifyOranNodePoolDoesNotExist(nodePool2, &waitGroup, ctx)
	waitGroup.Wait()
}
