package parser

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func testdataPath(parts ...string) string {
	_, file, _, _ := runtime.Caller(0)
	base := filepath.Join(filepath.Dir(file), "..", "..", "testdata")
	return filepath.Join(append([]string{base}, parts...)...)
}

func readTestdata(t *testing.T, parts ...string) string {
	t.Helper()
	p := testdataPath(parts...)
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to read testdata %s: %v", p, err)
	}
	return string(data)
}

// --- Spec Parsing Tests ---

func TestParseSpec_Title(t *testing.T) {
	spec := ParseSpec(readTestdata(t, "valid-spec", "basic.md"))
	if spec.Title != "My Feature Spec" {
		t.Errorf("got title %q, want %q", spec.Title, "My Feature Spec")
	}
}

func TestParseSpec_Purpose(t *testing.T) {
	spec := ParseSpec(readTestdata(t, "valid-spec", "basic.md"))
	if spec.Overview == "" {
		t.Error("overview is empty")
	}
	if len(spec.Overview) < 50 {
		t.Errorf("overview too short: %d chars", len(spec.Overview))
	}
}

func TestParseSpec_Requirements(t *testing.T) {
	spec := ParseSpec(readTestdata(t, "valid-spec", "basic.md"))
	if len(spec.Requirements) != 2 {
		t.Fatalf("got %d requirements, want 2", len(spec.Requirements))
	}
	if spec.Requirements[0].Name != "Data Processing" {
		t.Errorf("first requirement name = %q, want %q", spec.Requirements[0].Name, "Data Processing")
	}
	if spec.Requirements[1].Name != "Error Handling" {
		t.Errorf("second requirement name = %q, want %q", spec.Requirements[1].Name, "Error Handling")
	}
}

func TestParseSpec_Scenarios(t *testing.T) {
	spec := ParseSpec(readTestdata(t, "valid-spec", "basic.md"))
	if len(spec.Requirements) < 1 {
		t.Fatal("no requirements")
	}
	req := spec.Requirements[0]
	if len(req.Scenarios) != 2 {
		t.Fatalf("got %d scenarios for %q, want 2", len(req.Scenarios), req.Name)
	}
	if req.Scenarios[0].Name != "Happy path" {
		t.Errorf("scenario name = %q, want %q", req.Scenarios[0].Name, "Happy path")
	}
	if req.Scenarios[0].Text == "" {
		t.Error("scenario text is empty")
	}
}

func TestParseSpec_WhitespaceNormalization(t *testing.T) {
	content := `#  My Spec

##  Purpose
This is a purpose section that is long enough to pass validation checks for minimum length.

##  Requirements

###  Requirement:  Feature Name
The system SHALL do something.

####  Scenario:  Test Case
- **WHEN** something happens
- **THEN** result occurs
`
	spec := ParseSpec(content)
	if spec.Title != "My Spec" {
		t.Errorf("title = %q, want %q", spec.Title, "My Spec")
	}
	if len(spec.Requirements) != 1 {
		t.Fatalf("got %d requirements, want 1", len(spec.Requirements))
	}
	if spec.Requirements[0].Name != "Feature Name" {
		t.Errorf("requirement name = %q, want %q", spec.Requirements[0].Name, "Feature Name")
	}
	if len(spec.Requirements[0].Scenarios) != 1 {
		t.Fatalf("got %d scenarios, want 1", len(spec.Requirements[0].Scenarios))
	}
	if spec.Requirements[0].Scenarios[0].Name != "Test Case" {
		t.Errorf("scenario name = %q, want %q", spec.Requirements[0].Scenarios[0].Name, "Test Case")
	}
}

func TestParseSpec_EmptySections(t *testing.T) {
	content := `# Empty Spec

## Purpose

## Requirements
`
	spec := ParseSpec(content)
	if spec.Title != "Empty Spec" {
		t.Errorf("title = %q", spec.Title)
	}
	if spec.Overview != "" {
		t.Errorf("overview should be empty, got %q", spec.Overview)
	}
	if len(spec.Requirements) != 0 {
		t.Errorf("got %d requirements, want 0", len(spec.Requirements))
	}
}

// --- Change Parsing Tests ---

func TestParseChange_Title(t *testing.T) {
	change := ParseChange(readTestdata(t, "valid-change", "proposal.md"))
	if change.Title != "Add batch processing support" {
		t.Errorf("title = %q, want %q", change.Title, "Add batch processing support")
	}
}

func TestParseChange_Sections(t *testing.T) {
	change := ParseChange(readTestdata(t, "valid-change", "proposal.md"))
	if change.Why == "" {
		t.Error("Why section is empty")
	}
	if change.WhatChanges == "" {
		t.Error("What Changes section is empty")
	}
	if change.Impact == "" {
		t.Error("Impact section is empty")
	}
}

func TestParseChange_MissingSections(t *testing.T) {
	content := `# Change: Incomplete

## Why
Some reason that is long enough for the minimum length check.
`
	change := ParseChange(content)
	if change.Title != "Incomplete" {
		t.Errorf("title = %q", change.Title)
	}
	if change.WhatChanges != "" {
		t.Errorf("expected empty WhatChanges, got %q", change.WhatChanges)
	}
}

// --- Delta Parsing Tests ---

func TestParseDelta_Added(t *testing.T) {
	content := readTestdata(t, "valid-change", "specs", "data-processing", "spec.md")
	deltas := ParseDeltas(content)

	var added []int
	for i, d := range deltas {
		if d.Operation == "ADDED" {
			added = append(added, i)
		}
	}
	if len(added) == 0 {
		t.Fatal("no ADDED deltas found")
	}
	d := deltas[added[0]]
	if len(d.Requirements) != 1 || d.Requirements[0].Name != "Batch Processing" {
		t.Errorf("ADDED requirement = %+v", d.Requirements)
	}
}

func TestParseDelta_Modified(t *testing.T) {
	content := readTestdata(t, "valid-change", "specs", "data-processing", "spec.md")
	deltas := ParseDeltas(content)

	var modified []int
	for i, d := range deltas {
		if d.Operation == "MODIFIED" {
			modified = append(modified, i)
		}
	}
	if len(modified) == 0 {
		t.Fatal("no MODIFIED deltas found")
	}
}

func TestParseDelta_Removed(t *testing.T) {
	content := `## REMOVED Requirements

### Requirement: Deprecated Feature
This requirement is being removed because it is no longer needed.
`
	deltas := ParseDeltas(content)
	if len(deltas) != 1 {
		t.Fatalf("got %d deltas, want 1", len(deltas))
	}
	if deltas[0].Operation != "REMOVED" {
		t.Errorf("operation = %q", deltas[0].Operation)
	}
	if deltas[0].Requirements[0].Name != "Deprecated Feature" {
		t.Errorf("name = %q", deltas[0].Requirements[0].Name)
	}
}

func TestParseDelta_Renamed(t *testing.T) {
	content := `## RENAMED Requirements

### Requirement: Old Name
- **From**: Old Name
- **To**: New Name
`
	deltas := ParseDeltas(content)
	if len(deltas) != 1 {
		t.Fatalf("got %d deltas, want 1", len(deltas))
	}
	if deltas[0].Operation != "RENAMED" {
		t.Errorf("operation = %q", deltas[0].Operation)
	}
	if deltas[0].FromName != "Old Name" {
		t.Errorf("from = %q", deltas[0].FromName)
	}
	if deltas[0].ToName != "New Name" {
		t.Errorf("to = %q", deltas[0].ToName)
	}
}

func TestParseDelta_MultipleSections(t *testing.T) {
	content := readTestdata(t, "valid-change", "specs", "data-processing", "spec.md")
	deltas := ParseDeltas(content)
	if len(deltas) < 2 {
		t.Fatalf("got %d deltas, want at least 2 (ADDED + MODIFIED)", len(deltas))
	}
	ops := make(map[string]bool)
	for _, d := range deltas {
		ops[string(d.Operation)] = true
	}
	if !ops["ADDED"] || !ops["MODIFIED"] {
		t.Errorf("expected ADDED and MODIFIED operations, got %v", ops)
	}
}

func TestParseDelta_HeaderNormalization(t *testing.T) {
	content := `##  ADDED   Requirements

### Requirement: Normalized
The system SHALL work with normalized headers.

#### Scenario: Basic
- **WHEN** triggered
- **THEN** it works
`
	deltas := ParseDeltas(content)
	if len(deltas) != 1 {
		t.Fatalf("got %d deltas, want 1", len(deltas))
	}
	if deltas[0].Operation != "ADDED" {
		t.Errorf("operation = %q, want ADDED", deltas[0].Operation)
	}
}

// --- Task Progress Tests ---

func TestParseTaskProgress_Counts(t *testing.T) {
	status := ParseTaskProgress(readTestdata(t, "valid-change", "tasks.md"))
	if status.Total != 5 {
		t.Errorf("total = %d, want 5", status.Total)
	}
	if status.Completed != 3 {
		t.Errorf("completed = %d, want 3", status.Completed)
	}
}

func TestParseTaskProgress_CaseInsensitive(t *testing.T) {
	content := `- [X] Done uppercase
- [x] Done lowercase
- [ ] Not done
`
	status := ParseTaskProgress(content)
	if status.Total != 3 || status.Completed != 2 {
		t.Errorf("got total=%d completed=%d, want total=3 completed=2", status.Total, status.Completed)
	}
}

func TestParseTaskProgress_IgnoresNonTasks(t *testing.T) {
	content := `## Section Header
Some paragraph text.

- Regular list item
- [x] Actual task
- Another list item
* [x] Star task

Not a task: [x] inline checkbox
`
	status := ParseTaskProgress(content)
	if status.Total != 2 {
		t.Errorf("total = %d, want 2", status.Total)
	}
	if status.Completed != 2 {
		t.Errorf("completed = %d, want 2", status.Completed)
	}
}

// --- Malformed Input Tests ---

func TestParseSpec_EmptyFile(t *testing.T) {
	spec := ParseSpec("")
	if spec.Title != "" {
		t.Errorf("empty file should have empty title, got %q", spec.Title)
	}
	if len(spec.Requirements) != 0 {
		t.Errorf("empty file should have 0 requirements, got %d", len(spec.Requirements))
	}
}

func TestParseSpec_OnlyWhitespace(t *testing.T) {
	spec := ParseSpec("   \n\n  \t  \n")
	if spec.Title != "" || len(spec.Requirements) != 0 {
		t.Error("whitespace-only file should parse as empty")
	}
}

func TestParseSpec_BinaryData(t *testing.T) {
	// Should not panic on binary input
	spec := ParseSpec("\x00\x01\x02\xff\xfe\x80\x90")
	_ = spec // just verify no panic
}

func TestParseSpec_UnicodeHeaders(t *testing.T) {
	content := "# Spécification: 功能测试 🚀\n\n## Purpose\nThis spec covers unicode handling in headers and is long enough to pass the minimum length check.\n\n## Requirements\n\n### Requirement: Ünïcödé Fëätürë\nThe system SHALL handle unicode.\n\n#### Scenario: Ëmöjï 🎯\n- **WHEN** unicode input\n- **THEN** no crash\n"
	spec := ParseSpec(content)
	if spec.Title == "" {
		t.Error("should parse unicode title")
	}
	if len(spec.Requirements) != 1 {
		t.Fatalf("got %d requirements, want 1", len(spec.Requirements))
	}
	if spec.Requirements[0].Name != "Ünïcödé Fëätürë" {
		t.Errorf("requirement name = %q", spec.Requirements[0].Name)
	}
}

func TestParseSpec_CRLFLineEndings(t *testing.T) {
	content := "# CRLF Spec\r\n\r\n## Purpose\r\nThis spec uses Windows-style CRLF line endings and should still parse correctly as valid markdown.\r\n\r\n## Requirements\r\n\r\n### Requirement: CRLF Feature\r\nThe system SHALL handle CRLF.\r\n\r\n#### Scenario: Windows\r\n- **WHEN** CRLF\r\n- **THEN** works\r\n"
	spec := ParseSpec(content)
	if len(spec.Requirements) != 1 {
		t.Fatalf("CRLF: got %d requirements, want 1", len(spec.Requirements))
	}
}

func TestParseSpec_DeeplyNestedHeaders(t *testing.T) {
	content := "# Level 1\n## Level 2\n### Requirement: Level 3\nSHALL work.\n#### Scenario: Level 4\nsteps\n##### Level 5\n###### Level 6\n"
	spec := ParseSpec(content)
	if len(spec.Requirements) != 1 {
		t.Errorf("deeply nested: got %d requirements", len(spec.Requirements))
	}
}

func TestParseSpec_VeryLongLine(t *testing.T) {
	longLine := strings.Repeat("a", 100000)
	content := "# Long Spec\n\n## Purpose\n" + longLine + "\n\n## Requirements\n"
	spec := ParseSpec(content)
	if spec.Title != "Long Spec" {
		t.Errorf("title = %q", spec.Title)
	}
}

func TestParseSpec_HeadersWithNoContent(t *testing.T) {
	content := "# \n## \n### Requirement: \n#### Scenario: \n"
	spec := ParseSpec(content)
	// Should not panic; title should be empty string after trim
	_ = spec
}

func TestParseChange_EmptyFile(t *testing.T) {
	change := ParseChange("")
	if change.Title != "" || change.Why != "" {
		t.Error("empty file should produce empty change")
	}
}

func TestParseDelta_EmptyFile(t *testing.T) {
	deltas := ParseDeltas("")
	if len(deltas) != 0 {
		t.Errorf("empty file should produce 0 deltas, got %d", len(deltas))
	}
}

func TestParseDelta_InvalidOperations(t *testing.T) {
	content := "## INVALID Requirements\n\n### Requirement: Ghost\nSHALL haunt.\n"
	deltas := ParseDeltas(content)
	if len(deltas) != 0 {
		t.Errorf("invalid operation should produce 0 deltas, got %d", len(deltas))
	}
}

func TestParseDelta_CaseSensitiveOperations(t *testing.T) {
	content := "## added Requirements\n\n### Requirement: Lowercase\nSHALL work.\n"
	deltas := ParseDeltas(content)
	// "added" in lowercase should still match because parseDeltaHeader uppercases
	if len(deltas) != 1 {
		t.Errorf("lowercase 'added' should match, got %d deltas", len(deltas))
	}
}

func TestParseTaskProgress_EmptyFile(t *testing.T) {
	status := ParseTaskProgress("")
	if status.Total != 0 || status.Completed != 0 {
		t.Error("empty file should have 0 tasks")
	}
}

func TestParseTaskProgress_MalformedCheckboxes(t *testing.T) {
	content := "- [ x] Spaced\n- [x ] Trailing\n- [x Unclosed\n- x] No bracket\n- [X] Valid\n"
	status := ParseTaskProgress(content)
	// Only "- [X] Valid" should match
	if status.Total != 1 {
		t.Errorf("expected 1 valid task, got total=%d", status.Total)
	}
}

func TestParseSpec_ConsecutiveHeaders(t *testing.T) {
	// Headers with no content between them
	content := "# Title\n## Purpose\n## Requirements\n### Requirement: Empty Body\n#### Scenario: No Steps\n"
	spec := ParseSpec(content)
	if len(spec.Requirements) != 1 {
		t.Fatalf("got %d requirements, want 1", len(spec.Requirements))
	}
	if spec.Requirements[0].Name != "Empty Body" {
		t.Errorf("name = %q", spec.Requirements[0].Name)
	}
}

func TestParseSpec_DuplicateRequirementNames(t *testing.T) {
	content := "# Spec\n\n## Purpose\nLong enough purpose for the minimum length check requirement.\n\n## Requirements\n\n### Requirement: Same Name\nSHALL do first.\n\n#### Scenario: S1\n- steps\n\n### Requirement: Same Name\nSHALL do second.\n\n#### Scenario: S2\n- steps\n"
	spec := ParseSpec(content)
	if len(spec.Requirements) != 2 {
		t.Fatalf("should parse both duplicates, got %d", len(spec.Requirements))
	}
}
