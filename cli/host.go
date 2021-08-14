package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/larsks/oaitool/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdHostList(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "list --cluster <name_or_id>",
		Short:         "List hosts in the given cluster",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			clusterid, err := cmd.Flags().GetString("cluster")
			if err != nil {
				return err
			}

			cluster, err := ctx.api.FindCluster(clusterid)
			if err != nil {
				return err
			}
			if cluster == nil {
				return fmt.Errorf("unable to find cluster named \"%s\"", clusterid)
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			for _, host := range cluster.Hosts {
				inventory, err := host.GetInventory()
				if err != nil {
					return err
				}

				fmt.Fprintf(
					w,
					"%s\t%s\t%s\t%s\n",
					host.ID, host.RequestedHostname, host.Role,
					inventory.BmcAddress,
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

			clusterid, err := cmd.Flags().GetString("cluster")
			if err != nil {
				return err
			}

			cluster, err := ctx.api.FindCluster(clusterid)
			if err != nil {
				return err
			}
			if cluster == nil {
				return fmt.Errorf("unable to find cluster named \"%s\"", clusterid)
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
	)

	return &cmd
}
