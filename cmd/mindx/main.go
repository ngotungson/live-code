package main

import (
	"fmt"
	"mindx/internal/app/mindx"
	"os"
)

var (
	name    = "App"
	version = "0.0.0"
)

func main() {
	if err := mindx.Run(name, version); err != nil {
		fmt.Print("Run application ...")
		os.Exit(1)
	}
}
