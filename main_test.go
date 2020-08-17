package main

import (
	"testing"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
}

func TestExecuteMutator(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	event.Metrics = corev2.FixtureMetrics()
	ev, err := executeMutator(event)
	assert.NoError(err)
	assert.Equal(0, len(ev.Entity.Redact))
	assert.Equal(0, len(ev.Entity.System.Network.Interfaces))
	assert.Equal(0, len(ev.Entity.Subscriptions))
	assert.Equal(0, len(ev.Check.Handlers))
	assert.Equal(0, len(ev.Check.History))
	assert.Equal(0, len(ev.Check.RuntimeAssets))
	assert.Equal(0, len(ev.Check.Subscriptions))
}
