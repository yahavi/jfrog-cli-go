package log

import (
	"github.com/jfrog/jfrog-cli-go/utils/cliutils"
	"github.com/jfrog/jfrog-cli-go/utils/config"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func GetCliLogLevel() log.LevelType {
	switch os.Getenv(cliutils.LogLevel) {
	case "ERROR":
		return log.ERROR
	case "WARN":
		return log.WARN
	case "DEBUG":
		return log.DEBUG
	default:
		return log.INFO
	}
}

func SetDefaultLogger() {
	log.SetLogger(log.NewLogger(GetCliLogLevel(), nil))
}

func CreateLogFile() (*os.File, error) {
	logDir, err := config.CreateDirInJfrogHome("logs")
	if err != nil {
		return nil, err
	}

	currentTime := time.Now().Format("2006-01-02.15-04-05")
	pid := os.Getpid()

	fileName := filepath.Join(logDir, "jfrog-cli."+currentTime+"."+strconv.Itoa(pid)+".log")
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errorutils.WrapError(err)
	}

	return file, nil
}

// Closes the log file and resets to the default logger
func CloseLogFile(logFile *os.File) {
	if logFile != nil {
		SetDefaultLogger()
		err := logFile.Close()
		utils.WrapErrorWithMessage(err, "failed closing the log file")
	}
}
