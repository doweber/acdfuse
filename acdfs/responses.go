package acdfs

type MetadataList struct {
	Count     int
	NextToken string
	Data      []Metadata
}

type Metadata struct {
	Id       string
	ParentId string
	Name     string
	Kind     string
	Version  uint
	Status   string
	Parents  []string
}
