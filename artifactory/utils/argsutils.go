package utils

import (
	"strconv"

	"github.com/jfrog/jfrog-cli-go/utils/cliutils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

func ExtractNpmOptionsFromArgs(args []string) (threads int, cleanArgs []string, buildConfig *BuildConfiguration, err error) {
	threads = 3
	// Extract threads information from the args.
	flagIndex, valueIndex, numOfThreads, err := cliutils.FindFlag("--threads", args)
	if err != nil {
		return
	}
	cliutils.RemoveFlagFromCommand(&args, flagIndex, valueIndex)
	if numOfThreads != "" {
		threads, err = strconv.Atoi(numOfThreads)
		if err != nil {
			err = errorutils.WrapError(err)
			return
		}
	}

	cleanArgs, buildConfig, err = ExtractBuildDetailsFromArgs(args)
	return
}

func ExtractBuildDetailsFromArgs(args []string) (cleanArgs []string, buildConfig *BuildConfiguration, err error) {
	var flagIndex, valueIndex int
	buildConfig = &BuildConfiguration{}
	cleanArgs = append([]string(nil), args...)

	// Extract build-info information from the args.
	flagIndex, valueIndex, buildConfig.BuildName, err = cliutils.FindFlag("--build-name", cleanArgs)
	if err != nil {
		return
	}
	cliutils.RemoveFlagFromCommand(&cleanArgs, flagIndex, valueIndex)

	flagIndex, valueIndex, buildConfig.BuildNumber, err = cliutils.FindFlag("--build-number", cleanArgs)
	if err != nil {
		return
	}
	cliutils.RemoveFlagFromCommand(&cleanArgs, flagIndex, valueIndex)

	// Retreive build name and build number from env if both missing
	buildConfig.BuildName, buildConfig.BuildNumber = GetBuildNameAndNumber(buildConfig.BuildName, buildConfig.BuildNumber)
	flagIndex, valueIndex, buildConfig.Module, err = cliutils.FindFlag("--module", cleanArgs)
	if err != nil {
		return
	}
	cliutils.RemoveFlagFromCommand(&cleanArgs, flagIndex, valueIndex)
	err = ValidateBuildParams(buildConfig)
	return
}
