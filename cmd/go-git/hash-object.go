package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/shivamdixit/go-git/git"
	"github.com/shivamdixit/go-git/git/object"
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
	var repo *git.Repository
	if hashObjectWrite == true {
		r, err := git.Find(".")
		if err != nil {
			log.Fatal(err)
		}
		repo = r
	}

	data, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}

	var o object.Object
	switch hashObjectType {
	case object.TypeBlob:
		o = object.NewBlob(data, repo)
		break
	case object.TypeTree:
		o = object.NewTree(data, repo)
		break
	case object.TypeCommit:
		o = object.NewCommit(data, repo)
		break
	case object.TypeTag:
		o = object.NewTag(data, repo)
		break
	default:
		log.Fatalf("invalid object type: %s", hashObjectType)
	}

	_, hash, err := object.Hash(o)
	if err != nil {
		log.Fatal(err)
	}

	if hashObjectWrite {
		object.Write(o)
	} else {
		fmt.Print(hash)
	}
}
