package git

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

// Object interface represents a generic git object
type Object interface {
	name() string
	repository() *Repository
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

// Different types of objects supported by git
const (
	ObjectCommit = "commit"
	ObjectTree   = "tree"
	ObjectTag    = "tag"
	ObjectBlob   = "blob"
)

// Read reads a given git object in a repository
func Read(r *Repository, sha string) (Object, error) {
	// get the path of the object. First two bytes of the hash
	// are used to identify the directory, remaining are used as a file name.
	path, err := r.file(filepath.Join("objects", sha[:2], sha[2:]), false)
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
	case ObjectCommit:
		o = NewCommit(rawBytes[j:])
	case ObjectTree:
		o = NewTree(rawBytes[j:])
	case ObjectTag:
		o = NewTag(rawBytes[j:])
	case ObjectBlob:
		o = NewBlob(rawBytes[j:])
	default:
		return nil, fmt.Errorf("invalid object type :%s", objType)
	}

	return o, nil
}

// FindObj finds a git object as an object can be referenced by
// full hash, short hash, tags, etc.
func FindObj(r *Repository, name string) string {
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
	length := new(bytes.Buffer)
	err = binary.Write(length, binary.LittleEndian, len(data))
	if err != nil {
		return nil, err
	}

	// create headers of the object. Sample header:
	//
	// commit 1086.tree 29ff16c9c14e265 2b22f8b78bb08a5a
	// <type> <len><0x0><contents>
	result := append([]byte(o.name()), []byte{' '}...)
	result = append(result, length.Bytes()...)
	result = append(result, []byte{0x0}...)
	result = append(result, data...)

	return result, nil
}

// getHash returns SHA1 hash of a byte slice
func getHash(o []byte) (string, error) {
	sha := sha1.Sum(o)
	return hex.EncodeToString(sha[:]), nil
}

func Write(o Object) error {
	// calculate hash of the object and use it as path of the object
	r, err := raw(o)
	if err != nil {
		return err
	}

	hash, err := getHash(r)
	if err != nil {
		return err
	}

	p, err := o.repository().file(filepath.Join("objects", hash[:2], hash[2:]), true)
	if err != nil {
		return err
	}

	// compress the raw object and write it out to file
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer w.Close()

	w.Write(r)
	err = ioutil.WriteFile(p, b.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
