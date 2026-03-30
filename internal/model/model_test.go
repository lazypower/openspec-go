package model

import "testing"

func TestTaskStatus_Percent(t *testing.T) {
	tests := []struct {
		total, completed, want int
	}{
		{0, 0, 0},
		{5, 3, 60},
		{10, 10, 100},
		{3, 1, 33},
	}
	for _, tt := range tests {
		ts := TaskStatus{Total: tt.total, Completed: tt.completed}
		if got := ts.Percent(); got != tt.want {
			t.Errorf("TaskStatus{%d, %d}.Percent() = %d, want %d", tt.total, tt.completed, got, tt.want)
		}
	}
}

func TestValidationReport_HasErrors(t *testing.T) {
	report := &ValidationReport{
		Items: []ValidationItem{
			{ID: "a", Valid: true},
			{ID: "b", Valid: false},
		},
	}
	if !report.HasErrors() {
		t.Error("expected HasErrors to be true")
	}

	report2 := &ValidationReport{
		Items: []ValidationItem{
			{ID: "a", Valid: true},
		},
	}
	if report2.HasErrors() {
		t.Error("expected HasErrors to be false")
	}
}
