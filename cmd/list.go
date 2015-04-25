// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/codegangsta/cli"

	"github.com/DevMine/srctool/config"
	"github.com/DevMine/srctool/log"
)

// List command is used to list installed parsers or available parsers.
func List(c *cli.Context) {
	// this will create the config dir if it does not already exist
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if c.Bool("r") {
		remoteParsers := getRemoteParsers(cfg.DownloadServerURL)

		fmt.Println("available parsers:")
		for _, parser := range remoteParsers {
			fmt.Println("  * ", parser)
		}
		return
	}

	listInstalledParsers()
}

func listInstalledParsers() {
	fis, err := ioutil.ReadDir(filepath.Join(config.DataDir(), config.ParsersFolder))
	if err != nil {
		log.Fatal(err)
	}

	if len(fis) == 0 {
		fmt.Println("no parser installed")
		return
	}

	fmt.Println("available parsers:")
	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}

		if marched, err := filepath.Match("parser-*", fi.Name()); err != nil {
			log.Debug(err)
			continue
		} else if !marched {
			continue
		}

		fmt.Println("  * ", formatParserName(fi.Name()))
	}
}
