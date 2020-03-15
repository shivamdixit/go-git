package object

import "github.com/shivamdixit/go-git/git"

type Tree baseObject

func (t Tree) Name() string {
	return TypeTree
}

func (t Tree) repository() *git.Repository {
	return t.repo
}

func (t Tree) Serialize() ([]byte, error) {
	return nil, nil
}

func (t Tree) Deserialize(data []byte) error {
	return nil
}

func NewTree(data []byte) *Tree {
	return &Tree{}
}
