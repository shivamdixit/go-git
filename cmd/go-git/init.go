package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/shivamdixit/go-git/git"
)

var initCmdQuiet bool

func initCmdInit(initCmd *flag.FlagSet) {
	initCmd.Bool("quiet", false, "be quiet")

	// Define custom usage function for each subcommand
	initCmd.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: git init [<directory>]\n\n")
		initCmd.PrintDefaults()
	}
}

func initCmdExec(args []string) {
	// path is optional
	var dir string
	if len(args) == 0 {
		d, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		dir = d
	} else {
		dir = args[0]
	}

	err := git.Create(dir)
	if err != nil {
		fmt.Printf("failed: %s", err)
	}
}
