package ocloudparams

import (
	provisioningv1alpha1 "github.com/openshift-kni/oran-o2ims/api/provisioning/v1alpha1"
	hardwaremanagementv1alpha1 "github.com/openshift-kni/oran-o2ims/api/hardwaremanagement/v1alpha1"


	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/systemtestsparams"
	"github.com/openshift-kni/k8sreporter"
)

var (
	// Labels represents the range of labels that can be used for test cases selection.
	Labels = []string{systemtestsparams.Label, Label}

	// ReporterNamespacesToDump tells to the reporter from where to collect logs.
	ReporterNamespacesToDump = map[string]string{
		"oran-hwmgr-plugin":		"oran-hwmgr-plugin",
		"oran-o2ims": 				"oran-o2ims",
	}

	// ReporterCRDsToDump tells to the reporter what CRs to dump.
	ReporterCRDsToDump = []k8sreporter.CRData{
		{Cr: &provisioningv1alpha1.ClusterTemplateList{}},
		{Cr: &provisioningv1alpha1.ProvisioningRequestList{}},
		{Cr: &hardwaremanagementv1alpha1.HardwareTemplateList{}},
		{Cr: &hardwaremanagementv1alpha1.NodeList{}},
		{Cr: &hardwaremanagementv1alpha1.NodePoolList{}},
	}

	// TestNamespaceName is used for defining the namespace name where test resources are created.
	TestNamespaceName = "o-cloud-system-tests"
)
