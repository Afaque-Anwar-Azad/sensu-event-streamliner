package main

import (
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
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

func executeMutator(event *types.Event) (*types.Event, error) {

	event.Entity.Redact = nil
	event.Entity.System.Network.Interfaces = nil
	event.Entity.Subscriptions = nil
	event.Check.Handlers = nil
	event.Check.History = nil
	event.Check.RuntimeAssets = nil
	event.Check.Subscriptions = nil

	return event, nil
}
