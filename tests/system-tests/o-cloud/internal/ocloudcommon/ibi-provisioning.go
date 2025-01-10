package ocloudcommon

import (
	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
)

func VerifySuccessfulIbiSnoProvisioning() {
	It("should fail with a message", func() {
        Fail("Intentional failure for demonstration purposes H")
    })
}

func VerifyFailedIbiSnoProvisioning() {
	It("should fail with a message", func() {
        Fail("Intentional failure for demonstration purposes I")
    })
}