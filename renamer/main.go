package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var dry bool
	flag.BoolVar(&dry, "dry", true, "whether or not this should be a real or dry run")
	flag.Parse()

	dir := "fixtures"
	toRename := make(map[string][]string)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		currentDir := filepath.Dir(path)
		if mr, err := match(info.Name()); err == nil {
			key := filepath.Join(currentDir, fmt.Sprintf("%s.%s", mr.base, mr.ext))
			toRename[key] = append(toRename[key], info.Name())
		}
		return nil
	})

	for k, files := range toRename {
		dir := filepath.Dir(k)
		n := len(files)
		sort.Strings(files)
		for i, fileName := range files {
			m, _ := match(fileName)
			newFileName := fmt.Sprintf("%s - %d of %d.%s", m.base, (i + 1), n, m.ext)
			oldPath := filepath.Join(dir, fileName)
			newPath := filepath.Join(dir, newFileName)
			fmt.Printf("mv %s => %s\n", oldPath, newPath)
			if !dry {
				err := os.Rename(oldPath, newPath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "rename %s: %v\n", oldPath, err)
				}
			}
		}
	}
}

type matchResult struct {
	base, ext string
	index     int
}

// match returns the new file name, or an error if the file name
// didn't match our pattern.
func match(fileName string) (*matchResult, error) {
	pieces := strings.Split(fileName, ".")
	ext := pieces[len(pieces)-1]
	base := strings.Join(pieces[:len(pieces)-1], ".")
	pieces = strings.Split(base, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return nil, fmt.Errorf("%s didn't match our pattern", fileName)
	}
	return &matchResult{strings.Title(name), ext, number}, nil
}
