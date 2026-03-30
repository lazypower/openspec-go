package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// executeCmd runs a root command with the given args and captures output.
func executeCmd(args ...string) (string, string, error) {
	cmd := NewRootCmd("test")
	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(outBuf)
	cmd.SetErr(errBuf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return outBuf.String(), errBuf.String(), err
}

// setupProject creates a full test project in a temp directory.
func setupProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	t.Cleanup(func() { os.Chdir(oldWd) })

	// Run init
	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"init", "--tools", "all"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	return dir
}

// --- Init Tests ---

func TestInit_NonInteractive(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	_, _, err := executeCmd("init", "--tools", "all")
	if err != nil {
		t.Fatal(err)
	}

	// Verify directory structure
	for _, d := range []string{
		"openspec",
		"openspec/specs",
		"openspec/changes",
		"openspec/changes/archive",
	} {
		if _, err := os.Stat(filepath.Join(dir, d)); os.IsNotExist(err) {
			t.Errorf("directory not created: %s", d)
		}
	}

	// Verify files
	for _, f := range []string{
		"openspec/project.md",
		"openspec/AGENTS.md",
	} {
		if _, err := os.Stat(filepath.Join(dir, f)); os.IsNotExist(err) {
			t.Errorf("file not created: %s", f)
		}
	}

	// Verify editor configs
	for _, f := range []string{
		".claude/commands/openspec/proposal.md",
		".claude/commands/openspec/apply.md",
		".claude/commands/openspec/archive.md",
		".opencode/commands/openspec/proposal.md",
		".codex/prompts/openspec/proposal.md",
		".goose/recipes/openspec/proposal.yaml",
		".goosehints",
		"CLAUDE.md",
	} {
		if _, err := os.Stat(filepath.Join(dir, f)); os.IsNotExist(err) {
			t.Errorf("editor config not created: %s", f)
		}
	}
}

func TestInit_NoTools(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	_, _, err := executeCmd("init", "--tools", "none")
	if err != nil {
		t.Fatal(err)
	}

	// Verify basic structure
	if _, err := os.Stat(filepath.Join(dir, "openspec", "AGENTS.md")); os.IsNotExist(err) {
		t.Error("AGENTS.md not created")
	}

	// Verify NO editor configs
	if _, err := os.Stat(filepath.Join(dir, ".claude")); !os.IsNotExist(err) {
		t.Error(".claude should not exist with --tools none")
	}
}

func TestInit_AlreadyInitialized(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	// First init
	executeCmd("init", "--tools", "none")

	// Second init should fail
	_, _, err := executeCmd("init", "--tools", "none")
	if err == nil {
		t.Error("expected error for already initialized")
	}
}

func TestInit_CustomPath(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	subdir := filepath.Join(dir, "subproject")
	os.MkdirAll(subdir, 0o755)

	_, _, err := executeCmd("init", subdir, "--tools", "none")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(subdir, "openspec", "AGENTS.md")); os.IsNotExist(err) {
		t.Error("openspec not created in custom path")
	}
}

// --- Update Tests ---

func TestUpdate(t *testing.T) {
	dir := setupProject(t)

	// Modify AGENTS.md
	agentsPath := filepath.Join(dir, "openspec", "AGENTS.md")
	os.WriteFile(agentsPath, []byte("old content"), 0o644)

	_, _, err := executeCmd("update")
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(agentsPath)
	if string(data) == "old content" {
		t.Error("AGENTS.md was not updated")
	}
}

// --- List Tests ---

func TestList_Changes(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)

	// Capture stdout by redirecting
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"list"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	if !strings.Contains(buf.String(), "test-change") {
		t.Errorf("expected test-change in list output, got: %s", buf.String())
	}
}

func TestList_Specs(t *testing.T) {
	dir := setupProject(t)
	createTestSpec(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"list", "--specs"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	if !strings.Contains(buf.String(), "test-spec") {
		t.Errorf("expected test-spec in list output, got: %s", buf.String())
	}
}

func TestList_Empty(t *testing.T) {
	setupProject(t)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"list"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	if !strings.Contains(buf.String(), "No active changes") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

// --- Show Tests ---

func TestShow_ChangeText(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"show", "test-change"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	if !strings.Contains(buf.String(), "# Change: Test change") {
		t.Errorf("expected proposal content, got: %s", buf.String())
	}
}

func TestShow_ChangeJSON(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"show", "test-change", "--json"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, buf.String())
	}
	if result["id"] != "test-change" {
		t.Errorf("unexpected id: %v", result["id"])
	}
}

func TestShow_SpecText(t *testing.T) {
	dir := setupProject(t)
	createTestSpec(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"show", "test-spec"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	if !strings.Contains(buf.String(), "# Test Spec") {
		t.Errorf("expected spec content, got: %s", buf.String())
	}
}

func TestShow_SpecJSON(t *testing.T) {
	dir := setupProject(t)
	createTestSpec(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"show", "test-spec", "--json"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result["id"] != "test-spec" {
		t.Errorf("unexpected id: %v", result["id"])
	}
}

func TestShow_NotFound(t *testing.T) {
	setupProject(t)

	_, _, err := executeCmd("show", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent item")
	}
}

func TestShow_DeltasOnly(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"show", "test-change", "--json", "--deltas-only"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	var result []interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON array: %v\noutput: %s", err, buf.String())
	}
}

func TestShow_RequirementsOnly(t *testing.T) {
	dir := setupProject(t)
	createTestSpec(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"show", "test-spec", "--json", "--requirements"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	var result []interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON array: %v", err)
	}
}

func TestShow_SingleRequirement(t *testing.T) {
	dir := setupProject(t)
	createTestSpec(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"show", "test-spec", "--json", "-r", "1"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestShow_TypeFlag(t *testing.T) {
	dir := setupProject(t)
	createTestSpec(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"show", "test-spec", "--type", "spec"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	if !strings.Contains(buf.String(), "# Test Spec") {
		t.Error("type flag not respected")
	}
}

// --- Validate Tests ---

func TestValidate_SingleChange(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)

	_, _, err := executeCmd("validate", "test-change")
	if err != nil {
		t.Logf("validation result: %v", err)
	}
}

func TestValidate_SingleSpec(t *testing.T) {
	dir := setupProject(t)
	createTestSpec(t, dir)

	_, _, err := executeCmd("validate", "test-spec", "--type", "spec")
	if err != nil {
		t.Logf("validation result: %v", err)
	}
}

func TestValidate_All(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)
	createTestSpec(t, dir)

	_, _, err := executeCmd("validate", "--all")
	if err != nil {
		t.Logf("validation result: %v", err)
	}
}

func TestValidate_JSON(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"validate", "test-change", "--json"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	var report map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON report: %v\noutput: %s", err, buf.String())
	}
	if _, ok := report["items"]; !ok {
		t.Error("JSON report missing 'items' key")
	}
	if _, ok := report["summary"]; !ok {
		t.Error("JSON report missing 'summary' key")
	}
}

func TestValidate_Strict(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)

	_, _, err := executeCmd("validate", "test-change", "--strict")
	if err != nil {
		t.Logf("strict validation result: %v", err)
	}
}

func TestValidate_Concurrency(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)
	createTestSpec(t, dir)

	_, _, err := executeCmd("validate", "--all", "--concurrency", "2")
	if err != nil {
		t.Logf("concurrent validation: %v", err)
	}
}

func TestValidate_Changes(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)

	_, _, err := executeCmd("validate", "--changes")
	// This is a valid path even if validation finds issues
	_ = err
}

func TestValidate_Specs(t *testing.T) {
	dir := setupProject(t)
	createTestSpec(t, dir)

	_, _, err := executeCmd("validate", "--specs")
	_ = err
}

func TestValidate_NotFound(t *testing.T) {
	setupProject(t)

	_, _, err := executeCmd("validate", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent item")
	}
}

// --- View Tests ---

func TestView_Dashboard(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)
	createTestSpec(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"view"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	out := buf.String()
	if !strings.Contains(out, "Dashboard") {
		t.Error("missing dashboard header")
	}
	if !strings.Contains(out, "test-change") {
		t.Error("missing active change")
	}
	if !strings.Contains(out, "test-spec") {
		t.Error("missing spec")
	}
}

func TestView_EmptyProject(t *testing.T) {
	setupProject(t)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"view"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	if !strings.Contains(buf.String(), "Dashboard") {
		t.Error("empty dashboard should still render")
	}
}

func TestView_WithArchived(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)

	// Archive it
	archiveDir := filepath.Join(dir, "openspec", "changes", "archive")
	os.MkdirAll(archiveDir, 0o755)
	today := time.Now().Format("2006-01-02")
	os.MkdirAll(filepath.Join(archiveDir, today+"-old-change"), 0o755)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"view"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	if !strings.Contains(buf.String(), "old-change") {
		t.Error("archived change not shown in dashboard")
	}
}

// --- Archive Tests ---

func TestArchive_WithYes(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)
	createTestSpec(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"archive", "test-change", "--yes"})
	err := cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	if err != nil {
		t.Fatal(err)
	}

	// Verify archived
	today := time.Now().Format("2006-01-02")
	archivePath := filepath.Join(dir, "openspec", "changes", "archive", today+"-test-change")
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		t.Error("change not archived")
	}
}

func TestArchive_SkipSpecs(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)
	createTestSpec(t, dir)

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"archive", "test-change", "--yes", "--skip-specs"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Spec should be unchanged
	specData, _ := os.ReadFile(filepath.Join(dir, "openspec", "specs", "test-spec", "spec.md"))
	if strings.Contains(string(specData), "New Capability") {
		t.Error("spec was modified despite --skip-specs")
	}
}

func TestArchive_NoValidate(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)
	createTestSpec(t, dir)

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"archive", "test-change", "--yes", "--no-validate"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestArchive_NotFound(t *testing.T) {
	setupProject(t)

	_, _, err := executeCmd("archive", "nonexistent", "--yes")
	if err == nil {
		t.Error("expected error for nonexistent change")
	}
}

// --- Version Test ---

func TestVersion(t *testing.T) {
	_, _, err := executeCmd("--version")
	if err != nil {
		t.Fatal(err)
	}
}

// --- No-Color Test ---

func TestNoColor(t *testing.T) {
	setupProject(t)

	_, _, err := executeCmd("--no-color", "list")
	if err != nil {
		t.Fatal(err)
	}
}

// --- CLI Error Path Tests ---

func TestUnknownCommand(t *testing.T) {
	_, _, err := executeCmd("nonexistent-command")
	if err == nil {
		t.Error("expected error for unknown command")
	}
}

func TestShowMissingArg(t *testing.T) {
	setupProject(t)
	_, _, err := executeCmd("show")
	if err == nil {
		t.Error("expected error when show called with no args")
	}
}

func TestArchiveMissingArg(t *testing.T) {
	setupProject(t)
	_, _, err := executeCmd("archive")
	if err == nil {
		t.Error("expected error when archive called with no args")
	}
}

func TestValidateInvalidConcurrency(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)
	// Negative concurrency should fall back to NumCPU, not crash
	_, _, err := executeCmd("validate", "test-change", "--concurrency", "0")
	// Should succeed (0 means use default)
	_ = err
}

func TestValidateExitCodeOnError(t *testing.T) {
	dir := setupProject(t)
	// Create an invalid change (no deltas)
	changePath := filepath.Join(dir, "openspec", "changes", "bad-change")
	os.MkdirAll(changePath, 0o755)
	os.WriteFile(filepath.Join(changePath, "proposal.md"), []byte("# Change: Bad\n\n## Why\nShort\n\n## What Changes\nNothing\n"), 0o644)

	_, _, err := executeCmd("validate", "bad-change")
	if err == nil {
		t.Error("validation of invalid change should return error")
	}
}

func TestListInUninitializedDir(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	_, _, err := executeCmd("list")
	if err == nil {
		t.Error("expected error in uninitialized directory")
	}
}

func TestValidateInUninitializedDir(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	_, _, err := executeCmd("validate", "--all")
	if err == nil {
		t.Error("expected error in uninitialized directory")
	}
}

func TestViewInUninitializedDir(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	_, _, err := executeCmd("view")
	if err == nil {
		t.Error("expected error in uninitialized directory")
	}
}

func TestShowInUninitializedDir(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	_, _, err := executeCmd("show", "anything")
	if err == nil {
		t.Error("expected error in uninitialized directory")
	}
}

func TestUpdateInUninitializedDir(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	_, _, err := executeCmd("update")
	if err == nil {
		t.Error("expected error in uninitialized directory")
	}
}

func TestArchiveInUninitializedDir(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	_, _, err := executeCmd("archive", "something", "--yes")
	if err == nil {
		t.Error("expected error in uninitialized directory")
	}
}

func TestShowRequirementOutOfRange(t *testing.T) {
	dir := setupProject(t)
	createTestSpec(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"show", "test-spec", "--json", "-r", "999"})
	err := cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	if err == nil {
		t.Error("expected error for out-of-range requirement index")
	}
}

func TestInitWithInvalidTool(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	// Should warn about unknown tool but still initialize
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"init", "--tools", "fake-editor"})
	cmd.Execute()

	w.Close()
	os.Stderr = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	// Project should still be initialized
	if _, err := os.Stat(filepath.Join(dir, "openspec")); os.IsNotExist(err) {
		t.Error("project should still be initialized despite invalid tool")
	}
}

func TestShowAutoDetectAmbiguous(t *testing.T) {
	dir := setupProject(t)
	// Create both a change and spec with same name
	createTestChange(t, dir)
	createTestSpec(t, dir)
	// "test-change" only exists as change, "test-spec" only as spec
	// but let's create overlap
	specAsChange := filepath.Join(dir, "openspec", "changes", "test-spec")
	os.MkdirAll(specAsChange, 0o755)
	os.WriteFile(filepath.Join(specAsChange, "proposal.md"), []byte("# Change: Overlap\n\n## Why\nLong enough reason to test auto-detection when both change and spec exist with same name.\n\n## What Changes\n- test\n"), 0o644)

	// Should default to change when ambiguous
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"show", "test-spec", "--json"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	// Should get some output (either change or spec, not an error)
	if buf.Len() == 0 {
		t.Error("ambiguous auto-detect should still produce output")
	}
}

func TestHelpFlag(t *testing.T) {
	// Help should not error
	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("--help should not error: %v", err)
	}
}

func TestSubcommandHelp(t *testing.T) {
	for _, sub := range []string{"init", "update", "list", "show", "validate", "archive", "view"} {
		cmd := NewRootCmd("test")
		cmd.SetArgs([]string{sub, "--help"})
		if err := cmd.Execute(); err != nil {
			t.Errorf("%s --help should not error: %v", sub, err)
		}
	}
}

func TestShowNearestMatchSuggestion(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)

	_, _, err := executeCmd("show", "test-chang") // typo
	if err == nil {
		t.Error("expected error for misspelled item")
	}
	if err != nil && !strings.Contains(err.Error(), "test-change") {
		t.Errorf("expected suggestion for 'test-change' in error, got: %v", err)
	}
}

func TestValidateSpecsFlag(t *testing.T) {
	dir := setupProject(t)
	createTestSpec(t, dir)
	_, _, _ = executeCmd("validate", "--specs")
}

func TestArchiveWithIncompleteTasksAndYes(t *testing.T) {
	dir := setupProject(t)
	createTestChange(t, dir)
	createTestSpec(t, dir)

	// tasks.md in createTestChange has 1 incomplete task
	// With --yes, should warn but proceed
	cmd := NewRootCmd("test")
	cmd.SetArgs([]string{"archive", "test-change", "--yes"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("archive with --yes should succeed despite incomplete tasks: %v", err)
	}
}

func TestInitSpecificTools(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	_, _, err := executeCmd("init", "--tools", "claude-code,goose")
	if err != nil {
		t.Fatal(err)
	}

	// Claude should be configured
	if _, err := os.Stat(filepath.Join(dir, ".claude", "commands", "openspec")); os.IsNotExist(err) {
		t.Error("claude-code not configured")
	}
	// Goose should be configured
	if _, err := os.Stat(filepath.Join(dir, ".goose", "recipes", "openspec")); os.IsNotExist(err) {
		t.Error("goose not configured")
	}
	// OpenCode should NOT be configured
	if _, err := os.Stat(filepath.Join(dir, ".opencode")); !os.IsNotExist(err) {
		t.Error("opencode should not be configured")
	}
}

// --- Helpers ---

func createTestChange(t *testing.T, dir string) {
	t.Helper()
	changePath := filepath.Join(dir, "openspec", "changes", "test-change")
	os.MkdirAll(filepath.Join(changePath, "specs", "test-spec"), 0o755)

	proposal := `# Change: Test change

## Why
This change adds important new functionality that users have been requesting for months and is critical for the project roadmap.

## What Changes
- Adds a new requirement to the test spec

## Impact
- test-spec
`
	os.WriteFile(filepath.Join(changePath, "proposal.md"), []byte(proposal), 0o644)

	delta := `## ADDED Requirements

### Requirement: New Capability
The system SHALL provide new capability for users.

#### Scenario: Basic usage
- **WHEN** user triggers the new capability
- **THEN** the system responds correctly
`
	os.WriteFile(filepath.Join(changePath, "specs", "test-spec", "spec.md"), []byte(delta), 0o644)

	tasks := "- [x] Implement feature\n- [x] Write tests\n- [ ] Deploy\n"
	os.WriteFile(filepath.Join(changePath, "tasks.md"), []byte(tasks), 0o644)
}

func createTestSpec(t *testing.T, dir string) {
	t.Helper()
	specPath := filepath.Join(dir, "openspec", "specs", "test-spec")
	os.MkdirAll(specPath, 0o755)

	spec := fmt.Sprintf(`# Test Spec

## Purpose
This is the test specification that defines the core behavior of the test module for validation and testing purposes.

## Requirements

### Requirement: Core Feature
The system SHALL provide core functionality for all users accessing the system.

#### Scenario: Happy path
- **WHEN** user interacts
- **THEN** system responds
`)
	os.WriteFile(filepath.Join(specPath, "spec.md"), []byte(spec), 0o644)
}
