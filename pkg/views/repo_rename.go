package views

import (
	"fmt"
	"go-git-tk/pkg/gitlib"
	"go-git-tk/pkg/layouts"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	reporename_cancel key.Binding = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)
	reporename_accept key.Binding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "remove repo"),
	)
)

type reporename struct {
	parent   Richmodel
	repo     gitlib.Repo
	input    textinput.Model
	viewport viewport.Model
}

type rename_event struct {
	name string
}

func rename(name string) tea.Cmd {
	return func() tea.Msg {
		return rename_event{
			name,
		}
	}
}

func (rr reporename) Init() tea.Cmd {
	return nil
}

func (r reporename) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "esc":
			return r, ChangeView(r.parent)

		case "enter":
			//if r.input.Value() == path.Base(r.repo.repopath) {
			//return r, finalRemoveRepo(r.repo.repopath)
			//}
			return r, rename(r.input.Value())

		default:
			v, c := r.input.Update(msg)
			r.input = v
			return r, c

		}

	case rename_event:
		p := path.Dir(r.repo.Repopath)
		newPath := fmt.Sprintf("%s/%s.git", p, m.name)
		err := os.Rename(r.repo.Repopath, newPath)
		if err != nil {
			return r, ChangeView(errorview{
				Parent: r,
				Err:    err,
			})
		}
		r.repo.Repopath = newPath
		return r, ChangeView(MakeHomeList())
	}
	return r, nil
}

func (rr reporename) GetHeaderString() string {
	return "Rename Repository"
}

func (rr reporename) View() string {
	/*s := "Rename Repository\n\n"

	s += "Don't forget to update the Remote URL!\n\n"

	s += "Name:\n"
	s += rr.input.View()*/

	//return s
	return fmt.Sprintf("%s\n%s", rr.viewport.View(), rr.input.View())
}

func (rr reporename) GetKeymapString() []key.Binding {
	return []key.Binding{
		reporename_accept,
		reporename_cancel,
	}
}

func OpenRepoRename(parent Richmodel, repo gitlib.Repo) reporename {
	input := textinput.New()
	input.Focus()

	vp := viewport.New(layouts.GetContentInnerWidth(), layouts.GetContentInnerHeight()-3)
	vp.SetContent("Enter the new name of the new Repository and press enter.\n\nDont forget to update the Remote URL!")
	//vp.Style = layouts.ContentStyle

	basename := path.Base(repo.Repopath)
	input.SetValue(strings.TrimSuffix(basename, filepath.Ext(basename)))

	return reporename{
		parent:   parent,
		repo:     repo,
		input:    input,
		viewport: vp,
	}
}
