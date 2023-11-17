package main

import (
	"encoding/json"
	"fmt"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
	"io/ioutil"
	"net/http"
	"os"
)

// Config represents the mutator plugin config.
type Config struct {
	sensu.PluginConfig
}

var (
	mutatorConfig = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-event-streamliner",
			Short:    "Sensu Event Streamliner",
			Keyspace: "sensu.io/plugins/sensu-event-streamliner/config",
		},
	}
)

func main() {
	mutator := sensu.NewGoMutator(&mutatorConfig.PluginConfig, nil, checkArgs, executeMutator)
	mutator.Execute()
}

func checkArgs(_ *types.Event) error {
	return nil
}

func handleError(message string, err error) {
	fmt.Printf("%s: %s\n", message, err)
	os.Exit(1)
}
func executeMutator(event *types.Event) (*types.Event, error) {

	namespace, entityName := "default", event.Entity.Name
	url := fmt.Sprintf("%s/api/core/v2/namespaces/%s/entities/%s", os.Getenv("SENSU_API_URL"), namespace, entityName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		handleError("Error creating request", err)
	}
	req.Header.Add("Authorization", "Key "+os.Getenv("SENSU_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		handleError("Error making request", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError("Error reading response body", err)
	}

	var result map[string]json.RawMessage
	err = json.Unmarshal(body, &result)
	if err != nil {
		handleError("Error parsing JSON response", err)
	}

	var metadataDetails map[string]interface{}
	err = json.Unmarshal(result["metadata"], &metadataDetails)
	if err != nil {
		handleError("Error extracting metadata from response", err)
	}

	labels, ok := metadataDetails["labels"].(map[string]interface{})
	if !ok {
		handleError("Error extracting labels from metadata", fmt.Errorf("labels field not found"))
	}
	jsonLabels, err := json.Marshal(labels)
	if err != nil {
		handleError("Error encoding labels to JSON", err)
	}

	var finalLabels map[string]string
	err = json.Unmarshal([]byte(jsonLabels), &finalLabels)

	event.Entity.Labels = finalLabels

	annotations, ok := metadataDetails["annotations"].(map[string]interface{})
	if !ok {
		handleError("Error extracting labels from metadata", fmt.Errorf("labels field not found"))
	}
	jsonAnnotations, err := json.Marshal(annotations)
	if err != nil {
		handleError("Error encoding labels to JSON", err)
	}
	var finalAnnotations map[string]string
	err = json.Unmarshal([]byte(jsonAnnotations), &finalAnnotations)

	event.Entity.Annotations = finalAnnotations

	return event, nil
}
