package template

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	tmpl, err := Get("agents.md")
	if err != nil {
		t.Fatal(err)
	}
	if tmpl == nil {
		t.Error("template is nil")
	}
}

func TestGet_NotFound(t *testing.T) {
	_, err := Get("nonexistent.md")
	if err == nil {
		t.Error("expected error for nonexistent template")
	}
}

func TestRender(t *testing.T) {
	result, err := Render("agents.md", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "OpenSpec") {
		t.Error("rendered template missing expected content")
	}
}

func TestMustRender(t *testing.T) {
	result := MustRender("agents.md", nil)
	if !strings.Contains(result, "OpenSpec") {
		t.Error("rendered template missing expected content")
	}
}

func TestMustRender_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nonexistent template")
		}
	}()
	MustRender("nonexistent.md", nil)
}

func TestRaw(t *testing.T) {
	content, err := Raw("agents.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(content, "OpenSpec") {
		t.Error("raw template missing expected content")
	}
}

func TestRaw_NotFound(t *testing.T) {
	_, err := Raw("nonexistent.md")
	if err == nil {
		t.Error("expected error for nonexistent template")
	}
}

func TestAllTemplatesExist(t *testing.T) {
	templates := []string{"agents.md", "project.md", "claude.md", "proposal.md", "apply.md", "archive.md"}
	for _, name := range templates {
		if _, err := Raw(name); err != nil {
			t.Errorf("template %s not found: %v", name, err)
		}
	}
}
