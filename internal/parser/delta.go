package parser

import (
	"strings"

	"github.com/chuck/openspec-go/internal/model"
)

// ParseDeltas parses delta sections from a change spec file.
// Delta files contain sections like "## ADDED Requirements", "## MODIFIED Requirements", etc.
func ParseDeltas(content string) []model.Delta {
	var deltas []model.Delta
	lines := strings.Split(content, "\n")

	var (
		currentOp   model.DeltaOp
		currentReq  *model.Requirement
		currentScen *model.Scenario
		contentBuf  []string
	)

	flushContent := func() {
		text := strings.TrimSpace(strings.Join(contentBuf, "\n"))
		if currentScen != nil {
			currentScen.Text = text
		} else if currentReq != nil {
			currentReq.Text = text
		}
		contentBuf = nil
	}

	finishScenario := func() {
		if currentScen != nil && currentReq != nil {
			flushContent()
			currentReq.Scenarios = append(currentReq.Scenarios, *currentScen)
			currentScen = nil
		}
	}

	finishRequirement := func() {
		finishScenario()
		if currentReq != nil && currentOp != "" {
			// Only flush content to Text if no scenarios set it already
			if len(currentReq.Scenarios) == 0 {
				flushContent()
			}
			delta := model.Delta{
				Operation:    currentOp,
				Requirements: []model.Requirement{*currentReq},
			}
			// For RENAMED, parse FROM/TO from body text
			if currentOp == model.DeltaRenamed {
				from, to := parseRenamedPair(currentReq.Text)
				delta.FromName = from
				delta.ToName = to
			}
			deltas = append(deltas, delta)
			currentReq = nil
		}
	}

	for _, line := range lines {
		m := headerRe.FindStringSubmatch(line)
		if m == nil {
			contentBuf = append(contentBuf, line)
			continue
		}

		level := len(m[1])
		text := strings.TrimSpace(m[2])

		switch level {
		case 2:
			finishRequirement()
			op := parseDeltaHeader(text)
			if op != "" {
				currentOp = op
			} else {
				currentOp = ""
			}
			contentBuf = nil

		case 3:
			finishRequirement()
			name := extractAfterColon(text, "Requirement")
			if name != "" && currentOp != "" {
				currentReq = &model.Requirement{Name: name}
			}
			contentBuf = nil

		case 4:
			if currentScen != nil {
				finishScenario()
			} else if currentReq != nil {
				// Flush content between requirement header and first scenario as requirement text
				currentReq.Text = strings.TrimSpace(strings.Join(contentBuf, "\n"))
			}
			name := extractAfterColon(text, "Scenario")
			if name != "" && currentReq != nil {
				currentScen = &model.Scenario{Name: name}
			}
			contentBuf = nil

		default:
			contentBuf = append(contentBuf, line)
		}
	}

	finishRequirement()
	return deltas
}

// parseDeltaHeader detects "ADDED Requirements", "MODIFIED Requirements", etc.
func parseDeltaHeader(text string) model.DeltaOp {
	normalized := collapseWhitespace(strings.ToUpper(text))
	for _, op := range []model.DeltaOp{model.DeltaAdded, model.DeltaModified, model.DeltaRemoved, model.DeltaRenamed} {
		if normalized == string(op)+" REQUIREMENTS" {
			return op
		}
	}
	return ""
}

// parseRenamedPair extracts FROM/TO names from renamed requirement text.
func parseRenamedPair(text string) (from, to string) {
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "from:") || strings.HasPrefix(lower, "- from:") || strings.HasPrefix(lower, "- **from**:") || strings.HasPrefix(lower, "**from**:") {
			from = extractValue(line)
		}
		if strings.HasPrefix(lower, "to:") || strings.HasPrefix(lower, "- to:") || strings.HasPrefix(lower, "- **to**:") || strings.HasPrefix(lower, "**to**:") {
			to = extractValue(line)
		}
	}
	return
}

func extractValue(line string) string {
	idx := strings.Index(line, ":")
	if idx < 0 {
		return ""
	}
	return strings.TrimSpace(line[idx+1:])
}
