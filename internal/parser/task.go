package parser

import (
	"regexp"
	"strings"

	"github.com/chuck/openspec-go/internal/model"
)

var taskCheckboxRe = regexp.MustCompile(`^[-*]\s+\[([ xX])\]`)

// ParseTaskProgress counts completed and total tasks from a tasks.md file.
func ParseTaskProgress(content string) model.TaskStatus {
	var status model.TaskStatus
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		m := taskCheckboxRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		status.Total++
		if strings.EqualFold(m[1], "x") {
			status.Completed++
		}
	}
	return status
}
