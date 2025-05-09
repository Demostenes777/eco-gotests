package rdsmanagementcommon

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hashicorp/go-version"
	"github.com/openshift-kni/eco-goinfra/pkg/clusteroperator"
	"github.com/openshift-kni/eco-goinfra/pkg/clusterversion"
	"github.com/openshift-kni/eco-goinfra/pkg/mco"
	"github.com/openshift-kni/eco-goinfra/pkg/nodes"
	"github.com/openshift-kni/eco-gotests/tests/internal/cluster"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementinittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementparams"
)

// VerifieAllNodesAreReady waits for all the nodes in the cluster to report Ready state.
func VerifieAllNodesAreReady(ctx SpecContext) {
	By("Checking all nodes are Ready")

	Eventually(func(ctx SpecContext) bool {
		allNodes, err := nodes.List(APIClient, metav1.ListOptions{})
		if err != nil {
			glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof("Failed to list all nodes: %s", err)

			return false
		}

		for _, _node := range allNodes {
			glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof("Processing node %q", _node.Definition.Name)

			for _, condition := range _node.Object.Status.Conditions {
				if condition.Type == rdsmanagementparams.ConditionTypeReadyString {
					if condition.Status != rdsmanagementparams.ConstantTrueString {
						glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof("Node %q is notReady", _node.Definition.Name)
						glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof("  Reason: %s", condition.Reason)

						return false
					}
				}
			}
		}

		return true
	}).WithTimeout(25*time.Minute).WithPolling(15*time.Second).WithContext(ctx).Should(BeTrue(),
		"Some nodes are notReady")
}

// VerifyKubeletResourceReservationHasBeenIncreased assert system reserved memory for masters succeeded.
func VerifyKubeletResourceReservationHasBeenIncreased(ctx SpecContext) {
	glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof("Verify system reserved memory config for masters succeeded")

	systemReservedBuilder := mco.NewKubeletConfigBuilder(APIClient, rdsmanagementparams.KubeletConfigName).
		WithMCPoolSelector("pools.operator.machineconfiguration.openshift.io/master", "").
		WithSystemReserved(rdsmanagementparams.SystemReservedCPU, rdsmanagementparams.SystemReservedMemory)

	if !systemReservedBuilder.Exists() {
		glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof("Create system-reserved configuration")

		systemReserved, err := systemReservedBuilder.Create()
		Expect(err).ToNot(HaveOccurred(), "Failed to create %s kubeletConfig objects "+
			"with system-reserved definition", rdsmanagementparams.KubeletConfigName)

		Expect(systemReserved.Exists()).To(Equal(true),
			"Failed to setup master system reserved memory, %s kubeletConfig not found; %s",
			rdsmanagementparams.KubeletConfigName, err)
	}
}
func VerifyClusterIsOperational(ctx SpecContext) {
	glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof("Checking that the clusterversion is available")

	clusterVersion, err := cluster.GetOCPClusterVersion(APIClient)
	Expect(err).ToNot(HaveOccurred(), "error detecting clusterversion")
	ocpVersion, _ := version.NewVersion(clusterVersion.Definition.Status.Desired.Version)
	glog.V(rdsmanagementparams.RdsManagementLogLevel).Infof("Cluster Version: %s", ocpVersion)

	_, err = clusterversion.Pull(APIClient)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Error accessing csv: %v", err))

	var coBuilder []*clusteroperator.Builder
	coBuilder, err = clusteroperator.List(APIClient)
	Expect(err).To(BeNil(), fmt.Sprintf("ClusterOperator List not found: %v", err))
	Expect(len(coBuilder)).ToNot(Equal(0), "Empty clusterOperators list received")

	_, err = clusteroperator.WaitForAllClusteroperatorsAvailable(APIClient, 60*time.Second)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Error waiting for all available clusteroperators: %v", err))

}
