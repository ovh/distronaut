package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"encoding/json"
	"github.com/ovh/distronaut/pkg/distro"
	"os"
)

var sourceCmd = &cobra.Command{
	Use:   "source",
	Short: "List configured sources",
	Run: func(cmd *cobra.Command, args []string) {
		src, err := distro.ListSources(config, filter)
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
