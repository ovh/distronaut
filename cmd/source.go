package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ovh/distronaut/pkg/distro"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var sourceCmd = &cobra.Command{
	Use:   "source",
	Short: "List configured sources",
	Run: func(cmd *cobra.Command, args []string) {
		src, err := distro.ListSources(config, filter)
		if err != nil {
			log.Errorf("%s", err)
			os.Exit(1)
		}
		j, err := json.MarshalIndent(src, "", "  ")
		if err != nil {
			log.Errorf("%s", err)
			os.Exit(1)
		}
		fmt.Println(string(j))
	},
}
