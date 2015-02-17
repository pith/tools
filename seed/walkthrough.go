// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"io/ioutil"
	"path/filepath"
	"fmt"
)

func walkthroughDir(transPath string, dirPath string) []string {
	files, _ := ioutil.ReadDir(dirPath)	
	var res []string
	for _, f := range files {
		fpath := filepath.Join(dirPath, f.Name())
		if f.IsDir() {
			rres := walkthroughDir(transPath, fpath)
			for _, n := range rres {
				res = append(res, n)
			}
		} else if filepath.Clean(fpath) != filepath.Clean(transPath) {
			res = append(res, fpath)
		}
	}
	return res
}

func processFiles(files []string, transformations []Transformation) {
	first := true
	done := make(chan string, len(files))
	
	for _, f := range files {
		go func(filePath string) {
			origDat, data := processFile(filePath, transformations)
			if len(origDat) != len(data) {
				if first {
					fmt.Println("Apply transformations:")
					first = false
				}
				err := ioutil.WriteFile(filePath, data, 0644)
				if err != nil {
					fmt.Printf("Error writting file %s\n", filePath)
				}
				fmt.Printf("%s\n", filePath)
			}
			
			done <- "ok"
		}(f)
	}
	
	for _ = range files {
		<-done
	}
	fmt.Printf("Processed %v files\n", len(files))
}

func processFile(filePath string, transformations []Transformation) ([]byte, []byte) {
	var origDat []byte
	var data []byte
	for _, transf := range transformations {
		if checkFileName(filePath, transf) {
			if len(origDat) == 0 {
				dat, err := ioutil.ReadFile(filePath)
				if err != nil {
					fmt.Printf("Error reading file %s\n", filePath)
				}
				data = dat
				origDat = dat
			}
			
			if checkCondition(filePath, data, transf) {
				data = applyProcs(data, transf)
			}
		}
	}
	return origDat, data
}
