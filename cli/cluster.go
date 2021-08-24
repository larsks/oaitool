package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/larsks/oaitool/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getClusterFromArgs(ctx *Context, args []string) (*api.Cluster, error) {
	if len(args) < 1 || args[0] == "" {
		return nil, fmt.Errorf("missing cluster name")
	}
	clusterid := args[0]

	cluster, err := ctx.api.FindCluster(clusterid)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func NewCmdClusterList(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "list",
		Short:         "List available clusters",
		Args:          cobra.NoArgs,
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

func NewCmdClusterCreate(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "create <name_or_id>",
		Short:         "Create an assisted installer cluster",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var ps *api.PullSecret

			pspath, err := cmd.Flags().GetString("pull-secret")
			if err != nil {
				return err
			}

			if pspath != "" {
				ps, err = api.PullSecretFromFile(pspath)
			} else {
				ps, err = ctx.api.GetPullSecret()
			}

			if err != nil {
				return err
			}

			psjson, err := ps.ToJSON()
			if err != nil {
				return err
			}

			openshiftVersion, err := cmd.Flags().GetString("openshift-version")
			if err != nil {
				return err
			}

			baseDnsDomain, err := cmd.Flags().GetString("base-domain")
			if err != nil {
				return err
			}

			noDhcpAllocation, err := cmd.Flags().GetBool("no-dhcp-allocation")
			if err != nil {
				return err
			}

			apiVip, err := cmd.Flags().GetString("api-vip")
			if err != nil {
				return err
			}

			ingressVip, err := cmd.Flags().GetString("ingress-vip")
			if err != nil {
				return err
			}

			networkType, err := cmd.Flags().GetString("network-type")
			if err != nil {
				return err
			}

			if err = api.ValidateNetworkType(networkType); err != nil {
				return err
			}

			sshKeyFile, err := cmd.Flags().GetString("ssh-public-key")
			if err != nil {
				return err
			}

			var sshKey []byte
			if sshKeyFile != "" {
				log.Debugf("reading ssh key from %s", sshKeyFile)
				sshKey, err = ioutil.ReadFile(sshKeyFile)
				if err != nil {
					return err
				}
			}

			createParams := api.ClusterCreateParams{
				Name:              args[0],
				PullSecret:        string(psjson),
				OpenshiftVersion:  openshiftVersion,
				BaseDnsDomain:     baseDnsDomain,
				VipDhcpAllocation: !noDhcpAllocation,
				ApiVip:            apiVip,
				IngressVip:        ingressVip,
				NetworkType:       networkType,
				SshPublicKey:      string(sshKey),
			}

			log.Debugf("creating cluster with parameters: %+v", createParams)
			detail, err := ctx.api.CreateCluster(&createParams)
			if err != nil {
				return err
			}

			fmt.Printf("%s %s\n", detail.Name, detail.ID)

			return nil
		},
	}

	cmd.Flags().StringP("pull-secret", "s", "", "Read pull secret from a file")
	cmd.Flags().StringP("openshift-version", "o", "", "Set OpenShift version")
	cmd.Flags().StringP("base-domain", "b", "", "Base DNS Domain")
	cmd.Flags().Bool("no-dhcp-allocation", false, "Do not allocate VIP addresses using DHCP")
	cmd.Flags().StringP("api-vip", "a", "", "API VIP address")
	cmd.Flags().StringP("ingress-vip", "i", "", "Ingress VIP address")
	cmd.Flags().StringP("ssh-public-key", "k", "", "Public ssh key")
	cmd.Flags().StringP("network-type", "n", "", "Network type")

	return &cmd
}

func NewCmdClusterInstall(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "install <name_or_id> (--start | --cancel | --reset )",
		Short:         "Manage cluster install",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromArgs(ctx, args)
			if err != nil {
				return err
			}

			action := false
			for _, mode := range []string{"start", "cancel", "reset"} {
				flagval, err := cmd.Flags().GetBool(mode)
				if err != nil {
					return err
				}

				switch {
				case mode == "start" && flagval:
					log.Infof("starting install of cluster %s (%s)", cluster.Name, cluster.ID)
					err = ctx.api.InstallCluster(cluster.ID)
					action = true
				case mode == "cancel" && flagval:
					log.Infof("starting install of cluster %s (%s)", cluster.Name, cluster.ID)
					err = ctx.api.CancelCluster(cluster.ID)
					action = true
				case mode == "reset" && flagval:
					log.Infof("starting install of cluster %s (%s)", cluster.Name, cluster.ID)
					err = ctx.api.ResetCluster(cluster.ID)
					action = true
				}
			}

			if err != nil {
				return err
			}

			if !action {
				log.Warnf("no action specified")
			}

			return nil
		},
	}

	cmd.Flags().Bool("start", false, "path to config file")
	cmd.Flags().Bool("cancel", false, "path to config file")
	cmd.Flags().Bool("reset", false, "path to config file")

	return &cmd
}

func NewCmdClusterDelete(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "delete <name_or_id>",
		Short:         "Delete the specified cluster",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromArgs(ctx, args)
			if err != nil {
				return err
			}

			if err := ctx.api.DeleteCluster(cluster.ID); err != nil {
				return err
			}

			return nil
		},
	}

	return &cmd
}

func NewCmdClusterShow(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "show <name_or_id>",
		Short:         "Show details for a single cluster",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromArgs(ctx, args)
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

func NewCmdClusterStatus(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "status <name_or_id>",
		Short:         "Get cluster status",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromArgs(ctx, args)
			if err != nil {
				return err
			}

			fmt.Println(cluster.Status)
			return nil
		},
	}

	return &cmd
}

func NewCmdClusterGetImageUrl(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "get-image-url <name_or_id>",
		Short:         "Get discovery image download url",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromArgs(ctx, args)
			if err != nil {
				return err
			}

			if cluster.ImageInfo.DownloadUrl == "" {
				imageType, err := cmd.Flags().GetString("image-type")
				if err != nil {
					return err
				}

				log.Info("generating discovery image")
				cluster, err = ctx.api.CreateDiscoveryImage(cluster.ID, imageType, "")
				if err != nil {
					return err
				}
			}

			log.Debugf("image info: %+v", cluster.ImageInfo)
			fmt.Println(cluster.ImageInfo.DownloadUrl)
			return nil
		},
	}

	cmd.Flags().StringP("image-type", "T", "minimal-iso", "set discovery image type")

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
		NewCmdClusterStatus(ctx),
		NewCmdClusterDelete(ctx),
		NewCmdClusterInstall(ctx),
		NewCmdClusterCreate(ctx),
		NewCmdClusterGetImageUrl(ctx),
	)

	return &cmd
}
