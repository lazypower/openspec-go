package validator

import (
	"fmt"
	"strings"

	"github.com/chuck/openspec-go/internal/model"
)

// ValidateChange validates a change proposal and its deltas.
func ValidateChange(change model.Change) []model.Issue {
	var issues []model.Issue

	// Must have proposal sections
	if change.Why == "" {
		issues = append(issues, model.Issue{
			Level:   model.LevelError,
			Path:    change.ID + "/proposal.md",
			Message: "Missing required section: Why",
		})
	}
	if change.WhatChanges == "" {
		issues = append(issues, model.Issue{
			Level:   model.LevelError,
			Path:    change.ID + "/proposal.md",
			Message: "Missing required section: What Changes",
		})
	}

	// Why section length checks
	if change.Why != "" {
		if len(change.Why) < MinWhySectionLength {
			issues = append(issues, model.Issue{
				Level:   model.LevelWarning,
				Path:    change.ID + "/proposal.md",
				Message: "Why section is too short",
			})
		}
		if len(change.Why) > MaxWhySectionLength {
			issues = append(issues, model.Issue{
				Level:   model.LevelWarning,
				Path:    change.ID + "/proposal.md",
				Message: "Why section is too long",
			})
		}
	}

	// Must have at least one delta
	if len(change.Deltas) == 0 {
		issues = append(issues, model.Issue{
			Level:   model.LevelError,
			Path:    change.ID,
			Message: "Change must have at least one delta",
		})
	}

	// Delta count warning
	if len(change.Deltas) > MaxDeltasPerChange {
		issues = append(issues, model.Issue{
			Level:   model.LevelWarning,
			Path:    change.ID,
			Message: fmt.Sprintf("Change has %d deltas (threshold: %d)", len(change.Deltas), MaxDeltasPerChange),
		})
	}

	// Validate each delta
	seen := make(map[string]map[string]bool) // operation -> req name -> seen
	for _, delta := range change.Deltas {
		opKey := string(delta.Operation)
		if seen[opKey] == nil {
			seen[opKey] = make(map[string]bool)
		}

		for _, req := range delta.Requirements {
			// Duplicate check within section
			if seen[opKey][req.Name] {
				issues = append(issues, model.Issue{
					Level:   model.LevelError,
					Path:    change.ID,
					Message: fmt.Sprintf("Duplicate requirement name: %s", req.Name),
				})
			}
			seen[opKey][req.Name] = true

			// ADDED/MODIFIED must have scenarios
			if delta.Operation == model.DeltaAdded || delta.Operation == model.DeltaModified {
				if len(req.Scenarios) == 0 {
					issues = append(issues, model.Issue{
						Level:   model.LevelError,
						Path:    change.ID,
						Message: fmt.Sprintf("Requirement %q must have at least one scenario", req.Name),
					})
				}
				// Must use SHALL/MUST
				if !containsShallMust(req.Text) {
					issues = append(issues, model.Issue{
						Level:   model.LevelError,
						Path:    change.ID,
						Message: fmt.Sprintf("Requirement %q must contain SHALL or MUST keyword", req.Name),
					})
				}
			}
		}
	}

	return issues
}

func containsShallMust(text string) bool {
	upper := strings.ToUpper(text)
	return strings.Contains(upper, "SHALL") || strings.Contains(upper, "MUST")
}
