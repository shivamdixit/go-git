package object

import "github.com/shivamdixit/go-git/git"

type Commit baseObject

func (c Commit) Name() string {
	return TypeCommit
}

func (c Commit) repository() *git.Repository {
	return c.repo
}

func (c Commit) Serialize() ([]byte, error) {
	return nil, nil
}

func (c Commit) Deserialize(data []byte) error {
	return nil
}

func NewCommit(data []byte, r *git.Repository) *Commit {
	return &Commit{data: data, repo: r}
}
