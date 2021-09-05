package main

import (
	"flag"
	"fmt"

	"github.com/kfelter/sea/cmd/sea/docgen"
)

var root = flag.String("root", ".", "root of tree to scan")

func main() {
	flag.Parse()
	if err := docgen.GenDoc(*root); err != nil {
		fmt.Printf("err generating doc: %v\n", err)
	}
}
