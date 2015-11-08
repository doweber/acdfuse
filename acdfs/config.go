package acdfs

type Config struct {
	ContentUrl     string `json:"contentUrl"`
	CustomerExists bool   `json:"customerExists"`
	MetadataUrl    string `json:"metadataUrl"`
}
