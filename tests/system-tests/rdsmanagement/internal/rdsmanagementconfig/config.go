package rdsmanagementconfig

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kelseyhightower/envconfig"
	"github.com/openshift-kni/eco-gotests/tests/internal/config"

	"gopkg.in/yaml.v2"
)

const (
	// PathToDefaultRDSManagementParamsFile path to config file with default RDSManagement parameters.
	PathToDefaultRDSManagementParamsFile = "./default.yaml"
)

// IDMDetails structure to hold BMC details.
type IDMDetails struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	IPAddress     string `json:"ip"`
	FQDN          string `json:"fqdn"`
	IPAAdminPass  string `json:"ipaadminpass"`
	IPADirMgrPass string `json:"ipadirmgrpass"`
}

// ManagementConfig type keeps RDS Management configuration.
type ManagementConfig struct {
	*config.GeneralConfig
	// ClusterVersion is OCP version.
	ClusterVersion string `yaml:"rdsmanagement_cluster_version" envconfig:"ECO_RDSMANAGEMENT_CLUSTER_VERSION"`
	// MasterNodesCount is number of master nodes and master/worker nodes.
	MasterNodesCount int `yaml:"rdsmanagement_master_nodes_count" envconfig:"ECO_RDSMANAGEMENT_MASTER_NODES_COUNT"`
	// WorkerNodesCount is number of worker nodes.
	WorkerNodesCount int `yaml:"rdsmanagement_worker_nodes_count" envconfig:"ECO_RDSMANAGEMENT_WORKER_NODES_COUNT"`

	// AppsNS is the namespace where the applications are installed.
	AppsNS string `yaml:"rdsmanagement_apps_ns" envconfig:"ECO_RDSMANAGEMENT_APPS_NS"`
	// PerformanceAddonNamespace is the namespace of the Performance Addon operator.
	PerformanceAddonNS string `yaml:"rdsmanagement_performance_addon_ns" envconfig:"ECO_RDSMANAGEMENT_PERFORMANCE_ADDON_NS"`
	// OpenshiftVirtualizationNamespace is the namespace of the OpenShift Virtualization operator.
	OpenshiftVirtualizationNS string `yaml:"rdsmanagement_openshift_virtualization_ns" envconfig:"ECO_RDSMANAGEMENT_OPENSHIFT_VIRTUALIZATION_NS"`
	// QuayNamespace is the namespace of Quay.
	QuayNS string `yaml:"rdsmanagement_quay_ns" envconfig:"ECO_RDSMANAGEMENT_QUAY_NS"`
	// MetalLBNamespace is the namespace of MetalLB.
	MetalLBNS string `yaml:"rdsmanagement_metallb_ns" envconfig:"ECO_RDSMANAGEMENT_METALLB_NS"`
	// ACMNamespace is the namespace of ACM.
	AcmNS string `yaml:"rdsmanagement_acm_ns" envconfig:"ECO_RDSMANAGEMENT_ACM_NS"`
	// KafkaNamespace is the namespace of Kafka.
	KafkaNS string `yaml:"rdsmanagement_kafka_ns" envconfig:"ECO_RDSMANAGEMENT_KAFKA_NS"`
	// KafkaAdapterNS is the namespace of Kafka.
	KafkaAdapterNS string `yaml:"rdsmanagement_kafka_adapter_ns" envconfig:"ECO_RDSMANAGEMENT_KAFKA_NS"`
	// AnsibleNS is the namespace of Ansible Automation Platform.
	AnsibleNS string `yaml:"rdsmanagement_ansible_ns" envconfig:"ECO_RDSMANAGEMENT_ANSIBLE_NS"`
	// Amq7NS is the name of the namespace of AMQ7.
	Amq7NS string `yaml:"rdsmanagement_amq7_ns" envconfig:"ECO_RDSMANAGEMENT_AMQ7_NS"`
	// StfNS is the namespace of STF.
	StfNS string `yaml:"rdsmanagement_stf_ns" envconfig:"ECO_RDSMANAGEMENT_STF_NS"`

	// KubeletCPUAllocation is the CPU allocated by the kubelet.
	KubeletCPUAllocation string `yaml:"rdsmanagement_kubelet_cpu_allocation_ns" envconfig:"ECO_RDSMANAGEMENT_KUBELET_CPU_ALLOCATION_NS"`
	//KubeletMemoryAllocation is the memory allocated by the kubelet.
	KubeletMemoryAllocation string `yaml:"rdsmanagement_kubelet_memory_allocation_ns" envconfig:"ECO_RDSMANAGEMENT_KUBELET_MEMORY_ALLOCATION_NS"`

	// IDMDeployed indicates whether IDM has been deployed or not
	IDMDeployed bool `yaml:"rdsmanagement_idm_deployed" envconfig:"ECO_RDSMANAGEMENT_IDM_DEPLOYED"`
	// SatelliteDeployed indicates whether Satellite has been deployed or not
	SatelliteDeployed bool `yaml:"rdsmanagement_satellite_deployed" envconfig:"ECO_RDSMANAGEMENT_SATELLITE_DEPLOYED"`
	// StfDeployed indicates whether STF has been deployed or not
	StfDeployed                 bool `yaml:"rdsmanagement_stf_deployed" envconfig:"ECO_RDSMANAGEMENT_STF_DEPLOYED"`
	ControlPlaneLabelListOption metav1.ListOptions
	//Odf Max Device Count
	OdfMaxDeviceCount int    `yaml:"rdsmanagement_odf_max_device_count" envconfig:"ECO_RDSMANAGEMENT_ODF_MAX_DEVICE_COUNT"`
	OdfReplica        int    `yaml:"rdsmanagement_odf_replica" envconfig:"ECO_RDSMANAGEMENT_ODF_REPLICA"`
	OdfPortable       string `yaml:"rdsmanagement_odf_portable" envconfig:"ECO_RDSMANAGEMENT_ODF_PORTABLE"`
	OdfName           string `yaml:"rdsmanagement_odf_name" envconfig:"ECO_RDSMANAGEMENT_ODF_NAME"`

	IDMConfig struct {
		Username           string `yaml:"username" envconfig:"ECO_RDSMANAGEMENT_IDM_CONFIG_USERNAME"`
		Password           string `yaml:"password" envconfig:"ECO_RDSMANAGEMENT_IDM_CONFIG_PASSWORD"`
		IPAddress          string `yaml:"ip_address" envconfig:"ECO_RDSMANAGEMENT_IDM_CONFIG_IP_ADDRESS"`
		FQDN               string `yaml:"fqdn" envconfig:"ECO_RDSMANAGEMENT_IDM_CONFIG_FQDN"`
		IPAAdminPass       string `yaml:"ipa_admin_pass" envconfig:"ECO_RDSMANAGEMENT_IDM_CONFIG_IPA_ADMIN_PASS"`
		IPADirMgrPass      string `yaml:"ipa_dir_mgr_pass" envconfig:"ECO_RDSMANAGEMENT_IDM_CONFIG_IPA_DIR_MGR_PASS"`
		IDMOcpBindPassword string `yaml:"idm_ocp_bind_password" envconfig:"ECO_RDSMANAGEMENT_IDM_OCP_BIND_PASSWORD"`
	} `yaml:"rdsmanagement_idm_config"`
}

// NewManagementConfig returns instance of ManagementConfig config type.
func NewManagementConfig() *ManagementConfig {
	log.Print("Creating new ManagementConfig struct")

	var rdsManagementConf ManagementConfig
	rdsManagementConf.GeneralConfig = config.NewConfig()

	var confFile string

	if fileFromEnv, exists := os.LookupEnv("ECO_RDS_MANAGEMENT_CONFIG_FILE_PATH"); !exists {
		_, filename, _, _ := runtime.Caller(0)
		baseDir := filepath.Dir(filename)
		confFile = filepath.Join(baseDir, PathToDefaultRDSManagementParamsFile)
	} else {
		confFile = fileFromEnv
	}

	log.Printf("Open config file %s", confFile)

	err := readFile(&rdsManagementConf, confFile)
	if err != nil {
		log.Printf("Error to read config file %s", confFile)

		return nil
	}

	err = readEnv(&rdsManagementConf)

	if err != nil {
		log.Print("Error to read environment variables")

		return nil
	}

	return &rdsManagementConf
}

func readFile(rdsConfig *ManagementConfig, cfgFile string) error {
	openedCfgFile, err := os.Open(cfgFile)
	if err != nil {
		return err
	}

	defer func() {
		_ = openedCfgFile.Close()
	}()

	decoder := yaml.NewDecoder(openedCfgFile)
	err = decoder.Decode(&rdsConfig)

	if err != nil {
		return err
	}

	return nil
}

func readEnv(rdsConfig *ManagementConfig) error {
	err := envconfig.Process("", rdsConfig)
	if err != nil {
		return err
	}

	return nil
}
