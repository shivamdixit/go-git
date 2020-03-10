package git

type Commit baseObject

func (c Commit) name() string {
	return ObjectCommit
}

func (c Commit) repository() *Repository {
	return c.repo
}

func (c Commit) Serialize() ([]byte, error) {
	return nil, nil
}

func (c Commit) Deserialize(data []byte) error {
	return nil
}

func NewCommit(data []byte) Commit {
	return Commit{}
}
