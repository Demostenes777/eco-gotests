package ocloudconfig

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
	// PathToDefaultOCloudParamsFile path to config file with default o-cloud parameters.
	PathToDefaultOCloudParamsFile = "./default.yaml"
)

// OCloudConfig type keeps o-cloud configuration.
type OCloudConfig struct {
	*systemtestsconfig.SystemTestsConfig
	// todo
}

// NewOCloudConfig returns instance of OCloudConfig config type.
func NewOCloudConfig() *OCloudConfig {
	log.Print("Creating new OCloudConfig struct")

	var ocloudConf OCloudConfig
	ocloudConf.SystemTestsConfig = systemtestsconfig.NewSystemTestsConfig()

	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	confFile := filepath.Join(baseDir, PathToDefaultOCloudParamsFile)
	err := readFile(&ocloudConf, confFile)

	if err != nil {
		log.Printf("Error to read config file %s", confFile)

		return nil
	}

	err = readEnv(&ocloudConf)

	if err != nil {
		log.Print("Error to read environment variables")

		return nil
	}

	return &ocloudConf
}

func readFile(ocloudConfig *OCloudConfig, cfgFile string) error {
	openedCfgFile, err := os.Open(cfgFile)
	if err != nil {
		return err
	}

	defer func() {
		_ = openedCfgFile.Close()
	}()

	decoder := yaml.NewDecoder(openedCfgFile)
	err = decoder.Decode(&ocloudConfig)

	if err != nil {
		return err
	}

	return nil
}

func readEnv(ocloudConfig *OCloudConfig) error {
	err := envconfig.Process("", ocloudConfig)
	if err != nil {
		return err
	}

	return nil
}
