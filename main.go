package main

import (
	"github.com/ovh/distronaut/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	cmd.Execute()
}
