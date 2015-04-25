// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"

	"github.com/DevMine/srctool/config"
	"github.com/DevMine/srctool/log"
)

// Update command updates an installed parser.
// It expects only one command line argument: the parser name.
func Update(c *cli.Context) {
	if !c.Args().Present() {
		log.Fatal("expected 1 argument, found 0")
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	parserName := "parser-" + c.Args().First()

	if !isAlreadyInstalled(parserName) {
		log.Fatal(" parser not installed, install it first")
	}

	LocalChecksum, err := ioutil.ReadFile(config.LocalChecksumPath(parserName))
	if err != nil {
		log.Debug(err)
		log.Fatal("unable to read the MD5 file of the currently installed parser")
	}

	remoteCheckSum, err := fetchChecksum(cfg.DownloadServerURL, parserName)
	if err != nil {
		log.Fail(err)
	}

	if string(LocalChecksum) == remoteCheckSum {
		log.Info("latest version already installed")
		os.Exit(0)
	}

	if err = downloadParser(cfg.DownloadServerURL, parserName); err != nil {
		log.Fatal(err)
	}

	if err = uninstallParser(parserName, false); err != nil {
		log.Fatal(err)
	}

	if err = installParser(parserName); err != nil {
		log.Fatal(err)
	}

	log.Success("parser successfully updated")
}

func isAlreadyInstalled(parserName string) bool {
	parserPath := config.ParserPath(parserName)
	log.Debug("=> ", parserPath)

	if _, err := os.Stat(parserPath); err != nil {
		log.Debug(err)
		return false
	}

	return true
}
