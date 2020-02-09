package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var endpoint string
var token string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mservctl",
	Short: "CLI to control an Mserv instance",
	Long: `mservctl is a CLI application that enables listing and operating middleware in an Mserv instance.
Use a config file (by default at $HOME/.mservctl.yaml) in order to configure the Mserv to use with the CLI.
Alternatively pass the values with command line arguments, e.g.:
$ mservctl list -e https://remote.mserv:8989`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mservctl.yaml)")
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "mserv endpoint")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "mserv security token")

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
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".mservctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mservctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if endpoint != "" {
		viper.Set("endpoint", endpoint)
	}
	if token != "" {
		viper.Set("token", token)
	}
}
