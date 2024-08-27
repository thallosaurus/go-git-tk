package views

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thallosaurus/go-git-tk/pkg/gitlib"
)

type reporename struct {
	parent richmodel
	repo   gitlib.Repo
	input  textinput.Model
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
			return r, changeView(r.parent)

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
		err := os.Rename(r.repo.Repopath, fmt.Sprintf("%s/%s.git", p, m.name))
		if err != nil {
			return r, changeView(errorview{
				Parent: r,
				Err:    err,
			})
		}
		return r, changeView(NewHomeView())
	}
	return r, nil
}

func (rr reporename) View() string {
	s := "Rename Repository\n\n"

	s += "Don't forget to update the Remote URL!\n\n"

	s += "Name:\n"
	s += rr.input.View()

	return s
}

func (rr reporename) GetKeymapString() string {
	return "enter - confirm, esc - back"
}

func OpenRepoRename(parent richmodel, repo gitlib.Repo) reporename {
	input := textinput.New()
	input.Focus()

	basename := path.Base(repo.Repopath)

	input.SetValue(strings.TrimSuffix(basename, filepath.Ext(basename)))

	return reporename{
		parent: parent,
		repo:   repo,
		input:  input,
	}
}
