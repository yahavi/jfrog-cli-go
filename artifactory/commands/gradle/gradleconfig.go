package gradle

import (
	"github.com/jfrog/jfrog-cli-go/artifactory/utils"
	"github.com/jfrog/jfrog-cli-go/artifactory/utils/prompt"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	promptreader "github.com/jfrog/jfrog-client-go/utils/prompt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func CreateBuildConfig(configFilePath string) error {
	if err := prompt.VerifyConfigFile(configFilePath); err != nil {
		return err
	}

	configResult := &GradleBuildConfig{}
	configResult.Version = prompt.BUILD_CONF_VERSION
	configResult.ConfigType = utils.Gradle.String()

	vConfig, err := readGradleGlobalConfig()
	if err != nil {
		return err
	}
	configResult.UsePlugin = vConfig.GetBool(usePlugin)
	configResult.UseWrapper = vConfig.GetBool(useWrapper)

	vConfig, err = prompt.ReadArtifactoryServer("Resolve dependencies from Artifactory (y/n) [${default}]? ")
	if err != nil {
		return err
	}
	if vConfig.GetBool(prompt.USE_ARTIFACTORY) {
		err = configResult.Resolver.Server.Set(vConfig)
		if err != nil {
			return err
		}
		configResult.Resolver.Repo, err = prompt.ReadRepo("Set repository for dependencies resolution (press Tab for options): ", vConfig, utils.REMOTE, utils.VIRTUAL)
		if err != nil {
			return err
		}
	}

	vConfig, err = prompt.ReadArtifactoryServer("Deploy artifacts to Artifactory (y/n) [${default}]? ")
	if err != nil {
		return err
	}
	if vConfig.GetBool(prompt.USE_ARTIFACTORY) {
		err = configResult.Deployer.Server.Set(vConfig)
		if err != nil {
			return err
		}
		configResult.Deployer.Repo, err = prompt.ReadRepo("Set repository for artifacts deployment (press Tab for options): ", vConfig, utils.LOCAL, utils.VIRTUAL)
		if err != nil {
			return err
		}
		err = readDescriptors(&configResult.Deployer)
		if err != nil {
			return err
		}
	}

	resBytes, err := yaml.Marshal(&configResult)
	if err != nil {
		return errorutils.WrapError(err)
	}
	err = ioutil.WriteFile(configFilePath, resBytes, 0644)
	if err != nil {
		return errorutils.WrapError(err)
	}

	log.Info("Gradle build config successfully created.")
	return nil
}

func readGradleGlobalConfig() (*viper.Viper, error) {
	globalOptions := &promptreader.Array{
		Prompts: []promptreader.Prompt{
			&promptreader.YesNo{
				Msg:     "Is the Gradle Artifactory Plugin already applied in the build script (y/n) [${default}]? ",
				Default: "n",
				Label:   usePlugin,
			},
			&promptreader.YesNo{
				Msg:     "Use Gradle wrapper (y/n) [${default}]? ",
				Default: "n",
				Label:   useWrapper,
			},
		},
	}
	err := globalOptions.Read()
	if err != nil {
		return nil, errorutils.WrapError(err)
	}
	vConfig := globalOptions.GetResults()
	return vConfig, nil
}

func readDescriptors(deployer *GradleDeployer) error {
	descriptors := &promptreader.Array{
		Prompts: []promptreader.Prompt{
			&promptreader.YesNo{
				Msg:     "Deploy Maven descriptor (y/n) [${default}]? ",
				Default: "n",
				Label:   utils.MAVEN_DESCRIPTOR,
			},
			&promptreader.YesNo{
				Msg:     "Deploy Ivy descriptor (y/n) [${default}]? ",
				Default: "n",
				Label:   utils.IVY_DESCRIPTOR,
				Yes: &promptreader.Array{
					Prompts: []promptreader.Prompt{
						&promptreader.Simple{
							Msg:   "Set ivy pattern, [organization]/[module]/ivy-[revision].xml: ",
							Label: utils.IVY_PATTERN,
						},
						&promptreader.Simple{
							Msg:   "Set ivy artifact pattern, [organization]/[module]/[revision]/[artifact]-[revision](-[classifier]).[ext]: ",
							Label: utils.ARTIFACT_PATTERN,
						},
					},
				},
			},
		},
	}
	err := descriptors.Read()
	if err != nil {
		return errorutils.WrapError(err)
	}
	vConfig := descriptors.GetResults()
	deployer.DeployMavenDesc = vConfig.GetBool(utils.MAVEN_DESCRIPTOR)
	deployer.DeployIvyDesc = vConfig.GetBool(utils.IVY_DESCRIPTOR)
	if deployer.DeployIvyDesc {
		deployer.IvyPattern = vConfig.GetString(utils.IVY_PATTERN)
		deployer.ArtifactsPattern = vConfig.GetString(utils.ARTIFACT_PATTERN)
	}

	return nil
}

type GradleBuildConfig struct {
	prompt.CommonConfig `yaml:"common,inline"`
	UsePlugin           bool           `yaml:"usePlugin,omitempty"`
	UseWrapper          bool           `yaml:"useWrapper,omitempty"`
	Resolver            GradleRepo     `yaml:"resolver,omitempty"`
	Deployer            GradleDeployer `yaml:"deployer,omitempty"`
}
type GradleDeployer struct {
	GradleRepo       `yaml:"deployer,inline"`
	DeployMavenDesc  bool   `yaml:"deployMavenDescriptors,omitempty"`
	DeployIvyDesc    bool   `yaml:"deployIvyDescriptors,omitempty"`
	IvyPattern       string `yaml:"ivyPattern,omitempty"`
	ArtifactsPattern string `yaml:"artifactPattern,omitempty"`
}

type GradleRepo struct {
	Repo   string              `yaml:"repo,omitempty"`
	Server prompt.ServerConfig `yaml:"server,inline"`
}
