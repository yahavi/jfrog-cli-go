package cliutils

const (
	// General CLI constants
	CliVersion           = "1.32.4"
	ClientAgent          = "jfrog-cli-go"
	OnErrorPanic OnError = "panic"

	// CLI base commands constants:
	CmdArtifactory    = "rt"
	CmdBintray        = "bt"
	CmdMissionControl = "mc"
	CmdXray           = "xr"
	CmdCompletion     = "completion"

	// Download
	DownloadMinSplitKb    = 5120
	DownloadSplitCount    = 3
	DownloadMaxSplitCount = 15

	// Common
	Retries = 3

	// Env
	ReportUsage     = "JFROG_CLI_REPORT_USAGE"
	LogLevel        = "JFROG_CLI_LOG_LEVEL"
	OfferConfig     = "JFROG_CLI_OFFER_CONFIG"
	HomeDir         = "JFROG_CLI_HOME_DIR"
	TempDir         = "JFROG_CLI_TEMP_DIR"
	CI              = "CI"
	DependenciesDir = "JFROG_CLI_DEPENDENCIES_DIR"
	BuildName       = "JFROG_CLI_BUILD_NAME"
	BuildNumber     = "JFROG_CLI_BUILD_NUMBER"
	BuildUrl        = "JFROG_CLI_BUILD_URL"
	EnvExclude      = "JFROG_CLI_ENV_EXCLUDE"
	UserAgent       = "JFROG_CLI_USER_AGENT"
	// Deprecated:
	JfrogHomeEnv = "JFROG_CLI_HOME"
)
