package pip

import (
	"github.com/jfrog/jfrog-cli-go/artifactory/utils"
	"github.com/jfrog/jfrog-cli-go/artifactory/utils/prompt"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type PipBuildConfig struct {
	prompt.CommonConfig `yaml:"common,inline"`
	Resolver            utils.Repository `yaml:"resolver,omitempty"`
}

func CreateBuildConfig(global bool) error {
	projectDir, err := utils.GetProjectDir(global)
	if err != nil {
		return err
	}
	err = fileutils.CreateDirIfNotExist(projectDir)
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(projectDir, "pip.yaml")
	if err := prompt.VerifyConfigFile(configFilePath); err != nil {
		return err
	}

	var vConfig *viper.Viper
	configResult := &PipBuildConfig{}
	configResult.Version = prompt.BUILD_CONF_VERSION
	configResult.ConfigType = utils.Pip.String()
	configResult.Resolver.ServerId, vConfig, err = prompt.ReadServerId()
	if err != nil {
		return errorutils.WrapError(err)
	}
	configResult.Resolver.Repo, err = prompt.ReadRepo("Set repository for dependencies resolution (press Tab for options): ", vConfig, utils.LOCAL, utils.VIRTUAL)
	if err != nil {
		return errorutils.WrapError(err)
	}
	resBytes, err := yaml.Marshal(&configResult)
	if err != nil {
		return errorutils.WrapError(err)
	}
	err = ioutil.WriteFile(configFilePath, resBytes, 0644)
	if err != nil {
		return errorutils.WrapError(err)
	}
	log.Info("Pip build config successfully created.")

	return nil
}
