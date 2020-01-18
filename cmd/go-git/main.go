package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shivamdixit/go-git/git"
)

// List of commands supported by go-git
const (
	add        = "add"
	catFile    = "cat-file"
	checkout   = "checkout"
	commit     = "commit"
	hashObject = "hash-object"
	initialize = "init"
	log        = "log"
	lsTree     = "ls-tree"
	merge      = "merge"
	rebase     = "rebase"
	revParse   = "rev-parse"
	rm         = "rm"
	showRef    = "show-ref"
	tag        = "tag"
)

var initCmd *flag.FlagSet

func init() {
	// checkoutCmd := flag.NewFlagSet(checkout, flag.ExitOnError)
	// commitCmd := flag.NewFlagSet(commit, flag.ExitOnError)
	// hashObjectCmd := flag.NewFlagSet(hashObject, flag.ExitOnError)
	initCmd = flag.NewFlagSet(initialize, flag.ExitOnError)
	// logCmd := flag.NewFlagSet(log, flag.ExitOnError)
	// lsTreeCmd := flag.NewFlagSet(lsTree, flag.ExitOnError)
	// mergeCmd := flag.NewFlagSet(merge, flag.ExitOnError)
	// rebaseCmd := flag.NewFlagSet(rebase, flag.ExitOnError)
	// revParseCmd := flag.NewFlagSet(revParse, flag.ExitOnError)
	// rmCmd := flag.NewFlagSet(rm, flag.ExitOnError)
	// showRefCmd := flag.NewFlagSet(showRef, flag.ExitOnError)
	// tagCmd := flag.NewFlagSet(tag, flag.ExitOnError)
}

func initExec(path string) {
	err := git.Create(path)
	if err != nil {
		fmt.Printf("failed: %s", err)
	}
}

func main() {
	// providing a subcommand is must
	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch os.Args[1] {
	case initialize:
		initCmd.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if initCmd.Parsed() {
		// path is optional positional argument
		var dir string
		if len(os.Args) < 3 {
			d, err := os.Getwd()
			if err != nil {
				fmt.Errorf("cannot determine current directory %s", err)
				os.Exit(1)
			}
			dir = d
		} else {
			dir = os.Args[2]
		}

		initExec(dir)
	}
}
