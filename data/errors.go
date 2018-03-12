package data

// ParseError represents a parsing error e.g when parsing json
type ParseError struct {
	Base error
}

func (e ParseError) Error() string {
	return e.Base.Error()
}

// InvalidIDError represents an error where id is not a valid hex
type InvalidIDError struct {
	Base error
}

func (e InvalidIDError) Error() string {
	return e.Base.Error()
}

// ESError represents a parsing error e.g when parsing json
type ESError struct {
	Base error
}

func (e ESError) Error() string {
	return e.Base.Error()
}
