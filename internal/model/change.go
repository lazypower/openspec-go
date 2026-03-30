package model

// Change represents a parsed change proposal.
type Change struct {
	ID          string
	Title       string
	Why         string
	WhatChanges string
	Impact      string
	Deltas      []Delta
	Tasks       TaskStatus
}

// DeltaOp represents the type of delta operation.
type DeltaOp string

const (
	DeltaAdded    DeltaOp = "ADDED"
	DeltaModified DeltaOp = "MODIFIED"
	DeltaRemoved  DeltaOp = "REMOVED"
	DeltaRenamed  DeltaOp = "RENAMED"
)

// Delta represents a set of requirement changes targeting a spec.
type Delta struct {
	Operation    DeltaOp
	SpecName     string // target spec (derived from directory name)
	Requirements []Requirement
	// For RENAMED operations
	FromName string
	ToName   string
}

// TaskStatus tracks task completion from tasks.md.
type TaskStatus struct {
	Total     int
	Completed int
}

// Percent returns completion percentage (0-100).
func (t TaskStatus) Percent() int {
	if t.Total == 0 {
		return 0
	}
	return t.Completed * 100 / t.Total
}
