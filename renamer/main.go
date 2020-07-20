package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type file struct {
	name, path string
}

func main() {
	dir := "fixtures"
	var toRename []file
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path, info.IsDir())
		if info.IsDir() {
			return nil
		}
		if _, err := match(info.Name()); err == nil {
			toRename = append(toRename, file{
				name: info.Name(),
				path: path,
			})
		}
		return nil
	})

	for _, orig := range toRename {
		var nf file
		var err error
		nf.name, err = match(orig.name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "match %s: %v", orig.path, err.Error())
		}
		nf.path = filepath.Join(orig.path, nf.name)
		fmt.Printf("mv %s => %s\n", orig.path, nf.path)
		err = os.Rename(orig.path, nf.path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "rename %s: %v", orig.path, err)
		}
	}
}

// match returns the new file name, or an error if the file name
// didn't match our pattern.
func match(fileName string) (string, error) {
	pieces := strings.Split(fileName, ".")
	ext := pieces[len(pieces)-1]
	base := strings.Join(pieces[:len(pieces)-1], ".")
	pieces = strings.Split(base, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", fmt.Errorf("%s didn't match our pattern", fileName)
	}
	return fmt.Sprintf("%s - %d.%s", strings.Title(name), number, ext), nil
}
