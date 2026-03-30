package output

import (
	"fmt"
	"io"
	"strings"
)

// Table renders a simple aligned table.
type Table struct {
	Headers []string
	Rows    [][]string
}

// Render writes the table to w.
func (t *Table) Render(w io.Writer) {
	if len(t.Rows) == 0 && len(t.Headers) == 0 {
		return
	}

	// Calculate column widths
	widths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		widths[i] = len(h)
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print headers
	if len(t.Headers) > 0 {
		for i, h := range t.Headers {
			if i > 0 {
				fmt.Fprint(w, "  ")
			}
			fmt.Fprintf(w, "%-*s", widths[i], h)
		}
		fmt.Fprintln(w)
		// Separator
		for i, width := range widths {
			if i > 0 {
				fmt.Fprint(w, "  ")
			}
			fmt.Fprint(w, strings.Repeat("─", width))
		}
		fmt.Fprintln(w)
	}

	// Print rows
	for _, row := range t.Rows {
		for i, cell := range row {
			if i > 0 {
				fmt.Fprint(w, "  ")
			}
			if i < len(widths) {
				fmt.Fprintf(w, "%-*s", widths[i], cell)
			} else {
				fmt.Fprint(w, cell)
			}
		}
		fmt.Fprintln(w)
	}
}
