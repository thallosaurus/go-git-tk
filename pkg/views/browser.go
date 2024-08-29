package views

import (
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-git/go-billy/v5"
)

type repo_browser struct {
	fp     filepicker.Model
	Parent Richmodel
}

func (ev repo_browser) Init() tea.Cmd {
	return ev.fp.Init()
}

func (ev repo_browser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	ev.fp, cmd = ev.fp.Update(msg)
	return ev, cmd
}

func (ev repo_browser) GetKeymapString() []key.Binding {
	return []key.Binding{}
}
func (ev repo_browser) View() string {
	return ev.fp.View()
}

func MakeRepoBrowser(p Richmodel, repo billy.Filesystem) repo_browser {
	fp := filepicker.New()
	//fp.AllowedTypes = []string{".mod", ".sum", ".go", ".txt", ".md"}
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.Height = 2

	return repo_browser{
		Parent: p,
		fp:     fp,
	}
}
