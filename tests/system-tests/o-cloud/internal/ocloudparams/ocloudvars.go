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

	// PrName1 is the name of the first provisioning request
	PrName1 = "provisioning-request-1"

	// PrName2 is the name of the second provisioning request
	PrName2 = "provisioning-request-2"
	
	// TemplateVersion1 defines the version of the referenced ClusterTemplate used for the successful SNO provisioning
	TemplateVersion1 = "v4-18-0-ec3-1"

	// TemplateVersion2 defines the version of the referenced ClusterTemplate used for the failing SNO provisioning
	TemplateVersion2 = "v4-18-0-ec3-2"

	//nolint:lll
	// TemplateVersion3 defines the version of the referenced ClusterTemplate used for the multicluster provisioning with different templates
	TemplateVersion3 = "v4-18-0-ec3-3"

	// NodeClusterName1 is the name of the first ORAN Node Cluster
	NodeClusterName1 = "sno02"

	// NodeClusterName2 is the name of the second ORAN Node Cluster
	NodeClusterName2 = "sno03"

	// OCloudSiteId1 is the ID of the of the first ORAN O-Cloud Site
	OCloudSiteId1 = "sno02"

	// OCloudSiteId2 is the ID of the of the second ORAN O-Cloud Site
	OCloudSiteId2 = "sno03"
	
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
)
