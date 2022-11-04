/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "compose",
	Short: "Compose is a tool for defining and running multi-commands. With Compose, you use a YAML file to configure your application's services. Then, with a single command, you create and start all the services from your configuration.",
}

func Execute() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	const defaultCmd = "run"
	if cmd, _, err := rootCmd.Find(os.Args[1:]); err == nil && cmd == rootCmd {
		rootCmd.SetArgs(append([]string{defaultCmd}, os.Args[1:]...))
	}

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger, initConfig)

	rootCmd.PersistentFlags().StringP("file", "f", "", "Specify an alternate compose file")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")
}

func initConfig() {
	viper.SetConfigName("compose")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	if v, _ := rootCmd.Flags().GetString("file"); v != "" {
		viper.SetConfigFile(v)
	}

	if err := viper.ReadInConfig(); err == nil {
		log.Debug().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}
}

func initLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.InfoLevel)
	if ok, _ := rootCmd.Flags().GetBool("debug"); ok {
		log.Logger = log.Level(zerolog.DebugLevel)
	}
}
