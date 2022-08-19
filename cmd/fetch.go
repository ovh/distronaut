package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"encoding/json"
	"github.com/ovh/distronaut/pkg/distro"
	"os"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch available distribution links from configured sources (may take some time)",
	Run: func(cmd *cobra.Command, args []string) {
		src, err := distro.FetchSources(config, filter)
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}
		j, err := json.MarshalIndent(src, "", "  ")
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}
		fmt.Println(string(j))
	},
}
