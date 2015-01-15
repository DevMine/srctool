// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

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
		listRemoteParsers(cfg.DownloadServerURL)
		return
	}

	listInstalledParsers()
}

func listRemoteParsers(serverURL string) {
	md5sums, err := fetchChecksumsFile(serverURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("remotely available parsers:")
	for _, line := range strings.Split(md5sums, "\n") {
		tmp := strings.Split(line, " ")
		if len(tmp) != 2 {
			continue
		}

		path, parser := filepath.Split(tmp[1])
		if isSupported(path) {
			fmt.Println("  * ", formatParserName(parser))
		}
	}
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

// isSupported checks whether the current OS and architecture are supported.
func isSupported(path string) bool {
	suppPath := filepath.Join(runtime.GOOS, runtime.GOARCH) + string(filepath.Separator)
	return path == suppPath
}

func formatParserName(fileName string) string {
	return strings.Replace(removeExt(fileName), "parser-", "", -1)
}

// removeExt removes the extension of a given file name.
func removeExt(fileName string) string {
	ext := filepath.Ext(fileName)
	return fileName[0 : len(fileName)-len(ext)]
}
