package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestProgressBar(t *testing.T) {
	tests := []struct {
		pct, width int
		wantFilled int
		wantEmpty  int
	}{
		{0, 10, 0, 10},
		{50, 10, 5, 5},
		{100, 10, 10, 0},
		{75, 20, 15, 5},
	}

	for _, tt := range tests {
		bar := ProgressBar(tt.pct, tt.width)
		filled := strings.Count(bar, "█")
		empty := strings.Count(bar, "░")
		if filled != tt.wantFilled || empty != tt.wantEmpty {
			t.Errorf("ProgressBar(%d, %d) = %d filled + %d empty, want %d + %d",
				tt.pct, tt.width, filled, empty, tt.wantFilled, tt.wantEmpty)
		}
	}
}

func TestFormatProgress(t *testing.T) {
	if got := FormatProgress(3, 5); got != "3/5 (60%)" {
		t.Errorf("FormatProgress(3, 5) = %q", got)
	}
	if got := FormatProgress(0, 0); got != "0/0 (0%)" {
		t.Errorf("FormatProgress(0, 0) = %q", got)
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"key": "value"}
	if err := WriteJSON(&buf, data); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"key": "value"`) {
		t.Errorf("unexpected JSON: %s", buf.String())
	}
}

func TestSetNoColor(t *testing.T) {
	SetNoColor(true)
	if !NoColor {
		t.Error("NoColor should be true")
	}
	// Rendered text should have no ANSI codes
	result := Cyan("test")
	if strings.Contains(result, "\033[") {
		t.Error("expected no ANSI codes with NoColor=true")
	}

	SetNoColor(false)
	if NoColor {
		t.Error("NoColor should be false")
	}
}

func TestColorFunctions(t *testing.T) {
	// Just make sure they don't panic
	_ = Cyan("test")
	_ = Yellow("test")
	_ = Green("test")
	_ = Magenta("test")
	_ = Dim("test")
	_ = Bold("test")
}

func TestTable_Render(t *testing.T) {
	tbl := &Table{
		Headers: []string{"Name", "Count"},
		Rows: [][]string{
			{"alpha", "10"},
			{"beta", "20"},
		},
	}
	var buf bytes.Buffer
	tbl.Render(&buf)
	out := buf.String()
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "20") {
		t.Errorf("unexpected table output: %s", out)
	}
	if !strings.Contains(out, "─") {
		t.Error("missing separator line")
	}
}

func TestTable_Empty(t *testing.T) {
	tbl := &Table{}
	var buf bytes.Buffer
	tbl.Render(&buf)
	if buf.Len() != 0 {
		t.Error("empty table should produce no output")
	}
}
