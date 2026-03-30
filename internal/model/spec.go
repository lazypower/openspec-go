package model

// Spec represents a parsed specification file.
type Spec struct {
	Title        string
	Overview     string // Purpose section content
	Requirements []Requirement
	Sections     map[string]string // raw section name → content
}

// Requirement represents a single requirement within a spec.
type Requirement struct {
	Name      string
	Text      string // body text
	Scenarios []Scenario
}

// Scenario represents a test scenario within a requirement.
type Scenario struct {
	Name string
	Text string // raw bullet lines
}
