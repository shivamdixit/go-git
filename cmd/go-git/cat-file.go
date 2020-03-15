package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/shivamdixit/go-git/git"
	"github.com/shivamdixit/go-git/git/object"
)

var catFileType bool
var catFileSize bool
var catFilePretty bool

func catFileInit(catFileCmd *flag.FlagSet) {
	catFileCmd.BoolVar(&catFileType, "type", false, "show object type")
	catFileCmd.BoolVar(&catFileSize, "size", false, "show object size")
	catFileCmd.BoolVar(&catFilePretty, "pretty", false, "pretty print object's contents")

	catFileCmd.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: git cat-file (-type | -size | -pretty | <type>) <object>\n"+
			"<type> can be one of: blob, tree, commit, tag\n\n")
		catFileCmd.PrintDefaults()
	}
}

func catFileExec(catFileCmd *flag.FlagSet) {
	c := 0
	if catFileType {
		c += 1
	}
	if catFileSize {
		c += 1
	}
	if catFilePretty {
		c += 1
	}

	// exactly one option must be specified
	if c != 1 {
		catFileCmd.Usage()
		os.Exit(1)
	}

	// object SHA must be present
	if catFileCmd.NArg() != 1 {
		catFileCmd.Usage()
		os.Exit(1)
	}

	// fetch the repository from current directory
	r, err := git.Find(".")
	if err != nil {
		log.Fatal(err)
	}

	obj, err := object.Read(r, catFileCmd.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	if catFileType {
		fmt.Print(obj.Name())
	} else if catFileSize {
		// TODO: implement Size() method for objects
		//fmt.Print(obj.Size())
	} else {
		data, err := obj.Serialize()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(data))
	}
}
