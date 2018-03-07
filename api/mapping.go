package api

type MappingAttribute struct {
	Field string `json:"field"`
	Kind  string `json:"kind"`
}

type StreamMapping struct {
	Attributes []MappingAttribute `json:"attributes"`
}
