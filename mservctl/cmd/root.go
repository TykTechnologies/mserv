package cmd

import (
	"net/url"

	"github.com/TykTechnologies/mserv/mservclient/client"
	"github.com/TykTechnologies/mserv/util/logger"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var mservapi *client.Mserv

var log = logger.GetLogger("mservctl")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mservctl",
	Short: "CLI to control an Mserv instance",
	Long: `mservctl is a CLI application that enables listing and operating middleware in an Mserv instance.
Use a config file (by default at $HOME/.mservctl.yaml) in order to configure the Mserv to use with the CLI.
Alternatively pass the values with command line arguments, e.g.:

$ mservctl list -e https://remote.mserv:8989

Set TYK_MSERV_LOGLEVEL="debug" environment variable to see raw API requests and responses.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initMservApi)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mservctl.yaml)")
	rootCmd.PersistentFlags().StringP("endpoint", "e", "", "mserv endpoint")
	rootCmd.PersistentFlags().StringP("token", "t", "", "mserv security token")
	viper.BindPFlag("endpoint", rootCmd.Flag("endpoint"))
	viper.BindPFlag("token", rootCmd.Flag("token"))

	viper.SetDefault("endpoint", "http://localhost:8989")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".mservctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mservctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}
}

func initMservApi() {
	endpoint, err := parseEndpoint(viper.GetString("endpoint"))
	if err != nil {
		log.WithError(err).Fatal("Couldn't parse the mserv endpoint")
	}

	tr := httptransport.New(endpoint.Host, endpoint.Path, []string{endpoint.Scheme})
	tr.SetLogger(log)
	if log.Logger.GetLevel() == logrus.DebugLevel {
		tr.SetDebug(true)
	}

	mservapi = client.New(tr, nil)
}

func parseEndpoint(endpoint string) (*url.URL, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func defaultAuth() runtime.ClientAuthInfoWriter {
	return httptransport.APIKeyAuth("X-Api-Key", "header", viper.GetString("token"))
}
