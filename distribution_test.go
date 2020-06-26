package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-cli/inttestutils"
	"github.com/jfrog/jfrog-cli/utils/tests"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
)

const (
	bundleName    = "cli-test-bundle"
	bundleVersion = "10"
)

func InitDistributionTests() {
	*tests.RtDistributionUrl = utils.AddTrailingSlashIfNeeded(*tests.RtDistributionUrl)
	InitArtifactoryTests()
	inttestutils.SendGpgKeys(artHttpDetails)
}

func CleanDistributionTests() {
	inttestutils.DeleteGpgKeys(artHttpDetails)
	CleanArtifactoryTests()
}

func initDistributionTest(t *testing.T) {
	if !*tests.TestDistribution {
		t.Skip("Distribution is not being tested, skipping...")
	}
	// Delete old release bundle
	artifactoryCli.Exec("rbdel", bundleName, bundleVersion, "--site=*", "--delete-from-dist", "--quiet")
	inttestutils.WaitForDeletion(t, bundleName, bundleVersion, artHttpDetails)
}

func cleanDistributionTest(t *testing.T) {
	artifactoryCli.Exec("rbdel", bundleName, bundleVersion, "--site=*", "--delete-from-dist", "--quiet")
	inttestutils.WaitForDeletion(t, bundleName, bundleVersion, artHttpDetails)
	cleanArtifactoryTest()
}

func TestBundleDownload(t *testing.T) {
	initDistributionTest(t)

	// Upload files
	specFile, err := tests.CreateSpec(tests.SplitUploadSpecB)
	assert.NoError(t, err)
	artifactoryCli.Exec("u", "--spec="+specFile)

	// Create and distribute release bundle
	artifactoryCli.Exec("rbc", bundleName, bundleVersion, tests.Repo1+"/data/b1.in", "--sign")
	artifactoryCli.Exec("rbd", bundleName, bundleVersion, "--site=*")
	inttestutils.WaitForDistribution(t, bundleName, bundleVersion, artHttpDetails)

	// Download by bundle version, b2 and b3 should not be downloaded, b1 should
	artifactoryCli.Exec("dl "+tests.Repo1+"/data/* "+tests.Out+fileutils.GetFileSeparator()+"download"+fileutils.GetFileSeparator()+"simple_by_build"+fileutils.GetFileSeparator(), "--bundle="+bundleName+"/"+bundleVersion)

	// Validate files are downloaded by bundle version
	paths, _ := fileutils.ListFilesRecursiveWalkIntoDirSymlink(tests.Out, false)
	err = tests.ValidateListsIdentical(tests.GetBuildSimpleDownload(), paths)
	assert.NoError(t, err)

	// Cleanup
	cleanDistributionTest(t)
}

func TestBundleDownloadUsingSpec(t *testing.T) {
	initDistributionTest(t)

	// Upload files
	specFile, err := tests.CreateSpec(tests.SplitUploadSpecB)
	assert.NoError(t, err)
	artifactoryCli.Exec("u", "--spec="+specFile)
	inttestutils.WaitForDeletion(t, bundleName, bundleVersion, artHttpDetails)

	// Create release bundle
	distributionRules, err := tests.CreateSpec(tests.DistributionRules)
	assert.NoError(t, err)
	artifactoryCli.Exec("rbc", bundleName, bundleVersion, tests.Repo1+"/data/b1.in", "--sign")
	artifactoryCli.Exec("rbd", bundleName, bundleVersion, "--dist-rules="+distributionRules)
	inttestutils.WaitForDistribution(t, bundleName, bundleVersion, artHttpDetails)

	// Download by bundle version, b2 and b3 should not be downloaded, b1 should
	specFile, err = tests.CreateSpec(tests.BundleDownloadSpec)
	assert.NoError(t, err)
	artifactoryCli.Exec("dl", "--spec="+specFile)

	// Validate files are downloaded by bundle version
	paths, _ := fileutils.ListFilesRecursiveWalkIntoDirSymlink(tests.Out, false)
	err = tests.ValidateListsIdentical(tests.GetBuildSimpleDownload(), paths)
	assert.NoError(t, err)

	// Cleanup
	cleanDistributionTest(t)
}

func TestBundleDownloadNoPattern(t *testing.T) {
	initDistributionTest(t)

	// Upload files
	specFile, err := tests.CreateSpec(tests.SplitUploadSpecB)
	assert.NoError(t, err)
	artifactoryCli.Exec("u", "--spec="+specFile)

	// Create release bundle
	artifactoryCli.Exec("rbc", bundleName, bundleVersion, tests.Repo1+"/data/b1.in", "--sign")
	artifactoryCli.Exec("rbd", bundleName, bundleVersion, "--site=*")
	inttestutils.WaitForDistribution(t, bundleName, bundleVersion, artHttpDetails)

	// Download by bundle name and version with pattern "*", b2 and b3 should not be downloaded, b1 should
	artifactoryCli.Exec("dl", "*", "out/download/simple_by_build/data/", "--bundle="+bundleName+"/"+bundleVersion, "--flat")

	// Validate files are downloaded by bundle version
	paths, _ := fileutils.ListFilesRecursiveWalkIntoDirSymlink(tests.Out, false)
	err = tests.ValidateListsIdentical(tests.GetBuildSimpleDownload(), paths)
	assert.NoError(t, err)

	// Download by bundle name and version version without pattern, b2 and b3 should not be downloaded, b1 should
	tests.CleanFileSystem()
	specFile, err = tests.CreateSpec(tests.BundleDownloadSpecNoPattern)
	artifactoryCli.Exec("dl", "--spec="+specFile, "--flat")

	// Validate files are downloaded by bundle version
	paths, _ = fileutils.ListFilesRecursiveWalkIntoDirSymlink(tests.Out, false)
	err = tests.ValidateListsIdentical(tests.GetBuildSimpleDownload(), paths)
	assert.NoError(t, err)

	// Cleanup
	cleanDistributionTest(t)
}

func TestBundleExclusions(t *testing.T) {
	initDistributionTest(t)

	// Upload files
	specFile, err := tests.CreateSpec(tests.SplitUploadSpecB)
	assert.NoError(t, err)
	artifactoryCli.Exec("u", "--spec="+specFile)

	// Create release bundle. Include b1.in and b2.in. Exclude b3.in.
	artifactoryCli.Exec("rbc", bundleName, bundleVersion, tests.Repo1+"/data/b*.in", "--sign", "--exclusions=*b3.in")
	artifactoryCli.Exec("rbd", bundleName, bundleVersion, "--site=*")
	inttestutils.WaitForDistribution(t, bundleName, bundleVersion, artHttpDetails)

	// Download by bundle version, b2 and b3 should not be downloaded, b1 should
	artifactoryCli.Exec("dl "+tests.Repo1+"/data/* "+tests.Out+fileutils.GetFileSeparator()+"download"+fileutils.GetFileSeparator()+"simple_by_build"+fileutils.GetFileSeparator(), "--bundle="+bundleName+"/"+bundleVersion, "--exclusions=*b2.in")

	// Validate files are downloaded by bundle version
	paths, _ := fileutils.ListFilesRecursiveWalkIntoDirSymlink(tests.Out, false)
	err = tests.ValidateListsIdentical(tests.GetBuildSimpleDownload(), paths)
	assert.NoError(t, err)

	// Cleanup
	cleanDistributionTest(t)
}

func TestBundleCopy(t *testing.T) {
	initDistributionTest(t)

	// Upload files
	specFileA, err := tests.CreateSpec(tests.SplitUploadSpecA)
	assert.NoError(t, err)
	specFileB, err := tests.CreateSpec(tests.SplitUploadSpecB)
	assert.NoError(t, err)
	artifactoryCli.Exec("u", "--spec="+specFileB)
	artifactoryCli.Exec("u", "--spec="+specFileA)

	// Create release bundle
	artifactoryCli.Exec("rbc", bundleName, bundleVersion, tests.Repo1+"/data/a*", "--sign")
	artifactoryCli.Exec("rbd", bundleName, bundleVersion, "--site=*")
	inttestutils.WaitForDistribution(t, bundleName, bundleVersion, artHttpDetails)

	// Copy by bundle name and version
	specFile, err := tests.CreateSpec(tests.CopyByBundleSpec)
	assert.NoError(t, err)
	artifactoryCli.Exec("cp", "--spec="+specFile)

	// Validate files are moved by bundle version
	cpMvDlByBuildAssertSpec, err := tests.CreateSpec(tests.CpMvDlByBuildAssertSpec)
	assert.NoError(t, err)
	verifyExistInArtifactory(tests.GetBuildCopyExpected(), cpMvDlByBuildAssertSpec, t)

	// Cleanup
	cleanDistributionTest(t)
}

func TestSetPropsOnBundle(t *testing.T) {
	initDistributionTest(t)

	// Upload a file.
	artifactoryCli.Exec("u", "testsdata/a/a1.in", tests.Repo1+"/a.in")

	// Create release bundle
	artifactoryCli.Exec("rbc", bundleName, bundleVersion, tests.Repo1+"/a.in", "--sign")
	artifactoryCli.Exec("rbd", bundleName, bundleVersion, "--site=*")
	inttestutils.WaitForDistribution(t, bundleName, bundleVersion, artHttpDetails)

	// Set the 'prop=red' property to the file.
	artifactoryCli.Exec("sp", tests.Repo1+"/a.*", "prop=red", "--bundle="+bundleName+"/"+bundleVersion)
	// Now let's change the property value, by searching for the 'prop=red'.
	specFile, err := tests.CreateSpec(tests.SetDeletePropsSpec)
	assert.NoError(t, err)
	artifactoryCli.Exec("sp", "prop=green", "--spec="+specFile, "--bundle="+bundleName+"/"+bundleVersion)

	// Check that prop=green exist on a.in
	checkProperty(t, "a.in")

	cleanDistributionTest(t)
}

func TestCreateBundleWithProps(t *testing.T) {
	initDistributionTest(t)

	// Upload a file.
	artifactoryCli.Exec("u", "testsdata/a/a1.in", tests.Repo1+"/a.in")

	// Create release bundle
	artifactoryCli.Exec("rbc", bundleName, bundleVersion, tests.Repo1+"/a.in", "--sign", "--props=prop=green")
	artifactoryCli.Exec("rbd", bundleName, bundleVersion, "--site=*")
	inttestutils.WaitForDistribution(t, bundleName, bundleVersion, artHttpDetails)

	// Check that prop=green exist on a.in
	checkProperty(t, "a.in")
	cleanDistributionTest(t)
}

func TestSignReleaseBundle(t *testing.T) {
	initDistributionTest(t)

	// Upload files
	specFile, err := tests.CreateSpec(tests.SplitUploadSpecB)
	assert.NoError(t, err)
	artifactoryCli.Exec("u", "--spec="+specFile)

	// Create a release bundle without --sign and make sure it is not signed
	artifactoryCli.Exec("rbc", bundleName, bundleVersion, tests.Repo1+"/data/b1.in")
	distributableResponse := inttestutils.GetLocalBundle(t, bundleName, bundleVersion, artHttpDetails)
	assert.NotNil(t, distributableResponse)
	assert.Equal(t, inttestutils.Open, distributableResponse.State)

	// Sign the release bundle and make sure it is signed
	artifactoryCli.Exec("rbs", bundleName, bundleVersion)
	distributableResponse = inttestutils.GetLocalBundle(t, bundleName, bundleVersion, artHttpDetails)
	assert.NotNil(t, distributableResponse)
	assert.Equal(t, inttestutils.Signed, distributableResponse.State)

	// Cleanup
	cleanDistributionTest(t)
}

func TestBundleDeleteLocal(t *testing.T) {
	initDistributionTest(t)

	// Upload files
	specFile, err := tests.CreateSpec(tests.SplitUploadSpecB)
	assert.NoError(t, err)
	artifactoryCli.Exec("u", "--spec="+specFile)

	// Create a release bundle
	artifactoryCli.Exec("rbc", bundleName, bundleVersion, tests.Repo1+"/data/b1.in", "--sign")
	inttestutils.VerifyLocalBundleExistence(t, bundleName, bundleVersion, true, artHttpDetails)

	// Delete release bundle locally
	artifactoryCli.Exec("rbdel", bundleName, bundleVersion, "--site=*", "--delete-from-dist", "--quiet")
	inttestutils.VerifyLocalBundleExistence(t, bundleName, bundleVersion, false, artHttpDetails)

	// Cleanup
	cleanDistributionTest(t)
}

func TestUpdateReleaseBundle(t *testing.T) {
	initDistributionTest(t)

	// Upload files
	specFile, err := tests.CreateSpec(tests.SplitUploadSpecB)
	assert.NoError(t, err)
	artifactoryCli.Exec("u", "--spec="+specFile)

	// Create a release bundle with b2.in
	artifactoryCli.Exec("rbc", bundleName, bundleVersion, tests.Repo1+"/data/b2.in")
	inttestutils.VerifyLocalBundleExistence(t, bundleName, bundleVersion, true, artHttpDetails)

	// Update release bundle to have b1.in
	artifactoryCli.Exec("rbu", bundleName, bundleVersion, tests.Repo1+"/data/b1.in", "--sign", "--props=prop=green")

	// Distribute release bundle
	artifactoryCli.Exec("rbd", bundleName, bundleVersion, "--site=*")
	inttestutils.WaitForDistribution(t, bundleName, bundleVersion, artHttpDetails)

	// Download by bundle version, b2 and b3 should not be downloaded, b1 should
	artifactoryCli.Exec("dl "+tests.Repo1+"/data/* "+tests.Out+fileutils.GetFileSeparator()+"download"+fileutils.GetFileSeparator()+"simple_by_build"+fileutils.GetFileSeparator(), "--bundle="+bundleName+"/"+bundleVersion)

	// Validate files are downloaded by bundle version
	paths, _ := fileutils.ListFilesRecursiveWalkIntoDirSymlink(tests.Out, false)
	err = tests.ValidateListsIdentical(tests.GetBuildSimpleDownload(), paths)
	assert.NoError(t, err)

	// Check that prop=green exist on b1.in
	checkProperty(t, "b1.in")

	// Cleanup
	cleanDistributionTest(t)
}

func TestCreateBundleText(t *testing.T) {
	initDistributionTest(t)

	// Upload files
	specFile, err := tests.CreateSpec(tests.SplitUploadSpecB)
	assert.NoError(t, err)
	artifactoryCli.Exec("u", "--spec="+specFile)

	// Create a release bundle with release notes and description
	releaseNotesPath := filepath.Join(tests.GetTestResourcesPath(), "distribution", "releasenotes.md")
	description := "thisIsADescription"
	artifactoryCli.Exec("rbc", bundleName, bundleVersion, tests.Repo1+"/data/*", "--release-notes-path="+releaseNotesPath, "--desc="+description)

	// Validate release notes and description
	distributableResponse := inttestutils.GetLocalBundle(t, bundleName, bundleVersion, artHttpDetails)
	if distributableResponse != nil {
		assert.Equal(t, description, distributableResponse.Description)
		releaseNotes, err := ioutil.ReadFile(releaseNotesPath)
		assert.NoError(t, err)
		assert.Equal(t, string(releaseNotes), distributableResponse.ReleaseNotes.Content)
		assert.Equal(t, "markdown", distributableResponse.ReleaseNotes.Syntax)
	}

	cleanDistributionTest(t)
}

// Check that the artifact contains 'prop=green' property
func checkProperty(t *testing.T, artifactName string) {
	resultItems := searchItemsInArtifactory(t)
	assert.NotZero(t, len(resultItems), "No artifacts were found.")
	for _, item := range resultItems {
		properties := item.Properties
		if artifactName != item.Name {
			// If the result item is not what we search for, make sure it doesn't contain the property
			for _, prop := range properties {
				assert.NotEqual(t, "prop", prop.Key)
			}
			continue
		}
		assert.Equal(t, len(properties), 2, "Failed setting properties on item:", item.GetItemRelativePath())
		for _, prop := range properties {
			if prop.Key != "prop" {
				continue
			}
			assert.Equal(t, "prop", prop.Key, "Wrong property key")
			assert.Equal(t, "green", prop.Value, "Wrong property value")
		}
	}
}
