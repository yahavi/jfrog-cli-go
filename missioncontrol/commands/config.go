package commands

import (
	"errors"
	"github.com/jfrog/jfrog-cli-go/utils/config"
	"github.com/jfrog/jfrog-cli-go/utils/ioutils"
	"github.com/jfrog/jfrog-cli-go/utils/lock"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/url"
	"sync"
)

// Internal golang locking for the same process.
var mutex sync.Mutex

func GetConfig() (*config.MissionControlDetails, error) {
	return config.ReadMissionControlConf()
}

func ShowConfig() error {
	details, err := config.ReadMissionControlConf()
	if err != nil {
		return err
	}
	if details.Url != "" {
		log.Output("Url: " + details.Url)
	}
	if details.User != "" {
		log.Output("User: " + details.User)
	}
	if details.Password != "" {
		log.Output("Password: ***")
	}
	return nil
}

func ClearConfig() {
	config.SaveMissionControlConf(new(config.MissionControlDetails))
}

func Config(details, defaultDetails *config.MissionControlDetails, interactive bool) (conf *config.MissionControlDetails, err error) {
	mutex.Lock()
	lockFile, err := lock.CreateLock()
	defer mutex.Unlock()
	defer lockFile.Unlock()

	if err != nil {
		return nil, err
	}

	allowUsingSavedPassword := true
	conf = details
	if conf == nil {
		conf = new(config.MissionControlDetails)
	}
	if interactive {
		if defaultDetails == nil {
			defaultDetails, err = config.ReadMissionControlConf()
			if err != nil {
				return
			}
		}
		if conf.Url == "" {
			ioutils.ScanFromConsole("Mission Control URL", &conf.Url, defaultDetails.Url)
			var u *url.URL
			u, err = url.Parse(conf.Url)
			err = errorutils.WrapError(err)
			if err != nil {
				return
			}
			if u.Scheme != "http" && u.Scheme != "https" {
				err = errorutils.WrapError(errors.New("URL scheme is not valid " + u.Scheme))
				if err != nil {
					return
				}
			}
			allowUsingSavedPassword = false
		}
		ioutils.ReadCredentialsFromConsole(conf, defaultDetails, allowUsingSavedPassword)
	}
	conf.Url = utils.AddTrailingSlashIfNeeded(conf.Url)
	config.SaveMissionControlConf(conf)
	return
}

type ConfigFlags struct {
	MissionControlDetails *config.MissionControlDetails
	Interactive           bool
}
