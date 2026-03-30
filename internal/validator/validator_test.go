package validator

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chuck/openspec-go/internal/model"
)

func longString(n int) string {
	return strings.Repeat("a", n)
}

func validChange() model.Change {
	return model.Change{
		ID:          "test-change",
		Title:       "Test Change",
		Why:         longString(60),
		WhatChanges: "Some changes here",
		Impact:      "Some impact",
		Deltas: []model.Delta{
			{
				Operation: model.DeltaAdded,
				SpecName:  "test-spec",
				Requirements: []model.Requirement{
					{
						Name: "New Feature",
						Text: "The system SHALL do something.",
						Scenarios: []model.Scenario{
							{Name: "Happy path", Text: "- **WHEN** x\n- **THEN** y"},
						},
					},
				},
			},
		},
	}
}

// --- Change Validation Tests ---

func TestValidateChange_NoDeltasError(t *testing.T) {
	ch := validChange()
	ch.Deltas = nil
	issues := ValidateChange(ch)
	if !hasError(issues, "at least one delta") {
		t.Error("expected error about missing deltas")
	}
}

func TestValidateChange_NoScenariosError(t *testing.T) {
	ch := validChange()
	ch.Deltas[0].Requirements[0].Scenarios = nil
	issues := ValidateChange(ch)
	if !hasError(issues, "at least one scenario") {
		t.Error("expected error about missing scenarios")
	}
}

func TestValidateChange_NoShallMustError(t *testing.T) {
	ch := validChange()
	ch.Deltas[0].Requirements[0].Text = "The system does something."
	issues := ValidateChange(ch)
	if !hasError(issues, "SHALL or MUST") {
		t.Error("expected error about missing SHALL/MUST")
	}
}

func TestValidateChange_MissingWhyError(t *testing.T) {
	ch := validChange()
	ch.Why = ""
	issues := ValidateChange(ch)
	if !hasError(issues, "Missing required section: Why") {
		t.Error("expected error about missing Why section")
	}
}

func TestValidateChange_WhyTooShort(t *testing.T) {
	ch := validChange()
	ch.Why = "Short."
	issues := ValidateChange(ch)
	if !hasWarning(issues, "too short") {
		t.Error("expected warning about short Why section")
	}
}

func TestValidateChange_DuplicateRequirements(t *testing.T) {
	ch := validChange()
	ch.Deltas[0].Requirements = append(ch.Deltas[0].Requirements, model.Requirement{
		Name: "New Feature",
		Text: "The system SHALL do something else.",
		Scenarios: []model.Scenario{
			{Name: "Test", Text: "steps"},
		},
	})
	issues := ValidateChange(ch)
	if !hasError(issues, "Duplicate requirement") {
		t.Error("expected error about duplicate requirement")
	}
}

func TestValidateChange_ValidPasses(t *testing.T) {
	ch := validChange()
	issues := ValidateChange(ch)
	for _, iss := range issues {
		if iss.Level == model.LevelError {
			t.Errorf("unexpected error: %s", iss.Message)
		}
	}
}

// --- Spec Validation Tests ---

func TestValidateSpec_MissingPurposeError(t *testing.T) {
	spec := model.Spec{
		Title:    "Test",
		Sections: map[string]string{},
		Requirements: []model.Requirement{
			{Name: "R1", Text: "SHALL do x", Scenarios: []model.Scenario{{Name: "S1", Text: "steps"}}},
		},
	}
	issues := ValidateSpec(spec, "test.md")
	if !hasError(issues, "Purpose") {
		t.Error("expected error about missing Purpose")
	}
}

func TestValidateSpec_MissingRequirementsError(t *testing.T) {
	spec := model.Spec{
		Title:    "Test",
		Overview: longString(60),
		Sections: map[string]string{"Purpose": longString(60)},
	}
	issues := ValidateSpec(spec, "test.md")
	if !hasError(issues, "Requirements") {
		t.Error("expected error about missing Requirements")
	}
}

func TestValidateSpec_PurposeTooShortError(t *testing.T) {
	spec := model.Spec{
		Title:    "Test",
		Overview: "Short.",
		Sections: map[string]string{"Purpose": "Short."},
		Requirements: []model.Requirement{
			{Name: "R1", Text: "SHALL do x", Scenarios: []model.Scenario{{Name: "S1", Text: "steps"}}},
		},
	}
	issues := ValidateSpec(spec, "test.md")
	if !hasError(issues, "Purpose section too short") {
		t.Error("expected error about short Purpose")
	}
}

func TestValidateSpec_RequirementNoScenario(t *testing.T) {
	spec := model.Spec{
		Title:    "Test",
		Overview: longString(60),
		Sections: map[string]string{"Purpose": longString(60)},
		Requirements: []model.Requirement{
			{Name: "R1", Text: "SHALL do x"},
		},
	}
	issues := ValidateSpec(spec, "test.md")
	if !hasError(issues, "at least one scenario") {
		t.Error("expected error about missing scenario")
	}
}

func TestValidateSpec_ValidPasses(t *testing.T) {
	spec := model.Spec{
		Title:    "Test",
		Overview: longString(60),
		Sections: map[string]string{"Purpose": longString(60)},
		Requirements: []model.Requirement{
			{
				Name: "R1",
				Text: "The system SHALL do something.",
				Scenarios: []model.Scenario{
					{Name: "S1", Text: "steps"},
				},
			},
		},
	}
	issues := ValidateSpec(spec, "test.md")
	for _, iss := range issues {
		if iss.Level == model.LevelError {
			t.Errorf("unexpected error: %s", iss.Message)
		}
	}
}

// --- Strict Mode Tests ---

func TestValidateStrict_CrossSpecConflict(t *testing.T) {
	changes := []model.Change{
		{
			ID: "change-a",
			Deltas: []model.Delta{
				{Operation: model.DeltaModified, SpecName: "shared-spec", Requirements: []model.Requirement{{Name: "Feature X"}}},
			},
		},
		{
			ID: "change-b",
			Deltas: []model.Delta{
				{Operation: model.DeltaModified, SpecName: "shared-spec", Requirements: []model.Requirement{{Name: "Feature X"}}},
			},
		},
	}
	issues := ValidateStrict(changes, nil)
	if !hasWarning(issues, "multiple changes") {
		t.Error("expected warning about cross-spec conflict")
	}
}

func TestValidateStrict_MalformedRenamed(t *testing.T) {
	changes := []model.Change{
		{
			ID: "rename-change",
			Deltas: []model.Delta{
				{Operation: model.DeltaRenamed, FromName: "Old", ToName: ""},
			},
		},
	}
	issues := ValidateStrict(changes, nil)
	if !hasError(issues, "FROM and TO") {
		t.Error("expected error about malformed RENAMED")
	}
}

func TestValidateStrict_ModifiedCompleteness(t *testing.T) {
	specs := map[string]model.Spec{
		"my-spec": {
			Requirements: []model.Requirement{
				{Name: "Feature", Scenarios: []model.Scenario{{Name: "S1"}, {Name: "S2"}, {Name: "S3"}}},
			},
		},
	}
	changes := []model.Change{
		{
			ID: "mod-change",
			Deltas: []model.Delta{
				{
					Operation: model.DeltaModified,
					SpecName:  "my-spec",
					Requirements: []model.Requirement{
						{Name: "Feature", Scenarios: []model.Scenario{{Name: "S1"}}},
					},
				},
			},
		},
	}
	issues := ValidateStrict(changes, specs)
	if !hasWarning(issues, "fewer scenarios") {
		t.Error("expected warning about completeness")
	}
}

// --- Concurrent Validation Tests ---

func TestValidateConcurrent_BatchResults(t *testing.T) {
	items := []ValidationFunc{
		func() (model.ValidationItem, error) {
			return model.ValidationItem{ID: "b", Type: "change", Valid: true}, nil
		},
		func() (model.ValidationItem, error) {
			return model.ValidationItem{ID: "a", Type: "spec", Valid: false, Issues: []model.Issue{{Level: model.LevelError, Message: "bad"}}}, nil
		},
	}
	results := ValidateConcurrent(items, 2)
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	// Should be sorted by ID
	if results[0].ID != "a" || results[1].ID != "b" {
		t.Errorf("results not sorted: %v, %v", results[0].ID, results[1].ID)
	}
}

func TestValidateConcurrent_DeterministicOrder(t *testing.T) {
	var items []ValidationFunc
	for i := 0; i < 20; i++ {
		id := string(rune('a' + i%26))
		items = append(items, func() (model.ValidationItem, error) {
			return model.ValidationItem{ID: id, Type: "change", Valid: true}, nil
		})
	}
	results := ValidateConcurrent(items, 4)
	for i := 1; i < len(results); i++ {
		if results[i].ID < results[i-1].ID {
			t.Errorf("results not sorted at index %d: %s < %s", i, results[i].ID, results[i-1].ID)
		}
	}
}

func TestValidateConcurrent_SingleItemNoop(t *testing.T) {
	items := []ValidationFunc{
		func() (model.ValidationItem, error) {
			return model.ValidationItem{ID: "only", Type: "spec", Valid: true}, nil
		},
	}
	results := ValidateConcurrent(items, 1)
	if len(results) != 1 || results[0].ID != "only" {
		t.Errorf("unexpected result: %v", results)
	}
}

// helpers

func hasError(issues []model.Issue, substr string) bool {
	for _, iss := range issues {
		if iss.Level == model.LevelError && strings.Contains(iss.Message, substr) {
			return true
		}
	}
	return false
}

func hasWarning(issues []model.Issue, substr string) bool {
	for _, iss := range issues {
		if iss.Level == model.LevelWarning && strings.Contains(iss.Message, substr) {
			return true
		}
	}
	return false
}

// --- Boundary Condition Tests ---

func TestValidateChange_WhyExactlyMinLength(t *testing.T) {
	// 50 chars exactly — should NOT warn
	ch := validChange()
	ch.Why = strings.Repeat("a", 50)
	issues := ValidateChange(ch)
	if hasWarning(issues, "too short") {
		t.Error("50 chars should not trigger too-short warning")
	}
}

func TestValidateChange_WhyOneBelowMinLength(t *testing.T) {
	ch := validChange()
	ch.Why = strings.Repeat("a", 49)
	issues := ValidateChange(ch)
	if !hasWarning(issues, "too short") {
		t.Error("49 chars should trigger too-short warning")
	}
}

func TestValidateChange_WhyExactlyMaxLength(t *testing.T) {
	ch := validChange()
	ch.Why = strings.Repeat("a", 1000)
	issues := ValidateChange(ch)
	if hasWarning(issues, "too long") {
		t.Error("1000 chars should not trigger too-long warning")
	}
}

func TestValidateChange_WhyOneAboveMaxLength(t *testing.T) {
	ch := validChange()
	ch.Why = strings.Repeat("a", 1001)
	issues := ValidateChange(ch)
	if !hasWarning(issues, "too long") {
		t.Error("1001 chars should trigger too-long warning")
	}
}

func TestValidateSpec_PurposeExactlyMinLength(t *testing.T) {
	spec := model.Spec{
		Title:    "Test",
		Overview: strings.Repeat("a", 50),
		Sections: map[string]string{"Purpose": strings.Repeat("a", 50)},
		Requirements: []model.Requirement{
			{Name: "R1", Text: "The system SHALL do x.", Scenarios: []model.Scenario{{Name: "S1", Text: "steps"}}},
		},
	}
	issues := ValidateSpec(spec, "test.md")
	if hasError(issues, "Purpose section too short") {
		t.Error("50 chars should not trigger too-short error")
	}
}

func TestValidateSpec_PurposeOneBelowMinLength(t *testing.T) {
	spec := model.Spec{
		Title:    "Test",
		Overview: strings.Repeat("a", 49),
		Sections: map[string]string{"Purpose": strings.Repeat("a", 49)},
		Requirements: []model.Requirement{
			{Name: "R1", Text: "The system SHALL do x.", Scenarios: []model.Scenario{{Name: "S1", Text: "steps"}}},
		},
	}
	issues := ValidateSpec(spec, "test.md")
	if !hasError(issues, "Purpose section too short") {
		t.Error("49 chars should trigger too-short error")
	}
}

func TestValidateSpec_RequirementTextAtMaxLength(t *testing.T) {
	spec := model.Spec{
		Title:    "Test",
		Overview: longString(60),
		Sections: map[string]string{"Purpose": longString(60)},
		Requirements: []model.Requirement{
			{Name: "R1", Text: "SHALL " + strings.Repeat("a", 494), Scenarios: []model.Scenario{{Name: "S1", Text: "steps"}}},
		},
	}
	issues := ValidateSpec(spec, "test.md")
	if hasWarning(issues, "exceeds") {
		t.Error("500 chars should not trigger exceeds warning")
	}
}

func TestValidateSpec_RequirementTextAboveMaxLength(t *testing.T) {
	spec := model.Spec{
		Title:    "Test",
		Overview: longString(60),
		Sections: map[string]string{"Purpose": longString(60)},
		Requirements: []model.Requirement{
			{Name: "R1", Text: "SHALL " + strings.Repeat("a", 495), Scenarios: []model.Scenario{{Name: "S1", Text: "steps"}}},
		},
	}
	issues := ValidateSpec(spec, "test.md")
	if !hasWarning(issues, "exceeds") {
		t.Error("501 chars should trigger exceeds warning")
	}
}

func TestValidateChange_ExactlyMaxDeltas(t *testing.T) {
	ch := validChange()
	ch.Deltas = nil
	for i := 0; i < 10; i++ {
		ch.Deltas = append(ch.Deltas, model.Delta{
			Operation: model.DeltaAdded,
			SpecName:  fmt.Sprintf("spec-%d", i),
			Requirements: []model.Requirement{
				{Name: fmt.Sprintf("Req %d", i), Text: "SHALL do x.", Scenarios: []model.Scenario{{Name: "S", Text: "s"}}},
			},
		})
	}
	issues := ValidateChange(ch)
	if hasWarning(issues, "deltas") {
		t.Error("10 deltas should not trigger warning")
	}
}

func TestValidateChange_OneAboveMaxDeltas(t *testing.T) {
	ch := validChange()
	ch.Deltas = nil
	for i := 0; i < 11; i++ {
		ch.Deltas = append(ch.Deltas, model.Delta{
			Operation: model.DeltaAdded,
			SpecName:  fmt.Sprintf("spec-%d", i),
			Requirements: []model.Requirement{
				{Name: fmt.Sprintf("Req %d", i), Text: "SHALL do x.", Scenarios: []model.Scenario{{Name: "S", Text: "s"}}},
			},
		})
	}
	issues := ValidateChange(ch)
	if !hasWarning(issues, "deltas") {
		t.Error("11 deltas should trigger warning")
	}
}

func TestValidateChange_MissingWhatChanges(t *testing.T) {
	ch := validChange()
	ch.WhatChanges = ""
	issues := ValidateChange(ch)
	if !hasError(issues, "What Changes") {
		t.Error("expected error about missing What Changes section")
	}
}

func TestValidateChange_ShallCaseInsensitive(t *testing.T) {
	ch := validChange()
	ch.Deltas[0].Requirements[0].Text = "The system shall do something."
	issues := ValidateChange(ch)
	if hasError(issues, "SHALL or MUST") {
		t.Error("lowercase 'shall' should pass SHALL/MUST check")
	}
}

func TestValidateChange_MustKeyword(t *testing.T) {
	ch := validChange()
	ch.Deltas[0].Requirements[0].Text = "The system MUST do something."
	issues := ValidateChange(ch)
	if hasError(issues, "SHALL or MUST") {
		t.Error("MUST keyword should satisfy SHALL/MUST check")
	}
}

func TestValidateSpec_ShallInScenarioNotBody(t *testing.T) {
	// SHALL only in scenario text, not in requirement body
	spec := model.Spec{
		Title:    "Test",
		Overview: longString(60),
		Sections: map[string]string{"Purpose": longString(60)},
		Requirements: []model.Requirement{
			{
				Name: "R1",
				Text: "The system does something without keywords.",
				Scenarios: []model.Scenario{{Name: "S1", Text: "The system SHALL work"}},
			},
		},
	}
	issues := ValidateSpec(spec, "test.md")
	if !hasError(issues, "SHALL or MUST") {
		t.Error("SHALL in scenario but not body should still fail")
	}
}

func TestValidateChange_EmptyRequirementName(t *testing.T) {
	ch := validChange()
	ch.Deltas[0].Requirements[0].Name = ""
	issues := ValidateChange(ch)
	// Should still validate (empty name means duplicate detection key is "")
	_ = issues
}

func TestValidateStrict_NoChanges(t *testing.T) {
	issues := ValidateStrict(nil, nil)
	if len(issues) != 0 {
		t.Errorf("expected no issues for empty input, got %d", len(issues))
	}
}

func TestValidateStrict_RenamedBothEmpty(t *testing.T) {
	changes := []model.Change{
		{ID: "ch", Deltas: []model.Delta{{Operation: model.DeltaRenamed, FromName: "", ToName: ""}}},
	}
	issues := ValidateStrict(changes, nil)
	if !hasError(issues, "FROM and TO") {
		t.Error("both empty should trigger RENAMED error")
	}
}
