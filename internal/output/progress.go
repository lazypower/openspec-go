package output

import "fmt"

// FormatProgress returns a formatted progress string like "3/5 (60%)".
func FormatProgress(completed, total int) string {
	if total == 0 {
		return "0/0 (0%)"
	}
	pct := completed * 100 / total
	return fmt.Sprintf("%d/%d (%d%%)", completed, total, pct)
}
