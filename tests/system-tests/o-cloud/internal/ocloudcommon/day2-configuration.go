package ocloudcommon

import (
	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/o-cloud/internal/ocloudparams"
)

func SuccessfulOperatorUpgrade() {
	It("should fail with a message", func() {
        Fail("Intentional failure for demonstration purposes J")
    })
}

func FailedOperatorUpgrade() {
	It("should fail with a message", func() {
        Fail("Intentional failure for demonstration purposes K")
    })
}

func FailedPartialOperatorUpgrade() {
	It("should fail with a message", func() {
        Fail("Intentional failure for demonstration purposes L")
    })
}
