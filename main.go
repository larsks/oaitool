package main

import (
	"github.com/larsks/oaitool/cli"
	"github.com/spf13/cobra"
)

func main() {
	root := cli.NewCmdRoot()
	cobra.CheckErr(root.Execute())
}
