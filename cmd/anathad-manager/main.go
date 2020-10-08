package main

import (
	"fmt"
	"github.com/anathatech/cosmosd/lib"
	"os"
)


func main() {
	err := lib.Run(os.Args[1:])
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}