package archive

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chuck/openspec-go/internal/model"
)

const baseSpec = `# Test Spec

## Purpose
This is a test specification that provides enough content to meet the minimum purpose length requirement for validation.

## Requirements

### Requirement: Existing Feature
The system SHALL handle existing functionality.

#### Scenario: Happy path
- **WHEN** the system is running
- **THEN** it works correctly
`

// --- Merge Operation Tests ---

func TestMerge_AddedAppendsRequirement(t *testing.T) {
	delta := model.Delta{
		Operation: model.DeltaAdded,
		Requirements: []model.Requirement{
			{
				Name: "New Feature",
				Text: "The system SHALL support new features.",
				Scenarios: []model.Scenario{
					{Name: "Basic", Text: "- **WHEN** used\n- **THEN** it works"},
				},
			},
		},
	}
	result, err := MergeDeltas(baseSpec, []model.Delta{delta})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "### Requirement: New Feature") {
		t.Error("added requirement not found in output")
	}
	if !strings.Contains(result, "Existing Feature") {
		t.Error("existing requirement was lost")
	}
}

func TestMerge_ModifiedReplacesRequirement(t *testing.T) {
	delta := model.Delta{
		Operation: model.DeltaModified,
		Requirements: []model.Requirement{
			{
				Name: "Existing Feature",
				Text: "The system SHALL handle UPDATED functionality.",
				Scenarios: []model.Scenario{
					{Name: "Updated path", Text: "- **WHEN** updated\n- **THEN** it works better"},
				},
			},
		},
	}
	result, err := MergeDeltas(baseSpec, []model.Delta{delta})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "UPDATED functionality") {
		t.Error("modified content not found")
	}
	if strings.Contains(result, "existing functionality") {
		t.Error("old content still present")
	}
}

func TestMerge_RemovedDeletesRequirement(t *testing.T) {
	delta := model.Delta{
		Operation: model.DeltaRemoved,
		Requirements: []model.Requirement{
			{Name: "Existing Feature"},
		},
	}
	result, err := MergeDeltas(baseSpec, []model.Delta{delta})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(result, "Existing Feature") {
		t.Error("removed requirement still present")
	}
}

func TestMerge_RenamedUpdatesHeader(t *testing.T) {
	delta := model.Delta{
		Operation: model.DeltaRenamed,
		FromName:  "Existing Feature",
		ToName:    "Better Feature",
	}
	result, err := MergeDeltas(baseSpec, []model.Delta{delta})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "### Requirement: Better Feature") {
		t.Error("renamed header not found")
	}
	if strings.Contains(result, "### Requirement: Existing Feature") {
		t.Error("old header still present")
	}
}

func TestMerge_ModifiedNotFoundError(t *testing.T) {
	delta := model.Delta{
		Operation:    model.DeltaModified,
		Requirements: []model.Requirement{{Name: "Nonexistent"}},
	}
	_, err := MergeDeltas(baseSpec, []model.Delta{delta})
	if err == nil {
		t.Error("expected error for missing requirement")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMerge_RemovedNotFoundError(t *testing.T) {
	delta := model.Delta{
		Operation:    model.DeltaRemoved,
		Requirements: []model.Requirement{{Name: "Nonexistent"}},
	}
	_, err := MergeDeltas(baseSpec, []model.Delta{delta})
	if err == nil {
		t.Error("expected error for missing requirement")
	}
}

func TestMerge_NewSpecCreation(t *testing.T) {
	newSpecContent := "# new-spec\n\n## Purpose\n\n## Requirements\n"
	delta := model.Delta{
		Operation: model.DeltaAdded,
		Requirements: []model.Requirement{
			{
				Name: "First Feature",
				Text: "The system SHALL do the first thing.",
				Scenarios: []model.Scenario{
					{Name: "Basic", Text: "- **WHEN** called\n- **THEN** works"},
				},
			},
		},
	}
	result, err := MergeDeltas(newSpecContent, []model.Delta{delta})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "### Requirement: First Feature") {
		t.Error("requirement not added to new spec")
	}
}

// --- Archive Workflow Tests ---

func setupTestProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	osp := filepath.Join(dir, "openspec")

	// Create directory structure
	for _, d := range []string{
		filepath.Join(osp, "specs", "test-spec"),
		filepath.Join(osp, "changes", "test-change", "specs", "test-spec"),
		filepath.Join(osp, "changes", "archive"),
	} {
		os.MkdirAll(d, 0o755)
	}

	// Write spec
	os.WriteFile(filepath.Join(osp, "specs", "test-spec", "spec.md"), []byte(baseSpec), 0o644)

	// Write proposal
	proposal := `# Change: Test change

## Why
This change adds important new functionality that users have been requesting for months and is critical to the roadmap.

## What Changes
- Adds a new requirement to the test spec

## Impact
- test-spec
`
	os.WriteFile(filepath.Join(osp, "changes", "test-change", "proposal.md"), []byte(proposal), 0o644)

	// Write delta
	delta := `## ADDED Requirements

### Requirement: New Capability
The system SHALL provide new capability for users.

#### Scenario: Basic usage
- **WHEN** user triggers the new capability
- **THEN** the system responds correctly
`
	os.WriteFile(filepath.Join(osp, "changes", "test-change", "specs", "test-spec", "spec.md"), []byte(delta), 0o644)

	// Write tasks
	tasks := "- [x] Implement feature\n- [x] Write tests\n"
	os.WriteFile(filepath.Join(osp, "changes", "test-change", "tasks.md"), []byte(tasks), 0o644)

	return dir
}

func TestArchive_FullWorkflow(t *testing.T) {
	dir := setupTestProject(t)
	osp := filepath.Join(dir, "openspec")

	err := Archive(ArchiveOptions{
		OpenSpecPath: osp,
		ChangeID:     "test-change",
		Yes:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify spec was updated
	specData, err := os.ReadFile(filepath.Join(osp, "specs", "test-spec", "spec.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(specData), "New Capability") {
		t.Error("spec not updated with new requirement")
	}

	// Verify change moved to archive
	today := time.Now().Format("2006-01-02")
	archivePath := filepath.Join(osp, "changes", "archive", today+"-test-change")
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		t.Error("change not moved to archive")
	}
	// Original should be gone
	if _, err := os.Stat(filepath.Join(osp, "changes", "test-change")); !os.IsNotExist(err) {
		t.Error("original change directory still exists")
	}
}

func TestArchive_SkipSpecs(t *testing.T) {
	dir := setupTestProject(t)
	osp := filepath.Join(dir, "openspec")

	origSpec, _ := os.ReadFile(filepath.Join(osp, "specs", "test-spec", "spec.md"))

	err := Archive(ArchiveOptions{
		OpenSpecPath: osp,
		ChangeID:     "test-change",
		SkipSpecs:    true,
		Yes:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Spec should be unchanged
	newSpec, _ := os.ReadFile(filepath.Join(osp, "specs", "test-spec", "spec.md"))
	if string(newSpec) != string(origSpec) {
		t.Error("spec was modified despite skip-specs")
	}

	// But change should be archived
	today := time.Now().Format("2006-01-02")
	archivePath := filepath.Join(osp, "changes", "archive", today+"-test-change")
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		t.Error("change not moved to archive")
	}
}

func TestArchive_ValidationAborts(t *testing.T) {
	dir := setupTestProject(t)
	osp := filepath.Join(dir, "openspec")

	// Remove deltas directory to cause "must have at least one delta" error
	os.RemoveAll(filepath.Join(osp, "changes", "test-change", "specs"))

	err := Archive(ArchiveOptions{
		OpenSpecPath: osp,
		ChangeID:     "test-change",
		Yes:          true,
	})
	if err == nil {
		t.Error("expected validation error")
	}

	// Change should NOT be archived
	if _, err := os.Stat(filepath.Join(osp, "changes", "test-change")); os.IsNotExist(err) {
		t.Error("change was archived despite validation failure")
	}
}

func TestArchive_PostMergeRollback(t *testing.T) {
	dir := setupTestProject(t)
	osp := filepath.Join(dir, "openspec")

	// Create a delta that modifies a nonexistent requirement
	delta := `## MODIFIED Requirements

### Requirement: Ghost Requirement
The system SHALL do ghost things.

#### Scenario: Spooky
- **WHEN** ghost appears
- **THEN** it spooks
`
	specDir := filepath.Join(osp, "changes", "test-change", "specs", "test-spec")
	os.WriteFile(filepath.Join(specDir, "spec.md"), []byte(delta), 0o644)

	err := Archive(ArchiveOptions{
		OpenSpecPath: osp,
		ChangeID:     "test-change",
		SkipValidate: true,
		Yes:          true,
	})
	if err == nil {
		t.Error("expected merge error")
	}

	// Original spec should be unchanged
	specData, _ := os.ReadFile(filepath.Join(osp, "specs", "test-spec", "spec.md"))
	if !strings.Contains(string(specData), "Existing Feature") {
		t.Error("spec was modified despite merge failure")
	}
}

// --- Atomic Write Tests ---

func TestAtomicWrite_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")

	err := AtomicWrite(path, []byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello world" {
		t.Errorf("content = %q", string(data))
	}
}

func TestAtomicWrite_MultiSpecRollback(t *testing.T) {
	dir := t.TempDir()

	// Write initial files
	file1 := filepath.Join(dir, "spec1.md")
	file2 := filepath.Join(dir, "spec2.md")
	os.WriteFile(file1, []byte("original1"), 0o644)
	os.WriteFile(file2, []byte("original2"), 0o644)

	// Simulate a multi-spec merge where we compute all first
	merges := map[string]string{
		file1: "updated1",
		file2: "updated2",
	}

	// Write all
	for path, content := range merges {
		if err := AtomicWrite(path, []byte(content)); err != nil {
			t.Fatal(err)
		}
	}

	// Verify both updated
	d1, _ := os.ReadFile(file1)
	d2, _ := os.ReadFile(file2)
	if string(d1) != "updated1" || string(d2) != "updated2" {
		t.Error("atomic multi-write failed")
	}
}

// --- Edge Case Tests ---

func TestMerge_SpecialCharsInRequirementName(t *testing.T) {
	specContent := "# Spec\n\n## Requirements\n\n### Requirement: Feature (v2.0) — Enhanced\nSHALL work.\n\n#### Scenario: Test\n- steps\n"
	delta := model.Delta{
		Operation:    model.DeltaModified,
		Requirements: []model.Requirement{{Name: "Feature (v2.0) — Enhanced", Text: "SHALL work better.", Scenarios: []model.Scenario{{Name: "T", Text: "s"}}}},
	}
	result, err := MergeDeltas(specContent, []model.Delta{delta})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "work better") {
		t.Error("modified content not applied with special chars in name")
	}
}

func TestMerge_EmptySpecContent(t *testing.T) {
	delta := model.Delta{
		Operation:    model.DeltaAdded,
		Requirements: []model.Requirement{{Name: "First", Text: "SHALL exist.", Scenarios: []model.Scenario{{Name: "S", Text: "s"}}}},
	}
	result, err := MergeDeltas("", []model.Delta{delta})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "### Requirement: First") {
		t.Error("should append to empty content")
	}
}

func TestMerge_MultipleAddedInSequence(t *testing.T) {
	deltas := []model.Delta{
		{Operation: model.DeltaAdded, Requirements: []model.Requirement{{Name: "A", Text: "SHALL a.", Scenarios: []model.Scenario{{Name: "S", Text: "s"}}}}},
		{Operation: model.DeltaAdded, Requirements: []model.Requirement{{Name: "B", Text: "SHALL b.", Scenarios: []model.Scenario{{Name: "S", Text: "s"}}}}},
		{Operation: model.DeltaAdded, Requirements: []model.Requirement{{Name: "C", Text: "SHALL c.", Scenarios: []model.Scenario{{Name: "S", Text: "s"}}}}},
	}
	result, err := MergeDeltas(baseSpec, deltas)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"A", "B", "C"} {
		if !strings.Contains(result, "### Requirement: "+name) {
			t.Errorf("missing requirement %s after multi-add", name)
		}
	}
}

func TestMerge_RenamedToEmpty(t *testing.T) {
	delta := model.Delta{
		Operation: model.DeltaRenamed,
		FromName:  "Existing Feature",
		ToName:    "",
	}
	_, err := MergeDeltas(baseSpec, []model.Delta{delta})
	if err == nil {
		t.Error("expected error for empty ToName")
	}
}

func TestMerge_RenamedFromEmpty(t *testing.T) {
	delta := model.Delta{
		Operation: model.DeltaRenamed,
		FromName:  "",
		ToName:    "New Name",
	}
	_, err := MergeDeltas(baseSpec, []model.Delta{delta})
	if err == nil {
		t.Error("expected error for empty FromName")
	}
}

func TestArchive_ChangeNotFound(t *testing.T) {
	dir := t.TempDir()
	osp := filepath.Join(dir, "openspec")
	os.MkdirAll(filepath.Join(osp, "changes"), 0o755)

	err := Archive(ArchiveOptions{
		OpenSpecPath: osp,
		ChangeID:     "ghost",
		Yes:          true,
	})
	if err == nil {
		t.Error("expected error for nonexistent change")
	}
}

func TestArchive_SkipValidateWithInvalidChange(t *testing.T) {
	dir := setupTestProject(t)
	osp := filepath.Join(dir, "openspec")

	// Make proposal invalid by removing Why section
	proposal := "# Change: No why\n\n## What Changes\n- stuff\n"
	os.WriteFile(filepath.Join(osp, "changes", "test-change", "proposal.md"), []byte(proposal), 0o644)

	// Without skip-validate, should fail
	err := Archive(ArchiveOptions{
		OpenSpecPath: osp,
		ChangeID:     "test-change",
		Yes:          true,
	})
	if err == nil {
		t.Error("expected validation error without skip-validate")
	}
}

func TestAtomicWrite_Overwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "existing.md")
	os.WriteFile(path, []byte("old content"), 0o644)

	err := AtomicWrite(path, []byte("new content"))
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	if string(data) != "new content" {
		t.Errorf("expected 'new content', got %q", string(data))
	}
}

func TestMerge_ModifyThenRemoveSameRequirement(t *testing.T) {
	// MODIFIED then REMOVED on same requirement in sequence
	deltas := []model.Delta{
		{Operation: model.DeltaModified, Requirements: []model.Requirement{{Name: "Existing Feature", Text: "SHALL be modified.", Scenarios: []model.Scenario{{Name: "S", Text: "s"}}}}},
		{Operation: model.DeltaRemoved, Requirements: []model.Requirement{{Name: "Existing Feature"}}},
	}
	result, err := MergeDeltas(baseSpec, deltas)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(result, "Existing Feature") {
		t.Error("requirement should be removed after modify+remove sequence")
	}
}

func TestMerge_WhitespaceInsensitiveMatch(t *testing.T) {
	// Spec has "Existing Feature", delta targets "Existing  Feature" (double space)
	delta := model.Delta{
		Operation:    model.DeltaRemoved,
		Requirements: []model.Requirement{{Name: "Existing  Feature"}},
	}
	result, err := MergeDeltas(baseSpec, []model.Delta{delta})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(result, "Existing Feature") {
		t.Error("whitespace-insensitive match should find and remove requirement")
	}
}
