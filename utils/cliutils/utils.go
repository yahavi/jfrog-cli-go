package cliutils

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jfrog/jfrog-cli-go/utils/summary"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/pkg/errors"
	"github.com/ztrue/tracerr"
)

// Error modes (how should the application behave when the WrapError function is invoked):
type OnError string

var cliTempDir string
var cliUserAgent string

func init() {
	// Initialize the temp base-dir path of the CLI executions.
	cliTempDir = os.Getenv(TempDir)
	if cliTempDir == "" {
		cliTempDir = os.TempDir()
	}
	fileutils.SetTempDirBase(cliTempDir)

	// Initialize agent name and version.
	cliUserAgent = os.Getenv(UserAgent)
	if cliUserAgent == "" {
		cliUserAgent = ClientAgent + "/" + CliVersion
	}
}

// Exit codes:
type ExitCode struct {
	Code int
}

var ExitCodeNoError = ExitCode{0}
var ExitCodeError = ExitCode{1}
var ExitCodeFailNoOp = ExitCode{2}
var ExitCodeBuildScan = ExitCode{3}

type CliError struct {
	ExitCode
	TraceErr tracerr.Error
}

func (err CliError) Error() string {
	return err.TraceErr.Error()
}

func PanicOnError(err error) error {
	if err != nil {
		panic(err)
	}
	return err
}

func ExitOnErr(err error) {
	_, _, printStacktraceStr, _ := FindFlag("stacktrace", os.Args)
	printStacktrace, _ := strconv.ParseBool(printStacktraceStr)
	if err, ok := err.(CliError); ok {
		traceExit(err.ExitCode, err, printStacktrace)
	}
	if exitCode := GetExitCode(err, 0, 0, false); exitCode != ExitCodeNoError {
		traceExit(exitCode, err, printStacktrace)
	}
}

func GetCliError(err error, success, failed int, failNoOp bool) error {
	var traceError tracerr.Error
	switch GetExitCode(err, success, failed, failNoOp) {
	case ExitCodeError:
		{
			if err != nil {
				e, ok := err.(tracerr.Error)
				if ok {
					traceError = e
				} else {
					traceError = tracerr.CustomError(err, []tracerr.Frame{})
				}
			}
			return CliError{ExitCodeError, traceError}
		}
	case ExitCodeFailNoOp:
		traceError = tracerr.CustomError(errors.New("No errors, but also no files affected (fail-no-op flag)"), []tracerr.Frame{})
		return CliError{ExitCodeFailNoOp, traceError}
	default:
		return nil
	}
}

func ExitBuildScan(failBuild bool, err error) error {
	if failBuild && err != nil {
		traceError := tracerr.CustomError(errors.New("Build Scan Failed"), []tracerr.Frame{})
		return CliError{ExitCodeBuildScan, traceError}
	}

	return nil
}

func GetExitCode(err error, success, failed int, failNoOp bool) ExitCode {
	// Error occurred - Return 1
	if err != nil || failed > 0 {
		return ExitCodeError
	}
	// No errors, but also no files affected - Return 2 if failNoOp
	if success == 0 && failNoOp {
		return ExitCodeFailNoOp
	}
	// Otherwise - Return 0
	return ExitCodeNoError
}

func traceExit(exitCode ExitCode, err error, printStacktrace bool) {
	if err != nil && len(err.Error()) > 0 && printStacktrace {
		tracerr.PrintSourceColor(err)
	}
	os.Exit(exitCode.Code)
}

// func RunCommand(c *cli.Context, cmd func(c *cli.Context) error) error {
// 	err := cmd(c)
// 	if (err != nil && c.Bool("stacktrace")) {
// 		tracerr.PrintSourceColor(err)
// 	}
// 	return err
// }

// Print summary report.
// The given error will pass through and be returned as is if no other errors are raised.
func PrintSummaryReport(success, failed int, err error) error {
	summaryReport := summary.New(err)
	summaryReport.Totals.Success = success
	summaryReport.Totals.Failure = failed
	if err == nil && summaryReport.Totals.Failure != 0 {
		summaryReport.Status = summary.Failure
	}
	content, mErr := summaryReport.Marshal()
	if errorutils.WrapError(mErr) != nil {
		log.Error(mErr)
		return err
	}
	log.Output(utils.IndentJson(content))

	return err
}

func PrintHelpAndReturnError(msg string, context *cli.Context) error {
	log.Error(msg + " " + GetDocumentationMessage())
	cli.ShowCommandHelp(context, context.Command.Name)
	return errors.New(msg)
}

func InteractiveConfirm(message string) bool {
	var confirm string
	fmt.Print(message + " (y/n): ")
	fmt.Scanln(&confirm)
	return confirmAnswer(confirm)
}

func confirmAnswer(answer string) bool {
	answer = strings.ToLower(answer)
	return answer == "y" || answer == "yes"
}

func GetVersion() string {
	return CliVersion
}

func GetConfigVersion() string {
	return "1"
}

func GetDocumentationMessage() string {
	return "You can read the documentation at https://www.jfrog.com/confluence/display/CLI/JFrog+CLI"
}

func SumTrueValues(boolArr []bool) int {
	counter := 0
	for _, val := range boolArr {
		counter += utils.Bool2Int(val)
	}
	return counter
}

func SpecVarsStringToMap(rawVars string) map[string]string {
	if len(rawVars) == 0 {
		return nil
	}
	varCandidates := strings.Split(rawVars, ";")
	varsList := []string{}
	for _, v := range varCandidates {
		if len(varsList) > 0 && isEndsWithEscapeChar(varsList[len(varsList)-1]) {
			currentLastVar := varsList[len(varsList)-1]
			varsList[len(varsList)-1] = strings.TrimSuffix(currentLastVar, "\\") + ";" + v
			continue
		}
		varsList = append(varsList, v)
	}
	return varsAsMap(varsList)
}

func isEndsWithEscapeChar(lastVar string) bool {
	return strings.HasSuffix(lastVar, "\\")
}

func varsAsMap(vars []string) map[string]string {
	result := map[string]string{}
	for _, v := range vars {
		keyVal := strings.SplitN(v, "=", 2)
		if len(keyVal) != 2 {
			continue
		}
		result[keyVal[0]] = keyVal[1]
	}
	return result
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// Return the path of CLI temp dir.
// This path should be persistent, meaning - should not be cleared at the end of a CLI run.
func GetCliPersistentTempDirPath() string {
	return cliTempDir
}

func GetUserAgent() string {
	return cliUserAgent
}

type Credentials interface {
	SetUser(string)
	SetPassword(string)
	GetUser() string
	GetPassword() string
}
