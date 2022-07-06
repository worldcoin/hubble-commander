package main

import (
	"bytes"
	"sort"
	"strings"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/dgraph-io/badger/v3"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

//nolint:funlen  // one of the rare times when extracting functions reduces clarity
func auditDatabase(ctx *cli.Context) error {
	log.Info("Iterating the database")

	path := "/home/parallels/wc/hc/db/data/hubble"

	cfg := &config.BadgerConfig{
		Path: path,
	}

	// 1. open the database
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Opened the database: ", path)

	CONTINUE := false
	BREAK := true
	iteration := 0

	binToKeySize := make(map[string]uint64)
	binToValueSize := make(map[string]uint64)

	err = database.Iterator([]byte{}, db.PrefetchIteratorOpts, func(item *badger.Item) (bool, error) {
		key := item.Key()
		var value []byte

		err = item.Value(func(val []byte) error {
			value = val
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}

		if bytes.HasPrefix(key, []byte("_bhIndex:")) {
			asString := string(key)
			pieces := strings.Split(asString, ":")
			bin := strings.Join(pieces[:3], ":")

			binToKeySize[bin] += uint64(len(key))
			binToValueSize[bin] += uint64(len(value))
		} else if bytes.HasPrefix(key, []byte("bh_")) {
			asString := string(key)
			pieces := strings.Split(asString, ":")
			bin := pieces[0]

			binToKeySize[bin] += uint64(len(key))
			binToValueSize[bin] += uint64(len(value))
		} else {
			iteration += 1
			bin := "unknown"
			log.Infof("Unknown item. key=%q value=%x", key, value)
			binToKeySize[bin] += uint64(len(key))
			binToValueSize[bin] += uint64(len(value))
		}

		if iteration > 5 {
			return BREAK, nil
		} else {
			return CONTINUE, nil
		}
	})
	if err != nil && !errors.Is(err, db.ErrIteratorFinished) {
		log.Fatal(err)
	}

	totalKeys, totalValues := uint64(0), uint64(0)

	binToTotalSize := make(map[string]uint64)
	binKeys := make([]string, 0, len(binToKeySize))
	for bin := range binToKeySize {
		binKeys = append(binKeys, bin)
		binToTotalSize[bin] = binToKeySize[bin] + binToValueSize[bin]

		totalKeys += binToKeySize[bin]
		totalValues += binToValueSize[bin]
	}

	log.Infof(
		"Total accounted size: total=%s keys=%s values=%s",
		humanize.IBytes(totalKeys+totalValues),
		humanize.IBytes(totalKeys),
		humanize.IBytes(totalValues),
	)

	sort.SliceStable(binKeys, func(i, j int) bool {
		return binToTotalSize[binKeys[i]] > binToTotalSize[binKeys[j]]
	})

	log.Info("Largest bins: ", path)
	for _, bin := range binKeys {
		log.Infof(
			"- bin=%q total=%s keys=%s values=%s",
			bin,
			humanize.IBytes(binToTotalSize[bin]),
			humanize.IBytes(binToKeySize[bin]), humanize.IBytes(binToValueSize[bin]),
		)
	}

	return nil
}
