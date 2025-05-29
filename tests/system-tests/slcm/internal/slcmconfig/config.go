package slcmconfig

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kelseyhightower/envconfig"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/systemtestsconfig"
	"gopkg.in/yaml.v2"
)

const (
	// PathToDefaultSLCMParamsFile path to config file with SLCM parameters.
	PathToDefaultSLCMParamsFile = "./default.yaml"
)

// SLCMconfig type keeps the SLCM configuration.
type SLCMConfig struct {
	*systemtestsconfig.SystemTestsConfig

	// TargetKubeconfigPath is the path to the kubeconfig file.
	TargetKubeconfigPath string `yaml:"target_kubeconfig_path" envconfig:"SLCM_TARGET_KUBECONFIG_PATH"`

	// ArtefactOutputDir is the directory where artifacts will be output.
	ArtefactOutputDir string `yaml:"artefact_output_dir" envconfig:"SLCM_ARTEFACT_OUTPUT_DIR"`

	// FECNumAqsPerGroups is the number of AQs per group for FEC.
	FECNumAqsPerGroups string `yaml:"fec_num_aqs_per_groups" envconfig:"SLCM_FEC_NUM_AQS_PER_GROUPS"`

	// FECMaxQueueSize is the maximum queue size for FEC.
	FECMaxQueueSize string `yaml:"fec_max_queue_size" envconfig:"SLCM_FEC_MAX_QUEUE_SIZE"`

	// FECTDevice is the FEC device.
	FECTDevice string `yaml:"fec_device" envconfig:"SLCM_FEC_DEVICE"`

	// PAOPageCount is the page count for PAO.
	PAOPageCount string `yaml:"pao_page_count" envconfig:"SLCM_PAO_PAGE_COUNT"`

	// ReleaseVersion is the release version of OCP.
	ReleaseVersion string `yaml:"release_version" envconfig:"SLCM_RELEASE_VERSION"`

	// TALMNamespace is the namespace for the TALM operator.
	TALMNamespace string `yaml:"talm_namespace" envconfig:"SLCM_TALM_NAMESPACE"`

	// TALMNumPods is the number of pods for the TALM operator.
	TALMNumPods int `yaml:"talm_num_pods" envconfig:"SLCM_TALM_NUM_PODS"`

	// PCIDevice is the PCI device.
	PCIDevice string `yaml:"pci_device" envconfig:"SLCM_PCI_DEVICE"`

	// SRIOVVFSCount is the count of SR-IOV VFs.
	SRIOVVFSCount int `yaml:"sriov_vfs_count" envconfig:"SLCM_SRIOV_VFS_COUNT"`

	// PTPInterface is the PTP interface.
	PTPInterface string `yaml:"ptp_interface" envconfig:"SLCM_PTP_INTERFACE"`

	// ConfigPolicy is the configuration policy.
	ConfigPolicy string `yaml:"config_policy" envconfig:"SLCM_CONFIG_POLICY"`

	// TestBBDevOutputFile is the output file for test BBDev.
	TestBBDevOutputFile string `yaml:"test_bbdev_output_file" envconfig:"SLCM_TEST_BBDEV_OUTPUT_FILE"`

	// SiteName is the name of the site.
	SiteName string `yaml:"site_name" envconfig:"SLCM_SITE_NAME"`

	// SSHPort is the SSH port.
	SSHPort int `yaml:"ssh_port" envconfig:"SLCM_SSH_PORT"`

	// OCPIApiURL is the OCP API URL.
	OCPIApiURL string `yaml:"ocp_api_url" envconfig:"SLCM_OCP_API_URL"`

	// StabilityTestDuration is the duration of the stability test.
	StabilityTestDuration string `yaml:"stability_test_duration" envconfig:"SLCM_STABILITY_TEST_DURATION"`

	// CSVFile is the path to the CSV file.
	CSVFile string `yaml:"csv_file" envconfig:"SLCM_CSV_FILE"`

	// TrafficInjectionEnabled is a flag indicating if traffic injection is enabled.
	TrafficInjectionEnabled bool `yaml:"traffic_injection_enabled" envconfig:"SLCM_TRAFFIC_INJECTION_ENABLED"`

	// DUPolicyEnabled is a flag indicating if DU policy is enabled.
	DUPolicyEnabled bool `yaml:"du_policy_enabled" envconfig:"SLCM_DU_POLICY_ENABLED"`

	// PTPEnabled is a flag indicating if PTP is enabled.
	PTPEnabled bool `yaml:"ptp_enabled" envconfig:"SLCM_PTP_ENABLED"`

	// PTPLogTime is the log time for PTP.
	PTPLogTime string `yaml:"ptp_log_time" envconfig:"SLCM_PTP_LOG_TIME"`

	// PTPCfgFilePath is the path to the PTP configuration file.
	PTPCfgFilePath string `yaml:"ptp_cfg_file_path" envconfig:"SLCM_PTP_CFG_FILE_PATH"`

	// DUPTPOCConfigName is the name of the DU PTP OC configuration.
	DUPTPOCConfigName string `yaml:"du_ptp_oc_config_name" envconfig:"SLCM_DU_PTP_OC_CONFIG_NAME"`

	// DUTPPTBCConfigName is the name of the DU PTP BC configuration.
	DUTPPTBCConfigName string `yaml:"du_ptp_bc_config_name" envconfig:"SLCM_DU_PTP_BC_CONFIG_NAME"`

	// PTPJumpResultPass is a flag indicating if the PTP jump result is a pass.
	PTPJumpResultPass bool `yaml:"ptp_jump_result_pass" envconfig:"SLCM_PTP_JUMP_RESULT_PASS"`

	// PTPTimedoutResultPass is a flag indicating if the PTP timed out result is a pass.
	PTPTimedoutResultPass bool `yaml:"ptp_timedout_result_pass" envconfig:"SLCM_PTP_TIMEDOUT_RESULT_PASS"`

	// PTPEventsHealthPass is a flag indicating if the PTP events health is a pass.
	PTPEventsHealthPass bool `yaml:"ptp_events_health_pass" envconfig:"SLCM_PTP_EVENTS_HEALTH_PASS"`

	// PTPEventsFreerunPass is a flag indicating if the PTP events freerun is a pass.
	PTPEventsFreerunPass bool `yaml:"ptp_events_freerun_pass" envconfig:"SLCM_PTP_EVENTS_FREERUN_PASS"`

	// PTPEventsErrorPass is a flag indicating if the PTP events error is a pass.
	PTPEventsErrorPass bool `yaml:"ptp_events_error_pass" envconfig:"SLCM_PTP_EVENTS_ERROR_PASS"`

	// PTPEventsGetCurrentStatePass is a flag indicating if the PTP events get current state is a pass.
	PTPEventsGetCurrentStatePass bool `yaml:"ptp_events_getcurrentstate_pass" envconfig:"SLCM_PTP_EVENTS_GETCURRENTSTATE_PASS"`

	// DUStaticPVCS is the list of static PVCs for DU.
	DUStaticPVCS []string `yaml:"du_static_pvcs" envconfig:"SLCM_DU_STATIC_PVCS"`

	// DUNumPVs is the number of PVs for DU.
	DUNumPVs int `yaml:"du_num_pvs" envconfig:"SLCM_DU_NUM_PVS"`

	// DUExpectedBoundPVs is the expected number of bound PVs for DU.
	DUExpectedBoundPVs int `yaml:"du_expected_bound_pvs" envconfig:"SLCM_DU_EXPECTED_BOUND_PVS"`

	// DUNumPods is the number of pods for DU.
	DUNumPods int `yaml:"du_num_pods" envconfig:"SLCM_DU_NUM_PODS"`

	// DUNamespace is the namespace for DU.
	DUNamespace string `yaml:"du_namespace" envconfig:"SLCM_DU_NAMESPACE"`

	// DUSRIOVPods is the list of SR-IOV pods for DU.
	DUSRIOVPods []string `yaml:"du_sriov_pods" envconfig:"SLCM_DU_SRIOV_PODS"`

	// DUSRIOVNetworksName is the list of SR-IOV network names for DU.
	DUSRIOVNetworksName []string `yaml:"du_sriov_networks_name" envconfig:"SLCM_DU_SRIOV_NETWORKS_NAME"`

	// DatamodelPath is the path to the datamodel.
	DatamodelPath string `yaml:"datamodel_path" envconfig:"SLCM_DATAMODEL_PATH"`

	// SpirentVM is the hostname of the Spirent VM.
	SpirentVM string `yaml:"spirent_vm" envconfig:"SLCM_SPIRENT_VM"`

	// SpirentChassis is the hostname of the Spirent chassis.
	SpirentChassis string `yaml:"spirent_chassis" envconfig:"SLCM_SPIRENT_CHASSIS"`

	// SpirentPorts is the list of Spirent ports.
	SpirentPorts []int `yaml:"spirent_ports" envconfig:"SLCM_SPIRENT_PORTS"`

	// MaxPacketLoss is the maximum packet loss.
	MaxPacketLoss int `yaml:"max_packet_loss" envconfig:"SLCM_MAX_PACKET_LOSS"`

	// SwitchPingResult is the result of the switch ping.
	SwitchPingResult bool `yaml:"switch_ping_result" envconfig:"SLCM_SWITCH_PING_RESULT"`

	// LinkStatus is the status of the link.
	LinkStatus bool `yaml:"link_status" envconfig:"SLCM_LINK_STATUS"`

	// DeletedSitesList is the list of deleted sites.
	DeletedSitesList []string `yaml:"deleted_sites_list" envconfig:"SLCM_DELETED_SITES_LIST"`

	// SitesList is the list of sites.
	SitesList []string `yaml:"sites_list" envconfig:"SLCM_SITES_LIST"`

	// TotalPoliciesCount is the total count of policies.
	TotalPoliciesCount int `yaml:"total_policies_count" envconfig:"SLCM_TOTAL_POLICIES_COUNT"`

	// OCPVersion is the version of OCP.
	OCPVersion string `yaml:"ocp_version" envconfig:"SLCM_OCP_VERSION"`

	// BMHNamespace is the namespace for BMH.
	BMHNamespace string `yaml:"bmh_namespace" envconfig:"SLCM_BMH_NAMESPACE"`

	// BMHExpectedNumber is the expected number of BMHs.
	BMHExpectedNumber int `yaml:"bmh_expected_number" envconfig:"SLCM_BMH_EXPECTED_NUMBER"`

	// TREXServerAddress is the address of the TREX server.
	TREXServerAddress string `yaml:"trex_server_address" envconfig:"SLCM_TREX_SERVER_ADDRESS"`

	// TREXServerPort is the port of the TREX server.
	TREXServerPort int `yaml:"trex_server_port" envconfig:"SLCM_TREX_SERVER_PORT"`

	// MaxAvgLatency is the maximum average latency.
	MaxAvgLatency int `yaml:"max_avg_latency" envconfig:"SLCM_MAX_AVG_LATENCY"`

	// InjectionDuration is the duration of the injection.
	InjectionDuration float64 `yaml:"injection_duration" envconfig:"SLCM_INJECTION_DURATION"`

	// ExpectFailure is a flag indicating if failure is expected.
	ExpectFailure bool `yaml:"expect_failure" envconfig:"SLCM_EXPECT_FAILURE"`

	// MinPacketLoss is the minimum packet loss.
	MinPacketLoss float64 `yaml:"min_packet_loss" envconfig:"SLCM_MIN_PACKET_LOSS"`

	// TrafficInjector is the traffic injector.
	TrafficInjector string `yaml:"traffic_injector" envconfig:"SLCM_TRAFFIC_INJECTOR"`

	// CUMasterHostnames is the list of CU master hostnames.
	CUMasterHostnames []string `yaml:"cu_master_hostnames" envconfig:"SLCM_CU_MASTER_HOSTNAMES"`

	// HelmChartProfile is the Helm chart profile.
	HelmChartProfile string `yaml:"helm_chart_profile" envconfig:"SLCM_HELM_CHART_PROFILE"`
}

// NewSLCMConfig returns instance of SLCMConfig config type.
func NewSLCMConfig() *SLCMConfig {
	log.Print("Creating new SLCMConfig struct")

	var slcmConf SLCMConfig
	slcmConf.SystemTestsConfig = systemtestsconfig.NewSystemTestsConfig()

	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	confFile := filepath.Join(baseDir, PathToDefaultSLCMParamsFile)
	err := readFile(&slcmConf, confFile)

	if err != nil {
		log.Printf("Error to read config file %s", confFile)

		return nil
	}

	err = readEnv(&slcmConf)

	if err != nil {
		log.Print("Error to read environment variables")

		return nil
	}

	return &slcmConf
}

func readFile(slcmConfig *SLCMConfig, cfgFile string) error {
	openedCfgFile, err := os.Open(cfgFile)
	if err != nil {
		return err
	}

	defer func() {
		_ = openedCfgFile.Close()
	}()

	decoder := yaml.NewDecoder(openedCfgFile)
	err = decoder.Decode(&slcmConfig)

	if err != nil {
		return err
	}

	return nil
}

func readEnv(slcmConfig *SLCMConfig) error {
	err := envconfig.Process("", slcmConfig)
	if err != nil {
		return err
	}

	return nil
}
