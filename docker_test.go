package main

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	gofrogcmd "github.com/jfrog/gofrog/io"
	"github.com/jfrog/jfrog-cli/artifactory/commands/generic"
	"github.com/jfrog/jfrog-cli/artifactory/spec"
	"github.com/jfrog/jfrog-cli/artifactory/utils/docker"
	"github.com/jfrog/jfrog-cli/inttestutils"
	"github.com/jfrog/jfrog-cli/utils/cliutils"
	"github.com/jfrog/jfrog-cli/utils/tests"
)

func InitDockerTests() {
	initArtifactoryCli()
	cleanUpOldBuilds()
	inttestutils.CleanUpOldImages(artifactoryDetails, artHttpDetails)
	tests.AddTimestampToGlobalVars()
}

func initDockerTest(t *testing.T) {
	if !*tests.TestDocker {
		t.Skip("Skipping docker test. To run docker test add the '-test.docker=true' option.")
	}
}

func TestDockerPush(t *testing.T) {
	initDockerTest(t)
	runDockerPushTest(tests.DockerImageName, tests.DockerImageName+":1", false, t)
}

func TestDockerPushWithModuleName(t *testing.T) {
	initDockerTest(t)
	runDockerPushTest(tests.DockerImageName, ModuleNameJFrogTest, true, t)
}

func TestDockerPushWithMultipleSlash(t *testing.T) {
	initDockerTest(t)
	runDockerPushTest(tests.DockerImageName+"/multiple", "multiple:1", false, t)
}

// Run docker push to Artifactory
func runDockerPushTest(imageName, module string, withModule bool, t *testing.T) {
	imageTag := inttestutils.BuildTestDockerImage(imageName)
	buildNumber := "1"

	// Push docker image using docker client
	if withModule {
		artifactoryCli.Exec("docker-push", imageTag, *tests.DockerTargetRepo, "--build-name="+tests.DockerBuildName, "--build-number="+buildNumber, "--module="+module)
	} else {
		artifactoryCli.Exec("docker-push", imageTag, *tests.DockerTargetRepo, "--build-name="+tests.DockerBuildName, "--build-number="+buildNumber)
	}
	artifactoryCli.Exec("build-publish", tests.DockerBuildName, buildNumber)

	imagePath := path.Join(*tests.DockerTargetRepo, imageName, "1") + "/"
	validateDockerBuild(tests.DockerBuildName, buildNumber, imagePath, module, 7, 5, 7, t)
	inttestutils.DockerTestCleanup(artifactoryDetails, artHttpDetails, imageName, tests.DockerBuildName)

}
func TestDockerPushBuildNameNumberFromEnv(t *testing.T) {
	initDockerTest(t)
	imageTag := inttestutils.BuildTestDockerImage(tests.DockerImageName)
	buildNumber := "1"
	os.Setenv(cliutils.BuildName, tests.DockerBuildName)
	os.Setenv(cliutils.BuildNumber, buildNumber)
	defer os.Unsetenv(cliutils.BuildName)
	defer os.Unsetenv(cliutils.BuildNumber)

	// Push docker image using docker client
	artifactoryCli.Exec("docker-push", imageTag, *tests.DockerTargetRepo)
	artifactoryCli.Exec("build-publish")

	imagePath := path.Join(*tests.DockerTargetRepo, tests.DockerImageName, "1") + "/"
	validateDockerBuild(tests.DockerBuildName, buildNumber, imagePath, tests.DockerImageName+":1", 7, 5, 7, t)

	inttestutils.DockerTestCleanup(artifactoryDetails, artHttpDetails, tests.DockerImageName, tests.DockerBuildName)
}

func TestDockerPull(t *testing.T) {
	initDockerTest(t)

	imageTag := inttestutils.BuildTestDockerImage(tests.DockerImageName)

	// Push docker image using docker client
	artifactoryCli.Exec("docker-push", imageTag, *tests.DockerTargetRepo)

	buildNumber := "1"

	// Pull docker image using docker client
	artifactoryCli.Exec("docker-pull", imageTag, *tests.DockerTargetRepo, "--build-name="+tests.DockerBuildName, "--build-number="+buildNumber)
	artifactoryCli.Exec("build-publish", tests.DockerBuildName, buildNumber)

	imagePath := path.Join(*tests.DockerTargetRepo, tests.DockerImageName, "1") + "/"
	validateDockerBuild(tests.DockerBuildName, buildNumber, imagePath, tests.DockerImageName+":1", 0, 7, 7, t)

	buildNumber = "2"
	artifactoryCli.Exec("docker-pull", imageTag, *tests.DockerTargetRepo, "--build-name="+tests.DockerBuildName, "--build-number="+buildNumber, "--module="+ModuleNameJFrogTest)
	artifactoryCli.Exec("build-publish", tests.DockerBuildName, buildNumber)
	validateDockerBuild(tests.DockerBuildName, buildNumber, imagePath, ModuleNameJFrogTest, 0, 7, 7, t)

	inttestutils.DockerTestCleanup(artifactoryDetails, artHttpDetails, tests.DockerImageName, tests.DockerBuildName)
}

func dockerTestCleanup(imageName, buildName string) {
	// Remove build from Artifactory
	inttestutils.DeleteBuild(artifactoryDetails.Url, buildName, artHttpDetails)
	inttestutils.DockerTestCleanup(artifactoryDetails, artHttpDetails, tests.DockerImageName, tests.DockerBuildName)
}

func TestDockerClientApiVersionCmd(t *testing.T) {
	initDockerTest(t)

	// Run docker version command and expect no errors
	cmd := &docker.VersionCmd{}
	content, err := gofrogcmd.RunCmdOutput(cmd)
	assert.NoError(t, err)

	// Expect VersionRegex to match the output API version
	content = strings.TrimSpace(content)
	assert.True(t, docker.ApiVersionRegex.Match([]byte(content)))

	// Assert docker min API version
	assert.True(t, docker.IsCompatibleApiVersion(content))
}

func TestDockerFatManifestPull(t *testing.T) {
	initDockerTest(t)

	imageName := "traefik"
	imageTag := path.Join(*tests.DockerRepoDomain, imageName+":2.2")
	buildNumber := "1"

	// Pull docker image using docker client
	artifactoryCli.Exec("docker-pull", imageTag, *tests.DockerTargetRepo, "--build-name="+tests.DockerBuildName, "--build-number="+buildNumber)
	artifactoryCli.Exec("build-publish", tests.DockerBuildName, buildNumber)

	// Validate
	buildInfo, _ := inttestutils.GetBuildInfo(artifactoryDetails.Url, tests.DockerBuildName, buildNumber, t, artHttpDetails)
	validateBuildInfo(buildInfo, t, 6, 0, imageName+":2.2")

	inttestutils.DockerTestCleanup(artifactoryDetails, artHttpDetails, imageName, tests.DockerBuildName)
	inttestutils.DeleteTestDockerImage(imageTag)
}

func validateDockerBuild(buildName, buildNumber, imagePath, module string, expectedArtifacts, expectedDependencies, expectedItemsInArtifactory int, t *testing.T) {
	specFile := spec.NewBuilder().Pattern(imagePath + "*").BuildSpec()
	searchCmd := generic.NewSearchCommand()
	searchCmd.SetRtDetails(artifactoryDetails).SetSpec(specFile)
	assert.NoError(t, searchCmd.Search())
	assert.Len(t, searchCmd.SearchResult(), expectedItemsInArtifactory, "Docker build info was not pushed correctly")

	buildInfo, _ := inttestutils.GetBuildInfo(artifactoryDetails.Url, buildName, buildNumber, t, artHttpDetails)
	validateBuildInfo(buildInfo, t, expectedDependencies, expectedArtifacts, module)
}
