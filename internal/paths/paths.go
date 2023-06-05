package paths

// Path definitions

type Path interface {
	FullPath() string
	Exists() (bool, error)
	Delete() error
	Base() string
}

// Match definitions

type MatchType int

const (
	Mismatch MatchType = iota
	Match
)

func (m MatchType) String() string {
	switch m {
	case Match:
		return "Match"
	default:
		return "Mismatch"
	}
}
