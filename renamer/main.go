package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var re = regexp.MustCompile(`^(.+) ([0-9]{4}) \(([0-9]+) of ([0-9]+)\)\.(.+)$`)
var template = "$2 – $1 - $3 of $4.$5"

func main() {
	var dry bool
	flag.BoolVar(&dry, "dry", true, "whether or not this should be a real or dry run")
	flag.Parse()

	dir := "fixtures"
	var toRename []string

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if _, err := match(info.Name()); err == nil {
			toRename = append(toRename, path)
		}
		return nil
	})

	for _, oldPath := range toRename {
		dir := filepath.Dir(oldPath)
		fileName := filepath.Base(oldPath)
		newFileName, err := match(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
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

// match returns the new file name, or an error if the file name
// didn't match our pattern.
func match(fileName string) (string, error) {
	if !re.MatchString(fileName) {
		return "", fmt.Errorf("%s didn't match our pattern", fileName)
	}
	return re.ReplaceAllString(fileName, template), nil
}
