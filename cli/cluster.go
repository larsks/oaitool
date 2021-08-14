package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdClusterList(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "list",
		Short:         "List available clusters",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			clusters, err := ctx.api.ListClusters()
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			for _, cluster := range clusters {
				fmt.Fprintf(
					w,
					"%s\t%s\t%s\t%s\n",
					cluster.Name,
					cluster.BaseDNSDomain,
					cluster.ID,
					cluster.Status,
				)
			}
			w.Flush()

			return nil
		},
	}

	return &cmd
}

func NewCmdClusterShow(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "show <name_or_id>",
		Short:         "Show details for a single cluster",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("missing cluster name or id")
			}

			clusterid := args[0]
			log.Debugf("look up cluster %s", clusterid)
			cluster, err := ctx.api.FindCluster(clusterid)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			fmt.Fprintf(w, "Name\t%s\n", cluster.Name)
			fmt.Fprintf(w, "BaseDNSDomain\t%s\n", cluster.BaseDNSDomain)
			fmt.Fprintf(w, "ID\t%s\n", cluster.ID)
			fmt.Fprintf(w, "EnabledHostCount\t%d\n", cluster.EnabledHostCount)
			fmt.Fprintf(w, "APIVip\t%s\n", cluster.APIVip)
			fmt.Fprintf(w, "IngressVip\t%s\n", cluster.IngressVip)
			fmt.Fprintf(w, "OpenshiftVersion\t%s\n", cluster.OpenshiftVersion)
			fmt.Fprintf(w, "Status\t%s\n", cluster.Status)
			w.Flush()

			return nil
		},
	}

	return &cmd
}

func NewCmdCluster(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:   "cluster",
		Short: "Commands for interacting with clusters",
	}

	cmd.AddCommand(
		NewCmdClusterList(ctx),
		NewCmdClusterShow(ctx),
	)

	return &cmd
}
