// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/codegangsta/cli"

	"github.com/DevMine/srctool/config"
	"github.com/DevMine/srctool/log"
)

const (
	parserNameFmt = "parser-%s" // parser name format
)

// Install command installs a language parser.
// It expects only one command line argument: the parser name.
func Install(c *cli.Context) {
	if !c.Args().Present() {
		log.Fatal("expected 1 argument, found 0")
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	parserName := fmt.Sprintf(parserNameFmt, c.Args().First())

	if err := downloadParser(cfg.DownloadServerURL, parserName); err != nil {
		log.Fatal(err)
	}

	if err := installParser(parserName); err != nil {
		log.Fatal(err)
	}
}

func installParser(parserName string) error {
	if err := uncompressParser(parserName, config.ParsersDir()); err != nil {
		log.Debug(err)
		return errors.New("failed to uncompress the parser")
	}

	md5sum, err := checksum(config.TempPath(parserName))
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(config.LocalChecksumPath(parserName), []byte(md5sum), 0644); err != nil {
		log.Debug(err)
		return errors.New("failed to write MD5SUM file in the parser directory")
	}

	log.Success("parser successfully installed")
	return nil
}
