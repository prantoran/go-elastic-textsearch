package data

// Section represents a single law section
type Section struct {
	Detail string `json:"detail"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
}

// Ammendment represents a single law ammendment
type Ammendment struct {
	Ammendment string   `json:"ammendment"`
	Atags      []string `json:"atags,omitempty"`
	ID         int      `json:"id"`
}

// LawDocument encapsulates the generic structure of single law record
type LawDocument struct {
	CreatedAt   string       `json:"created_at,omitempty"`
	Sections    []Section    `json:"sections,omitempty"`
	Ammendments []Ammendment `json:"ammendments,omitempty"`
	Act         string       `json:"act,omitempty"`
	ID          string       `json:"id"`
	Preamble    []string     `json:"preamble,omitempty"`
	Title       string       `json:"title,omitempty"`
}
