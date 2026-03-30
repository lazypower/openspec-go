package validator

import (
	"fmt"

	"github.com/chuck/openspec-go/internal/model"
)

// ValidateStrict performs additional strict-mode checks across changes.
func ValidateStrict(changes []model.Change, specs map[string]model.Spec) []model.Issue {
	var issues []model.Issue

	// Cross-spec delta conflict detection: two changes modifying same requirement in same spec
	type reqKey struct{ spec, req string }
	modifiers := make(map[reqKey][]string) // key → list of change IDs

	for _, ch := range changes {
		for _, d := range ch.Deltas {
			if d.Operation == model.DeltaModified || d.Operation == model.DeltaRemoved {
				for _, r := range d.Requirements {
					k := reqKey{d.SpecName, r.Name}
					modifiers[k] = append(modifiers[k], ch.ID)
				}
			}
		}
	}
	for k, ids := range modifiers {
		if len(ids) > 1 {
			issues = append(issues, model.Issue{
				Level:   model.LevelWarning,
				Path:    k.spec,
				Message: fmt.Sprintf("Requirement %q modified by multiple changes: %v", k.req, ids),
			})
		}
	}

	// Well-formed RENAMED pairs
	for _, ch := range changes {
		for _, d := range ch.Deltas {
			if d.Operation == model.DeltaRenamed {
				if d.FromName == "" || d.ToName == "" {
					issues = append(issues, model.Issue{
						Level:   model.LevelError,
						Path:    ch.ID,
						Message: "RENAMED must specify both FROM and TO",
					})
				}
			}
		}
	}

	// MODIFIED completeness check
	for _, ch := range changes {
		for _, d := range ch.Deltas {
			if d.Operation != model.DeltaModified {
				continue
			}
			spec, ok := specs[d.SpecName]
			if !ok {
				continue
			}
			for _, deltaReq := range d.Requirements {
				for _, specReq := range spec.Requirements {
					if specReq.Name == deltaReq.Name {
						if len(deltaReq.Scenarios) < len(specReq.Scenarios) {
							issues = append(issues, model.Issue{
								Level:   model.LevelWarning,
								Path:    ch.ID,
								Message: fmt.Sprintf("MODIFIED requirement %q has fewer scenarios (%d) than original (%d)", deltaReq.Name, len(deltaReq.Scenarios), len(specReq.Scenarios)),
							})
						}
					}
				}
			}
		}
	}

	return issues
}
