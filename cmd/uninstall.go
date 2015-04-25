// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"os"

	"github.com/codegangsta/cli"

	"github.com/DevMine/srctool/config"
	"github.com/DevMine/srctool/log"
)

// Uninstall command uninstalls one or all language parser(s).
func Uninstall(c *cli.Context) {
	if !c.Args().Present() {
		uninstallAll(c.Bool("dry"))
	} else {
		if err := uninstall(genParserName(c.Args().First()), c.Bool("dry"), true); err != nil {
			log.Fatal(err)
		}
	}
}

func uninstallAll(dryMode bool) {
	parsers := getInstalledParsers()
	for _, parser := range parsers {
		if err := uninstall(genParserName(parser), dryMode, true); err != nil {
			log.Fail(err)
		}
	}
}

func uninstall(parserName string, dryMode bool, verbose bool) error {
	parserPath := config.ParserPath(parserName)

	if _, err := os.Stat(parserPath); os.IsNotExist(err) {
		log.Debug(err)
		return errors.New("the parser is not installed")
	}

	if dryMode {
		log.Info("parser path:", parserPath)
		return nil
	}

	log.Debug("removing ", parserPath)
	if err := os.RemoveAll(parserPath); err != nil {
		log.Debug(err)
		return errors.New("failed to remove the parser")
	}

	if verbose {
		log.Success("parser successfully removed")
	}
	return nil
}
