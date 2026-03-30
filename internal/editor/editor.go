package editor

import "fmt"

// Editor defines the interface all editor configurators implement.
type Editor interface {
	Name() string
	Configure(projectPath, openspecPath string) error
	UpdateExisting(projectPath, openspecPath string) error
	IsConfigured(projectPath string) bool
}

// Registry maps editor names to their configurators.
var Registry = map[string]Editor{
	"claude-code": &ClaudeCode{},
	"opencode":    &OpenCode{},
	"codex":       &Codex{},
	"goose":       &Goose{},
}

// AllNames returns all registered editor names.
func AllNames() []string {
	return []string{"claude-code", "opencode", "codex", "goose"}
}

// Get returns an editor by name.
func Get(name string) (Editor, error) {
	e, ok := Registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown editor: %s", name)
	}
	return e, nil
}
