package cli

import (
	"github.com/spf13/cobra"
)

func NewCmdClusterList() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "List available clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
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
