package ocloudcommon

import (
	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
)

func VerifyACM() {
	VerifyNamespaceExists(ocloudparams.AcmNamespace, nil)
	VerifyCsvSuccessful(ocloudparams.AcmSubscriptionName, ocloudparams.AcmNamespace)
	VerifyPodsRunning(ocloudparams.AcmNamespace)
}

func VerifyGitOpsOperator() {
	VerifyNamespaceExists(ocloudparams.OpenshiftGitOpsNamespace, nil)
	VerifyCsvSuccessful(ocloudparams.OpenshiftGitOpsSubscriptionName, ocloudparams.OpenshiftGitOpsNamespace)
}

func VerifySiteConfigOperator() {
	// todo - need to find how to verify it
}

func VerifyOranO2ImsOperator() {
	VerifyNamespaceExists(ocloudparams.OCloudO2ImsNamespace, nil)
	VerifyCsvSuccessful(ocloudparams.OCloudO2ImsSubscriptionName, ocloudparams.OCloudO2ImsNamespace)
}

func VerifyOranHardwareManagerPluginOperator() {
	VerifyNamespaceExists(ocloudparams.OCloudHardwareManagerPluginNamespace, nil)
	VerifyCsvSuccessful(ocloudparams.OCloudHardwareManagerPluginSubscriptionName, 
		ocloudparams.OCloudHardwareManagerPluginNamespace)
}

