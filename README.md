# oaitool

Oaitool is a tool for interacting with the OpenShift Assisted
Installer API.

## Building

To build this package you will need at least `go` version 1.15.

To build `oaitool`:

- Clone this repository
- Run `make`

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
  oaitool cluster [command]

Available Commands:
  create          Create an assisted installer cluster
  delete          Delete the specified cluster
  get-file        Get file from cluster
  get-image-url   Get discovery image download url
  get-kubeconfig  Get cluster kubeconfig
  install         Manage cluster install
  list            List available clusters
  set-vips        Create an assisted installer cluster
  show            Show details for a single cluster
  status          Get cluster status
  wait-for-status Wait until cluster reaches the named status

Flags:
      --cluster string   cluster id or name
  -h, --help             help for cluster

Global Flags:
  -u, --api-url string         set logging verbosity (default "https://api.openshift.com/api/assisted-install/v1")
  -f, --config-file string     path to config file
  -t, --offline-token string   offline api token
  -v, --verbose count          set logging verbosity

Use "oaitool cluster [command] --help" for more information about a command.
```

### Host commands

```
Commands for interacting with hosts in a cluster

Usage:
  oaitool host [command]

Available Commands:
  delete          Delete hosts from cluster
  find            Find hosts matching criteria
  list            List hosts in the given cluster
  set-name        Set cluster hostnames
  show            Show details for a single host
  wait-for-status Wait until hosts in cluster reach the named status

Flags:
      --cluster string   cluster id or name
  -h, --help             help for host

Global Flags:
  -u, --api-url string         set logging verbosity (default "https://api.openshift.com/api/assisted-install/v1")
  -f, --config-file string     path to config file
  -t, --offline-token string   offline api token
  -v, --verbose count          set logging verbosity

Use "oaitool host [command] --help" for more information about a command.
```
