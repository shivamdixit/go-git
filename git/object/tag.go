package object

import "github.com/shivamdixit/go-git/git"

type Tag baseObject

func (t Tag) Name() string {
	return TypeTag
}

func (t Tag) repository() *git.Repository {
	return t.repo
}

func (t Tag) Serialize() ([]byte, error) {
	return nil, nil
}

func (t Tag) Deserialize(data []byte) error {
	return nil
}

func NewTag(data []byte, r *git.Repository) *Tag {
	return &Tag{data: data, repo: r}
}
