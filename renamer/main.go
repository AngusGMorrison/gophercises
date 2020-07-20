package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	// fileName := "birthday_001.txt"
	// // => Birthday - 1 of 4.txt
	// newName, err := match(fileName, 4)
	// if err != nil {
	// 	fmt.Println("no match")
	// 	os.Exit(1)
	// }
	// fmt.Println(newName)
	dir := "fixtures"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var count int
	toRename := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
		} else {
			_, err := match(file.Name(), 4)
			if err == nil {
				count++
				toRename = append(toRename, file.Name())
			}
		}
	}

	for _, orig := range toRename {
		newFileName, err := match(orig, count)
		if err != nil {
			panic(err)
		}
		origPath := filepath.Join(dir, orig)
		newPath := filepath.Join(dir, newFileName)
		fmt.Printf("mv %s => %s\n", origPath, newPath)
	}
}

// match returns the new file name, or an error if the file name
// didn't match our pattern.
func match(fileName string, total int) (string, error) {
	pieces := strings.Split(fileName, ".")
	ext := pieces[len(pieces)-1]
	base := strings.Join(pieces[:len(pieces)-1], ".")
	pieces = strings.Split(base, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", fmt.Errorf("%s didn't match our pattern", fileName)
	}
	return fmt.Sprintf("%s - %d of %d.%s", strings.Title(name), number, total, ext), nil
}
