package main

import (
  log "github.com/sirupsen/logrus"
  "github.com/ovh/distronaut/cmd"
)

func main() {
  log.SetLevel(log.DebugLevel)
  cmd.Execute()
}