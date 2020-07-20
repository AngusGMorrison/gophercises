package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	dir := "fixtures"
	toRename := make(map[string][]string)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		currentDir := filepath.Dir(path)
		if _, err := match(info.Name()); err == nil {
			toRename[currentDir] = append(toRename[currentDir], info.Name())
		}
		return nil
	})

	for dir, files := range toRename {
		n := len(files)
		sort.Strings(files)
		for i, fileName := range files {
			m, _ := match(fileName)
			newFileName := fmt.Sprintf("%s - %d of %d.%s", m.base, (i + 1), n, m.ext)
			oldPath := filepath.Join(dir, fileName)
			newPath := filepath.Join(dir, newFileName)
			fmt.Printf("mv %s => %s\n", oldPath, newPath)
			err := os.Rename(oldPath, newPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "rename %s: %v\n", oldPath, err)
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
