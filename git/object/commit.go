package object

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/shivamdixit/go-git/git"
)

type Commit struct {
	baseObject
	headers map[string]string
}

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
	c.headers = make(map[string]string)

	r := bufio.NewReader(bytes.NewReader(data))
	// Read every line to identify key-value pairs
	for {
		line, err := r.ReadString(byte('\n'))
		// we reached EOF without commit message
		if err == io.EOF {
			return errors.New("invalid commit object")
		}
		if err != nil {
			return err
		}

		// blank line implies that the remainder message is commit
		if line == "\n" {
			msg, err := ioutil.ReadAll(r)
			if err != nil {
				return err
			}

			c.headers["body"] = string(msg)
			return nil
		}

		// if the next line starts with space, it implies
		// that the previous line is in continuation
		for next, _ := r.Peek(1); next[0] == byte(' '); {
			t, err := r.ReadString(byte('\n'))

			// we reached EOF without commit message
			if err == io.EOF {
				return errors.New("invalid commit object")
			}
			if err != nil {
				return err
			}

			line = line + t
			next, _ = r.Peek(1)
		}

		// Split the line to identify the key
		pair := strings.SplitN(line, " ", 2)
		c.headers[pair[0]] = pair[1]
	}

	return nil
}

func NewCommit(data []byte, r *git.Repository) *Commit {
	c := Commit{}

	c.data = data
	c.repo = r
	return &c
}
