package acdfs

const (
	FILE   = "FILE"
	FOLDER = "FOLDER"
)

type MetadataPage struct {
	Count     int
	NextToken string
	Data      []Metadata
}

type Metadata struct {
	Id       string
	ParentId string
	Name     string
	Kind     string // should match FILE or FOLDER const
	Version  uint
	Status   string
	Parents  []string
	Raw      string

	ContentProperties struct {
		Version     uint64
		Md5         string
		Size        uint64
		ContentType string
		Extension   string
	}
}
