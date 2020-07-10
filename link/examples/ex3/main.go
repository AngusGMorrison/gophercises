package main

import (
	"fmt"
	"os"

	"github.com/angusgmorrison/gophercises/link"
)

func main() {
	f, err := os.Open("fixtures/ex3.html")
	if err != nil {
		exit(err.Error())
	}
	defer f.Close()
	links, err := link.Parse(f)
	if err != nil {
		exit(err.Error())
	}
	fmt.Printf("%+v\n", links)
}

func exit(msg string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], msg)
	os.Exit(1)
}
