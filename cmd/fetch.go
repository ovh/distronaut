package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ovh/distronaut/pkg/distro"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch available distribution links from configured sources (may take some time)",
	Run: func(cmd *cobra.Command, args []string) {
		src, err := distro.FetchSources(config, filter)
		if err != nil {
			log.Errorf("%s", err)
			os.Exit(1)
		}
		j, err := json.MarshalIndent(src, "", "  ")
		if err != nil {
			log.Errorf("%s", err)
			os.Exit(1)
		}
		if outputFile != "" {
			if err := os.MkdirAll(filepath.Dir(outputFile), os.ModePerm); err != nil {
				log.Errorf("%s", err)
				os.Exit(1)
			}
			f, err := os.Create(outputFile)
			if err != nil {
				log.Errorf("%s", err)
				os.Exit(1)
			}
			defer f.Close()
			if _, err := f.Write(j); err != nil {
				log.Errorf("%s", err)
				os.Exit(1)
			}
		} else {
			fmt.Println(string(j))
		}
	},
}
