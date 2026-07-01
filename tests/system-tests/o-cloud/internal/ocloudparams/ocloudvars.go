package ocloudparams

import (
	"fmt"
	"strings"

	hardwaremanagementv1alpha1 "github.com/openshift-kni/oran-o2ims/api/hardwaremanagement/v1alpha1"
	provisioningv1alpha1 "github.com/openshift-kni/oran-o2ims/api/provisioning/v1alpha1"

	. "github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/o-cloud/internal/ocloudinittools"

	"github.com/openshift-kni/k8sreporter"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/system-tests/internal/systemtestsparams"
)

var (
	// Labels represents the range of labels that can be used for test cases selection.
	Labels = []string{systemtestsparams.Label, Label}

	// ReporterNamespacesToDump tells to the reporter from where to collect logs.
	ReporterNamespacesToDump = map[string]string{
		"oran-hwmgr-plugin": "oran-hwmgr-plugin",
		OranO2ImsNamespace:  OranO2ImsNamespace,
	}

	// ReporterCRDsToDump tells to the reporter what CRs to dump.
	ReporterCRDsToDump = []k8sreporter.CRData{
		{Cr: &provisioningv1alpha1.ClusterTemplateList{}},
		{Cr: &provisioningv1alpha1.ProvisioningRequestList{}},
		{Cr: &hardwaremanagementv1alpha1.AllocatedNodeList{}},
		{Cr: &hardwaremanagementv1alpha1.NodeAllocationRequestList{}},
	}

	// PolicyTemplateParameters defines the policy template parameters.
	PolicyTemplateParameters = map[string]any{}

	// ClusterInstanceParameters1 is the map with the cluster instance parameters for the first cluster.
	ClusterInstanceParameters1 = map[string]any{
		"clusterName": OCloudConfig.ClusterName1,
		"nodes": []map[string]any{
			{
				"hostName": OCloudConfig.HostName1,
				"nodeNetwork": map[string]any{
					"config": map[string]any{
						"interfaces": []map[string]any{
							{
								"ipv6": map[string]any{
									"address": []map[string]any{
										{
											"ip":            OCloudConfig.InterfaceIpv6_1,
											"prefix-length": 64,
										},
									},
								},
							},
						},
						"dns-resolver": map[string]any{
							"config": map[string]any{
								"server": []string{OCloudConfig.DNSIpv6},
							},
						},
						"routes": map[string]any{
							"config": []map[string]any{
								{
									"destination":        "::/0",
									"next-hop-interface": OCloudConfig.NextHopInterface,
									"next-hop-address":   OCloudConfig.NextHopIpv6,
								},
							},
						},
					},
				},
			},
		},
	}

	// ClusterInstanceParameters2 is the map with the cluster instance parameters for the second cluster.
	ClusterInstanceParameters2 = map[string]any{
		"clusterName": OCloudConfig.ClusterName2,
		"nodes": []map[string]any{
			{
				"hostName": OCloudConfig.HostName2,
				"nodeNetwork": map[string]any{
					"config": map[string]any{
						"interfaces": []map[string]any{
							{
								"ipv6": map[string]any{
									"address": []map[string]any{
										{
											"ip":            OCloudConfig.InterfaceIpv6_2,
											"prefix-length": 64,
										},
									},
								},
							},
						},
						"dns-resolver": map[string]any{
							"config": map[string]any{
								"server": []string{OCloudConfig.DNSIpv6},
							},
						},
						"routes": map[string]any{
							"config": []map[string]any{
								{
									"destination":        "::/0",
									"next-hop-interface": OCloudConfig.NextHopInterface,
									"next-hop-address":   OCloudConfig.NextHopIpv6,
								},
							},
						},
					},
				},
			},
		},
	}

	// skopeoRedhatOperatorsTemplate is the command template for tagging redhat-operators catalog images.
	skopeoRedhatOperatorsTemplate = "skopeo copy --authfile %s --tls-verify=false" +
		" docker://%s/olm/redhat-operators:v%s-%s docker://%s/olm/redhat-operators:v%s-day2"
	//nolint:lll
	// SnoKubeconfigCreate command to get the SNO kubeconfig file.
	SnoKubeconfigCreate = "oc -n %s get secret %s-admin-kubeconfig -o json | jq -r .data.kubeconfig | base64 -d > tmp/%s/auth/kubeconfig"
	//nolint:lll
	// CreateImageBasedInstallationConfig command to create the image based installation configuration template.
	CreateImageBasedInstallationConfig = "openshift-install image-based create image-config-template --dir tmp/ibi-iso-workdir"
	// CreateIsoImage command to create the ISO image.
	CreateIsoImage = "openshift-install image-based create image --dir tmp/ibi-iso-workdir"
	//nolint:lll
	// CheckIbiCompleted command to check that the image based installation has finished.
	CheckIbiCompleted = "journalctl -u install-rhcos-and-restore-seed.service | grep 'Finished SNO Image-based Installation.'"

	// SpokeSSHUser ssh user of the spoke cluster.
	SpokeSSHUser = "core"
	// SpokeSSHPasskeyPath path to the ssh key of the spoke cluster.
	SpokeSSHPasskeyPath = "/opt/id_rsa"
	// SeedGeneratorName name of the seedgenerator CR.
	SeedGeneratorName = "seedimage"
	// RegistryCertPath path to the registry certificate.
	RegistryCertPath = "/opt/registry.crt"
	// IbiConfigTemplate template for the image based installation configuration.
	IbiConfigTemplate = "/opt/ibi-config.yaml.tmpl"
	// IbiConfigTemplateYaml path to the YAML file with the image based installation configuration.
	IbiConfigTemplateYaml = "tmp/ibi-iso-workdir/image-based-installation-config.yaml"
	// IbiBasedImageSourcePath path to the base image.
	IbiBasedImageSourcePath = "tmp/ibi-iso-workdir/rhcos-ibi.iso"

	// PtpCPURequest is cpu request for the PTP container.
	PtpCPURequest = "50m"
	// PtpMemoryRequest is cpu request for the PTP container.
	PtpMemoryRequest = "100Mi"
	// PtpCPULimit is cpu limit for the PTP container.
	PtpCPULimit = "1m"
	// PtpMemoryLimit is cpu limit for the PTP container.
	PtpMemoryLimit = "1Mi"
)

// BuildSkopeoRedhatOperatorsUpgradeCmd constructs the skopeo command for upgrading redhat-operators catalog.
func BuildSkopeoRedhatOperatorsUpgradeCmd(authfilePath, registry, ocpVersion string) string {
	majorMinor := extractMajorMinor(ocpVersion)

	return fmt.Sprintf(skopeoRedhatOperatorsTemplate,
		authfilePath, registry, majorMinor, "new", registry, majorMinor)
}

// BuildSkopeoRedhatOperatorsDowngradeCmd constructs the skopeo command for downgrading redhat-operators catalog.
func BuildSkopeoRedhatOperatorsDowngradeCmd(authfilePath, registry, ocpVersion string) string {
	majorMinor := extractMajorMinor(ocpVersion)

	return fmt.Sprintf(skopeoRedhatOperatorsTemplate,
		authfilePath, registry, majorMinor, "old", registry, majorMinor)
}

func extractMajorMinor(version string) string {
	parts := strings.SplitN(version, ".", 3)
	if len(parts) < 2 {
		return version
	}

	return parts[0] + "." + parts[1]
}
