package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/larsks/oaitool/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdHostShow(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "show --cluster <cluster_name_or_id> <host_name_or_id>",
		Short:         "Show details for a single host",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
			}
			log.Debugf("found cluster %s", cluster.ID)

			host, err := ctx.api.FindHost(cluster.ID, args[0])
			if err != nil {
				return err
			}

			inventory, err := host.GetInventory()
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			fmt.Fprintf(w, "ID\t%s\n", host.ID)
			fmt.Fprintf(w, "Manufacturer\t%s\n", inventory.SystemVendor.Manufacturer)
			fmt.Fprintf(w, "Model\t%s\n", inventory.SystemVendor.ProductName)
			fmt.Fprintf(w, "Serial\t%s\n", inventory.SystemVendor.SerialNumber)
			fmt.Fprintf(w, "Role\t%s\n", host.Role)
			fmt.Fprintf(w, "Status\t%s\n", host.Status)
			fmt.Fprintf(w, "Stage\t%s\n", host.HostProgress.CurrentStage)
			fmt.Fprintf(w, "BMC Address\t%s\n", inventory.BmcAddress)
			fmt.Fprintf(w, "Architecture\t%s\n", inventory.CPU.Architecture)
			fmt.Fprintf(w, "CPU Model\t%s\n", inventory.CPU.ModelName)
			fmt.Fprintf(w, "Memory\t%d\n", inventory.Memory.PhysicalBytes/1024/1024/1024)
			fmt.Fprintf(w, "Interfaces\n")
			for _, iface := range inventory.Interfaces {
				var speed string

				if iface.SpeedMbps > 0 {
					speed = fmt.Sprintf("%d", iface.SpeedMbps)
				} else {
					speed = "-"
				}

				addresses := strings.Join(iface.Ipv4Addresses, " ")
				fmt.Fprintf(w,
					"\t%s\t%s\t%d\t%s\t%s\n",
					iface.Name, iface.MacAddress, iface.Mtu, speed, addresses)
			}

			fmt.Fprintf(w, "Disks\n")
			for _, disk := range inventory.Disks {
				if !disk.Bootable {
					continue
				}

				size := disk.SizeBytes / 1024 / 1024 / 1024
				fmt.Fprintf(w, "\t%s\t%s\t%s\t%s\t%d\n",
					disk.Name, disk.Serial, disk.Vendor, disk.Model, size)
			}

			w.Flush()

			return nil
		},
	}

	return &cmd
}

func listHosts(hosts []api.Host) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, host := range hosts {
		inventory, err := host.GetInventory()
		if err != nil {
			return err
		}

		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\t%s\n",
			host.ID, host.RequestedHostname, host.Role,
			inventory.BmcAddress, host.Status,
		)
	}
	w.Flush()

	return nil
}

func NewCmdHostList(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "list --cluster <name_or_id>",
		Short:         "List hosts in the given cluster",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
			}

			return listHosts(cluster.Hosts)
		},
	}

	return &cmd
}

func NewCmdHostDelete(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "delete --cluster <cluster_id> <host_id_orname> [...]",
		Short:         "Delete hosts from cluster",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no hostnames provided")
			}

			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
			}

			for _, name := range args {
				host, err := ctx.api.FindHost(cluster.ID, name)
				if err != nil {
					return err
				}

				log.Infof("deleting host %s (%s)", name, host.ID)
				if err := ctx.api.DeleteHost(cluster.ID, host.ID); err != nil {
					return err
				}
			}

			return nil
		},
	}

	return &cmd
}

func NewCmdHostSetName(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "set-name --cluster <cluster_id> [<host_id> <name> [...]]",
		Short:         "Set cluster hostnames",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no hostnames provided")
			}
			if len(args)%2 != 0 {
				return fmt.Errorf("wrong number of arguments")
			}

			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
			}

			pos := 0
			var hostnames []api.HostName
			for pos < len(args) {
				hostid := args[pos]
				hostname := args[pos+1]
				log.Infof("setting hostname %s = %s", hostid, hostname)
				pos += 2

				spec := api.HostName{
					ID:       hostid,
					HostName: hostname,
				}

				hostnames = append(hostnames, spec)
			}

			if err := ctx.api.SetHostnames(cluster.ID, hostnames); err != nil {
				return err
			}
			return nil
		},
	}

	return &cmd
}

func findHostByMac(hosts []api.Host, value string) []api.Host {
	var work []api.Host

	log.Debugf("searching for mac = %s", value)

	for _, host := range hosts {
		inventory, err := host.GetInventory()
		if err != nil {
			continue
		}

		for _, iface := range inventory.Interfaces {
			log.Debugf("want %s, have %s", value, iface.MacAddress)
			if iface.MacAddress == value {
				work = append(work, host)
				break
			}
		}
	}

	return work
}

func findHostByBmcAddress(hosts []api.Host, value string) []api.Host {
	var work []api.Host

	log.Debugf("searching for bmc address = %s", value)

	for _, host := range hosts {
		inventory, err := host.GetInventory()
		if err != nil {
			continue
		}

		if inventory.BmcAddress == value {
			work = append(work, host)
		}
	}

	return work
}

func findHostByVendor(hosts []api.Host, value string) []api.Host {
	var work []api.Host

	for _, host := range hosts {
		inventory, err := host.GetInventory()
		if err != nil {
			continue
		}

		if inventory.SystemVendor.Manufacturer == value {
			work = append(work, host)
		}
	}

	return work
}

func findHostByProduct(hosts []api.Host, value string) []api.Host {
	var work []api.Host

	for _, host := range hosts {
		inventory, err := host.GetInventory()
		if err != nil {
			continue
		}

		if inventory.SystemVendor.ProductName == value {
			work = append(work, host)
		}
	}

	return work
}

func NewCmdHostFind(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "find --cluster <cluster_id> -m name=value [-m ...]",
		Short:         "Find hosts matching criteria",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			match, err := cmd.Flags().GetStringArray("match")
			if err != nil {
				return err
			}

			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
			}

			selected := cluster.Hosts
			for _, spec := range match {
				parsed := strings.SplitN(spec, "=", 2)
				name := parsed[0]
				value := parsed[1]

				if len(parsed) != 2 {
					return fmt.Errorf("invalid match specification: %s",
						spec)
				}

				log.Debugf("searching for name=%s, val=%s\n",
					parsed[0], parsed[1])

				switch name {
				case "mac":
					selected = findHostByMac(selected, value)
				case "bmc_address", "bmc-address":
					selected = findHostByBmcAddress(selected, value)
				case "vendor":
					selected = findHostByVendor(selected, value)
				case "product":
					selected = findHostByProduct(selected, value)
				default:
					return fmt.Errorf("unsupported search key: %s", name)
				}
			}

			if len(selected) > 0 {
				return listHosts(selected)
			}

			return fmt.Errorf("no hosts matched your criteria")
		},
	}

	cmd.Flags().StringArrayP("match", "m", nil, "match criteria")

	return &cmd
}

func NewCmdHostWaitForStatus(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "wait-for-status [--hosts <count>] [--interval <seconds>] [--retries <retries>] [--timeout <seconds>] <status>",
		Short:         "Wait until hosts in cluster reach the named status",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
			}

			hostcount, err := cmd.Flags().GetInt("hosts")
			if err != nil {
				return err
			}
			if hostcount == 0 {
				hostcount = len(cluster.Hosts)
			}

			retries, err := cmd.Flags().GetInt("retries")
			if err != nil {
				return err
			}

			interval, err := cmd.Flags().GetInt("interval")
			if err != nil {
				return err
			}

			timeout, err := cmd.Flags().GetInt("timeout")
			if err != nil {
				return err
			}

			desired_status := args[0]
			if !api.ValidateHostStatus(desired_status) {
				return fmt.Errorf("invalid status")
			}

			log.Infof("waiting for %d hosts in cluster %s to reach status %s",
				hostcount, cluster.Name, desired_status)

			retry_count := 0
			time_start := time.Now()
			for {
				hosts_with_status := 0

				for _, host := range cluster.Hosts {
					log.Debugf("checking status for %s, have %s want %s",
						host.ID, host.Status, desired_status)

					if host.Status == desired_status {
						hosts_with_status++
					}
				}

				if hosts_with_status == hostcount {
					break
				}

				if timeout > 0 && time.Since(time_start).Seconds() > float64(timeout) {
					return fmt.Errorf("timed out waiting for status")
				}

				retry_count++
				if retries > 0 && retry_count > retries {
					return fmt.Errorf("too many retries waiting for status")
				}

				time.Sleep(time.Duration(interval) * time.Second)
				cluster, err = getClusterFromFlags(ctx, cmd)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().Int("retries", 0, "Number of times to check status")
	cmd.Flags().Int("interval", 5, "Number of seconds to sleep between retries")
	cmd.Flags().Int("timeout", 0, "Number of seconds after which we timeout")
	cmd.Flags().Int("hosts", 0, "Wait until this many hosts achieve target status")

	return &cmd
}

func NewCmdHost(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:   "host",
		Short: "Commands for interacting with hosts in a cluster",
	}

	cmd.PersistentFlags().String("cluster", "", "cluster id or name")

	cmd.AddCommand(
		NewCmdHostList(ctx),
		NewCmdHostSetName(ctx),
		NewCmdHostShow(ctx),
		NewCmdHostDelete(ctx),
		NewCmdHostFind(ctx),
		NewCmdHostWaitForStatus(ctx),
	)

	return &cmd
}
