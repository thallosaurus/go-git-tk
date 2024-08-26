package views

import tea "github.com/charmbracelet/bubbletea"

type key_mgmt struct{}

func (k key_mgmt) Init() tea.Cmd {
	return nil
}

func (k key_mgmt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return k, nil
}

func (k key_mgmt) View() string {
	return "key manager"
}

func (k key_mgmt) GetKeymapString() string {
	return ""
}
