package parser

import (
	"regexp"
	"strings"

	"github.com/chuck/openspec-go/internal/model"
)

var headerRe = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

// ParseSpec parses a spec markdown file into a Spec struct.
func ParseSpec(content string) model.Spec {
	spec := model.Spec{
		Sections: make(map[string]string),
	}
	lines := strings.Split(content, "\n")

	var (
		currentSection string
		currentReq     *model.Requirement
		currentScen    *model.Scenario
		contentBuf     []string
	)

	flushContent := func() {
		text := strings.TrimSpace(strings.Join(contentBuf, "\n"))
		if currentScen != nil {
			currentScen.Text = text
		} else if currentReq != nil {
			currentReq.Text = text
		} else if currentSection != "" {
			spec.Sections[currentSection] = text
			if strings.EqualFold(currentSection, "Purpose") {
				spec.Overview = text
			}
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
		if currentReq != nil {
			if len(currentReq.Scenarios) == 0 {
				flushContent()
			}
			spec.Requirements = append(spec.Requirements, *currentReq)
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
		case 1:
			flushContent()
			spec.Title = text

		case 2:
			finishRequirement()
			flushContent()
			currentSection = text

		case 3:
			finishRequirement()
			name := extractAfterColon(text, "Requirement")
			if name != "" {
				currentReq = &model.Requirement{Name: name}
			} else {
				currentReq = nil
				currentSection = text
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

	// flush remaining
	finishRequirement()
	flushContent()

	return spec
}

// extractAfterColon extracts the value after "prefix:" with whitespace normalization.
func extractAfterColon(text, prefix string) string {
	// Normalize whitespace in both prefix and text for matching
	normalized := collapseWhitespace(text)
	p := collapseWhitespace(prefix) + ":"
	if !strings.HasPrefix(normalized, p) {
		return ""
	}
	return strings.TrimSpace(normalized[len(p):])
}

func collapseWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
