package git

type Tree struct {
	repo *Repository
	data string
}

func (t Tree) name() string {
	return ObjectTree
}

func (t Tree) repository() *Repository {
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
