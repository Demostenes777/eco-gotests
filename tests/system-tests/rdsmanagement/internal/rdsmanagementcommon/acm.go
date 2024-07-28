package rdsmanagementcommon

import (
	"fmt"
	"time"

	"github.com/openshift-kni/eco-goinfra/pkg/deployment"
	"github.com/openshift-kni/eco-goinfra/pkg/pod"

	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/apiobjectshelper"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementinittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementparams"
)

var (
	acmOperatorDeployments = []string{"multicluster-operators", "management-ingress",
		"multiclusterhub", "klusterlet-addon-controller"}
)

func VerifyACMNamespace(ctx SpecContext) {
	// Verify namespace is exists
	err := apiobjectshelper.VerifyNamespaceExists(APIClient, rdsmanagementparams.OpenShiftVirtualizationNamespace, time.Second)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to pull %q namespace", rdsmanagementparams.ACMNameSpace))

}

func VerifyACMDeployment(ctx SpecContext) {
	// VerifyACMDeployment asserts ACM successfully installed.
	for _, operatorPod := range acmOperatorDeployments {
		glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof("Confirm that acm %s pod was deployed and running in %s namespace",
			operatorPod, rdsmanagementparams.ACMNameSpace)

		acmPods, err := pod.ListByNamePattern(APIClient, operatorPod, rdsmanagementparams.ACMNameSpace)
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("No %s pods were found in %s namespace; %v",
			operatorPod, rdsmanagementparams.ACMNameSpace, err))
		Expect(len(acmPods)).ToNot(Equal(0), fmt.Sprintf("No %s pods were found in %s namespace; %v",
			operatorPod, rdsmanagementparams.ACMNameSpace, err))

		acmPod := acmPods[0]
		acmPodName := acmPod.Object.Name

		err = acmPod.WaitUntilReady(time.Second)
		if err != nil {
			acmPodLog, _ := acmPod.GetLog(600*time.Second, operatorPod)
			glog.Fatalf("%s pod in %s namespace in a bad state: %s",
				acmPodName, rdsmanagementparams.ACMNameSpace, acmPodLog)

			for _, operatorDeployment := range acmOperatorDeployments {
				glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof("Confirm that %s deployment is running in %s namespace",
					operatorDeployment, rdsmanagementparams.ACMNameSpace)

				acmDeployment, err := deployment.Pull(APIClient, operatorDeployment, rdsmanagementparams.ACMNameSpace)
				Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("%s deployment not found in %s namespace; %v",
					operatorDeployment, rdsmanagementparams.ACMNameSpace, err))
				Expect(acmDeployment.IsReady(5*time.Second)).To(Equal(true),
					fmt.Sprintf("Bad state for %s deployment in %s namespace",
						operatorDeployment, rdsmanagementparams.ACMNameSpace))
			}
		}
	}
}
