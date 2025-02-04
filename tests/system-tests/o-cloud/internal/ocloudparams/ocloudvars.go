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

	// TemplateName defines the base name of the referenced ClusterTemplate.
	TemplateName = "sno-ran-du"


	PolicyTemplateParameters = map[string]any{}

	// PrSingleSnoAiSuccess is the name of the first provisioning request
	PrSingleSnoAiSuccess = "single-sno-ai-success"

	// PrSingleSnoAiFailure is the name of the first provisioning request
	PrSingleSnoAiFailure = "single-sno-ai-failure"

	// PrSingleSnoIbiSuccess is the name of the first provisioning request
	PrSingleSnoIbiSuccess = "single-sno-ibi-success"

	// PrSingleSnoIbiFailure is the name of the first provisioning request
	PrSingleSnoIbiFailure = "single-sno-ibi-failure"

	// PrMultipleSnoSameCTSno02 is the name of the second provisioning request
	PrMultipleSnoSameCTSno02 = "multiple-sno-same-cluster-template-sno02"

	// PrMultipleSnoSameCTSno03 is the name of the second provisioning request
	PrMultipleSnoSameCTSno03 = "multiple-sno-same-cluster-template-sno03"

	// PrMultipleSnoDifferentCTSno02 is the name of the second provisioning request
	PrMultipleSnoDifferentCTSno02 = "multiple-sno-different-cluster-template-sno02"

	// PrMultipleSnoDifferentCTSno03 is the name of the second provisioning request
	PrMultipleSnoDifferentCTSno03 = "multiple-sno-different-cluster-template-sno03"

	// PrMultipleDay2Sno02 is the name of the provisioning request for Day 2 operations in SNO 02
	PrMultipleDay2Sno02 = "multiple-day2-sno02"

	// PrMultipleDay2Sno03 is the name of the provisioning request for Day 2 operations in SNO 02
	PrMultipleDay2Sno03 = "multiple-day2-sno03"
	
	// TemplateVersion1 defines the version of the referenced ClusterTemplate used for the successful SNO provisioning using AI
	TemplateVersion1 = "v4-18-0-ec3-1"

	// TemplateVersion2 defines the version of the referenced ClusterTemplate used for the failing SNO provisioning using AI
	TemplateVersion2 = "v4-18-0-ec3-2"

	//nolint:lll
	// TemplateVersion3 defines the version of the referenced ClusterTemplate used for the multicluster provisioning with different templates
	TemplateVersion3 = "v4-18-0-ec3-3"

	//nolint:lll
	// TemplateVersion4 defines the version of the referenced ClusterTemplate used for the successful SNO provisioning using IBI
	TemplateVersion4 = "v4-18-0-ec3-4"

	//nolint:lll
	// TemplateVersion5 defines the version of the referenced ClusterTemplate used for the failing SNO provisioning using IBI
	TemplateVersion5 = "v4-18-0-ec3-5"

	//nolint:lll
	// TemplateVersion6 defines the version of the referenced ClusterTemplate used for the Day 2 operations
	TemplateVersion6 = "v4-18-0-ec3-6"

	// NodeClusterName1 is the name of the first ORAN Node Cluster
	NodeClusterName1 = "sno02"

	// NodeClusterName2 is the name of the second ORAN Node Cluster
	NodeClusterName2 = "sno03"

	// OCloudSiteId1 is the ID of the of the first ORAN O-Cloud Site
	OCloudSiteId1 = "sno02"

	// OCloudSiteId2 is the ID of the of the second ORAN O-Cloud Site
	OCloudSiteId2 = "sno03"

	// OCloudSiteId2 is the ID of the of the second ORAN O-Cloud Site
	HostName2 = "sno03.oran.telcoqe.eng.rdu2.dc.redhat.com"
	
	// ClusterInstanceParameters1 is the map with the cluster instance parameters for the first cluster
	ClusterInstanceParameters1 = map[string]any{
		"clusterName": "sno02",
		"nodes": []map[string]any{
			{
				"hostName": "sno02.oran.telcoqe.eng.rdu2.dc.redhat.comf",
			},
		},
	}

	// ClusterInstanceParameters2 is the map with the cluster instance parameters for the second cluster
	ClusterInstanceParameters2 = map[string]any{
		"clusterName": "sno03",
		"nodes": []map[string]any{
			{
				"hostName": "sno03.oran.telcoqe.eng.rdu2.dc.redhat.com",
				"nodeNetwork": map[string]any {
					"config": map[string]any {
						"interfaces": []map[string]any {
							{
								"name": "ens3f3",
								"type": "ethertype",
								"state": "up",
								"ipv6": map[string]any {
									"enabled": "true",
									"address": []map[string]any {
										{
											"ip": "2620:52:9:1698::6",
											"prefix-length": "64",
										},
									},
									"dhcp": "false",
									"autoconf": "false",
								},
								"ipv4": map[string]any {
									"enabled": "false",
								},
							},
						},
					},
				},
			},
		},
	}

	PTPVersionMajorOld uint64 = 4
	PTPVersionMinorOld uint64 = 18
	PTPVersionPatchOld uint64 = 0
	PTPVersionPrereleaseOld uint64 = 202411190136

	PTPVersionMajorNew uint64 = 4
	PTPVersionMinorNew uint64 = 18
	PTPVersionPatchNew uint64 = 0
	PTPVersionPrereleaseNew uint64 = 202412042342
)
