package cmd

import (
	"net/url"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/TykTechnologies/mserv/mservclient/client"
	"github.com/TykTechnologies/mserv/util/logger"
)

var (
	cfgFile  string
	mservapi *client.Mserv
)

var log = logger.GetLogger("mservctl")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mservctl",
	Short: "CLI to control an Mserv instance",
	Long: `mservctl is a CLI application that enables listing and operating middleware in an Mserv instance.
Use a config file (by default at $HOME/.mservctl.yaml) in order to configure the Mserv to use with the CLI.
Alternatively, pass the values with command line arguments.

Set TYK_MSERV_LOGLEVEL="debug" environment variable to see raw API requests and responses.
`,
	Example: "mservctl list -e https://remote.mserv:8989",
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
	rootCmd.PersistentFlags().BoolP("insecure-tls", "k", false, "allow insecure TLS for mserv client")

	viper.BindPFlag("endpoint", rootCmd.Flag("endpoint"))
	viper.BindPFlag("token", rootCmd.Flag("token"))
	viper.BindPFlag("insecure_tls", rootCmd.Flag("insecure-tls"))

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
	endpoint, err := url.Parse(viper.GetString("endpoint"))
	if err != nil {
		log.WithError(err).Fatal("Couldn't parse the mserv endpoint")
	}

	tlsOptions := httptransport.TLSClientOptions{InsecureSkipVerify: viper.GetBool("insecure_tls")}
	tlsClient, err := httptransport.TLSClient(tlsOptions)
	if err != nil {
		log.WithError(err).Fatal("Couldn't create client with TLS options")
	}

	tr := httptransport.NewWithClient(endpoint.Host, endpoint.Path, []string{endpoint.Scheme}, tlsClient)
	tr.SetLogger(log)
	tr.SetDebug(log.Logger.GetLevel() >= logrus.DebugLevel)

	mservapi = client.New(tr, nil)
}

func defaultAuth() runtime.ClientAuthInfoWriter {
	return httptransport.APIKeyAuth("X-Api-Key", "header", viper.GetString("token"))
}
