package acdfs

type MetadataList struct {
	Data []Metadata
}

type Metadata struct {
	Id      string
	Name    string
	Kind    string
	Version uint
	Status  string
	Parents []string
}
