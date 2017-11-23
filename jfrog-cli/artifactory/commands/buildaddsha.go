package commands

import (
	"github.com/jfrogdev/jfrog-cli-go/jfrog-cli/artifactory/utils/buildinfo"
	"github.com/jfrogdev/jfrog-cli-go/jfrog-cli/artifactory/utils"
	"github.com/jfrogdev/jfrog-cli-go/jfrog-client/utils/log"
)

func BuildAddArtifact(buildName, buildNumber, name, sha1 string) (err error) {
	if err = utils.SaveBuildGeneralDetails(buildName, buildNumber); err != nil {
		return
	}

	populateFunc := func(partial *buildinfo.Partial) {
		checksum := &buildinfo.Checksum{Sha1: sha1}
		partial.Artifacts = []buildinfo.Artifacts{{Name: name, Checksum: checksum}}
	}
	err = utils.SavePartialBuildInfo(buildName, buildNumber, populateFunc)
	log.Info("Successfully added", name, "to build info.")
	return
}
