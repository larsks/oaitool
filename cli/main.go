package cli

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/larsks/oaitool/api"
	"github.com/larsks/oaitool/version"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type (
	Context struct {
		api *api.ApiClient
	}
)

func getClusterFromFlags(ctx *Context, cmd *cobra.Command) (*api.Cluster, error) {
	clusterid, err := cmd.Flags().GetString("cluster")
	if err != nil {
		return nil, err
	}
	if clusterid == "" {
		return nil, fmt.Errorf("no cluster name provided")
	}

	cluster, err := ctx.api.FindCluster(clusterid)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func initLogging(cmd *cobra.Command) error {
	var loglevel log.Level

	verbose, err := cmd.Flags().GetCount("verbose")
	if err != nil {
		return err
	}

	switch verbose {
	case 0:
		loglevel = log.WarnLevel
	case 1:
		loglevel = log.InfoLevel
	default:
		loglevel = log.DebugLevel
	}

	log.SetLevel(loglevel)
	return nil
}

func NewConfig() *viper.Viper {
	config := viper.New()

	config.SetEnvPrefix("oai")
	config.AutomaticEnv()
	config.SetConfigName(".oiatool")

	replacer := strings.NewReplacer("-", "_")
	config.SetEnvKeyReplacer(replacer)

	log.Debugf("returning new config")

	return config
}

func initConfig(cmd *cobra.Command) error {
	cfgFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return err
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		viper.AddConfigPath(path.Join(home, ".config", "oaitool"))
		viper.AddConfigPath(path.Join(home, ".oaitool"))
		viper.AddConfigPath(".oaitool")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("oai")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	} else {
		log.Debugf("using config file %s", viper.ConfigFileUsed())
	}

	return nil
}

func initContext(cmd *cobra.Command, ctx *Context) error {
	offlinetoken := viper.GetString("offline-token")
	apiurl := viper.GetString("api-url")

	apiclient := api.NewApiClient(offlinetoken, apiurl)
	if apiclient == nil {
		return fmt.Errorf("failed to create api client")
	}

	ctx.api = apiclient

	return nil
}

func NewCmdVersion(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Build version: %s\n", version.BuildVersion)
			fmt.Printf("Build ref: %s\n", version.BuildRef)
			fmt.Printf("Build date: %s\n", version.BuildDate)
		},
	}

	return &cmd
}

func NewCmdRoot() *cobra.Command {
	ctx := &Context{}

	cmd := cobra.Command{
		Use:   "oaitool",
		Short: "A tool for interacting with the OpenShift Assisted Installer API",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := initLogging(cmd); err != nil {
				return err
			}

			if err := initConfig(cmd); err != nil {
				return err
			}

			if err := initContext(cmd, ctx); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringP("config-file", "f", "", "path to config file")
	cmd.PersistentFlags().StringP("offline-token", "t", "", "offline api token")
	cmd.PersistentFlags().CountP("verbose", "v", "set logging verbosity")
	cmd.PersistentFlags().StringP("api-url", "u", "https://api.openshift.com/api/assisted-install/v1", "set logging verbosity")

	viper.BindPFlag("offline-token", cmd.PersistentFlags().Lookup("offline-token"))
	viper.BindPFlag("api-url", cmd.PersistentFlags().Lookup("api-url"))

	cmd.AddCommand(
		NewCmdCluster(ctx),
		NewCmdHost(ctx),
		NewCmdVersion(ctx),
	)

	return &cmd
}
