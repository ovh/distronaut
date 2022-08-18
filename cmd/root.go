package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var config string
var filter string
var loglevel string

var rootCmd = &cobra.Command{
	Use:   "distronaut",
	Short: "Distronaut is a tool to fetch iso download links and metadata from across the web",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		lv, err := log.ParseLevel(loglevel)
		if err != nil {
			log.Warnf("failed to parse loglevel: %s (%s)", loglevel, err)
		}
		log.SetLevel(lv)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(sourceCmd)
	rootCmd.PersistentFlags().StringVarP(&loglevel, "log", "l", "warn", "log level")
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/sources.yml", "config file")
	rootCmd.PersistentFlags().StringVarP(&filter, "filter", "f", ".*", "filter regex")
}
