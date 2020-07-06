package inttestutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/buger/jsonparser"
	"github.com/jfrog/jfrog-cli/artifactory/spec"
	"github.com/jfrog/jfrog-cli/utils/cliutils"
	"github.com/jfrog/jfrog-cli/utils/config"
	"github.com/jfrog/jfrog-cli/utils/tests"
	"github.com/jfrog/jfrog-client-go/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/assert"
)

const (
	gpgKeyId                        = "234503"
	distributionGpgKeyCreatePattern = `{"public_key":"%s","private_key":"%s"}`
	artifactoryGpgKeyCreatePattern  = `{"alias":"cli tests distribution key","public_key":"%s"}`
)

type distributableDistributionStatus string
type receivedDistributionStatus string

const (
	Open                 distributableDistributionStatus = "OPEN"
	ReadyForDistribution distributableDistributionStatus = "READY_FOR_DISTRIBUTION"
	Signed               distributableDistributionStatus = "SIGNED"
	NotDistributed       receivedDistributionStatus      = "Not distributed"
	InProgress           receivedDistributionStatus      = "In progress"
	Completed            receivedDistributionStatus      = "Completed"
	Failed               receivedDistributionStatus      = "Failed"
)

// GET api/v1/release_bundle/:name/:version
// Retreive the status of a release bundle before distribution.
type distributableResponse struct {
	Name         string                          `json:"name,omitempty"`
	Version      string                          `json:"version,omitempty"`
	State        distributableDistributionStatus `json:"state,omitempty"`
	Description  string                          `json:"description,omitempty"`
	ReleaseNotes releaseNotesResponse            `json:"release_notes,omitempty"`
}

type releaseNotesResponse struct {
	Content string `json:"content,omitempty"`
	Syntax  string `json:"syntax,omitempty"`
}

// Get api/v1/release_bundle/:name/:version/distribution
// Retreive the status of a release bundle after distribution.
type receivedResponse struct {
	Id     string                     `json:"id,omitempty"`
	Status receivedDistributionStatus `json:"status,omitempty"`
}

type ReceivedResponses struct {
	receivedResponses []receivedResponse
}

// Send GPG keys to Distribution and Artifactory to allow signing of release bundles
func SendGpgKeys(artHttpDetails httputils.HttpClientDetails, distHttpDetails httputils.HttpClientDetails) {
	// Read gpg public and private keys
	keysDir := filepath.Join(tests.GetTestResourcesPath(), "distribution")
	publicKey, err := ioutil.ReadFile(filepath.Join(keysDir, "public.key"))
	cliutils.ExitOnErr(err)
	privateKey, err := ioutil.ReadFile(filepath.Join(keysDir, "private.key"))
	cliutils.ExitOnErr(err)

	// Create http client
	client, err := httpclient.ClientBuilder().Build()
	cliutils.ExitOnErr(err)

	// Send public and private keys to Distribution
	content := fmt.Sprintf(distributionGpgKeyCreatePattern, publicKey, privateKey)
	resp, body, err := client.SendPut(*tests.RtDistributionUrl+"api/v1/keys/pgp", []byte(content), distHttpDetails)
	cliutils.ExitOnErr(err)
	if resp.StatusCode != http.StatusOK {
		log.Error(resp.Status)
		log.Error(string(body))
		os.Exit(1)
	}

	// Send public key to Artifactory
	content = fmt.Sprintf(artifactoryGpgKeyCreatePattern, publicKey)
	resp, body, err = client.SendPost(*tests.RtUrl+"api/security/keys/trusted", []byte(content), artHttpDetails)
	cliutils.ExitOnErr(err)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		log.Error(resp.Status)
		log.Error(string(body))
		os.Exit(1)
	}
}

// Get a local release bundle
func GetLocalBundle(t *testing.T, bundleName, bundleVersion string, distHttpDetails httputils.HttpClientDetails) *distributableResponse {
	resp, body := getLocalBundle(t, bundleName, bundleVersion, distHttpDetails)
	if resp.StatusCode != http.StatusOK {
		t.Error(resp.Status)
		t.Error(string(body))
		return nil
	}
	response := &distributableResponse{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		t.Error(err)
		return nil
	}
	return response
}

// Return true if the release bundle exists locally on distribution
func VerifyLocalBundleExistence(t *testing.T, bundleName, bundleVersion string, expectExist bool, distHttpDetails httputils.HttpClientDetails) {
	for i := 0; i < 120; i++ {
		resp, body := getLocalBundle(t, bundleName, bundleVersion, distHttpDetails)
		switch resp.StatusCode {
		case http.StatusOK:
			if expectExist {
				return
			}
		case http.StatusNotFound:
			if !expectExist {
				return
			}
		default:
			t.Error(resp.Status)
			t.Error(string(body))
			return
		}
		t.Log("Waiting for " + bundleName + "/" + bundleVersion + "...")
		time.Sleep(time.Second)
	}
	t.Errorf("Release bundle %s/%s exist: %v unlike expected", bundleName, bundleVersion, expectExist)
}

// Wait for distribution of a release bundle
func WaitForDistribution(t *testing.T, bundleName, bundleVersion string, distHttpDetails httputils.HttpClientDetails) {
	client, err := httpclient.ClientBuilder().Build()
	assert.NoError(t, err)

	for i := 0; i < 120; i++ {
		resp, body, _, err := client.SendGet(*tests.RtDistributionUrl+"api/v1/release_bundle/"+bundleName+"/"+bundleVersion+"/distribution", true, distHttpDetails)
		assert.NoError(t, err)
		if resp.StatusCode != http.StatusOK {
			t.Error(resp.Status)
			t.Error(string(body))
			return
		}
		response := &ReceivedResponses{}
		err = json.Unmarshal(body, &response.receivedResponses)
		if err != nil {
			t.Error(err)
			return
		}
		if len(response.receivedResponses) == 0 {
			t.Error("Release bundle \"" + bundleName + "/" + bundleVersion + "\" not found")
			return
		}

		switch response.receivedResponses[0].Status {
		case Completed:
			return
		case Failed:
			t.Error("Distribution failed for " + bundleName + "/" + bundleVersion)
			return
		case InProgress, NotDistributed:
			// Wait
		}
		t.Log("Waiting for " + bundleName + "/" + bundleVersion + "...")
		time.Sleep(time.Second)
	}
	t.Error("Timeout for release bundle distribution " + bundleName + "/" + bundleVersion)
}

// Wait for deletion of a release bundle
func WaitForDeletion(t *testing.T, bundleName, bundleVersion string, distHttpDetails httputils.HttpClientDetails) {
	client, err := httpclient.ClientBuilder().Build()
	assert.NoError(t, err)

	for i := 0; i < 120; i++ {
		resp, body, _, err := client.SendGet(*tests.RtDistributionUrl+"api/v1/release_bundle/"+bundleName+"/"+bundleVersion+"/distribution", true, distHttpDetails)
		assert.NoError(t, err)
		if resp.StatusCode == http.StatusNotFound {
			return
		}
		if resp.StatusCode != http.StatusOK {
			t.Error(resp.Status)
			t.Error(string(body))
			return
		}
		t.Log("Waiting for distribution deletion " + bundleName + "/" + bundleVersion + "...")
		time.Sleep(time.Second)
	}
	t.Error("Timeout for release bundle deletion " + bundleName + "/" + bundleVersion)
}

func getLocalBundle(t *testing.T, bundleName, bundleVersion string, distHttpDetails httputils.HttpClientDetails) (*http.Response, []byte) {
	client, err := httpclient.ClientBuilder().Build()
	assert.NoError(t, err)

	resp, body, _, err := client.SendGet(*tests.RtDistributionUrl+"api/v1/release_bundle/"+bundleName+"/"+bundleVersion, true, distHttpDetails)
	assert.NoError(t, err)
	return resp, body
}

func CleanUpOldBundles(distHttpDetails httputils.HttpClientDetails, bundleVersion string, distributionCli *tests.JfrogCli) {
	getActualItems := func() ([]string, error) { return ListAllBundlesNames(distHttpDetails) }
	deleteItem := func(bundleName string) {
		distributionCli.Exec("rbdel", bundleName, bundleVersion, "--site=*", "--delete-from-dist", "--quiet")
		log.Info("Bundle", bundleName, "deleted.")
	}
	tests.CleanUpOldItems([]string{tests.BundleName}, getActualItems, deleteItem)
}

func ListAllBundlesNames(distHttpDetails httputils.HttpClientDetails) ([]string, error) {
	var bundlesNames []string

	// Build http client
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return nil, err
	}

	// Send get request
	resp, body, _, err := client.SendGet(*tests.RtDistributionUrl+"api/v1/release_bundle", true, distHttpDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body))
	}

	// Extract release bundle names from the json response
	var keyError error
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil || keyError != nil {
			return
		}
		bundleName, err := jsonparser.GetString(value, "name")
		if err != nil {
			keyError = err
			return
		}
		bundlesNames = append(bundlesNames, bundleName)
	})
	if keyError != nil {
		return nil, err
	}

	return bundlesNames, err
}

// Clean up 'cli-tests-dist1-<timestamp>' and 'cli-tests-dist2-<timestamp>' after running a distribution test
func CleanDistributionRepositories(t *testing.T, distributionDetails *config.ArtifactoryDetails) {
	deleteSpec := spec.NewBuilder().Pattern(tests.DistRepo1).BuildSpec()
	tests.DeleteFiles(deleteSpec, distributionDetails)
	deleteSpec = spec.NewBuilder().Pattern(tests.DistRepo1).BuildSpec()
	tests.DeleteFiles(deleteSpec, distributionDetails)
}
