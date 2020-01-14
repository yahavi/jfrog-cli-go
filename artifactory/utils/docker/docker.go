package docker

import (
	"errors"
	"fmt"
	gofrogcmd "github.com/jfrog/gofrog/io"
	"github.com/jfrog/jfrog-cli-go/artifactory/utils"
	"github.com/jfrog/jfrog-cli-go/utils/cliutils"
	"github.com/jfrog/jfrog-cli-go/utils/config"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io"
	"os/exec"
	"path"
	"strings"
)

// Docker login error message
const DockerLoginFailureMessage string = "Docker login failed for: %s.\nDocker image must be in the form: docker-registry-domain/path-in-repository/image-name:version."

func New(imageTag string) Image {
	return &image{tag: imageTag}
}

// Docker image
type Image interface {
	Push() error
	Id() (string, error)
	ParentId() (string, error)
	Tag() string
	Path() string
	Name() string
	Pull() error
}

// Internal implementation of docker image
type image struct {
	tag string
}

type DockerLoginConfig struct {
	ArtifactoryDetails *config.ArtifactoryDetails
}

// Push docker image
func (image *image) Push() error {
	cmd := &pushCmd{image: image}
	return gofrogcmd.RunCmd(cmd)
}

// Get docker image tag
func (image *image) Tag() string {
	return image.tag
}

// Get docker image ID
func (image *image) Id() (string, error) {
	cmd := &getImageIdCmd{image: image}
	content, err := gofrogcmd.RunCmdOutput(cmd)
	return strings.Trim(content, "\n"), err
}

// Get docker parent image ID
func (image *image) ParentId() (string, error) {
	cmd := &getParentId{image: image}
	content, err := gofrogcmd.RunCmdOutput(cmd)
	return strings.Trim(content, "\n"), err
}

// Get docker image relative path in Artifactory
func (image *image) Path() string {
	indexOfFirstSlash := strings.Index(image.tag, "/")
	indexOfLastColon := strings.LastIndex(image.tag, ":")

	if indexOfLastColon < 0 || indexOfLastColon < indexOfFirstSlash {
		return path.Join(image.tag[indexOfFirstSlash:], "latest")
	}
	return path.Join(image.tag[indexOfFirstSlash:indexOfLastColon], image.tag[indexOfLastColon+1:])
}

// Get docker image name
func (image *image) Name() string {
	indexOfLastSlash := strings.LastIndex(image.tag, "/")
	indexOfLastColon := strings.LastIndex(image.tag, ":")

	if indexOfLastColon < 0 || indexOfLastColon < indexOfLastSlash {
		return image.tag[indexOfLastSlash+1:] + ":latest"
	}
	return image.tag[indexOfLastSlash+1:]
}

// Pull docker image
func (image *image) Pull() error {
	cmd := &pullCmd{image: image}
	return gofrogcmd.RunCmd(cmd)
}

// Image push command
type pushCmd struct {
	image *image
}

func (pushCmd *pushCmd) GetCmd() *exec.Cmd {
	var cmd []string
	cmd = append(cmd, "docker")
	cmd = append(cmd, "push")
	cmd = append(cmd, pushCmd.image.tag)
	return exec.Command(cmd[0], cmd[1:]...)
}

func (pushCmd *pushCmd) GetEnv() map[string]string {
	return map[string]string{}
}

func (pushCmd *pushCmd) GetStdWriter() io.WriteCloser {
	return nil
}
func (pushCmd *pushCmd) GetErrWriter() io.WriteCloser {
	return nil
}

// Image get image id command
type getImageIdCmd struct {
	image *image
}

func (getImageId *getImageIdCmd) GetCmd() *exec.Cmd {
	var cmd []string
	cmd = append(cmd, "docker")
	cmd = append(cmd, "images")
	cmd = append(cmd, "--format", "{{.ID}}")
	cmd = append(cmd, "--no-trunc")
	cmd = append(cmd, getImageId.image.tag)
	return exec.Command(cmd[0], cmd[1:]...)
}

func (getImageId *getImageIdCmd) GetEnv() map[string]string {
	return map[string]string{}
}

func (getImageId *getImageIdCmd) GetStdWriter() io.WriteCloser {
	return nil
}

func (getImageId *getImageIdCmd) GetErrWriter() io.WriteCloser {
	return nil
}

// Image get parent image id command
type getParentId struct {
	image *image
}

func (getImageId *getParentId) GetCmd() *exec.Cmd {
	var cmd []string
	cmd = append(cmd, "docker")
	cmd = append(cmd, "inspect")
	cmd = append(cmd, "--format", "{{.Parent}}")
	cmd = append(cmd, getImageId.image.tag)
	return exec.Command(cmd[0], cmd[1:]...)
}

func (getImageId *getParentId) GetEnv() map[string]string {
	return map[string]string{}
}

func (getImageId *getParentId) GetStdWriter() io.WriteCloser {
	return nil
}

func (getImageId *getParentId) GetErrWriter() io.WriteCloser {
	return nil
}

// Get docker registry from tag
func ResolveRegistryFromTag(imageTag string) (string, error) {
	indexOfFirstSlash := strings.Index(imageTag, "/")
	if indexOfFirstSlash < 0 {
		err := errorutils.WrapError(errors.New("Invalid image tag received for pushing to Artifactory - tag does not include a slash."))
		return "", err
	}

	indexOfSecondSlash := strings.Index(imageTag[indexOfFirstSlash+1:], "/")
	// Reverse proxy Artifactory
	if indexOfSecondSlash < 0 {
		return imageTag[:indexOfFirstSlash], nil
	}
	// Can be reverse proxy or proxy-less Artifactory
	indexOfSecondSlash += indexOfFirstSlash + 1
	return imageTag[:indexOfSecondSlash], nil
}

// Login command
type LoginCmd struct {
	DockerRegistry string
	Username       string
	Password       string
}

func (loginCmd *LoginCmd) GetCmd() *exec.Cmd {
	if cliutils.IsWindows() {
		return exec.Command("cmd", "/C", "echo", "%DOCKER_PASS%|", "docker", "login", loginCmd.DockerRegistry, "--username", loginCmd.Username, "--password-stdin")
	}
	cmd := "echo $DOCKER_PASS " + fmt.Sprintf(`| docker login %s --username="%s" --password-stdin`, loginCmd.DockerRegistry, loginCmd.Username)
	return exec.Command("sh", "-c", cmd)
}

func (loginCmd *LoginCmd) GetEnv() map[string]string {
	return map[string]string{"DOCKER_PASS": loginCmd.Password}
}

func (loginCmd *LoginCmd) GetStdWriter() io.WriteCloser {
	return nil
}

func (loginCmd *LoginCmd) GetErrWriter() io.WriteCloser {
	return nil
}

// Image pull command
type pullCmd struct {
	image *image
}

func (pullCmd *pullCmd) GetCmd() *exec.Cmd {
	var cmd []string
	cmd = append(cmd, "docker")
	cmd = append(cmd, "pull")
	cmd = append(cmd, pullCmd.image.tag)
	return exec.Command(cmd[0], cmd[1:]...)
}

func (pullCmd *pullCmd) GetEnv() map[string]string {
	return map[string]string{}
}

func (pullCmd *pullCmd) GetStdWriter() io.WriteCloser {
	return nil
}

func (pullCmd *pullCmd) GetErrWriter() io.WriteCloser {
	return nil
}

func CreateServiceManager(artDetails *config.ArtifactoryDetails, threads int) (*artifactory.ArtifactoryServicesManager, error) {
	certPath, err := utils.GetJfrogSecurityDir()
	if err != nil {
		return nil, err
	}
	artAuth, err := artDetails.CreateArtAuthConfig()
	if err != nil {
		return nil, err
	}

	configBuilder := artifactory.NewConfigBuilder().
		SetArtDetails(artAuth).
		SetCertificatesPath(certPath).
		SetInsecureTls(artDetails.InsecureTls).
		SetThreads(threads)

	if threads != 0 {
		configBuilder.SetThreads(threads)
	}

	serviceConfig, err := configBuilder.Build()
	return artifactory.New(&artAuth, serviceConfig)
}

// First will try to login assuming a proxy-less tag (e.g. "registry-address/docker-repo/image:ver").
// If fails, we will try assuming a reverse proxy tag (e.g. "registry-address-docker-repo/image:ver").
func DockerLogin(imageTag string, config *DockerLoginConfig) error {
	imageRegistry, err := ResolveRegistryFromTag(imageTag)
	if err != nil {
		return err
	}

	username := config.ArtifactoryDetails.User
	password := config.ArtifactoryDetails.Password
	// If access-token exists, perform login with it.
	if config.ArtifactoryDetails.AccessToken != "" {
		log.Debug("Using access-token details in docker-login command.")
		username, err = auth.ExtractUsernameFromAccessToken(config.ArtifactoryDetails.AccessToken)
		if err != nil {
			return err
		}
		password = config.ArtifactoryDetails.AccessToken
	}

	// Perform login.
	cmd := &LoginCmd{DockerRegistry: imageRegistry, Username: username, Password: password}
	err = gofrogcmd.RunCmd(cmd)

	if exitCode := cliutils.GetExitCode(err, 0, 0, false); exitCode == cliutils.ExitCodeNoError {
		// Login succeeded
		return nil
	}
	log.Debug("Docker login while assuming proxy-less failed:", err)

	indexOfSlash := strings.Index(imageRegistry, "/")
	if indexOfSlash < 0 {
		return errorutils.WrapError(errors.New(fmt.Sprintf(DockerLoginFailureMessage, imageRegistry)))
	}

	cmd = &LoginCmd{DockerRegistry: imageRegistry[:indexOfSlash], Username: config.ArtifactoryDetails.User, Password: config.ArtifactoryDetails.Password}
	err = gofrogcmd.RunCmd(cmd)
	if err != nil {
		// Login failed for both attempts
		return errorutils.WrapError(errors.New(fmt.Sprintf(DockerLoginFailureMessage,
			fmt.Sprintf("%s, %s", imageRegistry, imageRegistry[:indexOfSlash])) + " " + err.Error()))
	}

	// Login succeeded
	return nil
}
