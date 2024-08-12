package rdsmanagementcommon

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/apiobjectshelper"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementinittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/rdsmanagement/internal/rdsmanagementparams"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func VerifyArgoCD(ctx SpecContext) {
	err := apiobjectshelper.VerifyNamespaceExists(APIClient, rdsmanagementparams.ArgoCDNamespace, time.Second)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to pull %q namespace", rdsmanagementparams.ArgoCDNamespace))
	// Verify all pods are running
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=argocd",
	})

	Expect(err).ToNot(HaveOccurred(), "failed to list pods")
	for _, pod := range pods.Items {
		Expect(pod.Status.Phase).To(Equal(corev1.PodRunning), "pod %s is not running", pod.Name)
		// Check CSV exists
		csv, err := clientset.AppsV1().ControllerRevisions(rdsmanagementparams.ArgoCDNamespace).Get(context.TODO(), csvName, metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred(), "failed to get CSV %s", csvName)
		Expect(csv).ToNot(BeNil(), "CSV %s not found", csvName)

	}
}
