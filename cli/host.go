package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

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

func NewCmdHost(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:   "host",
		Short: "Commands for interacting with hosts in a cluster",
	}

	cmd.AddCommand(
		NewCmdHostList(ctx),
	)

	return &cmd
}
