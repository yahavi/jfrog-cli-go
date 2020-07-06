package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/jfrog/jfrog-cli/artifactory/commands/generic"
	"github.com/jfrog/jfrog-cli/artifactory/utils"
	"github.com/jfrog/jfrog-cli/inttestutils"
	"github.com/jfrog/jfrog-cli/utils/cliutils"
	"github.com/jfrog/jfrog-cli/utils/tests"
	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	rtutils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ModuleNameJFrogTest = "jfrog-test"

func TestBuildAddDependenciesFromHomeDir(t *testing.T) {
	initArtifactoryTest(t)
	// Clean old build tests if exists
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)

	fileName := "cliTestFile.txt"
	testFileRelPath, testFileAbs := createFileInHomeDir(t, fileName)

	test := buildAddDepsBuildInfoTestParams{description: "'rt bad' from home dir", commandArgs: []string{testFileRelPath, "--recursive=false"}, expectedDependencies: []string{fileName}, buildName: tests.RtBuildName1, buildNumber: "1"}
	collectDepsAndPublishBuild(test, false, t)
	validateBuildAddDepsBuildInfo(t, test)

	os.Remove(testFileAbs)
	cleanArtifactoryTest()
}

func TestBuildPromote(t *testing.T) {
	initArtifactoryTest(t)
	buildNumberA := "10"

	// Upload files with buildName and buildNumber
	specFileA, err := tests.CreateSpec(tests.SplitUploadSpecA)
	assert.NoError(t, err)
	artifactoryCli.Exec("upload", "--spec="+specFileA, "--build-name="+tests.RtBuildName1, "--build-number="+buildNumberA)
	artifactoryCli.Exec("build-publish", tests.RtBuildName1, buildNumberA)

	key1 := "key"
	value1 := "v1,v2"
	key2 := "another"
	value2 := "property"

	// Promote build to Repo1 using build name and build number as args.
	artifactoryCli.Exec("build-promote", tests.RtBuildName1, buildNumberA, tests.RtRepo1, fmt.Sprintf("--props=%s=%s;%s=%s", key1, value1, key2, value2))
	buildInfo, _ := inttestutils.GetBuildInfo(artifactoryDetails.Url, tests.RtBuildName1, buildNumberA, t, artHttpDetails)
	resultItems := getResultItemsFromArtifactory(tests.SearchAllRepo1, t)

	assert.Equal(t, len(buildInfo.Modules[0].Artifacts), len(resultItems), "Incorrect number of artifacts were uploaded")

	// Promote the same build to Repo2 using build name and build number as env vars.
	os.Setenv(cliutils.BuildName, tests.RtBuildName1)
	os.Setenv(cliutils.BuildNumber, buildNumberA)
	defer os.Unsetenv(cliutils.BuildName)
	defer os.Unsetenv(cliutils.BuildNumber)
	artifactoryCli.Exec("build-promote", tests.RtRepo2, fmt.Sprintf("--props=%s=%s;%s=%s", key1, value1, key2, value2))
	buildInfo, _ = inttestutils.GetBuildInfo(artifactoryDetails.Url, tests.RtBuildName1, buildNumberA, t, artHttpDetails)
	resultItems = getResultItemsFromArtifactory(tests.SearchRepo2, t)

	assert.Equal(t, len(buildInfo.Modules[0].Artifacts), len(resultItems), "Incorrect number of artifacts were uploaded")

	propsMap := map[string]string{
		"build.name":   buildInfo.Name,
		"build.number": buildInfo.Number,
		key1:           value1,
		key2:           value2,
	}

	validateArtifactsProperties(resultItems, t, propsMap)
	cleanArtifactoryTest()
}

// Returns the artifacts found by the provided spec
func getResultItemsFromArtifactory(specName string, t *testing.T) []rtutils.ResultItem {
	searchGoSpecFile, err := tests.CreateSpec(specName)
	assert.NoError(t, err)
	spec, flags := getSpecAndCommonFlags(searchGoSpecFile)
	flags.SetArtifactoryDetails(artAuth)
	var resultItems []rtutils.ResultItem
	for i := 0; i < len(spec.Files); i++ {
		searchParams, err := generic.GetSearchParams(spec.Get(i))
		assert.NoError(t, err)

		currentResultItems, err := services.SearchBySpecFiles(searchParams, flags, rtutils.ALL)
		assert.NoError(t, err, "Failed Searching files")
		resultItems = append(resultItems, currentResultItems...)
	}
	return resultItems
}

// This function validates the properties on the provided artifacts. Every property within the provided map should be attached to the artifact.
func validateArtifactsProperties(resultItems []rtutils.ResultItem, t *testing.T, propsMap map[string]string) {
	for _, item := range resultItems {
		properties := item.Properties
		assert.GreaterOrEqual(t, len(properties), 1, "Failed finding properties on item:", item.GetItemRelativePath())
		propertiesMap := tests.ConvertSliceToMap(properties)

		for key, value := range propsMap {
			valueFromArtifact, contains := propertiesMap[key]
			assert.True(t, contains, "Failed finding %s property on %s", key, item.Name)
			assert.Equalf(t, value, valueFromArtifact, "Wrong value for %s property on %s.", key, item.Name)
		}
	}
}

func TestBuildAddDependenciesDryRun(t *testing.T) {
	initArtifactoryTest(t)
	// Clean old build tests if exists
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)
	err := utils.RemoveBuildDir(tests.RtBuildName1, "1")
	assert.NoError(t, err)

	wd, err := os.Getwd()
	assert.NoError(t, err)

	defer os.Chdir(wd)
	err = os.Chdir("testsdata")
	assert.NoError(t, err)

	noCredsCli := tests.NewJfrogCli(execMain, "jfrog rt", "")
	// Execute the bad command
	noCredsCli.Exec("bad", tests.RtBuildName1, "1", "a/*", "--dry-run=true")
	buildDir, err := utils.GetBuildDir(tests.RtBuildName1, "1")
	assert.NoError(t, err)

	files, _ := ioutil.ReadDir(buildDir)
	assert.Zero(t, len(files), "'rt bad' command with dry-run failed. The dry-run option has no effect.")

	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)
	os.Chdir(wd)
	cleanArtifactoryTest()
}

func TestBuildPublishDryRun(t *testing.T) {
	initArtifactoryTest(t)
	buildNumber := "11"
	// Clean old build tests if exists.
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)
	assert.NoError(t, utils.RemoveBuildDir(tests.RtBuildName1, buildNumber))

	// Upload files with build name & number.
	specFile, err := tests.CreateSpec(tests.UploadFlatRecursive)
	assert.NoError(t, err)
	artifactoryCli.Exec("upload", "--spec="+specFile, "--build-name="+tests.RtBuildName1, "--build-number="+buildNumber)
	// Verify build dir is not empty
	assert.NotEmpty(t, getFilesFromBuildDir(t, tests.RtBuildName1, buildNumber))

	// Execute the bp command with dry run
	artifactoryCli.Exec("bp", tests.RtBuildName1, buildNumber, "--dry-run=true")
	// Verify build dir is not empty
	assert.NotEmpty(t, getFilesFromBuildDir(t, tests.RtBuildName1, buildNumber))
	// Verify build was not published
	_, found := inttestutils.GetBuildInfo(artifactoryDetails.Url, tests.RtBuildName1, buildNumber, t, artHttpDetails)
	assert.False(t, found)

	// Execute the bp command without dry run
	artifactoryCli.Exec("bp", tests.RtBuildName1, buildNumber)
	// Verify build dir is empty
	assert.Empty(t, getFilesFromBuildDir(t, tests.RtBuildName1, buildNumber))
	// Verify build was published
	buildInfo, _ := inttestutils.GetBuildInfo(artifactoryDetails.Url, tests.RtBuildName1, buildNumber, t, artHttpDetails)
	validateBuildInfo(buildInfo, t, 0, 9, tests.RtBuildName1)

	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)
	cleanArtifactoryTest()
}

func getFilesFromBuildDir(t *testing.T, buildName, buildNumber string) []os.FileInfo {
	buildDir, err := utils.GetBuildDir(buildName, buildNumber)
	assert.NoError(t, err)

	files, err := ioutil.ReadDir(buildDir)
	assert.NoError(t, err)
	return files
}

func TestBuildAddDependencies(t *testing.T) {
	initArtifactoryTest(t)
	// Clean old build tests if exists
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)

	allFiles := []string{"a1.in", "a2.in", "a3.in", "b1.in", "b2.in", "b3.in", "c1.in", "c2.in", "c3.in"}
	var badTests = []buildAddDepsBuildInfoTestParams{
		{description: "'rt bad' simple cli", commandArgs: []string{"testsdata/a/*"}, expectedDependencies: allFiles},
		{description: "'rt bad' single file", commandArgs: []string{"testsdata/a/a1.in"}, expectedDependencies: []string{"a1.in"}},
		{description: "'rt bad' none recursive", commandArgs: []string{"testsdata/a/*", "--recursive=false"}, expectedDependencies: []string{"a1.in", "a2.in", "a3.in"}},
		{description: "'rt bad' special chars recursive", commandArgs: []string{getSpecialCharFilePath()}, expectedDependencies: []string{"a1.in"}},
		{description: "'rt bad' exclude command line wildcards", commandArgs: []string{"testsdata/a/*", "--exclude-patterns=*a2*;*a3.in"}, expectedDependencies: []string{"a1.in", "b1.in", "b2.in", "b3.in", "c1.in", "c2.in", "c3.in"}},
		{description: "'rt bad' spec", commandArgs: []string{"--spec=" + tests.GetFilePathForArtifactory(tests.BuildAddDepsSpec)}, expectedDependencies: allFiles},
		{description: "'rt bad' two specFiles", commandArgs: []string{"--spec=" + tests.GetFilePathForArtifactory(tests.BuildAddDepsDoubleSpec)}, expectedDependencies: []string{"a1.in", "a2.in", "a3.in", "b1.in", "b2.in", "b3.in"}},
		{description: "'rt bad' exclude command line regexp", commandArgs: []string{"testsdata/a/a(.*)", "--exclude-patterns=(.*)a2.*;.*a3.in", "--regexp=true"}, expectedDependencies: []string{"a1.in"}},
	}

	// Tests compatibility to file paths with windows separators.
	if cliutils.IsWindows() {
		var compatibilityTests = []buildAddDepsBuildInfoTestParams{
			{description: "'rt bad' win compatibility by arguments", commandArgs: []string{"testsdata\\\\a\\\\a1.in"}, expectedDependencies: []string{"a1.in"}},
			{description: "'rt bad' win compatibility by spec", commandArgs: []string{"--spec=" + tests.GetFilePathForArtifactory(tests.WinBuildAddDepsSpec)}, expectedDependencies: allFiles},
		}
		badTests = append(badTests, compatibilityTests...)
	}

	for i, badTest := range badTests {
		badTest.buildName = tests.RtBuildName1
		badTest.buildNumber = strconv.Itoa(i + 1)

		collectDepsAndPublishBuild(badTest, true, t)
		validateBuildAddDepsBuildInfo(t, badTest)
		utils.RemoveBuildDir(badTest.buildName, badTest.buildNumber)

		collectDepsAndPublishBuild(badTest, false, t)
		validateBuildAddDepsBuildInfo(t, badTest)
		utils.RemoveBuildDir(badTest.buildName, badTest.buildNumber)
	}
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)
}

// Test publish build info without --build-url
func TestArtifactoryPublishBuildInfo(t *testing.T) {
	initArtifactoryTest(t)
	buildNumber := "10"
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)

	body := uploadFilesAndGetBuildInfo(t, tests.RtBuildName1, buildNumber, "")

	// Validate no build url
	_, _, _, err := jsonparser.Get(body, "buildInfo", "url")
	assert.Error(t, err, "Build url is expected to be empty")

	// Cleanup
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)
	cleanArtifactoryTest()
}

// Test publish build info with --build-url
func TestArtifactoryPublishBuildInfoBuildUrl(t *testing.T) {
	initArtifactoryTest(t)
	buildNumber := "11"
	buildUrl := "http://example.ci.com"
	os.Setenv(cliutils.BuildUrl, "http://override-me.ci.com")
	defer os.Unsetenv(cliutils.BuildUrl)
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)

	body := uploadFilesAndGetBuildInfo(t, tests.RtBuildName1, buildNumber, buildUrl)

	// Validate correctness of build url
	actualBuildUrl, err := jsonparser.GetString(body, "buildInfo", "url")
	assert.NoError(t, err)
	assert.Equal(t, buildUrl, actualBuildUrl)

	// Cleanup
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)
	cleanArtifactoryTest()
}

// Test publish build info with JFROG_CLI_BUILD_URL env
func TestArtifactoryPublishBuildInfoBuildUrlFromEnv(t *testing.T) {
	initArtifactoryTest(t)
	buildNumber := "11"
	buildUrl := "http://example-env.ci.com"
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)
	os.Setenv(cliutils.BuildUrl, buildUrl)
	defer os.Unsetenv(cliutils.BuildUrl)
	body := uploadFilesAndGetBuildInfo(t, tests.RtBuildName1, buildNumber, "")

	// Validate correctness of build url
	actualBuildUrl, err := jsonparser.GetString(body, "buildInfo", "url")
	assert.NoError(t, err)
	assert.Equal(t, buildUrl, actualBuildUrl)

	// Cleanup
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)
	cleanArtifactoryTest()
}

func TestArtifactoryCleanBuildInfo(t *testing.T) {
	initArtifactoryTest(t)
	buildNumber := "11"
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)

	// Upload files with buildName and buildNumber
	specFile, err := tests.CreateSpec(tests.SplitUploadSpecA)
	assert.NoError(t, err)
	artifactoryCli.Exec("upload", "--spec="+specFile, "--build-name="+tests.RtBuildName1, "--build-number="+buildNumber)

	// Cleanup buildInfo with the same buildName and buildNumber
	artifactoryCli.WithoutCredentials().Exec("build-clean", tests.RtBuildName1, buildNumber)

	// Upload different files with the same buildName and buildNumber
	specFile, err = tests.CreateSpec(tests.SplitUploadSpecB)
	assert.NoError(t, err)
	artifactoryCli.Exec("upload", "--spec="+specFile, "--build-name="+tests.RtBuildName1, "--build-number="+buildNumber)

	// Publish buildInfo
	artifactoryCli.Exec("build-publish", tests.RtBuildName1, buildNumber)

	// Download by build and verify that only artifacts uploaded after clean are downloaded
	outputDir := filepath.Join(tests.Out, "clean-build")
	artifactoryCli.Exec("download", tests.RtRepo1, outputDir+fileutils.GetFileSeparator(), "--build="+tests.RtBuildName1+"/"+buildNumber)
	paths, _ := fileutils.ListFilesRecursiveWalkIntoDirSymlink(outputDir, false)
	tests.VerifyExistLocally(tests.GetCleanBuild(), paths, t)

	// Cleanup
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)
	cleanArtifactoryTest()
}

func TestArtifactoryBuildCollectEnv(t *testing.T) {
	initArtifactoryTest(t)
	buildNumber := "12"

	// Build collect env
	os.Setenv("DONT_COLLECT", "foo")
	os.Setenv("COLLECT", "bar")
	artifactoryCli.WithoutCredentials().Exec("bce", tests.RtBuildName1, buildNumber)

	// Publish build info
	artifactoryCli.Exec("bp", tests.RtBuildName1, buildNumber, "--env-exclude=*password*;*psw*;*secret*;*key*;*token*;DONT_COLLECT")
	buildInfo, _ := inttestutils.GetBuildInfo(artifactoryDetails.Url, tests.RtBuildName1, buildNumber, t, artHttpDetails)

	// Make sure no sensitive data in build env
	for k := range buildInfo.Properties {
		assert.NotContains(t, k, "password")
		assert.NotContains(t, k, "psw")
		assert.NotContains(t, k, "secret")
		assert.NotContains(t, k, "key")
		assert.NotContains(t, k, "token")
		assert.NotContains(t, k, "DONT_COLLECT")
	}

	// Make sure "COLLECT" env appear in build env
	assert.Contains(t, buildInfo.Properties, "buildInfo.env.COLLECT")

	// Cleanup
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)
	cleanArtifactoryTest()
}

func TestBuildAddGit(t *testing.T) {
	testBuildAddGit(t, false)
}

func TestBuildAddGitEnvBuildNameAndNumber(t *testing.T) {
	testBuildAddGit(t, true)
}

func testBuildAddGit(t *testing.T, useEnvBuildNameAndNumber bool) {
	initArtifactoryTest(t)
	gitCollectCliRunner := tests.NewJfrogCli(execMain, "jfrog rt", "")
	buildNumber := "13"

	// Populate cli config with 'default' server
	oldHomeDir := os.Getenv(cliutils.HomeDir)
	createJfrogHomeConfig(t, true)

	// Create .git folder for this test
	originalFolder := "buildaddgit_.git_suffix"
	baseDir, dotGitPath := tests.PrepareDotGitDir(t, originalFolder, "testsdata")

	// Get path for build-add-git config file
	pwd, _ := os.Getwd()
	configPath := filepath.Join(pwd, "testsdata", "buildaddgit_config.yaml")

	// Run build-add-git
	var err error
	if useEnvBuildNameAndNumber {
		os.Setenv(cliutils.BuildName, tests.RtBuildName1)
		os.Setenv(cliutils.BuildNumber, buildNumber)
		defer os.Unsetenv(cliutils.BuildName)
		defer os.Unsetenv(cliutils.BuildNumber)
		err = gitCollectCliRunner.Exec("build-add-git", baseDir, "--config="+configPath)
	} else {
		err = gitCollectCliRunner.Exec("build-add-git", tests.RtBuildName1, buildNumber, baseDir, "--config="+configPath)
	}
	defer cleanBuildAddGitTest(t, baseDir, originalFolder, oldHomeDir, dotGitPath)
	if err != nil {
		t.Fatal(err)
	}
	// Clear previous build if exists and publish build-info
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)
	artifactoryCli.Exec("build-publish", tests.RtBuildName1, buildNumber)

	// Fetch the published build-info for validation
	buildInfo, _ := inttestutils.GetBuildInfo(artifactoryDetails.Url, tests.RtBuildName1, buildNumber, t, artHttpDetails)
	if t.Failed() {
		t.FailNow()
	}
	require.NotNil(t, buildInfo.Vcs, "Received build-info with empty VCS.")

	// Validate results
	expectedVcsUrl := "https://github.com/jfrog/jfrog-cli-go.git"
	expectedVcsRevision := "b033a0e508bdb52eee25654c9e12db33ff01b8ff"
	buildInfoVcsUrl := buildInfo.Vcs.Url
	buildInfoVcsRevision := buildInfo.Vcs.Revision
	assert.Equal(t, expectedVcsRevision, buildInfoVcsRevision, "Wrong revision")
	assert.Equal(t, expectedVcsUrl, buildInfoVcsUrl, "Wrong url")
	assert.False(t, buildInfo.Issues == nil || len(buildInfo.Issues.AffectedIssues) != 4,
		"Wrong issues number, expected 4 issues, received: %+v", *buildInfo.Issues)
}

func cleanBuildAddGitTest(t *testing.T, baseDir, originalFolder, oldHomeDir, dotGitPath string) {
	tests.RenamePath(dotGitPath, filepath.Join(baseDir, originalFolder), t)
	inttestutils.DeleteBuild(artifactoryDetails.Url, tests.RtBuildName1, artHttpDetails)
	os.Setenv(cliutils.HomeDir, oldHomeDir)
	cleanArtifactoryTest()
}

func TestReadGitConfig(t *testing.T) {
	dotGitPath := getCliDotGitPath(t)
	gitManager := clientutils.NewGitManager(dotGitPath)
	err := gitManager.ReadConfig()
	assert.NoError(t, err, "Failed to read .git config file.")

	workingDir, err := os.Getwd()
	assert.NoError(t, err, "Failed to get current dir.")
	gitExecutor := tests.GitExecutor(workingDir)
	revision, _, err := gitExecutor.GetRevision()
	require.NoError(t, err)
	assert.Equal(t, revision, gitManager.GetRevision(), "Wrong revision")

	url, _, err := gitExecutor.GetUrl()
	require.NoError(t, err)
	if !strings.HasSuffix(url, ".git") {
		url += ".git"
	}

	assert.Equal(t, url, gitManager.GetUrl(), "Wrong url")
}

func uploadFilesAndGetBuildInfo(t *testing.T, buildName, buildNumber, buildUrl string) []byte {
	uploadFiles(t, "upload", "--build-name="+buildName, "--build-number="+buildNumber)

	//publish buildInfo
	publishBuildInfoArgs := []string{"build-publish", buildName, buildNumber}
	if buildUrl != "" {
		publishBuildInfoArgs = append(publishBuildInfoArgs, "--build-url="+buildUrl)
	}
	artifactoryCli.Exec(publishBuildInfoArgs...)

	//validate files are uploaded with the build info name and number
	props := fmt.Sprintf("build.name=%v;build.number=%v", buildName, buildNumber)
	verifyExistInArtifactoryByProps(tests.GetSimpleUploadExpectedRepo1(), tests.RtRepo1+"/*", props, t)

	//download build info
	buildInfoUrl := fmt.Sprintf("%vapi/build/%v/%v", artifactoryDetails.Url, buildName, buildNumber)
	client, err := httpclient.ClientBuilder().Build()
	assert.NoError(t, err)
	_, body, _, err := client.SendGet(buildInfoUrl, false, artHttpDetails)
	assert.NoError(t, err)
	return body
}

func uploadFiles(t *testing.T, args ...string) {
	// Upload files with buildName and buildNumber
	specFile, err := tests.CreateSpec(tests.UploadFlatRecursive)
	assert.NoError(t, err)
	args = append(args, "--spec="+specFile)
	artifactoryCli.Exec(args...)
}

func downloadFiles(t *testing.T, args ...string) {
	// Download files with buildName and buildNumber
	specFile, err := tests.CreateSpec(tests.DownloadAllRepo1TestResources)
	assert.NoError(t, err)
	args = append(args, "--spec="+specFile)
	artifactoryCli.Exec(args...)
}

func TestModuleName(t *testing.T) {
	initArtifactoryTest(t)
	buildName := tests.RtBuildName1
	type command struct {
		execFunc func(t *testing.T, args ...string)
		args     []string
	}

	tests := []struct {
		testName             string
		buildNumber          string
		moduleName           string
		expectedDependencies int
		expectedArtifacts    int
		execCommands         []command
	}{
		{"uploadWithModuleChange", "9", ModuleNameJFrogTest, 0, 9, []command{{uploadFiles, []string{"upload", "--build-name=" + buildName, "--module=" + ModuleNameJFrogTest}}}},
		{"uploadWithoutChange", "10", buildName, 0, 9, []command{{uploadFiles, []string{"upload", "--build-name=" + buildName}}}},
		{"downloadWithModuleChange", "11", ModuleNameJFrogTest, 9, 0, []command{{downloadFiles, []string{"download", "--build-name=" + buildName, "--module=" + ModuleNameJFrogTest}}}},
		{"downloadWithoutModuleChange", "12", buildName, 9, 0, []command{{downloadFiles, []string{"download", "--build-name=" + buildName}}}},
		{"uploadAndDownloadAggregationWithModuleChange", "13", ModuleNameJFrogTest, 9, 9, []command{{uploadFiles, []string{"upload", "--build-name=" + buildName, "--module=" + ModuleNameJFrogTest}}, {downloadFiles, []string{"download", "--build-name=" + buildName, "--module=" + ModuleNameJFrogTest}}}},
		{"uploadAndDownloadAggregationWithoutModuleChange", "14", buildName, 9, 9, []command{{uploadFiles, []string{"upload", "--build-name=" + buildName}}, {downloadFiles, []string{"download", "--build-name=" + buildName}}}},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			for _, exeCommand := range test.execCommands {
				exeCommand.args = append(exeCommand.args, "--build-number="+test.buildNumber)
				exeCommand.execFunc(t, exeCommand.args...)
			}
			artifactoryCli.Exec("bp", buildName, test.buildNumber)
			buildInfo, _ := inttestutils.GetBuildInfo(artifactoryDetails.Url, buildName, test.buildNumber, t, artHttpDetails)
			validateBuildInfo(buildInfo, t, test.expectedDependencies, test.expectedArtifacts, test.moduleName)
		})
	}

	inttestutils.DeleteBuild(artifactoryDetails.Url, buildName, artHttpDetails)
}

func collectDepsAndPublishBuild(badTest buildAddDepsBuildInfoTestParams, useEnvBuildNameAndNumber bool, t *testing.T) {
	noCredsCli := tests.NewJfrogCli(execMain, "jfrog rt", "")
	// Remove old tests data from fs if exists
	err := utils.RemoveBuildDir(badTest.buildName, badTest.buildNumber)
	assert.NoError(t, err)

	command := []string{"bad"}
	if useEnvBuildNameAndNumber {
		os.Setenv(cliutils.BuildName, badTest.buildName)
		os.Setenv(cliutils.BuildNumber, badTest.buildNumber)
		defer os.Unsetenv(cliutils.BuildName)
		defer os.Unsetenv(cliutils.BuildNumber)
	} else {
		command = append(command, badTest.buildName, badTest.buildNumber)
	}

	// Execute tha bad command
	noCredsCli.Exec(append(command, badTest.commandArgs...)...)
	artifactoryCli.Exec("bp", badTest.buildName, badTest.buildNumber)
}

func validateBuildAddDepsBuildInfo(t *testing.T, buildInfoTestParams buildAddDepsBuildInfoTestParams) {
	buildInfo, _ := inttestutils.GetBuildInfo(artifactoryDetails.Url, buildInfoTestParams.buildName, buildInfoTestParams.buildNumber, t, artHttpDetails)
	if buildInfo.Modules == nil || len(buildInfo.Modules) == 0 {
		buildInfoString, _ := json.Marshal(buildInfo)
		// Case no module was not created
		assert.Failf(t, "%s test with the command: \nrt bad %s \nexpected to have module with the following dependencies: \n%s \nbut has no modules: \n%s",
			buildInfoTestParams.description, buildInfoTestParams.commandArgs, buildInfoTestParams.expectedDependencies, buildInfoString)
	}
	// The checksums are ignored when comparing the actual and the expected
	assert.Equalf(t, len(buildInfoTestParams.expectedDependencies), len(buildInfo.Modules[0].Dependencies),
		"%s test with the command: \nrt bad %s  \nexpected to have the following dependencies: \n%s \nbut has: \n%s",
		buildInfoTestParams.description, buildInfoTestParams.commandArgs, buildInfoTestParams.expectedDependencies, dependenciesToPrintableArray(buildInfo.Modules[0].Dependencies))

	for _, expectedDependency := range buildInfoTestParams.expectedDependencies {
		found := false
		for _, actualDependency := range buildInfo.Modules[0].Dependencies {
			if actualDependency.Id == expectedDependency {
				found = true
				break
			}
		}
		// The checksums are ignored when comparing the actual and the expected
		assert.Truef(t, found, "%s test with the command: \nrt bad %s \nexpected to have the following dependencies: \n%s \nbut has: \n%s",
			buildInfoTestParams.description, buildInfoTestParams.commandArgs, buildInfoTestParams.expectedDependencies, dependenciesToPrintableArray(buildInfo.Modules[0].Dependencies))
	}
}

func dependenciesToPrintableArray(dependencies []buildinfo.Dependency) []string {
	ids := []string{}
	for _, dependency := range dependencies {
		ids = append(ids, dependency.Id)
	}
	return ids
}

type buildAddDepsBuildInfoTestParams struct {
	description          string
	commandArgs          []string
	expectedDependencies []string
	buildName            string
	buildNumber          string
	validationFunc       func(*testing.T, buildAddDepsBuildInfoTestParams)
}
