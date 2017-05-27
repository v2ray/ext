package build

import "net/http"
import "io/ioutil"

func getGCloudMetadata(name string, scope string) (string, error) {
	response, err := http.Get("http://metadata.google.internal/computeMetadata/v1/" + scope + "/attributes/" + name)
	if err != nil {
		return "", newError("failed to get gcloud attribute: ", name, " in ", scope).Base(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", newError("failed to read gcloud attribute: ", name, " in ", scope).Base(err)
	}
	return string(body), nil
}

func GetGCloudInstanceMetadata(name string) (string, error) {
	return getGCloudMetadata(name, "instance")
}

func GetGCloudProjectMetadata(name string) (string, error) {
	return getGCloudMetadata(name, "project")
}
