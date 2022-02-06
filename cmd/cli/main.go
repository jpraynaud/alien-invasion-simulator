package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/jpraynaud/alien-invasion-simulator/cmd/cli/cmd"
)

func main() {
	cmd.Execute()
}

func init() {
	log.SetLevel(log.WarnLevel)
}
