package views

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type repo_import struct{}

func (k repo_import) Init() tea.Cmd {
	return nil
}

func (k repo_import) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return k, nil
}

func (k repo_import) View() string {
	return "Import Repo"
}

func (k repo_import) GetKeymapString() []key.Binding {
	return []key.Binding{}
}
