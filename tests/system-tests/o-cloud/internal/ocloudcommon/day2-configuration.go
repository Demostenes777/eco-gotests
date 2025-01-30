package ocloudcommon

import (
	//"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	//. "github.com/onsi/gomega"
	//. "github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"
	//"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
)

func SuccessfulOperatorUpgrade() {
	Fail("Intentional failure for demonstration purposes J")
	// Deploy SNO 02 and SNO 03 
	// Once they are completed export their kufeconfig files to environment variables
	// Set SNO 02 and SNO 03 api client kubeconfig path
	// Verify CSV version
}

func FailedOperatorUpgrade() {
	Fail("Intentional failure for demonstration purposes K")
}

func FailedPartialOperatorUpgrade() {
	Fail("Intentional failure for demonstration purposes L")
}
