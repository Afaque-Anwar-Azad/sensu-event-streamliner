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

	event.Entity.Redact = nil
	event.Entity.System.Network.Interfaces = nil
	event.Entity.Subscriptions = []string{
		"mutator_is_working",
	}
	event.Check.Handlers = nil
	event.Check.History = nil
	event.Check.RuntimeAssets = nil
	event.Check.Subscriptions = nil

	apiURL, namespace, token, entityName := "http://34.244.175.180:8080", "default", "09da2702-d9b3-436f-b1b4-e84949467adc", event.Entity.Name
	url := fmt.Sprintf("%s/api/core/v2/namespaces/%s/entities/%s", apiURL, namespace, entityName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		handleError("Error creating request", err)
	}
	req.Header.Add("Authorization", "Key "+token)

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

	var metadataLabels map[string]interface{}
	err = json.Unmarshal(result["metadata"], &metadataLabels)
	if err != nil {
		handleError("Error extracting metadata from response", err)
	}

	labels, ok := metadataLabels["labels"].(map[string]interface{})
	if !ok {
		handleError("Error extracting labels from metadata", fmt.Errorf("labels field not found"))
	}
	jsonLabels, err := json.Marshal(labels)
	if err != nil {
		handleError("Error encoding labels to JSON", err)
	}

	var finallabels map[string]string
	err = json.Unmarshal([]byte(jsonLabels), &finallabels)
	event.Entity.Labels = finallabels

	return event, nil
}
