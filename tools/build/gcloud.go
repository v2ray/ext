package build

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func getGCloudMetadata(name string, scope string) (string, error) {
	request, err := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/"+scope+"/attributes/"+name, nil)
	if err != nil {
		return "", newError("failed to create http request").Base(err)
	}
	request.Header.Set("Metadata-Flavor", "Google")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", newError("failed to get gcloud attribute: ", name, " in ", scope).Base(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", newError("failed to read gcloud attribute: ", name, " in ", scope).Base(err)
	}
	return strings.TrimSpace(string(body)), nil
}

func GetGCloudInstanceMetadata(name string) (string, error) {
	return getGCloudMetadata(name, "instance")
}

func GetGCloudProjectMetadata(name string) (string, error) {
	return getGCloudMetadata(name, "project")
}
