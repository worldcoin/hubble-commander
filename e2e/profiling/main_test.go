// +build profiling

package profiling

import (
	"log"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
)

func TestCommander(t *testing.T) {
	cfg := config.GetConfig()
	cfg.Rollup.Prune = true

	cmd := commander.NewCommander(cfg)

	err := cmd.Start()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	time.Sleep(60 * time.Second)

	err = cmd.Stop()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
