package npm

import (
	gofrogcmd "github.com/jfrog/gofrog/io"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"io"
	"io/ioutil"
)

func Version(executablePath string) ([]byte, error) {

	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()
	defer pipeWriter.Close()
	var npmError error

	configListCmdConfig := createVersionCmdConfig(executablePath, pipeWriter)
	go func() {
		npmError = gofrogcmd.RunCmd(configListCmdConfig)
	}()

	data, err := ioutil.ReadAll(pipeReader)
	if err != nil {
		return nil, errorutils.WrapError(err)
	}

	if npmError != nil {
		return nil, errorutils.WrapError(npmError)
	}

	return data, nil
}

func createVersionCmdConfig(executablePath string, pipeWriter *io.PipeWriter) *NpmConfig {
	return &NpmConfig{
		Npm:       executablePath,
		Command:   []string{"-version"},
		StrWriter: pipeWriter,
		ErrWriter: nil,
	}
}
