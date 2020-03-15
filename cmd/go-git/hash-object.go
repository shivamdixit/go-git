package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/shivamdixit/go-git/git"
)

var hashObjectType string
var hashObjectWrite bool

func hashObjectInit(hashObjectCmd *flag.FlagSet) {
	hashObjectCmd.StringVar(&hashObjectType, "type", "blob", "specify the object type")
	hashObjectCmd.BoolVar(&hashObjectWrite, "write", false, "actually write the object into the object database")

	hashObjectCmd.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: git hash-object [-type <type>] [-write] <file>...\n\n")
		hashObjectCmd.PrintDefaults()
	}
}

func hashObjectExec(args []string) {
	if hashObjectWrite == true {
		_, err := git.Find(".")
		if err != nil {
			log.Fatal(err)
		}
	}

	data, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}

	var o git.Object
	switch hashObjectType {
	case git.ObjectBlob:
		o = git.NewBlob(data)
		break
	case git.ObjectTree:
		o = git.NewTree(data)
		break
	case git.ObjectCommit:
		o = git.NewCommit(data)
		break
	case git.ObjectTag:
		o = git.NewTag(data)
		break
	default:
		log.Fatalf("invalid object type: %s", hashObjectType)
	}

	_, hash, err := git.Hash(o)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(hash)
}
