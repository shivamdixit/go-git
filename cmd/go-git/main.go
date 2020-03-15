package main

import (
	"flag"
	"fmt"
	"os"
)

// List of commands supported by go-git
const (
	add        = "add"
	catFile    = "cat-file"
	checkout   = "checkout"
	commit     = "commit"
	hashObject = "hash-object"
	initialize = "init"
	gitLog     = "log" // to not confuse with log package
	lsTree     = "ls-tree"
	merge      = "merge"
	rebase     = "rebase"
	revParse   = "rev-parse"
	rm         = "rm"
	showRef    = "show-ref"
	tag        = "tag"
)

var initCmd *flag.FlagSet
var catFileCmd *flag.FlagSet
var hashObjectCmd *flag.FlagSet

var defaultMessage = `
These are common Git commands used in various situations:

start a working area (see also: git help tutorial)
   clone     Clone a repository into a new directory
   init      Create an empty Git repository or reinitialize an existing one

work on the current change (see also: git help everyday)
   add       Add file contents to the index
   mv        Move or rename a file, a directory, or a symlink
   restore   Restore working tree files
   rm        Remove files from the working tree and from the index

examine the history and state (see also: git help revisions)
   bisect    Use binary search to find the commit that introduced a bug
   diff      Show changes between commits, commit and working tree, etc
   grep      Print lines matching a pattern
   log       Show commit logs
   show      Show various types of objects
   status    Show the working tree status

grow, mark and tweak your common history
   branch    List, create, or delete branches
   commit    Record changes to the repository
   merge     Join two or more development histories together
   rebase    Reapply commits on top of another base tip
   reset     Reset current HEAD to the specified state
   switch    Switch branches
   tag       Create, list, delete or verify a tag object signed with GPG

collaborate (see also: git help workflows)
   fetch     Download objects and refs from another repository
   pull      Fetch from and integrate with another repository or a local branch
   push      Update remote refs along with associated objects

low level plumbing commands
   cat-file    Print contents of an object
   hash-object Create a git object from a file
`

func init() {
	catFileCmd = flag.NewFlagSet(catFile, flag.ExitOnError)
	// checkoutCmd := flag.NewFlagSet(checkout, flag.ExitOnError)
	// commitCmd := flag.NewFlagSet(commit, flag.ExitOnError)
	hashObjectCmd = flag.NewFlagSet(hashObject, flag.ExitOnError)
	initCmd = flag.NewFlagSet(initialize, flag.ExitOnError)
	// logCmd := flag.NewFlagSet(log, flag.ExitOnError)
	// lsTreeCmd := flag.NewFlagSet(lsTree, flag.ExitOnError)
	// mergeCmd := flag.NewFlagSet(merge, flag.ExitOnError)
	// rebaseCmd := flag.NewFlagSet(rebase, flag.ExitOnError)
	// revParseCmd := flag.NewFlagSet(revParse, flag.ExitOnError)
	// rmCmd := flag.NewFlagSet(rm, flag.ExitOnError)
	// showRefCmd := flag.NewFlagSet(showRef, flag.ExitOnError)
	// tagCmd := flag.NewFlagSet(tag, flag.ExitOnError)

	// Define the default usage for the git command
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s: missing subcommand.\n%s", os.Args[0], defaultMessage)

		hashObjectCmd.PrintDefaults()
	}

	// Initialize all the subcommands
	initCmdInit(initCmd)
	hashObjectInit(hashObjectCmd)
	catFileInit(catFileCmd)
}

func main() {
	// providing a subcommand is must
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case initialize:
		initCmd.Parse(os.Args[2:])
		break
	case catFile:
		catFileCmd.Parse(os.Args[2:])
		break
	case hashObject:
		hashObjectCmd.Parse(os.Args[2:])
		break
	default:
		flag.Usage()
		os.Exit(1)
	}

	if initCmd.Parsed() {
		initCmdExec(initCmd.Args())
	}

	if catFileCmd.Parsed() {
		if len(catFileCmd.Args()) < 1 {
			catFileCmd.Usage()
			os.Exit(1)
		}

		catFileExec(catFileCmd)
	}

	if hashObjectCmd.Parsed() {
		if len(hashObjectCmd.Args()) < 1 {
			hashObjectCmd.Usage()
			os.Exit(1)
		}

		hashObjectExec(hashObjectCmd.Args())
	}
}
