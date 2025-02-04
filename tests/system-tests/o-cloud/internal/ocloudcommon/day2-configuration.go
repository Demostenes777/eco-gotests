package ocloudcommon

import (
	"fmt"
	"os"

	"sync"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/openshift-kni/eco-goinfra/pkg/oran"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/olm"

	. "github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/csv"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/shell"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
)

func SuccessfulOperatorUpgrade(ctx SpecContext) {
	DowngradeImages()
	// Deploy SNO 02 and SNO 03
	prName1 := ocloudparams.PrMultipleDay2Sno02
	prName2 := ocloudparams.PrMultipleDay2Sno02

	var wg sync.WaitGroup
	wg.Add(2)
	go ProvisionSnoCluster(
		prName1,
		ocloudparams.TemplateName,
		ocloudparams.TemplateVersion6,
		ocloudparams.NodeClusterName1,
		ocloudparams.OCloudSiteId1,
		ocloudparams.PolicyTemplateParameters,
		ocloudparams.ClusterInstanceParameters1,
		&wg)
	go ProvisionSnoCluster(
		prName2,
		ocloudparams.TemplateName,
		ocloudparams.TemplateVersion6,
		ocloudparams.NodeClusterName2,
		ocloudparams.OCloudSiteId2,
		ocloudparams.PolicyTemplateParameters,
		ocloudparams.ClusterInstanceParameters2,
		&wg)
	wg.Wait()

	node1, nodePool1, ns1, ci1 := VerifyAndRetrieveAssociatedCRsForAI(
		prName1,
		ocloudparams.NodeClusterName1,
		ocloudparams.NodeClusterName1,
		ocloudparams.NodeClusterName1,
		ocloudparams.NodeClusterName1,
		ctx)

	node2, nodePool2, ns2, ci1 := VerifyAndRetrieveAssociatedCRsForAI(
		prName2,
		ocloudparams.NodeClusterName2,
		ocloudparams.NodeClusterName2,
		ocloudparams.NodeClusterName2,
		ocloudparams.NodeClusterName2,
		ctx)

	var wg2 sync.WaitGroup
	wg2.Add(2)
	go VerifyAllPoliciesInNamespaceAreCompliant(ns1.Object.Name, ctx, &wg2)
	go VerifyAllPoliciesInNamespaceAreCompliant(ns2.Object.Name, ctx, &wg2)
	wg2.Wait()

 	prSno1, err := oran.PullPR(HubAPIClient, prName1)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Failed to retrieve PR %s", prName1))

	VerifyProvisioningRequestState(prSno1, prName1, "fulfilled")
	glog.V(ocloudparams.OCloudLogLevel).Infof("Provisioning request %s is fulfilled", prName1)

	sno1ApiClient := CreateSnoApiClient(ocloudparams.NodeClusterName1)
	VerifyPtpOperatorVersionInSno(
		sno1ApiClient, 
		ocloudparams.PTPVersionMajorOld, 
		ocloudparams.PTPVersionMinorOld, 
		ocloudparams.PTPVersionPatchOld, 
		ocloudparams.PTPVersionPrereleaseOld)

	prSno2, err := oran.PullPR(HubAPIClient, prName2)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Failed to retrieve PR %s", prName2))

	VerifyProvisioningRequestState(prSno2, prName2, "fulfilled")
	glog.V(ocloudparams.OCloudLogLevel).Infof("Provisioning request %s is fulfilled", prName2)
	
	sno2ApiClient := CreateSnoApiClient(ocloudparams.NodeClusterName2)
	VerifyPtpOperatorVersionInSno(
		sno2ApiClient, 
		ocloudparams.PTPVersionMajorOld, 
		ocloudparams.PTPVersionMinorOld, 
		ocloudparams.PTPVersionPatchOld, 
		ocloudparams.PTPVersionPrereleaseOld)

	UpgradeImages()

	var wg3 sync.WaitGroup
	wg3.Add(2)
	go VerifyNotAllPoliciesInNamespaceAreCompliant(ocloudparams.NodeClusterName1, ctx, &wg3)
	go VerifyNotAllPoliciesInNamespaceAreCompliant(ocloudparams.NodeClusterName2, ctx, &wg3)
	wg3.Wait()


	go VerifyProvisioningRequestState(prSno1, prName1, "progressing")

	go VerifyProvisioningRequestState(prSno2, prName2, "progressing")

	//path = fmt.Sprintf("tmp/%s/auth", ocloudparams.NodeClusterName2)
	//err = os.MkdirAll(path, 0750)
	//Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error creating directory %s", path))

	//path = fmt.Sprintf("tmp/%s", ocloudparams.NodeClusterName1)
	//err = os.RemoveAll(path)
	//Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error removing directory %s", path))

	//


	
	

	// Once they are completed export their kufeconfig files to environment variables
	// Set SNO 02 and SNO 03 api client kubeconfig path
	// Verify CSV version

	

	//DeprovisionAiSnoCluster(pr, ns, ci, node, nodePool, ctx)
}

func VerifyPtpOperatorVersionInSno(sno1ApiClient *clients.Settings, 
	major uint64, minor uint64, patch uint64, prerelease uint64) {
	csvName, err := csv.GetCurrentCSVNameFromSubscription(sno1ApiClient, 
		ocloudparams.PtpOperatorSubscriptionName, ocloudparams.PtpNamespace)
	Expect(err).NotTo(HaveOccurred(), 
		fmt.Sprintf("csv %s not found in namespace %s", csvName, ocloudparams.PtpNamespace))

	csvObj, err := olm.PullClusterServiceVersion(sno1ApiClient, csvName, ocloudparams.PtpNamespace)
	Expect(err).NotTo(HaveOccurred(), 
		fmt.Sprintf("failed to pull %q csv from the %s namespace", csvName, ocloudparams.PtpNamespace))

	versionOk := false
	ptpVersion := csvObj.Object.Spec.Version
	if ptpVersion.Major == major &&
		ptpVersion.Minor == minor &&
		ptpVersion.Patch == patch {
		for _, pre := range csvObj.Object.Spec.Version.Pre {
			if pre.VersionNum == prerelease {
				versionOk = true
			}
		}
	}

	Expect(versionOk).To(BeTrue(), fmt.Sprintf("PTP version %s is not the expected one", ptpVersion))
}

func CreateSnoApiClient(nodeName string) *clients.Settings {
	path := fmt.Sprintf("tmp/%s/auth", nodeName)
	err := os.MkdirAll(path, 0750)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error creating directory %s", path))

	createSnoKubeconfig := fmt.Sprintf(ocloudparams.SnoKubeconfigCreate, nodeName, nodeName, nodeName)
	_, err = shell.ExecuteCmd(createSnoKubeconfig)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error creating %s kubeconfig", nodeName))

	snoKubeconfigPath := fmt.Sprintf("tmp/%s/auth/kubeconfig", nodeName) 
	snoApiClient := clients.New(snoKubeconfigPath)
	return snoApiClient
}


func FailedOperatorUpgrade() {
	Fail("Intentional failure for demonstration purposes K")
}

func FailedPartialOperatorUpgrade() {
	Fail("Intentional failure for demonstration purposes L")
}
