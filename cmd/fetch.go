package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ovh/distronaut/pkg/distro"
	"encoding/json"
	"os"

)

var config string
var filter string

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch available distribution links from sources",
	Run: func(cmd *cobra.Command, args []string) {
		jd, err := json.MarshalIndent(distro.Fetch(config, filter), "", "  ")
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}
		fmt.Println(string(jd))
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&config, "config", "c", "config/sources.yml", "Config file")
	fetchCmd.Flags().StringVarP(&filter, "filter", "f", ".*", "Filter regex")
}
