package model

// IssueLevel represents the severity of a validation issue.
type IssueLevel string

const (
	LevelError   IssueLevel = "ERROR"
	LevelWarning IssueLevel = "WARNING"
	LevelInfo    IssueLevel = "INFO"
)

// Issue represents a single validation finding.
type Issue struct {
	Level   IssueLevel `json:"level"`
	Path    string     `json:"path"`
	Message string     `json:"message"`
	Line    int        `json:"line,omitempty"`
}

// ValidationItem represents the validation result for a single item (change or spec).
type ValidationItem struct {
	ID     string     `json:"id"`
	Type   string     `json:"type"` // "change" or "spec"
	Valid  bool       `json:"valid"`
	Issues []Issue    `json:"issues"`
}

// ValidationReport is the structured output from a validation run.
type ValidationReport struct {
	Items   []ValidationItem  `json:"items"`
	Summary ValidationSummary `json:"summary"`
}

// ValidationSummary provides aggregate counts.
type ValidationSummary struct {
	Totals SummaryTotals     `json:"totals"`
	ByType map[string]int    `json:"byType"`
}

// SummaryTotals contains total/passed/failed counts.
type SummaryTotals struct {
	Items  int `json:"items"`
	Passed int `json:"passed"`
	Failed int `json:"failed"`
}

// HasErrors returns true if any item has ERROR-level issues.
func (r *ValidationReport) HasErrors() bool {
	for _, item := range r.Items {
		if !item.Valid {
			return true
		}
	}
	return false
}
