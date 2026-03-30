package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/chuck/openspec-go/internal/model"
	"github.com/chuck/openspec-go/internal/parser"
	"github.com/chuck/openspec-go/internal/validator"
)

var reqHeaderRe = regexp.MustCompile(`(?i)^###\s+Requirement:\s+(.+)$`)

// MergeDeltas applies deltas to a spec's markdown content and returns the new content.
func MergeDeltas(specContent string, deltas []model.Delta) (string, error) {
	result := specContent

	for _, delta := range deltas {
		var err error
		switch delta.Operation {
		case model.DeltaAdded:
			result, err = applyAdded(result, delta)
		case model.DeltaModified:
			result, err = applyModified(result, delta)
		case model.DeltaRemoved:
			result, err = applyRemoved(result, delta)
		case model.DeltaRenamed:
			result, err = applyRenamed(result, delta)
		}
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

func applyAdded(content string, delta model.Delta) (string, error) {
	for _, req := range delta.Requirements {
		block := renderRequirement(req)
		// Append after the last content
		content = strings.TrimRight(content, "\n") + "\n\n" + block + "\n"
	}
	return content, nil
}

func applyModified(content string, delta model.Delta) (string, error) {
	for _, req := range delta.Requirements {
		start, end, found := findRequirementBlock(content, req.Name)
		if !found {
			return "", fmt.Errorf("requirement %q not found for MODIFIED", req.Name)
		}
		replacement := renderRequirement(req)
		content = content[:start] + replacement + content[end:]
	}
	return content, nil
}

func applyRemoved(content string, delta model.Delta) (string, error) {
	for _, req := range delta.Requirements {
		start, end, found := findRequirementBlock(content, req.Name)
		if !found {
			return "", fmt.Errorf("requirement %q not found for REMOVED", req.Name)
		}
		// Remove the block and any trailing blank line
		removed := content[:start] + content[end:]
		content = strings.ReplaceAll(removed, "\n\n\n", "\n\n")
	}
	return content, nil
}

func applyRenamed(content string, delta model.Delta) (string, error) {
	if delta.FromName == "" || delta.ToName == "" {
		return "", fmt.Errorf("RENAMED requires both FROM and TO names")
	}
	lines := strings.Split(content, "\n")
	found := false
	for i, line := range lines {
		m := reqHeaderRe.FindStringSubmatch(line)
		if m != nil && normalizeWhitespace(m[1]) == normalizeWhitespace(delta.FromName) {
			lines[i] = "### Requirement: " + delta.ToName
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("requirement %q not found for RENAMED", delta.FromName)
	}
	return strings.Join(lines, "\n"), nil
}

// findRequirementBlock locates a requirement block by name, returning start/end byte offsets.
func findRequirementBlock(content, name string) (int, int, bool) {
	lines := strings.Split(content, "\n")
	normalizedName := normalizeWhitespace(name)
	start := -1
	startLine := -1

	for i, line := range lines {
		m := reqHeaderRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		if normalizeWhitespace(m[1]) == normalizedName {
			startLine = i
			break
		}
	}

	if startLine < 0 {
		return 0, 0, false
	}

	// Calculate byte offset for start
	start = 0
	for i := 0; i < startLine; i++ {
		start += len(lines[i]) + 1 // +1 for newline
	}

	// Find end: next ### header or end of content
	end := len(content)
	bytePos := start
	for i := startLine; i < len(lines); i++ {
		if i > startLine {
			m := reqHeaderRe.FindStringSubmatch(lines[i])
			if m != nil {
				end = bytePos
				break
			}
			// Also stop at ## headers
			if strings.HasPrefix(lines[i], "## ") {
				end = bytePos
				break
			}
		}
		bytePos += len(lines[i]) + 1
	}

	return start, end, true
}

func renderRequirement(req model.Requirement) string {
	var b strings.Builder
	b.WriteString("### Requirement: ")
	b.WriteString(req.Name)
	b.WriteString("\n")
	if req.Text != "" {
		b.WriteString(req.Text)
		b.WriteString("\n")
	}
	for _, sc := range req.Scenarios {
		b.WriteString("\n#### Scenario: ")
		b.WriteString(sc.Name)
		b.WriteString("\n")
		if sc.Text != "" {
			b.WriteString(sc.Text)
			b.WriteString("\n")
		}
	}
	return b.String()
}

func normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

// AtomicWrite writes content to a file using write-then-rename.
func AtomicWrite(path string, content []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".openspec-tmp-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(content); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("closing temp file: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("renaming temp file: %w", err)
	}
	return nil
}

// ArchiveOptions configures the archive operation.
type ArchiveOptions struct {
	OpenSpecPath string
	ChangeID     string
	SkipSpecs    bool
	SkipValidate bool
	Yes          bool
}

// Archive performs the full archive workflow: validate → merge → move.
func Archive(opts ArchiveOptions) error {
	changePath := filepath.Join(opts.OpenSpecPath, "changes", opts.ChangeID)
	if _, err := os.Stat(changePath); os.IsNotExist(err) {
		return fmt.Errorf("change %q not found", opts.ChangeID)
	}

	// Parse the change
	proposalContent, err := os.ReadFile(filepath.Join(changePath, "proposal.md"))
	if err != nil {
		return fmt.Errorf("reading proposal: %w", err)
	}
	change := parser.ParseChange(string(proposalContent))
	change.ID = opts.ChangeID

	// Parse deltas from specs/ subdirectory
	specsDir := filepath.Join(changePath, "specs")
	if entries, err := os.ReadDir(specsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			specFile := filepath.Join(specsDir, entry.Name(), "spec.md")
			deltaContent, err := os.ReadFile(specFile)
			if err != nil {
				continue
			}
			deltas := parser.ParseDeltas(string(deltaContent))
			for i := range deltas {
				deltas[i].SpecName = entry.Name()
			}
			change.Deltas = append(change.Deltas, deltas...)
		}
	}

	// Parse tasks
	tasksContent, err := os.ReadFile(filepath.Join(changePath, "tasks.md"))
	if err == nil {
		change.Tasks = parser.ParseTaskProgress(string(tasksContent))
	}

	// Pre-archive validation
	if !opts.SkipValidate {
		issues := validator.ValidateChange(change)
		for _, iss := range issues {
			if iss.Level == model.LevelError {
				return fmt.Errorf("validation failed: %s", iss.Message)
			}
		}
	}

	// Merge deltas into specs
	if !opts.SkipSpecs {
		// Compute all merges first (all-or-nothing)
		type mergeResult struct {
			path    string
			content string
		}
		var merges []mergeResult

		for _, delta := range change.Deltas {
			specPath := filepath.Join(opts.OpenSpecPath, "specs", delta.SpecName, "spec.md")
			var specContent string
			if data, err := os.ReadFile(specPath); err == nil {
				specContent = string(data)
			} else if delta.Operation == model.DeltaAdded {
				// Create new spec
				specContent = fmt.Sprintf("# %s\n\n## Purpose\n\n## Requirements\n", delta.SpecName)
			} else {
				return fmt.Errorf("spec file not found: %s", specPath)
			}

			merged, err := MergeDeltas(specContent, []model.Delta{delta})
			if err != nil {
				return fmt.Errorf("merge failed for %s: %w", delta.SpecName, err)
			}

			merges = append(merges, mergeResult{path: specPath, content: merged})
		}

		// Post-merge validation
		if !opts.SkipValidate {
			for _, mr := range merges {
				spec := parser.ParseSpec(mr.content)
				issues := validator.ValidateSpec(spec, mr.path)
				for _, iss := range issues {
					if iss.Level == model.LevelError {
						return fmt.Errorf("post-merge validation failed for %s: %s", mr.path, iss.Message)
					}
				}
			}
		}

		// Write all merged specs atomically
		for _, mr := range merges {
			dir := filepath.Dir(mr.path)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("creating spec directory: %w", err)
			}
			if err := AtomicWrite(mr.path, []byte(mr.content)); err != nil {
				return fmt.Errorf("writing merged spec: %w", err)
			}
		}
	}

	// Move to archive
	archiveDir := filepath.Join(opts.OpenSpecPath, "changes", "archive")
	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		return fmt.Errorf("creating archive directory: %w", err)
	}
	archiveName := time.Now().Format("2006-01-02") + "-" + opts.ChangeID
	archivePath := filepath.Join(archiveDir, archiveName)
	if err := os.Rename(changePath, archivePath); err != nil {
		return fmt.Errorf("moving to archive: %w", err)
	}

	return nil
}
