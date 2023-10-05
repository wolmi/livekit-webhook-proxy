/*
Copyright © 2023 Miquel Ramon Ortega Tido
*/
package cmd

import (
	"context"
	"fmt"
	"livekit-webhook-proxy/types"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "livekit-webhook-proxy",
	Short: "livekit webhook proxy allows sending events to GCP PubSub",
	Long: `The Livekit webhook proxy is designed to allow configuring it as webhook endpoint
and manages to publish them in a GCP PubSub topic.`,

	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		proxy := types.Proxy{}

		proxy.Init(ctx)

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.livekit-webhook-proxy.yaml)")

	rootCmd.PersistentFlags().Bool("debug", false, "Debug mode")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")) // nolint

	rootCmd.Flags().IntP("port", "p", 8080, "Port to start the proxy")
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port")) // nolint

	rootCmd.Flags().StringP("topic", "t", "", "Topic in the GCP PubSub to publish events")
	viper.BindPFlag("topic", rootCmd.Flags().Lookup("topic")) // nolint

	rootCmd.Flags().StringP("project-id", "P", "", "GCP ProjectId")
	viper.BindPFlag("project-id", rootCmd.Flags().Lookup("project-id")) // nolint
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find current directory.
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in current directory with name ".config" (without extension).
		viper.AddConfigPath(pwd)

		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("WEBHOOK_PROXY")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
