package cli

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func configureLogging(cmd *cobra.Command) error {
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

		viper.AddConfigPath(home)
		viper.SetConfigName(".oaitool")
	}

	viper.SetEnvPrefix("oai")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return nil
}

func NewCmdRoot() *cobra.Command {
	cmd := cobra.Command{
		Use:   "oai",
		Short: "A tool for interacting with the OpenShift Assisted Installer API",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := configureLogging(cmd); err != nil {
				return err
			}

			if err := initConfig(cmd); err != nil {
				return err
			}

			token := viper.GetString("offline-token")

			log.Debugf("using config file: %s\n", viper.ConfigFileUsed())
			log.Debugf("found token: %s\n", token)

			return nil
		},
	}

	cmd.PersistentFlags().StringP("config-file", "f", "", "path to config file")
	cmd.PersistentFlags().StringP("offline-token", "t", "", "offline api token")
	cmd.PersistentFlags().StringP("access-token", "T", "", "api access token")
	cmd.PersistentFlags().StringP("output-format", "o", "text", "select output format")
	cmd.PersistentFlags().CountP("verbose", "v", "set logging verbosity")

	viper.BindPFlag("offline-token", cmd.PersistentFlags().Lookup("offline-token"))

	cmd.AddCommand(
		NewCmdCluster(),
	)

	return &cmd
}
