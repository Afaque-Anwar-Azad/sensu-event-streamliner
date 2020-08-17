[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/nixwiz/sensu-event-streamliner)
![Go Test](https://github.com/nixwiz/sensu-event-streamliner/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/nixwiz/sensu-event-streamliner/workflows/goreleaser/badge.svg)

# Sensu Event Streamliner

## Table of Contents
- [Overview](#overview)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Mutator definition](#mutator-definition)
- [Installation from source](#installation-from-source)
- [Contributing](#contributing)

## Overview

The Sensu Event Streamliner is a [Sensu Mutator][2] that removes certain redundant event information
that may not be needed in certain situations (e.g. sending off to an event indexer).

The event fields removed are:

* event.Entity.Redact
* event.Entity.System.Network.Interfaces
* event.Entity.Subscriptions
* event.Check.Handlers
* event.Check.History
* event.Check.RuntimeAssets
* event.Check.Subscriptions

My anecdotal testing has shown that this reduces the event payload by between one and two KiB per event.

## Usage examples

There are no arguments to this mutator, so the usage is quite simple.

```
Sensu Event Streamliner

Usage:
  sensu-event-streamliner [flags]
  sensu-event-streamliner [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -h, --help   help for sensu-event-streamliner

Use "sensu-event-streamliner [command] --help" for more information about a command.
```

## Configuration

### Asset registration

[Sensu Assets][4] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add nixwiz/sensu-event-streamliner
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][3]

### Mutator definition

```yml
---
type: Mutator
api_version: core/v2
metadata:
  name: sensu-event-streamliner
  namespace: default
spec:
  command: sensu-event-streamliner
  runtime_assets:
  - nixwiz/sensu-event-streamliner
```

### Handler definition
```yml
---
type: Handler
api_version: core/v2
metadata:
  name: pushover
  namespace: default
spec:
  command: sensu-go-pushover-handler
  env_vars: null
  filters:
  - is_incident
  - not_silenced
  - fatigue_check
  handlers: null
  mutator: sensu-event-streamliner
  runtime_assets:
  - nixwiz/sensu-go-pushover-handler
  timeout: 10
  type: pipe
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable binary from this source.

From the local path of the sensu-event-streamliner repository:

```
go build
```

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://docs.sensu.io/sensu-go/latest/reference/mutators/
[3]: https://bonsai.sensu.io/assets/nixwiz/sensu-event-streamliner
[4]: https://docs.sensu.io/sensu-go/latest/reference/assets/
