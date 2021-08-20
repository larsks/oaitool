# oaitool

Oaitool is a tool for interacting with the OpenShift Assisted
Installer API.

## Authentication

You'll need to acquire an offline token from
https://console.redhat.com/openshift/token/show.

Put your API token into the file `~/.config/oaitool/config.yml`, like
this:

```
offline-token: "...token goes here..."
```

## Commands

### Cluster commands

```
Commands for interacting with clusters

Usage:
  oai cluster [command]

Available Commands:
  delete      Delete the specified cluster
  install     Manage cluster install
  list        List available clusters
  show        Show details for a single cluster

Flags:
  -h, --help   help for cluster

Global Flags:
  -u, --api-url string         set logging verbosity (default "https://api.openshift.com/api/assisted-install/v1")
  -f, --config-file string     path to config file
  -t, --offline-token string   offline api token
  -v, --verbose count          set logging verbosity

Use "oai cluster [command] --help" for more information about a command.
```

### Host commands

```
Commands for interacting with hosts in a cluster

Usage:
  oai host [command]

Available Commands:
  delete      Delete hosts from cluster
  find        Find hosts matching criteria
  list        List hosts in the given cluster
  set-name    Set cluster hostnames
  show        Show details for a single host

Flags:
  -c, --cluster string   cluster id or name
  -h, --help             help for host

Global Flags:
  -u, --api-url string         set logging verbosity (default "https://api.openshift.com/api/assisted-install/v1")
  -f, --config-file string     path to config file
  -t, --offline-token string   offline api token
  -v, --verbose count          set logging verbosity

Use "oai host [command] --help" for more information about a command.
```
