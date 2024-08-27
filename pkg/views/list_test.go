package views

import tea "github.com/charmbracelet/bubbletea"

type ListTest struct {
}

func (l ListTest) Init() tea.Cmd {
	return nil
}
func (l ListTest) View() string {
	return ""
}
func (l ListTest) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return l, nil
}
func (l ListTest) GetKeymapString() string {
	return ""
}

func MakeListTest() ListTest {
	return ListTest{}
}
