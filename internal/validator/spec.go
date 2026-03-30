package validator

import (
	"fmt"
	"strings"

	"github.com/chuck/openspec-go/internal/model"
)

// ValidateSpec validates a parsed spec.
func ValidateSpec(spec model.Spec, path string) []model.Issue {
	var issues []model.Issue

	// Must have Purpose section
	_, hasPurpose := spec.Sections["Purpose"]
	if !hasPurpose && spec.Overview == "" {
		issues = append(issues, model.Issue{
			Level:   model.LevelError,
			Path:    path,
			Message: "Missing required section: Purpose",
		})
	}

	// Purpose minimum length
	if spec.Overview != "" && len(spec.Overview) < MinPurposeLength {
		issues = append(issues, model.Issue{
			Level:   model.LevelError,
			Path:    path,
			Message: "Purpose section too short",
		})
	}

	// Must have Requirements section
	if len(spec.Requirements) == 0 {
		// Check if the section header exists but has no requirements
		_, hasReqSection := spec.Sections["Requirements"]
		if !hasReqSection {
			issues = append(issues, model.Issue{
				Level:   model.LevelError,
				Path:    path,
				Message: "Missing required section: Requirements",
			})
		}
	}

	// Validate each requirement
	for _, req := range spec.Requirements {
		if len(req.Scenarios) == 0 {
			issues = append(issues, model.Issue{
				Level:   model.LevelError,
				Path:    path,
				Message: fmt.Sprintf("Requirement %q must have at least one scenario", req.Name),
			})
		}
		if !strings.Contains(strings.ToUpper(req.Text), "SHALL") && !strings.Contains(strings.ToUpper(req.Text), "MUST") {
			issues = append(issues, model.Issue{
				Level:   model.LevelError,
				Path:    path,
				Message: fmt.Sprintf("Requirement %q must contain SHALL or MUST keyword", req.Name),
			})
		}
		if len(req.Text) > MaxRequirementTextLength {
			issues = append(issues, model.Issue{
				Level:   model.LevelWarning,
				Path:    path,
				Message: fmt.Sprintf("Requirement %q text exceeds %d characters", req.Name, MaxRequirementTextLength),
			})
		}
	}

	return issues
}
