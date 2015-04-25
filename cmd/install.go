// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"io/ioutil"

	"github.com/codegangsta/cli"

	"github.com/DevMine/srctool/config"
	"github.com/DevMine/srctool/log"
)

// Install command installs one or all language parser(s).
func Install(c *cli.Context) {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if !c.Args().Present() {
		installAll(cfg)
	} else {
		if err := install(cfg, genParserName(c.Args().First()), true); err != nil {
			log.Fatal(err)
		}
	}
}

func installAll(cfg *config.Config) {
	parsers := getRemoteParsers(cfg.DownloadServerURL)
	for _, parser := range parsers {
		if err := install(cfg, genParserName(parser), true); err != nil {
			log.Fail(err)
		}
	}
}

func install(cfg *config.Config, parserName string, verbose bool) error {
	if err := downloadParser(cfg.DownloadServerURL, parserName, verbose); err != nil {
		return err
	}

	if err := uncompressParser(parserName, config.ParsersDir()); err != nil {
		log.Debug(err)
		return errors.New("failed to uncompress the " + parserName + " archive")
	}

	md5sum, err := checksum(config.TempPath(parserName))
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(config.LocalChecksumPath(parserName), []byte(md5sum), 0644); err != nil {
		log.Debug(err)
		return errors.New("failed to write MD5SUM file in the parser directory")
	}

	if verbose {
		log.Success(parserName, " successfully installed")
	}
	return nil
}
