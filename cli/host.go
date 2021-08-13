package cli

import (
	"github.com/spf13/cobra"
)

func NewCmdHost() *cobra.Command {
	cmd := cobra.Command{
		Use:   "host",
		Short: "Commands for interacting with hosts in a cluster",
	}

	return &cmd
}
