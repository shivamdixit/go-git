package git

type Tag struct {
	repo *Repository
	data string
}

func (t Tag) name() string {
	return ObjectTag
}

func (t Tag) repository() *Repository {
	return t.repo
}

func (t Tag) Serialize() ([]byte, error) {
	return nil, nil
}

func (t Tag) Deserialize(data []byte) error {
	return nil
}

func NewTag(data []byte) *Tag {
	return &Tag{}
}
