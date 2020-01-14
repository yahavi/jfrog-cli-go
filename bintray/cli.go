package bintray

import (
	"errors"
	"github.com/jfrog/jfrog-cli-go/utils/ioutils"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jfrog/jfrog-cli-go/bintray/commands"
	accesskeysdoc "github.com/jfrog/jfrog-cli-go/docs/bintray/accesskeys"
	configdocs "github.com/jfrog/jfrog-cli-go/docs/bintray/config"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/downloadfile"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/downloadver"
	entitlementsdocs "github.com/jfrog/jfrog-cli-go/docs/bintray/entitlements"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/gpgsignfile"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/gpgsignver"
	logsdocs "github.com/jfrog/jfrog-cli-go/docs/bintray/logs"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/packagecreate"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/packagedelete"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/packageshow"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/packageupdate"
	streamdocs "github.com/jfrog/jfrog-cli-go/docs/bintray/stream"
	uploaddocs "github.com/jfrog/jfrog-cli-go/docs/bintray/upload"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/urlsign"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/versioncreate"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/versiondelete"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/versionpublish"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/versionshow"
	"github.com/jfrog/jfrog-cli-go/docs/bintray/versionupdate"
	"github.com/jfrog/jfrog-cli-go/docs/common"
	"github.com/jfrog/jfrog-cli-go/utils/cliutils"
	"github.com/jfrog/jfrog-cli-go/utils/config"
	"github.com/jfrog/jfrog-client-go/bintray"
	"github.com/jfrog/jfrog-client-go/bintray/auth"
	"github.com/jfrog/jfrog-client-go/bintray/services"
	"github.com/jfrog/jfrog-client-go/bintray/services/accesskeys"
	"github.com/jfrog/jfrog-client-go/bintray/services/entitlements"
	"github.com/jfrog/jfrog-client-go/bintray/services/packages"
	"github.com/jfrog/jfrog-client-go/bintray/services/url"
	"github.com/jfrog/jfrog-client-go/bintray/services/utils"
	"github.com/jfrog/jfrog-client-go/bintray/services/versions"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
)

func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:         "config",
			Flags:        getConfigFlags(),
			Aliases:      []string{"c"},
			Usage:        configdocs.Description,
			HelpName:     common.CreateUsage("bt config", configdocs.Description, configdocs.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return configure(c)
			},
		},
		{
			Name:         "upload",
			Flags:        getUploadFlags(),
			Aliases:      []string{"u"},
			Usage:        uploaddocs.Description,
			HelpName:     common.CreateUsage("bt upload", uploaddocs.Description, uploaddocs.Usage),
			UsageText:    uploaddocs.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return upload(c)
			},
		},
		{
			Name:         "download-file",
			Flags:        getDownloadFileFlags(),
			Aliases:      []string{"dlf"},
			Usage:        downloadfile.Description,
			HelpName:     common.CreateUsage("bt download-file", downloadfile.Description, downloadfile.Usage),
			UsageText:    downloadfile.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return downloadFile(c)
			},
		},
		{
			Name:         "download-ver",
			Flags:        getDownloadVersionFlags(),
			Aliases:      []string{"dlv"},
			Usage:        downloadver.Description,
			HelpName:     common.CreateUsage("bt download-ver", downloadver.Description, downloadver.Usage),
			UsageText:    downloadver.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return downloadVersion(c)
			},
		},
		{
			Name:         "package-show",
			Flags:        getFlags(),
			Aliases:      []string{"ps"},
			Usage:        packageshow.Description,
			HelpName:     common.CreateUsage("bt package-show", packageshow.Description, packageshow.Usage),
			UsageText:    packageshow.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return showPackage(c)
			},
		},
		{
			Name:         "package-create",
			Flags:        getCreateAndUpdatePackageFlags(),
			Aliases:      []string{"pc"},
			Usage:        packagecreate.Description,
			HelpName:     common.CreateUsage("bt package-create", packagecreate.Description, packagecreate.Usage),
			UsageText:    packagecreate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return createPackage(c)
			},
		},
		{
			Name:         "package-update",
			Flags:        getCreateAndUpdatePackageFlags(),
			Aliases:      []string{"pu"},
			Usage:        packageupdate.Description,
			HelpName:     common.CreateUsage("bt package-update", packageupdate.Description, packageupdate.Usage),
			UsageText:    packageupdate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return updatePackage(c)
			},
		},
		{
			Name:         "package-delete",
			Flags:        getDeletePackageAndVersionFlags(),
			Aliases:      []string{"pd"},
			Usage:        packagedelete.Description,
			HelpName:     common.CreateUsage("bt package-delete", packagedelete.Description, packagedelete.Usage),
			UsageText:    packagedelete.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return deletePackage(c)
			},
		},
		{
			Name:         "version-show",
			Flags:        getFlags(),
			Aliases:      []string{"vs"},
			Usage:        versionshow.Description,
			HelpName:     common.CreateUsage("bt version-show", versionshow.Description, versionshow.Usage),
			UsageText:    versionshow.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return showVersion(c)
			},
		},
		{
			Name:         "version-create",
			Flags:        getCreateAndUpdateVersionFlags(),
			Aliases:      []string{"vc"},
			Usage:        versioncreate.Description,
			HelpName:     common.CreateUsage("bt version-create", versioncreate.Description, versioncreate.Usage),
			UsageText:    versioncreate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return createVersion(c)
			},
		},
		{
			Name:         "version-update",
			Flags:        getCreateAndUpdateVersionFlags(),
			Aliases:      []string{"vu"},
			Usage:        versionupdate.Description,
			HelpName:     common.CreateUsage("bt version-update", versionupdate.Description, versionupdate.Usage),
			UsageText:    versionupdate.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return updateVersion(c)
			},
		},
		{
			Name:         "version-delete",
			Flags:        getDeletePackageAndVersionFlags(),
			Aliases:      []string{"vd"},
			Usage:        versiondelete.Description,
			HelpName:     common.CreateUsage("bt version-delete", versiondelete.Description, versiondelete.Usage),
			UsageText:    versiondelete.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return deleteVersion(c)
			},
		},
		{
			Name:         "version-publish",
			Flags:        getFlags(),
			Aliases:      []string{"vp"},
			Usage:        versionpublish.Description,
			HelpName:     common.CreateUsage("bt version-publish", versionpublish.Description, versionpublish.Usage),
			UsageText:    versionpublish.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return publishVersion(c)
			},
		},
		{
			Name:         "entitlements",
			Flags:        getEntitlementsFlags(),
			Aliases:      []string{"ent"},
			Usage:        entitlementsdocs.Description,
			HelpName:     common.CreateUsage("bt entitlements", entitlementsdocs.Description, entitlementsdocs.Usage),
			UsageText:    entitlementsdocs.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return handleEntitlements(c)
			},
		},
		{
			Name:         "access-keys",
			Flags:        getAccessKeysFlags(),
			Aliases:      []string{"acc-keys"},
			Usage:        accesskeysdoc.Description,
			HelpName:     common.CreateUsage("bt access-keys", accesskeysdoc.Description, accesskeysdoc.Usage),
			UsageText:    accesskeysdoc.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return accessKeys(c)
			},
		},
		{
			Name:         "url-sign",
			Flags:        getUrlSigningFlags(),
			Aliases:      []string{"us"},
			Usage:        urlsign.Description,
			HelpName:     common.CreateUsage("bt url-sign", urlsign.Description, urlsign.Usage),
			UsageText:    urlsign.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return signUrl(c)
			},
		},
		{
			Name:         "gpg-sign-file",
			Flags:        getGpgSigningFlags(),
			Aliases:      []string{"gsf"},
			Usage:        gpgsignfile.Description,
			HelpName:     common.CreateUsage("bt gpg-sign-file", gpgsignfile.Description, gpgsignfile.Usage),
			UsageText:    gpgsignfile.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return gpgSignFile(c)
			},
		},
		{
			Name:         "gpg-sign-ver",
			Flags:        getGpgSigningFlags(),
			Aliases:      []string{"gsv"},
			Usage:        gpgsignver.Description,
			HelpName:     common.CreateUsage("bt gpg-sign-ver", gpgsignver.Description, gpgsignver.Usage),
			UsageText:    gpgsignver.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return gpgSignVersion(c)
			},
		},
		{
			Name:         "logs",
			Flags:        getFlags(),
			Aliases:      []string{"l"},
			Usage:        logsdocs.Description,
			HelpName:     common.CreateUsage("bt logs", logsdocs.Description, logsdocs.Usage),
			UsageText:    logsdocs.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return logs(c)
			},
		},
		{
			Name:         "stream",
			Flags:        getStreamFlags(),
			Aliases:      []string{"st"},
			Usage:        streamdocs.Description,
			HelpName:     common.CreateUsage("bt stream", streamdocs.Description, streamdocs.Usage),
			UsageText:    streamdocs.Arguments,
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(c *cli.Context) error {
				return stream(c)
			},
		},
	}
}

func getFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "[Optional] Bintray username. If not set, the subject sent as part of the command argument is used for authentication.` `",
		},
		cli.StringFlag{
			Name:  "key",
			Value: "",
			Usage: "[Mandatory] Bintray API key` `",
		},
		cli.BoolFlag{
			Name:  "stacktrace",
			Usage: "[Default: false] Set to true to print the error stacktrace for all errors.` `",
		},
	}
}

func getStreamFlags() []cli.Flag {
	return append(getFlags(), []cli.Flag{
		cli.StringFlag{
			Name:  "include",
			Value: "",
			Usage: "[Optional] List of events type in the form of \"value1;value2;...\" leave empty to include all.` `",
		},
	}...)
}

func getConfigFlags() []cli.Flag {
	flags := []cli.Flag{
		cli.BoolTFlag{
			Name:  "interactive",
			Usage: "[Default: true] Set to false if you do not want the config command to be interactive.",
		},
	}
	flags = append(flags, getFlags()...)
	return append(flags, cli.StringFlag{
		Name:  "licenses",
		Value: "",
		Usage: "[Optional] Default package licenses in the form of Apache-2.0,GPL-3.0...` `",
	})
}

func getPackageFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "licenses",
			Value: "",
			Usage: "[Mandatory for OSS] Package licenses in the form of Apache-2.0,GPL-3.0...` `",
		},
		cli.StringFlag{
			Name:  "vcs-url",
			Value: "",
			Usage: "[Mandatory for OSS] Package VCS URL.` `",
		},
		cli.BoolFlag{
			Name:  "pub-dn",
			Usage: "[Default: false] Public download numbers.",
		},
		cli.BoolTFlag{
			Name:  "pub-stats",
			Usage: "[Default: true] Public statistics.",
		},
		cli.StringFlag{
			Name:  "desc",
			Value: "",
			Usage: "[Optional] Package description.` `",
		},
		cli.StringFlag{
			Name:  "labels",
			Value: "",
			Usage: "[Optional] Package lables in the form of \"lable11\",\"lable2\"...` `",
		},
		cli.StringFlag{
			Name:  "cust-licenses",
			Value: "",
			Usage: "[Optional] Package custom licenses in the form of \"my-license-1\",\"my-license-2\"...` `",
		},
		cli.StringFlag{
			Name:  "website-url",
			Value: "",
			Usage: "[Optional] Package web site URL.` `",
		},
		cli.StringFlag{
			Name:  "issuetracker-url",
			Value: "",
			Usage: "[Optional] Package Issues Tracker URL.` `",
		},
		cli.StringFlag{
			Name:  "github-repo",
			Value: "",
			Usage: "[Optional] Package Github repository.` `",
		},
		cli.StringFlag{
			Name:  "github-rel-notes",
			Value: "",
			Usage: "[Optional] Github release notes file.` `",
		},
	}
}

func getVersionFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "github-tag-rel-notes",
			Usage: "[Default: false] Set to true if you wish to use a Github tag release notes.",
		},
		cli.StringFlag{
			Name:  "desc",
			Value: "",
			Usage: "[Optional] Version description.` `",
		},
		cli.StringFlag{
			Name:  "released",
			Value: "",
			Usage: "[Optional] Release date in ISO8601 format (yyyy-MM-dd'T'HH:mm:ss.SSSZ)` `",
		},
		cli.StringFlag{
			Name:  "github-rel-notes",
			Value: "",
			Usage: "[Optional] Github release notes file.` `",
		},
		cli.StringFlag{
			Name:  "vcs-tag",
			Value: "",
			Usage: "[Optional] VCS tag.` `",
		},
	}
}

func getCreateAndUpdatePackageFlags() []cli.Flag {
	return append(getFlags(), getPackageFlags()...)
}

func getCreateAndUpdateVersionFlags() []cli.Flag {
	return append(getFlags(), getVersionFlags()...)
}

func getDeletePackageAndVersionFlags() []cli.Flag {
	return append(getFlags(), cli.BoolFlag{
		Name:  "quiet",
		Usage: "[Default: false] Set to true to skip the delete confirmation message.",
	})
}

func getDownloadFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "flat",
			Usage: "[Default: false] Set to true if you do not wish to have the Bintray path structure created locally for your downloaded files.",
		},
		cli.StringFlag{
			Name:  "min-split",
			Value: "",
			Usage: "[Default: 5120] Minimum file size in KB to split into ranges when downloading. Set to -1 for no splits.` `",
		},
		cli.StringFlag{
			Name:  "split-count",
			Value: "",
			Usage: "[Default: 3] Number of parts to split a file when downloading. Set to 0 for no splits.` `",
		},
		cli.BoolFlag{
			Name:  "unpublished",
			Usage: "[Default: false] Download both published and unpublished files.",
		},
	}
}

func getDownloadFileFlags() []cli.Flag {
	return append(getFlags(), getDownloadFlags()...)
}

func getDownloadVersionFlags() []cli.Flag {
	flags := append(getFlags(), cli.StringFlag{
		Name:  "threads",
		Value: "",
		Usage: "[Default: 3] Number of artifacts to download in parallel.` `",
	})
	return append(flags, getDownloadFlags()...)
}

func getUploadFlags() []cli.Flag {
	return append(getFlags(), []cli.Flag{
		cli.BoolTFlag{
			Name:  "recursive",
			Usage: "[Default: true] Set to false if you do not wish to collect files in sub-folders to be uploaded to Bintray.",
		},
		cli.BoolTFlag{
			Name:  "flat",
			Usage: "[Default: true] If set to false, files are uploaded according to their file system hierarchy.",
		},
		cli.BoolFlag{
			Name:  "regexp",
			Usage: "[Default: false] Set to true to use a regular expression instead of wildcards expression to collect files to upload.",
		},
		cli.BoolFlag{
			Name:  "publish",
			Usage: "[Default: false] Set to true to publish the uploaded files.",
		},
		cli.BoolFlag{
			Name:  "override",
			Usage: "[Default: false] Set to true to enable overriding existing published files.",
		},
		cli.BoolFlag{
			Name:  "explode",
			Usage: "[Default: false] Set to true to explode archived files after upload.",
		},
		cli.StringFlag{
			Name:  "threads",
			Value: "",
			Usage: "[Default: 3] Number of artifacts to upload in parallel.` `",
		},
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "[Default: false] Set to true to disable communication with Bintray.",
		},
		cli.StringFlag{
			Name:  "deb",
			Value: "",
			Usage: "[Optional] Used for Debian packages in the form of distribution/component/architecture.` `",
		},
	}...)
}

func getEntitlementsFlags() []cli.Flag {
	return append(getFlags(), []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[Optional] Entitlement ID. Used for entitlements update.` `",
		},
		cli.StringFlag{
			Name:  "access",
			Usage: "[Optional] Entitlement access. Used for entitlements creation and update.` `",
		},
		cli.StringFlag{
			Name:  "keys",
			Usage: "[Optional] Used for entitlements creation and update. List of Access Keys in the form of \"key1\",\"key2\"...` `",
		},
		cli.StringFlag{
			Name:  "path",
			Usage: "[Optional] Entitlement path. Used for entitlements creating and update.` `",
		},
	}...)
}

func getAccessKeysFlags() []cli.Flag {
	return append(getFlags(), []cli.Flag{
		cli.StringFlag{
			Name:  "org",
			Usage: "[Optional] Bintray organization` `",
		},
		cli.StringFlag{
			Name:  "password",
			Usage: "[Optional] Access Key password.` `",
		},
		cli.StringFlag{
			Name:  "expiry",
			Usage: "[Optional] Access Key expiry (required for 'jfrog bt acc-keys show/create/update/delete'` `",
		},
		cli.StringFlag{
			Name:  "ex-check-url",
			Usage: "[Optional] You can optionally provide an existence check directive, in the form of a callback URL, to verify whether the source identity of the Access Key still exists.` `",
		},
		cli.StringFlag{
			Name:  "ex-check-cache",
			Usage: "[Optional] You can optionally provide the period in seconds for the callback URL results cache.` `",
		},
		cli.StringFlag{
			Name:  "white-cidrs",
			Usage: "[Optional] Specifying white CIDRs in the form of 127.0.0.1/22,193.5.0.1/92 will allow access only for those IPs that exist in that address range.` `",
		},
		cli.StringFlag{
			Name:  "black-cidrs",
			Usage: "[Optional] Specifying black CIDRs in the form of 127.0.0.1/22,193.5.0.1/92 will block access for all IPs that exist in the specified range.` `",
		},
		cli.BoolTFlag{
			Name:  "api-only",
			Usage: "[Default: true] You can set api_only to false to allow access keys access to Bintray UI as well as to the API.` `",
		},
	}...)
}

func getUrlSigningFlags() []cli.Flag {
	return append(getFlags(), []cli.Flag{
		cli.StringFlag{
			Name:  "expiry",
			Usage: "[Optional] An expiry date for the URL, in Unix epoch time in milliseconds, after which the URL will be invalid. By default, expiry date will be 24 hours.` `",
		},
		cli.StringFlag{
			Name:  "valid-for",
			Usage: "[Optional] The number of seconds since generation before the URL expires. Mutually exclusive with the --expiry option.` `",
		},
		cli.StringFlag{
			Name:  "callback-id",
			Usage: "[Optional] An applicative identifier for the request. This identifier appears in download logs and is used in email and download webhook notifications.` `",
		},
		cli.StringFlag{
			Name:  "callback-email",
			Usage: "[Optional] An email address to send mail to when a user has used the download URL. This requires a callback_id. The callback-id will be included in the mail message.` `",
		},
		cli.StringFlag{
			Name:  "callback-url",
			Usage: "[Optional] A webhook URL to call when a user has used the download URL.` `",
		},
		cli.StringFlag{
			Name:  "callback-method",
			Usage: "[Optional] HTTP method to use for making the callback. Will use POST by default. Supported methods are: GET, POST, PUT and HEAD.` `",
		},
	}...)
}

func getGpgSigningFlags() []cli.Flag {
	return append(getFlags(), cli.StringFlag{
		Name:  "passphrase",
		Usage: "[Optional] GPG passphrase.` `",
	})
}

func configure(c *cli.Context) error {
	if c.NArg() > 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	} else if c.NArg() == 1 {
		if c.Args().Get(0) == "show" {
			commands.ShowConfig()
		} else if c.Args().Get(0) == "clear" {
			commands.ClearConfig()
		} else {
			return errors.New("Unknown argument '" + c.Args().Get(0) + "'. Available arguments are 'show' and 'clear'.")
		}
	} else {
		interactive := c.BoolT("interactive")
		if !interactive {
			if c.String("user") == "" || c.String("key") == "" {
				return errors.New("The --user and --key options are mandatory when the --interactive option is set to false")
			}
		}
		bintrayDetails, err := createBintrayDetails(c, false)
		if err != nil {
			return err
		}

		cliBtDetails := &config.BintrayDetails{
			User:              bintrayDetails.GetUser(),
			Key:               bintrayDetails.GetKey(),
			ApiUrl:            bintrayDetails.GetApiUrl(),
			DownloadServerUrl: bintrayDetails.GetDownloadServerUrl(),
			DefPackageLicense: bintrayDetails.GetDefPackageLicense(),
		}
		commands.Config(cliBtDetails, nil, interactive)
	}

	return nil
}

func showPackage(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	packagePath, err := packages.CreatePath(c.Args().Get(0))
	if err != nil {
		return err
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}

	return commands.ShowPackage(btConfig, packagePath)
}

func showVersion(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	versionPath, err := versions.CreatePath(c.Args().Get(0))
	if err != nil {
		return err
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}

	return commands.ShowVersion(btConfig, versionPath)
}

func createPackage(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	packageParams, err := createPackageParams(c)
	if err != nil {
		return err
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}

	return commands.CreatePackage(btConfig, packageParams)
}

func createVersion(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	versionParams, err := createVersionParams(c)
	if err != nil {
		return err
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}

	return commands.CreateVersion(btConfig, versionParams)
}

func updateVersion(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	versionParams, err := createVersionParams(c)
	if err != nil {
		return err
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}

	return commands.UpdateVersion(btConfig, versionParams)
}

func updatePackage(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	packageParams, err := createPackageParams(c)
	if err != nil {
		return err
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}

	return commands.UpdatePackage(btConfig, packageParams)
}

func deletePackage(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	packagePath, err := packages.CreatePath(c.Args().Get(0))
	if err != nil {
		return err
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}
	if !c.Bool("quiet") {
		confirmed := cliutils.InteractiveConfirm("Delete package " + packagePath.Package + "?")
		if !confirmed {
			return nil
		}
	}
	return commands.DeletePackage(btConfig, packagePath)
}

func deleteVersion(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	versionPath, err := versions.CreatePath(c.Args().Get(0))
	if err != nil {
		return err
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}
	if !c.Bool("quiet") {
		confirmed := cliutils.InteractiveConfirm("Delete version " + versionPath.Version +
			" of package " + versionPath.Package + "?")
		if !confirmed {
			return nil
		}
	}
	return commands.DeleteVersion(btConfig, versionPath)
}

func publishVersion(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	versionPath, err := versions.CreatePath(c.Args().Get(0))
	if err != nil {
		return err
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}
	return commands.PublishVersion(btConfig, versionPath)
}

func downloadVersion(c *cli.Context) error {
	if c.NArg() < 1 || c.NArg() > 2 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	var err error
	params := services.NewDownloadVersionParams()
	params.IncludeUnpublished = c.Bool("unpublished")
	params.Path, err = services.CreateVersionDetailsForDownloadVersion(c.Args().Get(0))
	if err != nil {
		return err
	}

	params.TargetPath = c.Args().Get(1)
	if strings.HasPrefix(params.TargetPath, "/") {
		params.TargetPath = params.TargetPath[1:]
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}
	downloaded, failed, err := commands.DownloadVersion(btConfig, params)
	err = cliutils.PrintSummaryReport(downloaded, failed, err)
	if err != nil {
		return err
	}
	if failed > 0 {
		return errors.New("")
	}
	return nil
}

func upload(c *cli.Context) error {
	if c.NArg() < 2 || c.NArg() > 3 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	params := services.NewUploadParams()
	params.Pattern = c.Args().Get(0)
	if cliutils.IsWindows() {
		params.Pattern = ioutils.UnixToWinPathSeparator(params.Pattern)
	}

	var err error
	params.Path, err = versions.CreatePath(c.Args().Get(1))
	if err != nil {
		return err
	}

	params.TargetPath = c.Args().Get(2)
	if strings.HasPrefix(params.TargetPath, "/") {
		params.TargetPath = params.TargetPath[1:]
	}

	params.Deb = c.String("deb")
	if params.Deb != "" && len(strings.Split(params.Deb, "/")) != 3 {
		return errors.New("The --deb option should be in the form of distribution/component/architecture")
	}

	params.Recursive = c.BoolT("recursive")
	params.Flat = c.BoolT("flat")
	params.Publish = c.Bool("publish")
	params.Override = c.Bool("override")
	params.Explode = c.Bool("explode")
	params.UseRegExp = c.Bool("regexp")

	uploadConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}
	uploaded, failed, err := commands.Upload(uploadConfig, params)
	err = cliutils.PrintSummaryReport(uploaded, failed, err)
	if err != nil {
		return err
	}
	if failed > 0 {
		return errors.New("")
	}

	return nil
}

func downloadFile(c *cli.Context) error {
	if c.NArg() < 1 || c.NArg() > 2 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	var err error
	params := services.NewDownloadFileParams()
	params.Flat = c.Bool("flat")
	params.IncludeUnpublished = c.Bool("unpublished")
	params.MinSplitSize, err = getMinSplitFlag(c)
	if err != nil {
		return nil
	}
	params.SplitCount, err = getSplitCountFlag(c)
	if err != nil {
		return err
	}
	params.PathDetails, err = utils.CreatePathDetails(c.Args().Get(0))
	if err != nil {
		return err
	}

	params.TargetPath = c.Args().Get(1)
	if strings.HasPrefix(params.TargetPath, "/") {
		params.TargetPath = params.TargetPath[1:]
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}
	downloaded, failed, err := commands.DownloadFile(btConfig, params)
	return cliutils.PrintSummaryReport(downloaded, failed, err)
}

func signUrl(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}
	signUrlParams, err := createUrlSigningFlags(c)
	if err != nil {
		return err
	}

	return commands.SignVersion(btConfig, signUrlParams)
}

func gpgSignFile(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	pathDetails, err := utils.CreatePathDetails(c.Args().Get(0))
	if err != nil {
		return err
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}

	return commands.GpgSignFile(btConfig, pathDetails, c.String("passphrase"))
}

func logs(c *cli.Context) error {
	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}

	if c.NArg() == 1 {
		versionPath, err := versions.CreatePath(c.Args().Get(0))
		if err != nil {
			return err
		}
		err = commands.LogsList(btConfig, versionPath)
	} else if c.NArg() == 3 && c.Args().Get(0) == "download" {
		versionPath, err := versions.CreatePath(c.Args().Get(1))
		if err != nil {
			return err
		}
		err = commands.DownloadLog(btConfig, versionPath, c.Args().Get(2))
	} else {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	return err
}

func stream(c *cli.Context) error {
	bintrayDetails, err := createBintrayDetails(c, true)
	if err != nil {
		return err
	}
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}

	streamDetails := &commands.StreamDetails{
		BintrayDetails: bintrayDetails,
		Subject:        c.Args().Get(0),
		Include:        c.String("include"),
	}
	commands.Stream(streamDetails, os.Stdout)

	return nil
}

func gpgSignVersion(c *cli.Context) error {
	if c.NArg() != 1 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	versionDetails, err := versions.CreatePath(c.Args().Get(0))
	if err != nil {
		return err
	}

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}

	return commands.GpgSignVersion(btConfig, versionDetails, c.String("passphrase"))
}

func accessKeys(c *cli.Context) error {
	var err error
	org := c.String("org")

	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}
	if c.NArg() == 0 {
		err = commands.ShowAllAccessKeys(btConfig, org)
		return err
	}
	if c.NArg() != 2 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	keyId := c.Args().Get(1)

	switch c.Args().Get(0) {
	case "show":
		err = commands.ShowAccessKey(btConfig, org, keyId)
	case "delete":
		err = commands.DeleteAccessKey(btConfig, org, keyId)
	case "create":
		params, err := createAccessKeysParams(c, org, keyId)
		if err != nil {
			return err
		}
		err = commands.CreateAccessKey(btConfig, params)
	case "update":
		params, err := createAccessKeysParams(c, org, keyId)
		if err != nil {
			return err
		}
		err = commands.UpdateAccessKey(btConfig, params)
	default:
		return errors.New("Expecting show, create, update or delete before the key argument. Got " + c.Args().Get(0))
	}
	return err
}

func handleEntitlements(c *cli.Context) error {
	if c.NArg() == 0 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	btConfig, err := newBintrayConfig(c)
	if err != nil {
		return err
	}

	if c.NArg() == 1 {
		details, err := entitlements.CreateVersionDetails(c.Args().Get(0))
		if err != nil {
			return err
		}
		return commands.ShowAllEntitlements(btConfig, details)
	}

	if c.NArg() != 2 {
		return cliutils.PrintHelpAndReturnError("Wrong number of arguments.", c)
	}
	versionPath, err := entitlements.CreateVersionDetails(c.Args().Get(1))
	if err != nil {
		return nil
	}

	switch c.Args().Get(0) {
	case "show":
		id := c.String("id")
		if id == "" {
			return errors.New("Please add the --id option")
		}
		err = commands.ShowEntitlement(btConfig, id, versionPath)
	case "create":
		params, err := createEntitlementFlagsForCreate(c, versionPath)
		if err != nil {
			return err
		}
		err = commands.CreateEntitlement(btConfig, params)
	case "update":
		params, err := createEntitlementFlagsForUpdate(c, versionPath)
		if err != nil {
			return err
		}
		err = commands.UpdateEntitlement(btConfig, params)
	case "delete":
		id := c.String("id")
		if id == "" {
			return errors.New("Please add the --id option")
		}
		err = commands.DeleteEntitlement(btConfig, id, versionPath)
	default:
		return errors.New("Expecting show, create, update or delete before " + c.Args().Get(1) + ". Got " + c.Args().Get(0))
	}
	return err
}

func createPackageParams(c *cli.Context) (*packages.Params, error) {
	licenses := c.String("licenses")
	if licenses == "" {
		confDetails, err := commands.GetConfig()
		if err != nil {
			return nil, err
		}
		licenses = confDetails.DefPackageLicense
	}

	packagePath, err := packages.CreatePath(c.Args().Get(0))
	if err != nil {
		return nil, err
	}

	params := packages.NewPackageParams()
	params.Path = packagePath
	params.Desc = c.String("desc")
	params.Labels = c.String("labels")
	params.Licenses = licenses
	params.CustomLicenses = c.String("cust-licenses")
	params.VcsUrl = c.String("vcs-url")
	params.WebsiteUrl = c.String("website-url")
	params.IssueTrackerUrl = c.String("issuetracker-url")
	params.GithubRepo = c.String("github-repo")
	params.GithubReleaseNotesFile = c.String("github-rel-notes")
	params.PublicDownloadNumbers = c.Bool("pub-dn")
	params.PublicStats = c.BoolT("pub-stats")

	return params, nil
}

func newBintrayConfig(c *cli.Context) (bintray.Config, error) {
	btDetails, err := createBintrayDetails(c, true)
	if err != nil {
		return nil, err
	}
	threads, err := getThreadsOptionValue(c)
	if err != nil {
		return nil, err
	}
	btConfig := bintray.NewConfigBuilder().
		SetBintrayDetails(btDetails).
		SetDryRun(c.Bool("dry-run")).
		SetThreads(threads).
		Build()
	return btConfig, nil
}

func createVersionParams(c *cli.Context) (*versions.Params, error) {
	versionDetails, err := versions.CreatePath(c.Args().Get(0))
	if err != nil {
		return nil, err
	}

	params := versions.NewVersionParams()
	params.Path = versionDetails
	params.Desc = c.String("desc")
	params.VcsTag = c.String("vcs-tag")
	params.Released = c.String("released")
	params.GithubReleaseNotesFile = c.String("github-rel-notes")
	params.GithubUseTagReleaseNotes = c.Bool("github-tag-rel-notes")

	return params, nil
}

func createUrlSigningFlags(c *cli.Context) (*url.Params, error) {
	if c.String("valid-for") != "" {
		_, err := strconv.ParseInt(c.String("valid-for"), 10, 64)
		if err != nil {
			return nil, errors.New("The '--valid-for' option should have a numeric value.")
		}
	}
	urlSigningDetails, err := utils.CreatePathDetails(c.Args().Get(0))
	if err != nil {
		return nil, err
	}

	var expiry int64
	if c.String("expiry") != "" {
		var err error
		expiry, err = strconv.ParseInt(c.String("expiry"), 10, 64)
		if err != nil {
			return nil, errors.New("The --expiry option should have a numeric value.")
		}
	}

	params := url.NewURLParams()
	params.PathDetails = urlSigningDetails
	params.Expiry = expiry
	params.ValidFor = c.Int("valid-for")
	params.CallbackId = c.String("callback-id")
	params.CallbackEmail = c.String("callback-email")
	params.CallbackUrl = c.String("callback-url")
	params.CallbackMethod = c.String("callback-method")

	return params, nil
}

func getThreadsOptionValue(c *cli.Context) (threads int, err error) {
	if c.String("threads") == "" {
		threads = 3
		err = nil
	} else {
		var err error
		threads, err = strconv.Atoi(c.String("threads"))
		if err != nil || threads < 1 {
			return 0, errors.New("The '--threads' option should have a numeric positive value.")
		}
	}
	return
}

func createEntitlementFlagsForCreate(c *cli.Context, path *versions.Path) (*entitlements.Params, error) {
	if c.String("access") == "" {
		return nil, errors.New("Please add the --access option")
	}

	params := entitlements.NewEntitlementsParams()
	params.VersionPath = path
	params.Path = c.String("path")
	params.Access = c.String("access")
	params.Keys = c.String("keys")

	return params, nil
}

func createEntitlementFlagsForUpdate(c *cli.Context, path *versions.Path) (*entitlements.Params, error) {
	if c.String("id") == "" {
		return nil, errors.New("Please add the --id option")
	}
	if c.String("access") == "" {
		return nil, errors.New("Please add the --access option")
	}

	params := entitlements.NewEntitlementsParams()
	params.VersionPath = path
	params.Id = c.String("id")
	params.Path = c.String("path")
	params.Access = c.String("access")
	params.Keys = c.String("keys")

	return params, nil
}

func createAccessKeysParams(c *cli.Context, org, keyId string) (*accesskeys.Params, error) {
	var cachePeriod int
	if c.String("ex-check-cache") != "" {
		var err error
		cachePeriod, err = strconv.Atoi(c.String("ex-check-cache"))
		if err != nil {
			return nil, errors.New("The --ex-check-cache option should have a numeric value.")
		}
	}

	var expiry int64
	if c.String("expiry") != "" {
		var err error
		expiry, err = strconv.ParseInt(c.String("expiry"), 10, 64)
		if err != nil {
			return nil, errors.New("The --expiry option should have a numeric value.")
		}
	}

	params := accesskeys.NewAccessKeysParams()
	params.Id = keyId
	params.Password = c.String("password")
	params.Org = org
	params.Expiry = expiry
	params.ExistenceCheckUrl = c.String("ex-check-url")
	params.ExistenceCheckCache = cachePeriod
	params.WhiteCidrs = c.String("white-cidrs")
	params.BlackCidrs = c.String("black-cidrs")
	params.ApiOnly = c.BoolT("api-only")

	return params, nil
}

func offerConfig(c *cli.Context) (*config.BintrayDetails, error) {
	exists, err := config.IsBintrayConfExists()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, nil
	}
	val, err := clientutils.GetBoolEnvValue(cliutils.OfferConfig, true)
	if err != nil {
		return nil, err
	}
	if !val {
		config.SaveBintrayConf(new(config.BintrayDetails))
		return nil, nil
	}
	msg := "Some CLI commands require the following common options:\n" +
		"- User\n" +
		"- API Key\n" +
		"- Default Package Licenses\n" +
		"Configuring JFrog CLI with these parameters now will save you having to include them as command options.\n" +
		"You can also configure these parameters later using the 'config' command.\n" +
		"Configure now?"
	confirmed := cliutils.InteractiveConfirm(msg)
	if !confirmed {
		config.SaveBintrayConf(new(config.BintrayDetails))
		return nil, nil
	}
	bintrayDetails, err := createBintrayDetails(c, false)
	if err != nil {
		return nil, err
	}
	cliBtDetails := &config.BintrayDetails{
		ApiUrl:            bintrayDetails.GetApiUrl(),
		DownloadServerUrl: bintrayDetails.GetDownloadServerUrl(),
		User:              bintrayDetails.GetUser(),
		Key:               bintrayDetails.GetKey(),
		DefPackageLicense: bintrayDetails.GetDefPackageLicense()}

	details, err := commands.Config(nil, cliBtDetails, true)
	if err != nil {
		return nil, err
	}
	details.ApiUrl = bintrayDetails.GetApiUrl()
	details.DownloadServerUrl = bintrayDetails.GetDownloadServerUrl()
	return details, nil
}

func createBintrayDetails(c *cli.Context, includeConfig bool) (auth.BintrayDetails, error) {
	if includeConfig {
		bintrayDetails, err := offerConfig(c)
		if err != nil {
			return nil, err
		}
		if bintrayDetails != nil {
			btDetails := auth.NewBintrayDetails()
			btDetails.SetApiUrl(bintrayDetails.ApiUrl)
			btDetails.SetDownloadServerUrl(bintrayDetails.DownloadServerUrl)
			btDetails.SetUser(bintrayDetails.User)
			btDetails.SetKey(bintrayDetails.Key)
			btDetails.SetDefPackageLicense(bintrayDetails.DefPackageLicense)
			return btDetails, nil
		}
	}
	user := c.String("user")
	key := c.String("key")
	defaultPackageLicenses := c.String("licenses")
	if includeConfig && (user == "" || key == "" || defaultPackageLicenses == "") {
		confDetails, err := commands.GetConfig()
		if err != nil {
			return nil, err
		}
		if user == "" {
			user = confDetails.User
		}
		if key == "" {
			key = confDetails.Key
		}
		if key == "" {
			return nil, errors.New("Please set your Bintray API key using the config command or send it as the --key option.")
		}
		if defaultPackageLicenses == "" {
			defaultPackageLicenses = confDetails.DefPackageLicense
		}
	}
	btDetails := auth.NewBintrayDetails()
	apiUrl := os.Getenv("JFROG_CLI_BINTRAY_API_URL")
	if apiUrl != "" {
		apiUrl = clientutils.AddTrailingSlashIfNeeded(apiUrl)
		btDetails.SetApiUrl(apiUrl)
	}
	downloadServerUrl := os.Getenv("JFROG_CLI_BINTRAY_DOWNLOAD_URL")
	if downloadServerUrl != "" {
		downloadServerUrl = clientutils.AddTrailingSlashIfNeeded(downloadServerUrl)
		btDetails.SetDownloadServerUrl(downloadServerUrl)
	}

	btDetails.SetUser(user)
	btDetails.SetKey(key)
	btDetails.SetDefPackageLicense(defaultPackageLicenses)
	return btDetails, nil
}

func getMinSplitFlag(c *cli.Context) (int64, error) {
	if c.String("min-split") == "" {
		return 5120, nil
	}
	minSplit, err := strconv.ParseInt(c.String("min-split"), 10, 64)
	if err != nil {
		return 0, cliutils.PrintHelpAndReturnError("The '--min-split' option should have a numeric value.", c)
	}
	return minSplit, nil
}

func getSplitCountFlag(c *cli.Context) (int, error) {
	if c.String("split-count") == "" {
		return 3, nil
	}
	splitCount, err := strconv.Atoi(c.String("split-count"))
	if err != nil {
		return 0, cliutils.PrintHelpAndReturnError("The '--split-count' option should have a numeric value.", c)
	}
	if splitCount > 15 {
		return 0, errors.New("The '--split-count' option value is limitted to a maximum of 15.")
	}
	if splitCount < 0 {
		return 0, errors.New("The '--split-count' option cannot have a negative value.")
	}
	return splitCount, nil
}
