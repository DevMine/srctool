// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/DevMine/srcanlzr/src"
	"github.com/codegangsta/cli"

	"github.com/DevMine/srctool/config"
	"github.com/DevMine/srctool/log"
)

// Parse command runs all installed parsers on a project, merges the resulting
// JSON and outputs the result in JSON to stdout.
// It expects only one command line argument: the directory of a project.
func Parse(ctx *cli.Context) {
	if !ctx.Args().Present() {
		log.Fatal("expected 1 argument, found 0")
	}

	projectPath := ctx.Args().First()
	parsersPath := filepath.Join(config.DataDir(), config.ParsersFolder)

	fis, err := ioutil.ReadDir(parsersPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(fis) == 0 {
		log.Fatal("no parser installed")
		return
	}

	totalWaits := 0
	c := make(chan []byte)

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

		totalWaits++
		go cmdRoutine(parsersPath, fi.Name(), projectPath, c)
	}

	var prjs []*src.Project

	for ; totalWaits > 0; totalWaits-- {
		select {
		case bs := <-c:
			prj, err := src.Unmarshal(bs)
			if err != nil {
				log.Fail(err)
			}

			prjs = append(prjs, prj)
		}
	}

	close(c)

	log.Info("merging JSON outputs")
	var prj *src.Project
	if prj, err = src.MergeAll(prjs...); err != nil {
		log.Debug(err)
		log.Fatal("failed to merge all JSON")
	}

	bs, err := json.Marshal(prj)
	if err != nil {
		log.Debug(err)
		log.Fatal("unable to marshal the final JSON")
	}

	fmt.Println(string(bs))

	log.Success("done parsing")
}

// cmdRoutine runs a language parser on a project.
func cmdRoutine(parsersPath, parserName, projectPath string, c chan []byte) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	if err := os.Chdir(filepath.Join(parsersPath, parserName)); err != nil {
		log.Fatal(err)
	}

	execPath := "./parser"

	cmd := exec.Command(execPath, projectPath)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	log.Debug("command: ", strings.Join(cmd.Args, " "))

	if err := cmd.Run(); err != nil {
		log.Debug("debug:", err)
		log.Fatal(fmt.Sprintf("failed to parse with the %s parser", parserName))
	}

	if errBuf.Len() > 0 {
		log.Info(fmt.Sprintf("parser %s errors:", parserName))
		log.Fail(errBuf.String())
	}

	if outBuf.Len() == 0 {
		log.Fatal(fmt.Sprintf("the %s parser did not produce any output", parserName))
	}

	c <- outBuf.Bytes()
}
