package slcmcommon

import (
	"fmt"
	"time"

	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/pod"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/slcm/internal/slcminittools"
)

// TestDUPodsCount verifies the number of DU pods matches the expected count.
func TestDUPodsCount() {
	Eventually(func() (int, error) {
		pods, err := pod.List(APIClient, SLCMConfig.DUNamespace)
		if err != nil {
			return 0, err
		}
		return len(pods), nil
	}).WithTimeout(1500*time.Second).WithPolling(10*time.Second).Should(Equal(SLCMConfig.DUNumPods),
		fmt.Sprintf("The number of DU pods does not match the expected count"))
}

// TestDUPodsStatus verifies all DU pods are in the Running state.
func TestDUPodsStatus() {
	Eventually(func() (bool, error) {
		pods, err := pod.List(APIClient, SLCMConfig.DUNamespace)
		if err != nil {
			return false, err
		}
		if len(pods) == 0 {
			return false, fmt.Errorf("no pods found in namespace %s", SLCMConfig.DUNamespace)
		}

		var failedPods []string
		for _, p := range pods {
			if p.Object.Status.Phase != "Running" {
				failedPods = append(failedPods, p.Object.Name)
			}
			for _, condition := range p.Object.Status.Conditions {
				if condition.Type == "Ready" && condition.Status != "True" {
					failedPods = append(failedPods, p.Object.Name)
				}
			}
		}
		return len(failedPods) == 0, nil
	}).WithTimeout(600*time.Second).WithPolling(10*time.Second).Should(BeTrue(),
		fmt.Sprintf("There are failed pods"))
}
