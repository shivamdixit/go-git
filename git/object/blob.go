package object

import "github.com/shivamdixit/go-git/git"

type Blob baseObject

func (b Blob) Name() string {
	return TypeBlob
}

func (b Blob) repository() *git.Repository {
	return b.repo
}

func (b Blob) Serialize() ([]byte, error) {
	// blob is a raw data structure and it has no specific format
	// therefore, just return the raw byte data
	return b.data, nil
}

func (b Blob) Deserialize(data []byte) error {
	b.data = data

	return nil
}

func NewBlob(data []byte) *Blob {
	return &Blob{data: data}
}
