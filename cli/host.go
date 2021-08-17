package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/larsks/oaitool/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getClusterFromFlags(ctx *Context, cmd *cobra.Command) (*api.ClusterDetail, error) {
	clusterid, err := cmd.Flags().GetString("cluster")
	if err != nil {
		return nil, err
	}

	cluster, err := ctx.api.FindCluster(clusterid)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func NewCmdHostShow(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "show --cluster <cluster_name_or_id> <host_name_or_id>",
		Short:         "List hosts in the given cluster",
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
				fmt.Fprintf(w, "\t%s\t%s\t%s\t%s\n", iface.Name, iface.MacAddress, speed, addresses)
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

	cmd.Flags().StringP("cluster", "c", "", "cluster id or name")
	cmd.MarkFlagRequired("cluster")

	return &cmd
}

func NewCmdHostList(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "list --cluster <name_or_id>",
		Short:         "List hosts in the given cluster",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			for _, host := range cluster.Hosts {
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
		},
	}

	cmd.Flags().StringP("cluster", "c", "", "cluster id or name")
	cmd.MarkFlagRequired("cluster")

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

	cmd.Flags().StringP("cluster", "c", "", "cluster id or name")
	cmd.MarkFlagRequired("cluster")

	return &cmd
}

func NewCmdHost(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:   "host",
		Short: "Commands for interacting with hosts in a cluster",
	}

	cmd.AddCommand(
		NewCmdHostList(ctx),
		NewCmdHostSetName(ctx),
		NewCmdHostShow(ctx),
	)

	return &cmd
}
