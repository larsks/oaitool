package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/larsks/oaitool/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdClusterList() *cobra.Command {
	var apiclient *api.ApiClient

	cmd := cobra.Command{
		Use:   "list",
		Short: "List available clusters",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			offlinetoken := viper.GetString("offline-token")
			apiurl := viper.GetString("api-url")

			apiclient = api.NewApiClient(offlinetoken, apiurl)
			if apiclient == nil {
				return fmt.Errorf("failed to create api client")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			clusters, err := apiclient.ListClusters()
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			for _, cluster := range clusters {
				fmt.Fprintf(
					w,
					"%s\t%s\t%s\n",
					cluster.Name,
					cluster.BaseDNSDomain,
					cluster.ID,
				)
			}
			w.Flush()

			return nil
		},
	}

	return &cmd
}

func NewCmdCluster() *cobra.Command {
	cmd := cobra.Command{
		Use:   "cluster",
		Short: "Commands for interacting with clusters",
	}

	cmd.AddCommand(
		NewCmdClusterList(),
	)

	return &cmd
}
