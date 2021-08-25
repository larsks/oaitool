package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"
	"time"

	"github.com/larsks/oaitool/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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

			networkType, err := cmd.Flags().GetString("network-type")
			if err != nil {
				return err
			}

			if !api.ValidateNetworkType(networkType) {
				return fmt.Errorf("invalid network type")
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
				Name:             args[0],
				PullSecret:       string(psjson),
				OpenshiftVersion: openshiftVersion,
				BaseDnsDomain:    baseDnsDomain,
				SshPublicKey:     string(sshKey),
				NetworkType:      networkType,
			}

			log.Infof("creating cluster %s", createParams.Name)
			log.Debugf("creating cluster with parameters: %+v", createParams)
			cluster, err := ctx.api.CreateCluster(&createParams)
			if err != nil {
				return err
			}

			fmt.Printf("%s %s\n", cluster.Name, cluster.ID)

			return nil
		},
	}

	cmd.Flags().String("pull-secret", "", "Read pull secret from a file")
	cmd.Flags().String("openshift-version", "", "Set OpenShift version")
	cmd.Flags().String("base-domain", "", "Base DNS Domain")
	cmd.Flags().String("ssh-public-key", "", "Public ssh key")
	cmd.Flags().String("network-type", "OpenShiftSDN", "Network type")

	return &cmd
}

func NewCmdClusterSetVips(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "set-vips --api-vip a.b.c.d --ingress-vip w.x.y.z",
		Short:         "Create an assisted installer cluster",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
			}

			if cluster.MachineNetworkCidr == "" {
				return fmt.Errorf("cluster does not have a machine network defined")
			}

			apiVip, err := cmd.Flags().GetString("api-vip")
			if err != nil {
				return err
			}

			ingressVip, err := cmd.Flags().GetString("ingress-vip")
			if err != nil {
				return err
			}

			networkPatch := api.ClusterNetworkPatch{
				ApiVip:            apiVip,
				IngressVip:        ingressVip,
				VipDhcpAllocation: false,
			}
			log.Debugf("patching cluster network configuration: %+v", networkPatch)
			_, err = ctx.api.PatchCluster(cluster.ID, &networkPatch)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().String("api-vip", "", "API VIP")
	cmd.Flags().String("ingress-vip", "", "Ingress VIP")
	if err := cmd.MarkFlagRequired("api-vip"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("ingress-vip"); err != nil {
		panic(err)
	}

	return &cmd
}

func NewCmdClusterInstall(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "install (--start | --cancel | --reset )",
		Short:         "Manage cluster install",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromFlags(ctx, cmd)
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

				if err != nil {
					return err
				}
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
		Use:           "delete",
		Short:         "Delete the specified cluster",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
			}

			log.Infof("deleting cluster %s", cluster.Name)
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
		Use:           "show",
		Short:         "Show details for a single cluster",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
			}

			use_json, err := cmd.Flags().GetBool("json")
			if err != nil {
				return err
			}

			if use_json {
				clusterJson, err := cluster.ToJSON()
				if err != nil {
					return err
				}

				os.Stdout.Write(clusterJson)
			} else {

				w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
				fmt.Fprintf(w, "Name\t%s\n", cluster.Name)
				fmt.Fprintf(w, "BaseDNSDomain\t%s\n", cluster.BaseDNSDomain)
				fmt.Fprintf(w, "ID\t%s\n", cluster.ID)
				fmt.Fprintf(w, "EnabledHostCount\t%d\n", cluster.EnabledHostCount)
				fmt.Fprintf(w, "ApiVip\t%s\n", cluster.ApiVip)
				fmt.Fprintf(w, "IngressVip\t%s\n", cluster.IngressVip)
				fmt.Fprintf(w, "OpenshiftVersion\t%s\n", cluster.OpenshiftVersion)
				fmt.Fprintf(w, "Status\t%s\n", cluster.Status)
				w.Flush()
			}

			return nil
		},
	}

	cmd.Flags().BoolP("json", "j", false, "Show full JSON data")

	return &cmd
}

func NewCmdClusterStatus(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "status",
		Short:         "Get cluster status",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromFlags(ctx, cmd)
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
		Use:           "get-image-url",
		Short:         "Get discovery image download url",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
			}

			if cluster.ImageInfo.DownloadUrl == "" {
				imageType, err := cmd.Flags().GetString("image-type")
				if err != nil {
					return err
				}
				if !api.ValidateImageType(imageType) {
					return fmt.Errorf("invalid image type")
				}

				log.Info("generating discovery image")
				cluster, err = ctx.api.CreateDiscoveryImage(cluster.ID, imageType, "")
				if err != nil {
					return err
				}

				if cluster.ImageInfo.DownloadUrl == "" {
					return fmt.Errorf("failed to retrieve discovery image url")
				}
			}

			log.Debugf("image info: %+v", cluster.ImageInfo)
			fmt.Println(cluster.ImageInfo.DownloadUrl)
			return nil
		},
	}

	cmd.Flags().String("image-type", "minimal-iso", "set discovery image type")

	return &cmd
}

func NewCmdClusterGetKubeconfig(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "get-kubeconfig <name_or_id>",
		Short:         "Get cluster kubeconfig",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
			}

			kubeconfig, err := ctx.api.GetKubeconfig(cluster.ID)
			if err != nil {
				return err
			}

			os.Stdout.Write(kubeconfig)

			return nil
		},
	}

	return &cmd
}

func NewCmdClusterGetFile(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "get-file <filename>",
		Short:         "Get file from cluster",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
			}

			if !api.ValidateDownloadFile(args[0]) {
				return fmt.Errorf("invalid filename")
			}

			content, err := ctx.api.GetFile(cluster.ID, args[0])
			if err != nil {
				return err
			}

			os.Stdout.Write(content)

			return nil
		},
	}

	return &cmd
}

func NewCmdClusterWaitForStatus(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:           "wait-for-status [--interval <seconds>] [--retries <retries>] [--timeout <seconds>] <status>",
		Short:         "Wait until cluster reaches the named status",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := getClusterFromFlags(ctx, cmd)
			if err != nil {
				return err
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
			if !api.ValidateClusterStatus(desired_status) {
				return fmt.Errorf("invalid status")
			}

			log.Infof("waiting for cluster %s to reach status %s",
				cluster.Name, desired_status)

			retry_count := 0
			time_start := time.Now()
			for {
				if cluster.Status == desired_status {
					break
				}

				log.Debugf("checking status, have %s want %s",
					cluster.Status, desired_status)

				if timeout > 0 && time.Since(time_start).Seconds() > float64(timeout) {
					return fmt.Errorf("timed out waiting for status")
				}

				retry_count++
				if retries > 0 && retry_count > retries {
					return fmt.Errorf("too many retries waiting for status")
				}

				time.Sleep(time.Duration(interval) * time.Second)
				cluster, err = ctx.api.GetCluster(cluster.ID)
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

	return &cmd
}

func NewCmdCluster(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:   "cluster",
		Short: "Commands for interacting with clusters",
	}

	cmd.PersistentFlags().String("cluster", "", "cluster id or name")

	cmd.AddCommand(
		NewCmdClusterList(ctx),
		NewCmdClusterShow(ctx),
		NewCmdClusterStatus(ctx),
		NewCmdClusterDelete(ctx),
		NewCmdClusterInstall(ctx),
		NewCmdClusterCreate(ctx),
		NewCmdClusterSetVips(ctx),
		NewCmdClusterGetImageUrl(ctx),
		NewCmdClusterGetKubeconfig(ctx),
		NewCmdClusterGetFile(ctx),
		NewCmdClusterWaitForStatus(ctx),
	)

	return &cmd
}
