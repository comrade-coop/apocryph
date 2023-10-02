package main

import (
	"os"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		println(err)
		os.Exit(1)
	}
}
