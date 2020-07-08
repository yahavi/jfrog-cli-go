package tests

import (
	"path/filepath"

	"github.com/jfrog/jfrog-cli/artifactory/commands/generic"
	clientutils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

const (
	ArchiveEntriesDownload                 = "archive_entries_download_spec.json"
	ArchiveEntriesUpload                   = "archive_entries_upload_spec.json"
	BuildAddDepsDoubleSpec                 = "build_add_deps_double_spec.json"
	BuildAddDepsSpec                       = "build_add_deps_simple_spec.json"
	BuildDownloadSpec                      = "build_download_spec.json"
	BuildDownloadSpecNoBuildNumber         = "build_download_spec_no_build_number.json"
	BuildDownloadSpecNoBuildNumberWithSort = "build_download_spec_no_build_number_with_sort.json"
	BuildDownloadSpecNoPattern             = "build_download_spec_no_pattern.json"
	BundleDownloadSpec                     = "bundle_download_spec.json"
	BundleDownloadSpecNoPattern            = "bundle_download_spec_no_pattern.json"
	CopyByBuildPatternAllSpec              = "move_copy_delete_by_build_pattern_all_spec.json"
	CopyByBuildSpec                        = "move_copy_delete_by_build_spec.json"
	CopyByBundleSpec                       = "copy_by_bundle_spec.json"
	CopyByBundleAssertSpec                 = "copy_by_bundle_assert_spec.json"
	CopyItemsSpec                          = "copy_items_spec.json"
	CopyMoveSimpleSpec                     = "copy_move_simple.json"
	CpMvDlByBuildAssertSpec                = "copy_by_build_assert_spec.json"
	DebianTestRepositoryConfig             = "debian_test_repository_config.json"
	DebianUploadSpec                       = "upload_debian_spec.json"
	DeleteSimpleSpec                       = "delete_simple_spec.json"
	DeleteSpec                             = "delete_spec.json"
	DeleteSpecWildcardInRepo               = "delete_spec_wildcard.json"
	DelSpecExclude                         = "delete_spec_exclude.json"
	DelSpecExclusions                      = "delete_spec_exclusions.json"
	DistributionRepoConfig1                = "dist_repository_config1.json"
	DistributionRepoConfig2                = "dist_repository_config2.json"
	DistributionRules                      = "distribution_rules.json"
	DistributionSetDeletePropsSpec         = "dist_set_delete_props_spec.json"
	DistributionUploadSpecA                = "dist_upload_spec_a.json"
	DistributionUploadSpecB                = "dist_upload_spec_b.json"
	DownloadAllRepo1TestResources          = "download_all_repo1_test_resources.json"
	DownloadEmptyDirs                      = "download_empty_dir_spec.json"
	DownloadModFileGo                      = "downloadmodfile_go.json"
	DownloadModOfDependencyGo              = "downloadmodofdependency_go.json"
	DownloadSpecExclude                    = "download_spec_exclude.json"
	DownloadSpecExclusions                 = "download_spec_exclusions.json"
	DownloadWildcardRepo                   = "download_wildcard_repo.json"
	GitLfsAssertSpec                       = "git_lfs_assert_spec.json"
	GitLfsTestRepositoryConfig             = "git_lfs_test_repository_config.json"
	GoLocalRepositoryConfig                = "go_local_repository_config.json"
	GradleConfig                           = "gradle.yaml"
	GradleRemoteRepositoryConfig           = "gradle_remote_repository_config.json"
	GradleRepositoryConfig                 = "gradle_repository_config.json"
	GradleServerIDConfig                   = "gradle_server_id.yaml"
	GradleServerIDUsesPluginConfig         = "gradle_server_id_uses_plugin.yaml"
	GradleUsernamePasswordTemplate         = "gradle_user_pass_template.yaml"
	HttpsProxyEnvVar                       = "PROXY_HTTPS_PORT"
	MavenConfig                            = "maven.yaml"
	MavenRemoteRepositoryConfig            = "maven_remote_repository_config.json"
	MavenRepositoryConfig1                 = "maven_repository_config1.json"
	MavenRepositoryConfig2                 = "maven_repository_config2.json"
	MavenServerIDConfig                    = "maven_server_id.yaml"
	MavenUsernamePasswordTemplate          = "maven_user_pass_template.yaml"
	MoveCopySpecExclude                    = "move_copy_spec_exclude.json"
	MoveCopySpecExclusions                 = "move_copy_spec_exclusions.json"
	Repo2RepositoryConfig                  = "repo2_repository_config.json"
	NpmLocalRepositoryConfig               = "npm_local_repository_config.json"
	NpmRemoteRepositoryConfig              = "npm_remote_repository_config.json"
	NugetRemoteRepo                        = "jfrog-cli-tests-nuget-remote-repo"
	Out                                    = "out"
	PypiRemoteRepositoryConfig             = "pypi_remote_repository_config.json"
	PypiVirtualRepositoryConfig            = "pypi_virtual_repository_config.json"
	RepoDetailsUrl                         = "api/repositories/"
	RtServerId                             = "rtTestServerId"
	SearchAllGradle                        = "search_all_gradle.json"
	SearchAllMaven                         = "search_all_maven.json"
	SearchAllRepo1                         = "search_all_repo1.json"
	SearchGo                               = "search_go.json"
	SearchDistRepoByInSuffix               = "search_dist_repo_by_in_suffix.json"
	SearchRepo1ByInSuffix                  = "search_repo1_by_in_suffix.json"
	SearchRepo1TestResources               = "search_repo1_test_resources.json"
	SearchRepo2                            = "search_repo2.json"
	SearchSimplePlaceholders               = "search_simple_placeholders.json"
	SearchTargetInRepo2                    = "search_target_in_repo2.json"
	SearchTxt                              = "search_txt.json"
	SetDeletePropsSpec                     = "set_delete_props_spec.json"
	Repo1RepositoryConfig                  = "repo1_repository_config.json"
	SplitUploadSpecA                       = "upload_split_spec_a.json"
	SplitUploadSpecB                       = "upload_split_spec_b.json"
	Temp                                   = "tmp"
	UploadEmptyDirs                        = "upload_empty_dir_spec.json"
	UploadFileWithParenthesesSpec          = "upload_file_with_parentheses.json"
	UploadFlatNonRecursive                 = "upload_flat_non_recursive.json"
	UploadFlatRecursive                    = "upload_flat_recursive.json"
	UploadMultipleFileSpecs                = "upload_multiple_file_specs.json"
	UploadSimplePlaceholders               = "upload_simple_placeholders.json"
	UploadSpecExclude                      = "upload_spec_exclude.json"
	UploadSpecExcludeRegex                 = "upload_spec_exclude_regex.json"
	UploadTempWildcard                     = "upload_temp_wildcard.json"
	UploadWithPropsSpec                    = "upload_with_props_spec.json"
	VirtualRepositoryConfig                = "specs_virtual_repository_config.json"
	WinBuildAddDepsSpec                    = "win_simple_build_add_deps_spec.json"
	WinSimpleDownloadSpec                  = "win_simple_download_spec.json"
	WinSimpleUploadSpec                    = "win_simple_upload_spec.json"
	ReplicationTempCreate                  = "replication_push_create.json"
)

var (
	// Repositories
	DistRepo1        = "cli-tests-dist1"
	DistRepo2        = "cli-tests-dist2"
	GoRepo           = "cli-tests-go"
	GradleRepo       = "cli-tests-gradle"
	MvnRemoteRepo    = "cli-tests-mvn-remote"
	GradleRemoteRepo = "cli-tests-gradle-remote"
	MvnRepo1         = "cli-tests-mvn1"
	MvnRepo2         = "cli-tests-mvn2"
	NpmRepo          = "cli-tests-npm"
	NpmRemoteRepo    = "cli-tests-npm-remote"
	PypiRemoteRepo   = "cli-tests-pypi-remote"
	PypiVirtualRepo  = "cli-tests-pypi-virtual"
	RtDebianRepo     = "cli-tests-debian"
	RtLfsRepo        = "cli-tests-lfs"
	RtRepo1          = "cli-tests-rt1"
	RtRepo2          = "cli-tests-rt2"
	RtVirtualRepo    = "cli-tests-rt-virtual"
	// These are not actual repositories. These patterns are meant to be used in both Repo1 and Repo2.
	RtRepo1And2            = "cli-tests-rt*"
	RtRepo1And2Placeholder = "cli-tests-rt(*)"

	BundleName      = "cli-tests-dist-bundle"
	DockerBuildName = "cli-tests-docker-build"
	DockerImageName = "cli-tests-docker-image"
	DotnetBuildName = "cli-tests-dotnet-build"
	GoBuildName     = "cli-tests-go-build"
	GradleBuildName = "cli-tests-gradle-build"
	NpmBuildName    = "cli-tests-npm-build"
	NuGetBuildName  = "cli-tests-nuget-build"
	PipBuildName    = "cli-tests-pip-build"
	RtBuildName1    = "cli-tests-rt-build1"
	RtBuildName2    = "cli-tests-rt-build2"
)

func GetTxtUploadExpectedRepo1() []string {
	return []string{
		RtRepo1 + "/cliTestFile.txt",
	}
}

func GetSimpleUploadExpectedRepo1() []string {
	return []string{
		RtRepo1 + "/test_resources/a3.in",
		RtRepo1 + "/test_resources/a1.in",
		RtRepo1 + "/test_resources/a2.in",
		RtRepo1 + "/test_resources/b2.in",
		RtRepo1 + "/test_resources/b3.in",
		RtRepo1 + "/test_resources/b1.in",
		RtRepo1 + "/test_resources/c2.in",
		RtRepo1 + "/test_resources/c1.in",
		RtRepo1 + "/test_resources/c3.in",
	}
}

func GetSimpleWildcardUploadExpectedRepo1() []string {
	return []string{
		RtRepo1 + "/upload_simple_wildcard/github.com/github.in",
	}
}

func GetSimpleUploadSpecialCharNoRegexExpectedRepo1() []string {
	return []string{
		RtRepo1 + "/a1.in",
	}
}

func GetSimpleUploadSpecialCharNoRegexExpectedRepo2() []string {
	return []string{
		RtRepo2 + "/a1.in",
	}
}

func GetSimpleUploadSpecialCharNoRegexExpected2filesRepo1() []string {
	return []string{
		RtRepo1 + "/a1.in",
		RtRepo1 + "/a2.in",
	}
}

func GetSimpleUploadSpecialCharNoRegexExpected2filesRepo2() []string {
	return []string{
		RtRepo2 + "/a1.in",
		RtRepo2 + "/a2.in",
	}
}

func GetUploadSpecExcludeRepo1() []string {
	return []string{
		RtRepo1 + "/a1.in",
		RtRepo1 + "/b1.in",
		RtRepo1 + "/c2.in",
		RtRepo1 + "/c3.in",
	}
}

func GetUploadDebianExpected() []string {
	return []string{
		RtDebianRepo + "/data/a1.in",
		RtDebianRepo + "/data/a2.in",
		RtDebianRepo + "/data/a3.in",
	}
}

func GetSingleFileCopy() []string {
	return []string{
		RtRepo2 + "/path/a1.in",
	}
}

func GetSingleFileCopyFullPath() []string {
	return []string{
		RtRepo2 + "/path/inner/a1.in",
	}
}

func GetSingleInnerFileCopyFullPath() []string {
	return []string{
		RtRepo2 + "/path/path/inner/a1.in",
	}
}

func GetFolderCopyTwice() []string {
	return []string{
		RtRepo2 + "/path/inner/a1.in",
		RtRepo2 + "/path/path/inner/a1.in",
	}
}

func GetFolderCopyIntoFolder() []string {
	return []string{
		RtRepo2 + "/path/path/inner/a1.in",
	}
}

func GetSingleDirectoryCopyFlat() []string {
	return []string{
		RtRepo2 + "/inner/a1.in",
	}
}

func GetAnyItemCopy() []string {
	return []string{
		RtRepo2 + "/path/inner/a1.in",
		RtRepo2 + "/someFile",
	}
}

func GetAnyItemCopyRecursive() []string {
	return []string{
		RtRepo2 + "/a/b/a1.in",
		RtRepo2 + "/aFile",
	}
}

func GetCopyFolderRename() []string {
	return []string{
		RtRepo2 + "/newPath/inner/a1.in",
	}
}

func GetAnyItemCopyUsingSpec() []string {
	return []string{
		RtRepo2 + "/a1.in",
	}
}

func GetExplodeUploadExpectedRepo1() []string {
	return []string{
		RtRepo1 + "/a/a3.in",
		RtRepo1 + "/a/a1.in",
		RtRepo1 + "/a/a2.in",
		RtRepo1 + "/a/b/b1.in",
		RtRepo1 + "/a/b/b2.in",
		RtRepo1 + "/a/b/b3.in",
		RtRepo1 + "/a/b/c/c1.in",
		RtRepo1 + "/a/b/c/c2.in",
		RtRepo1 + "/a/b/c/c3.in",
	}
}

func GetCopyFileNameWithParentheses() []string {
	return []string{
		RtRepo2 + "/testsdata/b/(/(.in",
		RtRepo2 + "/testsdata/b/(b/(b.in",
		RtRepo2 + "/testsdata/b/)b/)b.in",
		RtRepo2 + "/testsdata/b/b(/b(.in",
		RtRepo2 + "/testsdata/b/b)/b).in",
		RtRepo2 + "/testsdata/b/(b)/(b).in",
		RtRepo2 + "/testsdata/b/)b)/)b).in",
		RtRepo2 + "/(/b(.in",
		RtRepo2 + "/()/(b.in",
		RtRepo2 + "/()/testsdata/b/(b)/(b).in",
		RtRepo2 + "/(/testsdata/b/(/(.in.zip",
		RtRepo2 + "/(/in-b(",
		RtRepo2 + "/(/b(.in-up",
		RtRepo2 + "/c/(.in.zip",
	}
}
func GetUploadFileNameWithParentheses() []string {
	return []string{
		RtRepo1 + "/(.in",
		RtRepo1 + "/(b.in",
		RtRepo1 + "/)b.in",
		RtRepo1 + "/b(.in",
		RtRepo1 + "/b).in",
		RtRepo1 + "/(b).in",
		RtRepo1 + "/)b).in",
		RtRepo1 + "/(new)/testsdata/b/(/(.in",
		RtRepo1 + "/(new)/testsdata/b/(b/(b.in",
		RtRepo1 + "/(new)/testsdata/b/b(/b(.in",
		RtRepo1 + "/new)/testsdata/b/b)/b).in",
		RtRepo1 + "/new)/testsdata/b/(b)/(b).in",
		RtRepo1 + "/(new/testsdata/b/)b)/)b).in",
		RtRepo1 + "/(new/testsdata/b/)b/)b.in",
	}
}

func GetMoveCopySpecExpected() []string {
	return []string{
		RtRepo2 + "/copy_move_target/a1.in",
		RtRepo2 + "/copy_move_target/a2.in",
		RtRepo2 + "/copy_move_target/a3.in",
		RtRepo2 + "/copy_move_target/b/b1.in",
		RtRepo2 + "/copy_move_target/b/b2.in",
		RtRepo2 + "/copy_move_target/b/b3.in",
		RtRepo2 + "/copy_move_target/b/c/c1.in",
		RtRepo2 + "/copy_move_target/b/c/c2.in",
		RtRepo2 + "/copy_move_target/b/c/c3.in",
	}
}

func GetRepo1TestResourcesExpected() []string {
	return []string{
		RtRepo1 + "/test_resources/a1.in",
		RtRepo1 + "/test_resources/a2.in",
		RtRepo1 + "/test_resources/a3.in",
		RtRepo1 + "/test_resources/b/b1.in",
		RtRepo1 + "/test_resources/b/b2.in",
		RtRepo1 + "/test_resources/b/b3.in",
		RtRepo1 + "/test_resources/b/c/c1.in",
		RtRepo1 + "/test_resources/b/c/c2.in",
		RtRepo1 + "/test_resources/b/c/c3.in",
	}
}

func GetBuildBeforeCopyExpected() []string {
	return GetBuildBeforeMoveExpected()
}

func GetBuildCopyExpected() []string {
	return []string{
		RtRepo1 + "/data/a1.in",
		RtRepo1 + "/data/a2.in",
		RtRepo1 + "/data/a3.in",
		RtRepo1 + "/data/b1.in",
		RtRepo1 + "/data/b2.in",
		RtRepo1 + "/data/b3.in",
		RtRepo2 + "/data/a1.in",
		RtRepo2 + "/data/a2.in",
		RtRepo2 + "/data/a3.in",
	}
}

func GetBundleCopyExpected() []string {
	return []string{
		DistRepo1 + "/data/a1.in",
		DistRepo1 + "/data/a2.in",
		DistRepo1 + "/data/a3.in",
		DistRepo1 + "/data/b1.in",
		DistRepo1 + "/data/b2.in",
		DistRepo1 + "/data/b3.in",
		DistRepo2 + "/data/a1.in",
		DistRepo2 + "/data/a2.in",
		DistRepo2 + "/data/a3.in",
	}
}

func GetGitLfsExpected() []string {
	return []string{
		RtLfsRepo + "/objects/4b/f4/4bf4c8c0fef3f5c8cf6f255d1c784377138588c0a9abe57e440bce3ccb350c2e",
	}
}

func GetBuildBeforeMoveExpected() []string {
	return []string{
		RtRepo1 + "/data/b1.in",
		RtRepo1 + "/data/b2.in",
		RtRepo1 + "/data/b3.in",
		RtRepo1 + "/data/a1.in",
		RtRepo1 + "/data/a2.in",
		RtRepo1 + "/data/a3.in",
	}
}

func GetBuildMoveExpected() []string {
	return []string{
		RtRepo1 + "/data/b1.in",
		RtRepo1 + "/data/b2.in",
		RtRepo1 + "/data/b3.in",
		RtRepo2 + "/data/a1.in",
		RtRepo2 + "/data/a2.in",
		RtRepo2 + "/data/a3.in",
	}
}

func GetBuildCopyExclude() []string {
	return []string{
		RtRepo1 + "/data/a1.in",
		RtRepo1 + "/data/a2.in",
		RtRepo1 + "/data/a3.in",
		RtRepo1 + "/data/b1.in",
		RtRepo1 + "/data/b2.in",
		RtRepo1 + "/data/b3.in",
		RtRepo2 + "/data/a1.in",
		RtRepo2 + "/data/a2.in",
		RtRepo2 + "/data/a3.in",
	}
}

func GetBuildDeleteExpected() []string {
	return []string{
		RtRepo1 + "/data/b1.in",
		RtRepo1 + "/data/b2.in",
		RtRepo1 + "/data/b3.in",
	}
}

func GetExtractedDownload() []string {
	return []string{
		filepath.Join(Out, "randFile"),
		filepath.Join(Out, "concurrent.tar.gz"),
	}
}

func GetFileWithParenthesesDownload() []string {
	return []string{
		filepath.Join(Out, "testsdata"),
		filepath.Join(Out, "testsdata/b"),
		filepath.Join(Out, "testsdata/b/("),
		filepath.Join(Out, "testsdata/b/(/(.in"),
		filepath.Join(Out, "testsdata/b/(b"),
		filepath.Join(Out, "testsdata/b/(b/(b.in"),
		filepath.Join(Out, "testsdata/b/(b)"),
		filepath.Join(Out, "testsdata/b/(b)/(b).in"),
		filepath.Join(Out, "testsdata/b/)b"),
		filepath.Join(Out, "testsdata/b/)b/)b.in"),
		filepath.Join(Out, "testsdata/b/)b)"),
		filepath.Join(Out, "testsdata/b/)b)/)b).in"),
		filepath.Join(Out, "testsdata/b/b("),
		filepath.Join(Out, "testsdata/b/b(/b(.in"),
		filepath.Join(Out, "testsdata/b/b)"),
		filepath.Join(Out, "testsdata/b/b)/b).in"),
	}
}

func GetVirtualDownloadExpected() []string {
	return []string{
		filepath.Join(Out, "a/a1.in"),
		filepath.Join(Out, "a/a2.in"),
		filepath.Join(Out, "a/a3.in"),
		filepath.Join(Out, "a/b/b1.in"),
		filepath.Join(Out, "a/b/b2.in"),
		filepath.Join(Out, "a/b/b3.in"),
		filepath.Join(Out, "a/b/c/c1.in"),
		filepath.Join(Out, "a/b/c/c2.in"),
		filepath.Join(Out, "a/b/c/c3.in"),
	}
}

func GetExpectedSyncDeletesDownloadStep2() []string {
	localPathPrefix := filepath.Join("syncDir", "testsdata", "a")
	return []string{
		filepath.Join(Out, localPathPrefix, "a1.in"),
		filepath.Join(Out, localPathPrefix, "a2.in"),
		filepath.Join(Out, localPathPrefix, "a3.in"),
		filepath.Join(Out, localPathPrefix, "b/b1.in"),
		filepath.Join(Out, localPathPrefix, "b/b2.in"),
		filepath.Join(Out, localPathPrefix, "b/b3.in"),
		filepath.Join(Out, localPathPrefix, "b/c/c1.in"),
		filepath.Join(Out, localPathPrefix, "b/c/c2.in"),
		filepath.Join(Out, localPathPrefix, "b/c/c3.in"),
	}
}

func GetExpectedSyncDeletesDownloadStep3() []string {
	return []string{
		filepath.Join(Out, "a1.in"),
		filepath.Join(Out, "a2.in"),
		filepath.Join(Out, "a3.in"),
		filepath.Join(Out, "b1.in"),
		filepath.Join(Out, "b2.in"),
		filepath.Join(Out, "b3.in"),
		filepath.Join(Out, "c1.in"),
		filepath.Join(Out, "c2.in"),
		filepath.Join(Out, "c3.in"),
	}
}

func GetExpectedSyncDeletesDownloadStep4() []string {
	return []string{
		filepath.Join(Out, "a2.in"),
		filepath.Join(Out, "b2.in"),
		filepath.Join(Out, "c2.in"),
	}
}

func GetSyncExpectedDeletesDownloadStep5() []string {
	localPathPrefix := filepath.Join("syncDir", "testsdata", "a")
	return []string{
		filepath.Join(Out, localPathPrefix, "a1.in"),
		filepath.Join(Out, localPathPrefix, "a2.in"),
		filepath.Join(Out, localPathPrefix, "a3.in"),
		filepath.Join(Out, localPathPrefix, "b/b1.in"),
		filepath.Join(Out, localPathPrefix, "b/b2.in"),
		filepath.Join(Out, localPathPrefix, "b/b3.in"),
	}
}

func GetSyncExpectedDeletesDownloadStep6() []string {
	localPathPrefix := "/syncDir/testsdata/archives/"
	return []string{
		RtRepo1 + localPathPrefix + "a.zip",
		RtRepo1 + localPathPrefix + "b.zip",
		RtRepo1 + localPathPrefix + "c.zip",
		RtRepo1 + localPathPrefix + "d.zip",
	}
}

func GetSyncExpectedDeletesDownloadStep7() []string {
	localPathPrefix := filepath.Join("syncDir", "testsdata", "archives")
	return []string{
		filepath.Join(Out, localPathPrefix, "a.zip"),
		filepath.Join(Out, localPathPrefix, "b.zip"),
		filepath.Join(Out, localPathPrefix, "c.zip"),
		filepath.Join(Out, localPathPrefix, "d.zip"),
	}
}

func GetDownloadWildcardRepo() []string {
	return []string{
		RtRepo1 + "/path/a1.in",
		RtRepo2 + "/path/a2.in",
	}
}

func GetDownloadUnicode() []string {
	return []string{
		RtRepo1 + "/testsdata/unicode/dirλrec၇tory/文件.in",
		RtRepo1 + "/testsdata/unicode/dirλrec၇tory/aȩ.ȥ1",
		RtRepo1 + "/testsdata/unicode/Ԙחלص⩥.in",
	}
}

func GetBuildDownload() []string {
	return []string{
		Out,
		filepath.Join(Out, "download"),
		filepath.Join(Out, "download", "aql_by_build"),
		filepath.Join(Out, "download", "aql_by_build", "data"),
		filepath.Join(Out, "download", "aql_by_build", "data", "a1.in"),
		filepath.Join(Out, "download", "aql_by_build", "data", "a2.in"),
		filepath.Join(Out, "download", "aql_by_build", "data", "a3.in"),
	}
}

func GetBuildDownloadDoesntExist() []string {
	return []string{
		Out,
	}
}

func GetBuildDownloadByShaAndBuild() []string {
	return []string{
		Out,
		filepath.Join(Out, "download"),
		filepath.Join(Out, "download", "aql_by_build"),
		filepath.Join(Out, "download", "aql_by_build", "data"),
		filepath.Join(Out, "download", "aql_by_build", "data", "a10.in"),
	}
}

func GetBuildDownloadByShaAndBuildName() []string {
	return []string{
		Out,
		filepath.Join(Out, "download"),
		filepath.Join(Out, "download", "aql_by_build"),
		filepath.Join(Out, "download", "aql_by_build", "data"),
		filepath.Join(Out, "download", "aql_by_build", "data", "a11.in"),
	}
}

func GetBuildSimpleDownload() []string {
	return []string{
		Out,
		filepath.Join(Out, "download"),
		filepath.Join(Out, "download", "simple_by_build"),
		filepath.Join(Out, "download", "simple_by_build", "data"),
		filepath.Join(Out, "download", "simple_by_build", "data", "b1.in"),
	}
}

func GetBuildSimpleDownloadNoPattern() []string {
	return []string{
		Out,
		filepath.Join(Out, "download"),
		filepath.Join(Out, "download", "simple_by_build"),
		filepath.Join(Out, "download", "simple_by_build", "data"),
		filepath.Join(Out, "download", "simple_by_build", "data", "a1.in"),
		filepath.Join(Out, "download", "simple_by_build", "data", "a2.in"),
		filepath.Join(Out, "download", "simple_by_build", "data", "a3.in"),
	}
}

func GetBuildExcludeDownload() []string {
	return []string{
		Out,
		filepath.Join(Out, "download"),
		filepath.Join(Out, "download", "aql_by_artifacts"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data", "a3.in"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data", "b1.in"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data", "b2.in"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data", "b3.in"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data", "c1.in"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data", "c3.in"),
	}
}

func GetBuildExcludeDownloadBySpec() []string {
	return []string{
		Out,
		filepath.Join(Out, "download"),
		filepath.Join(Out, "download", "aql_by_artifacts"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data", "a2.in"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data", "b1.in"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data", "b2.in"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data", "b3.in"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data", "c1.in"),
		filepath.Join(Out, "download", "aql_by_artifacts", "data", "c3.in"),
	}
}

func GetCleanBuild() []string {
	return []string{
		filepath.Join(Out, "clean-build"),
		filepath.Join(Out, "clean-build", "data"),
		filepath.Join(Out, "clean-build", "data", "b1.in"),
		filepath.Join(Out, "clean-build", "data", "b2.in"),
		filepath.Join(Out, "clean-build", "data", "b3.in"),
	}
}

func GetMultipleFileSpecs() []string {
	return []string{
		RtRepo1 + "/multiple/a1.out",
		RtRepo1 + "/multiple/properties/testsdata/a/b/b2.in",
	}
}

func GetSimplePlaceholders() []string {
	return []string{
		RtRepo1 + "/simple_placeholders/a-in.out",
		RtRepo1 + "/simple_placeholders/b/b-in.out",
		RtRepo1 + "/simple_placeholders/b/c/c-in.out",
	}
}

func GetSimpleDelete() []string {
	return []string{
		RtRepo1 + "/test_resources/a1.in",
		RtRepo1 + "/test_resources/a2.in",
		RtRepo1 + "/test_resources/a3.in",
	}
}

func GetDeleteFolderWithWildcard() []string {
	return []string{
		RtRepo1 + "/test_resources/a1.in",
		RtRepo1 + "/test_resources/a2.in",
		RtRepo1 + "/test_resources/a3.in",
		RtRepo1 + "/test_resources/b/b1.in",
		RtRepo1 + "/test_resources/b/b2.in",
		RtRepo1 + "/test_resources/b/b3.in",
	}
}

func GetSearchIncludeDirsFiles() []generic.SearchResult {
	return []generic.SearchResult{
		{
			Path:  RtRepo1 + "/",
			Type:  "folder",
			Props: make(map[string][]string, 0),
			Size:  0,
		},
		{
			Path:  RtRepo1 + "/data",
			Type:  "folder",
			Props: make(map[string][]string, 0),
			Size:  0,
		},
		{
			Path:  RtRepo1 + "/data/testsdata",
			Type:  "folder",
			Props: make(map[string][]string, 0),
			Size:  0,
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a",
			Type:  "folder",
			Props: make(map[string][]string, 0),
			Size:  0,
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/a1.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  7,
			Sha1:  "507ac63c6b0f650fb6f36b5621e70ebca3b0965c",
			Md5:   "65298e78fe5883eee82056bc6d0d7f4c",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/a2.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  7,
			Sha1:  "de2f31d77e2c2b1039a806f21b0c5f3243e45588",
			Md5:   "28f9732cb82a2d11760e38614246ad6d",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/a3.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  7,
			Sha1:  "29d38faccfe74dee60d0142a716e8ea6fad67b49",
			Md5:   "73c046196302ff7218d47046cf3c0501",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b",
			Type:  "folder",
			Props: make(map[string][]string, 0),
			Size:  0,
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b/b1.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  9,
			Sha1:  "954cf8f3f75c41f18540bb38460910b4f0074e6f",
			Md5:   "4f5561d29422374e40bd97d28b12cf35",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b/b2.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  9,
			Sha1:  "3b60b837e037568856bedc1dd4952d17b3f06972",
			Md5:   "6931271be1e5f98e36bdc7a05097407b",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b/b3.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  9,
			Sha1:  "ec6420d2b5f708283619b25e68f9ddd351f555fe",
			Md5:   "305b21db102cf3a3d2d8c3f7e9584dba",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b/c",
			Type:  "folder",
			Props: make(map[string][]string, 0),
			Size:  0,
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b/c/c1.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  11,
			Sha1:  "063041114949bf19f6fe7508aef639640e7edaac",
			Md5:   "e53098d3d8ee1f5eb38c2ec3c783ef3d",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b/c/c2.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  11,
			Sha1:  "a4f912be11e7d1d346e34c300e6d4b90e136896e",
			Md5:   "82b6d565393a3fd1cc4778b1d53c0664",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b/c/c3.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  11,
			Sha1:  "2d6ee506188db9b816a6bfb79c5df562fc1d8658",
			Md5:   "d8020b86244956f647cf1beff5acdb90",
		},
	}
}

func GetSearchNotIncludeDirsFiles() []generic.SearchResult {
	return []generic.SearchResult{
		{
			Path:  RtRepo1 + "/data/testsdata/a/a1.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  7,
			Sha1:  "507ac63c6b0f650fb6f36b5621e70ebca3b0965c",
			Md5:   "65298e78fe5883eee82056bc6d0d7f4c",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/a2.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  7,
			Sha1:  "de2f31d77e2c2b1039a806f21b0c5f3243e45588",
			Md5:   "28f9732cb82a2d11760e38614246ad6d",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/a3.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  7,
			Sha1:  "29d38faccfe74dee60d0142a716e8ea6fad67b49",
			Md5:   "73c046196302ff7218d47046cf3c0501",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b/b1.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  9,
			Sha1:  "954cf8f3f75c41f18540bb38460910b4f0074e6f",
			Md5:   "4f5561d29422374e40bd97d28b12cf35",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b/b2.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  9,
			Sha1:  "3b60b837e037568856bedc1dd4952d17b3f06972",
			Md5:   "6931271be1e5f98e36bdc7a05097407b",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b/b3.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  9,
			Sha1:  "ec6420d2b5f708283619b25e68f9ddd351f555fe",
			Md5:   "305b21db102cf3a3d2d8c3f7e9584dba",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b/c/c1.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  11,
			Sha1:  "063041114949bf19f6fe7508aef639640e7edaac",
			Md5:   "e53098d3d8ee1f5eb38c2ec3c783ef3d",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b/c/c2.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  11,
			Sha1:  "a4f912be11e7d1d346e34c300e6d4b90e136896e",
			Md5:   "82b6d565393a3fd1cc4778b1d53c0664",
		},
		{
			Path:  RtRepo1 + "/data/testsdata/a/b/c/c3.in",
			Type:  "file",
			Props: make(map[string][]string, 0),
			Size:  11,
			Sha1:  "2d6ee506188db9b816a6bfb79c5df562fc1d8658",
			Md5:   "d8020b86244956f647cf1beff5acdb90",
		},
	}
}

func GetSearchPropsStep1() []generic.SearchResult {
	return []generic.SearchResult{
		{
			Path: RtRepo1 + "/a/a3.in",
			Type: "file",
			Size: 7,
			Sha1: "29d38faccfe74dee60d0142a716e8ea6fad67b49",
			Md5:  "73c046196302ff7218d47046cf3c0501",
			Props: map[string][]string{
				"a": {"1"},
				"b": {"3"},
				"c": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/b2.in",
			Type: "file",
			Size: 9,
			Sha1: "3b60b837e037568856bedc1dd4952d17b3f06972",
			Md5:  "6931271be1e5f98e36bdc7a05097407b",
			Props: map[string][]string{
				"b": {"1"},
				"c": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/b3.in",
			Type: "file",
			Size: 9,
			Sha1: "ec6420d2b5f708283619b25e68f9ddd351f555fe",
			Md5:  "305b21db102cf3a3d2d8c3f7e9584dba",
			Props: map[string][]string{
				"a": {"1"},
				"b": {"2"},
				"c": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/c/c2.in",
			Type: "file",
			Size: 11,
			Sha1: "a4f912be11e7d1d346e34c300e6d4b90e136896e",
			Md5:  "82b6d565393a3fd1cc4778b1d53c0664",
			Props: map[string][]string{
				"c": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/c/c3.in",
			Type: "file",
			Size: 11,
			Sha1: "2d6ee506188db9b816a6bfb79c5df562fc1d8658",
			Md5:  "d8020b86244956f647cf1beff5acdb90",
			Props: map[string][]string{
				"c": {"3"},
			},
		},
	}
}

func GetSearchPropsStep2() []generic.SearchResult {
	return []generic.SearchResult{
		{
			Path: RtRepo1 + "/a/a1.in",
			Type: "file",
			Size: 7,
			Sha1: "507ac63c6b0f650fb6f36b5621e70ebca3b0965c",
			Md5:  "65298e78fe5883eee82056bc6d0d7f4c",
			Props: map[string][]string{
				"a": {"2"},
				"b": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/a2.in",
			Type: "file",
			Size: 7,
			Sha1: "de2f31d77e2c2b1039a806f21b0c5f3243e45588",
			Md5:  "28f9732cb82a2d11760e38614246ad6d",
			Props: map[string][]string{
				"a": {"1"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/b1.in",
			Type: "file",
			Size: 9,
			Sha1: "954cf8f3f75c41f18540bb38460910b4f0074e6f",
			Md5:  "4f5561d29422374e40bd97d28b12cf35",
			Props: map[string][]string{
				"a": {"1"},
				"c": {"5"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/c/c1.in",
			Type: "file",
			Size: 11,
			Sha1: "063041114949bf19f6fe7508aef639640e7edaac",
			Md5:  "e53098d3d8ee1f5eb38c2ec3c783ef3d",
			Props: map[string][]string{
				"b": {"1"},
			},
		},
	}
}

func GetSearchPropsStep3() []generic.SearchResult {
	return []generic.SearchResult{
		{
			Path: RtRepo1 + "/a/a1.in",
			Type: "file",
			Size: 7,
			Sha1: "507ac63c6b0f650fb6f36b5621e70ebca3b0965c",
			Md5:  "65298e78fe5883eee82056bc6d0d7f4c",
			Props: map[string][]string{
				"a": {"2"},
				"b": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/a2.in",
			Type: "file",
			Size: 7,
			Sha1: "de2f31d77e2c2b1039a806f21b0c5f3243e45588",
			Md5:  "28f9732cb82a2d11760e38614246ad6d",
			Props: map[string][]string{
				"a": {"1"},
			},
		},
		{
			Path: RtRepo1 + "/a/a3.in",
			Type: "file",
			Size: 7,
			Sha1: "29d38faccfe74dee60d0142a716e8ea6fad67b49",
			Md5:  "73c046196302ff7218d47046cf3c0501",
			Props: map[string][]string{
				"a": {"1"},
				"b": {"3"},
				"c": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/b1.in",
			Type: "file",
			Size: 9,
			Sha1: "954cf8f3f75c41f18540bb38460910b4f0074e6f",
			Md5:  "4f5561d29422374e40bd97d28b12cf35",
			Props: map[string][]string{
				"a": {"1"},
				"c": {"5"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/b2.in",
			Type: "file",
			Size: 9,
			Sha1: "3b60b837e037568856bedc1dd4952d17b3f06972",
			Md5:  "6931271be1e5f98e36bdc7a05097407b",
			Props: map[string][]string{
				"b": {"1"},
				"c": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/c/c1.in",
			Type: "file",
			Size: 11,
			Sha1: "063041114949bf19f6fe7508aef639640e7edaac",
			Md5:  "e53098d3d8ee1f5eb38c2ec3c783ef3d",
			Props: map[string][]string{
				"b": {"1"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/c/c2.in",
			Type: "file",
			Size: 11,
			Sha1: "a4f912be11e7d1d346e34c300e6d4b90e136896e",
			Md5:  "82b6d565393a3fd1cc4778b1d53c0664",
			Props: map[string][]string{
				"c": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/c/c3.in",
			Type: "file",
			Size: 11,
			Sha1: "2d6ee506188db9b816a6bfb79c5df562fc1d8658",
			Md5:  "d8020b86244956f647cf1beff5acdb90",
			Props: map[string][]string{
				"c": {"3"},
			},
		},
	}
}

func GetSearchPropsStep4() []generic.SearchResult {
	return []generic.SearchResult{
		{
			Path: RtRepo1 + "/a/a3.in",
			Type: "file",
			Size: 7,
			Sha1: "29d38faccfe74dee60d0142a716e8ea6fad67b49",
			Md5:  "73c046196302ff7218d47046cf3c0501",
			Props: map[string][]string{
				"a": {"1"},
				"b": {"3"},
				"c": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/b2.in",
			Type: "file",
			Size: 9,
			Sha1: "3b60b837e037568856bedc1dd4952d17b3f06972",
			Md5:  "6931271be1e5f98e36bdc7a05097407b",
			Props: map[string][]string{
				"b": {"1"},
				"c": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/c/c2.in",
			Type: "file",
			Size: 11,
			Sha1: "a4f912be11e7d1d346e34c300e6d4b90e136896e",
			Md5:  "82b6d565393a3fd1cc4778b1d53c0664",
			Props: map[string][]string{
				"c": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/c/c3.in",
			Type: "file",
			Size: 11,
			Sha1: "2d6ee506188db9b816a6bfb79c5df562fc1d8658",
			Md5:  "d8020b86244956f647cf1beff5acdb90",
			Props: map[string][]string{
				"c": {"3"},
			},
		},
	}
}

func GetSearchPropsStep5() []generic.SearchResult {
	return make([]generic.SearchResult, 0)
}

func GetSearchPropsStep6() []generic.SearchResult {
	return []generic.SearchResult{
		{
			Path: RtRepo1 + "/a/b/c/c1.in",
			Type: "file",
			Size: 11,
			Sha1: "063041114949bf19f6fe7508aef639640e7edaac",
			Md5:  "e53098d3d8ee1f5eb38c2ec3c783ef3d",
			Props: map[string][]string{
				"b": {"1"},
			},
		},
	}
}

func GetSearchResultAfterDeleteByPropsStep1() []generic.SearchResult {
	return []generic.SearchResult{
		{
			Path: RtRepo1 + "/a/a1.in",
			Type: "file",
			Size: 7,
			Sha1: "507ac63c6b0f650fb6f36b5621e70ebca3b0965c",
			Md5:  "65298e78fe5883eee82056bc6d0d7f4c",
			Props: map[string][]string{
				"a": {"2"},
				"b": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/a2.in",
			Type: "file",
			Size: 7,
			Sha1: "de2f31d77e2c2b1039a806f21b0c5f3243e45588",
			Md5:  "28f9732cb82a2d11760e38614246ad6d",
			Props: map[string][]string{
				"a": {"1"},
			},
		},
		{
			Path: RtRepo1 + "/a/a3.in",
			Type: "file",
			Size: 7,
			Sha1: "29d38faccfe74dee60d0142a716e8ea6fad67b49",
			Md5:  "73c046196302ff7218d47046cf3c0501",
			Props: map[string][]string{
				"a": {"1"},
				"b": {"3"},
				"c": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/b2.in",
			Type: "file",
			Size: 9,
			Sha1: "3b60b837e037568856bedc1dd4952d17b3f06972",
			Md5:  "6931271be1e5f98e36bdc7a05097407b",
			Props: map[string][]string{
				"b": {"1"},
				"c": {"3"},
				"D": {"5"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/b3.in",
			Type: "file",
			Size: 9,
			Sha1: "ec6420d2b5f708283619b25e68f9ddd351f555fe",
			Md5:  "305b21db102cf3a3d2d8c3f7e9584dba",
			Props: map[string][]string{
				"a": {"1"},
				"b": {"2"},
				"c": {"3"},
				"D": {"5"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/c/c2.in",
			Type: "file",
			Size: 11,
			Sha1: "a4f912be11e7d1d346e34c300e6d4b90e136896e",
			Md5:  "82b6d565393a3fd1cc4778b1d53c0664",
			Props: map[string][]string{
				"c": {"3"},
				"D": {"2"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/c/c3.in",
			Type: "file",
			Size: 11,
			Sha1: "2d6ee506188db9b816a6bfb79c5df562fc1d8658",
			Md5:  "d8020b86244956f647cf1beff5acdb90",
			Props: map[string][]string{
				"c": {"3"},
				"D": {"2"},
			},
		},
	}
}

func GetSearchResultAfterDeleteByPropsStep2() []generic.SearchResult {
	return []generic.SearchResult{
		{
			Path: RtRepo1 + "/a/a1.in",
			Type: "file",
			Size: 7,
			Sha1: "507ac63c6b0f650fb6f36b5621e70ebca3b0965c",
			Md5:  "65298e78fe5883eee82056bc6d0d7f4c",
			Props: map[string][]string{
				"a": {"2"},
				"b": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/a2.in",
			Type: "file",
			Size: 7,
			Sha1: "de2f31d77e2c2b1039a806f21b0c5f3243e45588",
			Md5:  "28f9732cb82a2d11760e38614246ad6d",
			Props: map[string][]string{
				"a": {"1"},
			},
		},
		{
			Path: RtRepo1 + "/a/a3.in",
			Type: "file",
			Size: 7,
			Sha1: "29d38faccfe74dee60d0142a716e8ea6fad67b49",
			Md5:  "73c046196302ff7218d47046cf3c0501",
			Props: map[string][]string{
				"a": {"1"},
				"b": {"3"},
				"c": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/b/b3.in",
			Type: "file",
			Size: 9,
			Sha1: "ec6420d2b5f708283619b25e68f9ddd351f555fe",
			Md5:  "305b21db102cf3a3d2d8c3f7e9584dba",
			Props: map[string][]string{
				"a": {"1"},
				"b": {"2"},
				"c": {"3"},
				"D": {"5"},
			},
		},
	}
}

func GetSearchResultAfterDeleteByPropsStep3() []generic.SearchResult {
	return []generic.SearchResult{
		{
			Path: RtRepo1 + "/a/a1.in",
			Type: "file",
			Size: 7,
			Sha1: "507ac63c6b0f650fb6f36b5621e70ebca3b0965c",
			Md5:  "65298e78fe5883eee82056bc6d0d7f4c",
			Props: map[string][]string{
				"a": {"2"},
				"b": {"3"},
			},
		},
		{
			Path: RtRepo1 + "/a/a3.in",
			Type: "file",
			Size: 7,
			Sha1: "29d38faccfe74dee60d0142a716e8ea6fad67b49",
			Md5:  "73c046196302ff7218d47046cf3c0501",
			Props: map[string][]string{
				"a": {"1"},
				"b": {"3"},
				"c": {"3"},
			},
		},
	}
}

func GetMavenDeployedArtifacts() []string {
	return []string{
		MvnRepo1 + "/org/jfrog/cli-test/1.0/cli-test-1.0.jar",
		MvnRepo1 + "/org/jfrog/cli-test/1.0/cli-test-1.0.pom",
	}
}

func GetGradleDeployedArtifacts() []string {
	return []string{
		GradleRepo + "/minimal-example/1.0/minimal-example-1.0.jar",
	}
}

func GetNpmDeployedScopedArtifacts() []string {
	return []string{
		NpmRepo + "/@jscope/jfrog-cli-tests/-/jfrog-cli-tests-1.0.0.tgz",
	}
}
func GetNpmDeployedArtifacts() []string {
	return []string{
		NpmRepo + "/jfrog-cli-tests/-/jfrog-cli-tests-1.0.0.tgz",
	}
}

func GetSortAndLimit() []string {
	return []string{
		Out,
		filepath.Join(Out, "download"),
		filepath.Join(Out, "download", "sort_limit"),
		filepath.Join(Out, "download", "sort_limit", "data"),
		filepath.Join(Out, "download", "sort_limit", "data", "a1.in"),
		filepath.Join(Out, "download", "sort_limit", "data", "b"),
		filepath.Join(Out, "download", "sort_limit", "data", "b", "c"),
		filepath.Join(Out, "download", "sort_limit", "data", "b", "c", "c1.in"),
		filepath.Join(Out, "download", "sort_limit", "data", "b", "c", "c2.in"),
		filepath.Join(Out, "download", "sort_limit", "data", "b", "c", "c3.in"),
	}
}

func GetBuildDownloadByShaAndBuildNameWithSort() []string {
	return []string{
		filepath.Join(Out, "download", "sort_limit_by_build"),
		filepath.Join(Out, "download", "sort_limit_by_build", "data"),
		filepath.Join(Out, "download", "sort_limit_by_build", "data", "a11.in"),
	}
}

func GetBuildArchiveEntriesDownloadCli() []string {
	return []string{
		Out,
		filepath.Join(Out, "a.zip"),
		filepath.Join(Out, "b.zip"),
		filepath.Join(Out, "c.zip"),
	}
}

func GetBuildArchiveEntriesSpecificPathDownload() []string {
	return []string{
		Out,
		filepath.Join(Out, "b.zip"),
	}
}

func GetBuildArchiveEntriesDownloadSpec() []string {
	return []string{
		Out,
		filepath.Join(Out, "d.zip"),
	}
}

func GetWinCompatibility() []string {
	return []string{
		Out,
		filepath.Join(Out, "win"),
		filepath.Join(Out, "win", "a2.in"),
		filepath.Join(Out, "win", "b1.in"),
		filepath.Join(Out, "win", "b3.in"),
	}
}

func GetUploadExpectedRepo1SyncDeleteStep1() []string {
	return []string{
		RtRepo1 + "/syncDir/testsdata/a/a3.in",
		RtRepo1 + "/syncDir/testsdata/a/a1.in",
		RtRepo1 + "/syncDir/testsdata/a/a2.in",
		RtRepo1 + "/syncDir/testsdata/a/b/b1.in",
		RtRepo1 + "/syncDir/testsdata/a/b/b2.in",
		RtRepo1 + "/syncDir/testsdata/a/b/b3.in",
		RtRepo1 + "/syncDir/testsdata/a/b/c/c1.in",
		RtRepo1 + "/syncDir/testsdata/a/b/c/c2.in",
		RtRepo1 + "/syncDir/testsdata/a/b/c/c3.in",
	}
}

func GetUploadExpectedRepo1SyncDeleteStep2() []string {
	return []string{
		RtRepo1 + "/syncDir/testsdata/a/a3.in",
		RtRepo1 + "/syncDir/testsdata/a/a1.in",
		RtRepo1 + "/syncDir/testsdata/a/a2.in",
		RtRepo1 + "/syncDir/testsdata/a/b/b1.in",
		RtRepo1 + "/syncDir/testsdata/a/b/c/c1.in",
	}
}

func GetUploadExpectedRepo1SyncDeleteStep3() []string {
	return []string{
		RtRepo1 + "/syncDir/a.zip",
		RtRepo1 + "/syncDir/b.zip",
		RtRepo1 + "/syncDir/c.zip",
		RtRepo1 + "/syncDir/d.zip",
	}
}

func GetReplicationConfig() []clientutils.ReplicationParams {
	return []clientutils.ReplicationParams{
		{
			Url:                    *RtUrl,
			Username:               "admin",
			Password:               "password",
			CronExp:                "0 0 12 * * ?",
			RepoKey:                RtRepo1,
			EnableEventReplication: false,
			SocketTimeoutMillis:    15000,
			Enabled:                true,
			SyncDeletes:            true,
			SyncProperties:         true,
			SyncStatistics:         false,
			PathPrefix:             "/my/path",
		},
	}
}
