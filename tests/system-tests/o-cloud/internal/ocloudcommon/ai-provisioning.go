package ocloudcommon

import (
	"fmt"
	"sync"

	//"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"

	//"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
	"github.com/openshift-kni/eco-goinfra/pkg/oran"
)

func VerifySuccessfulSnoProvisioning() {
	// todo
}

func VerifySuccessfulSnoDeprovisioning() {
	Fail("Intentional failure for demonstration purposes B")
}

func VerifyFailedSnoProvisioning() {
	Fail("Intentional failure for demonstration purposes C")
}

func VerifySimultaneousSnoProvisioningSameClusterTemplate() {
	Fail("Intentional failure for demonstration purposes D")
}

func VerifySimultaneousSnoDeprovisioningSameClusterTemplate(ctx SpecContext) {
	pr1name := "multiple-sno-same-cluster-template-sno02"
	pr1, err := oran.PullPR(HubAPIClient, pr1name)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to get PR %s", pr1name))
	VerifyProvisioningRequestState(pr1, pr1name, "fulfilled")

	pr2name := "multiple-sno-same-cluster-template-sno03"
	pr2, err := oran.PullPR(HubAPIClient, pr2name)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to get PR %s", pr2name))
	VerifyProvisioningRequestState(pr2, pr2name, "fulfilled")

	By(fmt.Sprintf("Verify that %s PR and %s PR are using the same template version", pr1name, pr2name))
	pr1TemplateVersion := pr1.Object.Spec.TemplateName + pr1.Object.Spec.TemplateVersion
	pr2TemplateVersion := pr2.Object.Spec.TemplateName + pr1.Object.Spec.TemplateVersion
	Expect(pr1TemplateVersion).To(Equal(pr2TemplateVersion),
		fmt.Sprintf("PR %s and %s are not using the same cluster template", pr1name, pr2name))

	ci1ns, ci1name, oranNode1 := "sno02", "sno02", "sno02" // This should be part of the config
	VerifyNamespaceExists(ci1ns)
	VerifyClusterInstanceCompleted(pr1name, ci1ns, ci1name)
	VerifyPoliciesInNamespace(ci1ns)
	VerifyOranNode(oranNode1)
	VerifyOranNodePool(oranNode1)

	ci2ns, ci2name, oranNode2 := "sno03", "sno03", "sno03" // This should be part of the config
	VerifyNamespaceExists(ci2ns)
	VerifyClusterInstanceCompleted(pr2name, ci2ns, ci2name)
	VerifyPoliciesInNamespace(ci2ns)
	VerifyOranNode(oranNode2)
	VerifyOranNodePool(oranNode2)

	// todo - delete PRs

	var waitGroup1 sync.WaitGroup
	waitGroup1.Add(2)
	
	go VerifyProvisioningRequestIsDeleted(pr1, &waitGroup1, ctx)
	go VerifyProvisioningRequestIsDeleted(pr2, &waitGroup1, ctx)

	waitGroup1.Wait()
}

func VerifySimultaneousSnoProvisioningDifferentClusterTemplate() {
	Fail("Intentional failure for demonstration purposes F")
}

func VerifySimultaneousSnoDeprovisioningDifferentClusterTemplate() {
	Fail("Intentional failure for demonstration purposes G")
}
