// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

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

	installedParsers := getInstalledParsers()
	if len(installedParsers) == 0 {
		fmt.Println("no parser installed")
		return
	}

	fmt.Println("installed parsers:")
	for _, parser := range installedParsers {
		fmt.Println("  * ", parser)
	}
}
