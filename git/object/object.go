package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/shivamdixit/go-git/git"
)

// baseObject is the type for all high level objects like
// commit, tag, blob, etc.
type baseObject struct {
	repo *git.Repository
	data []byte
}

// Object interface represents a generic high level git object
type Object interface {
	Name() string
	repository() *git.Repository
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

// Different types of objects supported by git
const (
	TypeCommit = "commit"
	TypeTree   = "tree"
	TypeTag    = "tag"
	TypeBlob   = "blob"
)

// Read reads a given git object in a repository
func Read(r *git.Repository, sha string) (Object, error) {
	// get the path of the object. First two bytes of the hash
	// are used to identify the directory, remaining are used as a file name.
	path, err := r.File(filepath.Join("objects", sha[:2], sha[2:]), false)
	if err != nil {
		return nil, err
	}

	// read the object file
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	reader, err := zlib.NewReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	var raw bytes.Buffer
	io.Copy(&raw, reader)

	rawBytes := raw.Bytes()
	// find the first delimiter char ' ' (space)
	i := bytes.Index(rawBytes, []byte{' '})
	objType := rawBytes[:i]

	// next delimiter is the null character
	j := bytes.Index(rawBytes, []byte{0x0})
	size := string(rawBytes[i+1 : j])

	// convert and validate the value of size
	s, err := strconv.Atoi(size)
	if err != nil {
		return nil, err
	}
	if s != (len(rawBytes) - j - 1) {
		return nil, fmt.Errorf("malformed object: %s, bad length", sha)
	}

	var o Object
	switch string(objType) {
	case TypeCommit:
		o = NewCommit(rawBytes[j:], r)
	case TypeTree:
		o = NewTree(rawBytes[j:], r)
	case TypeTag:
		o = NewTag(rawBytes[j:], r)
	case TypeBlob:
		o = NewBlob(rawBytes[j:], r)
	default:
		return nil, fmt.Errorf("invalid object type :%s", objType)
	}

	return o, nil
}

// FindObj finds a git object as an object can be referenced by
// full hash, short hash, tags, etc.
func FindObj(r *git.Repository, name string) string {
	// TODO: implement find logic
	return name
}

// raw returns the raw object with headers
func raw(o Object) ([]byte, error) {
	data, err := o.Serialize()
	if err != nil {
		return nil, err
	}

	// calculate length of data (content of the object) as byte slice
	length := []byte(strconv.Itoa(len(data)))

	// create headers of the object. Sample header:
	//
	// commit 1086.tree 29ff16c9c14e265 2b22f8b78bb08a5a
	// <type> <len><0x0><contents>
	result := append([]byte(o.Name()), []byte{' '}...)
	result = append(result, length...)
	result = append(result, []byte{0x0}...)
	result = append(result, data...)

	return result, nil
}

// Hash returns the "raw" object along with its hash
func Hash(o Object) ([]byte, string, error) {
	r, err := raw(o)
	if err != nil {
		return nil, "", err
	}

	sha := sha1.Sum(r)
	return r, hex.EncodeToString(sha[:]), nil
}

func Write(o Object) error {
	// create a raw object with all headers and create its hash
	// Hash is used as a path for the object
	raw, hash, err := Hash(o)
	if err != nil {
		return err
	}

	p, err := o.repository().File(filepath.Join("objects", hash[:2], hash[2:]), true)
	if err != nil {
		return err
	}

	// compress the raw object and write it out to file
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(raw)
	// must be closed to flush the buffer
	w.Close()

	err = ioutil.WriteFile(p, b.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
