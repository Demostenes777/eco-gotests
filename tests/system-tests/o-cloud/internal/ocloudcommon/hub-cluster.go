package ocloudcommon

import (
	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
)

// VerifyACM verifies that ACM is installed in the Hub cluster.
func VerifyACM() {
	VerifyNamespaceExists(ocloudparams.AcmNamespace, nil)
	VerifyCsvSuccessful(ocloudparams.AcmSubscriptionName, ocloudparams.AcmNamespace)
	VerifyAllPodsRunningInNamespace(ocloudparams.AcmNamespace)
}

// VerifyGitOpsOperator verifies that the GitOps operator is installed in the Hub cluster.
func VerifyGitOpsOperator() {
	VerifyNamespaceExists(ocloudparams.OpenshiftGitOpsNamespace, nil)
	VerifyCsvSuccessful(ocloudparams.OpenshiftGitOpsSubscriptionName, ocloudparams.OpenshiftGitOpsNamespace)
}

// VerifySiteConfigOperator verifies that SiteConfig is installed in the hub cluster.
func VerifySiteConfigOperator() {
	// todo - need to find how to verify it
}

// VerifyOranO2ImsOperator verifies that the O-Cloud Manager operator is installed in the hub cluster.
func VerifyOranO2ImsOperator() {
	VerifyNamespaceExists(ocloudparams.OCloudO2ImsNamespace, nil)
	VerifyCsvSuccessful(ocloudparams.OCloudO2ImsSubscriptionName, ocloudparams.OCloudO2ImsNamespace)
}

// VerifyOranHardwareManagerPluginOperator verifies that the O-Cloud Hardware Manager Plugin operator
// is installed in the Hub cluster.
func VerifyOranHardwareManagerPluginOperator() {
	VerifyNamespaceExists(ocloudparams.OCloudHardwareManagerPluginNamespace, nil)
	VerifyCsvSuccessful(ocloudparams.OCloudHardwareManagerPluginSubscriptionName,
		ocloudparams.OCloudHardwareManagerPluginNamespace)
}
