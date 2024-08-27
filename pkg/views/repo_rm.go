package views

import (
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thallosaurus/go-git-tk/pkg/gitlib"
)

type reporemove struct {
	parent  richmodel
	repo    gitlib.Repo
	confirm textinput.Model
}

func (r reporemove) Init() tea.Cmd {
	return nil
}

func finalRemoveRepo(path string) tea.Cmd {
	return func() tea.Msg {
		return removeRepo{
			path,
		}
	}
}

func (r reporemove) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "esc":
			return r, changeView(r.parent)

		case "enter":
			if r.confirm.Value() == path.Base(r.repo.Repopath) {
				return r, finalRemoveRepo(r.repo.Repopath)
			}

		default:
			v, c := r.confirm.Update(msg)
			r.confirm = v
			return r, c

		}

	case removeRepo:
		err := os.RemoveAll(m.path)

		if err != nil {
			return r, changeView(errorview{
				Parent: r,
				Err:    err,
			})
		} else {
			return r, changeView(NewHomeView())
		}
	}

	return r, nil
}

func (r reporemove) View() string {
	var sb strings.Builder

	sb.WriteString("Are you sure to remove the Repo?\n\n\n")
	sb.WriteString(dangerousStyle.Render("You will lose everything in this repository! No undo!"))
	sb.WriteString("\n\n\n")
	sb.WriteString("Type the name of the repository below.\n\n")

	sb.WriteString(r.confirm.View())

	return sb.String()
}

func (r reporemove) GetKeymapString() string {
	return "enter - yes, esc - no"
}

func ConfirmRepoRemove(parent richmodel, repo gitlib.Repo) reporemove {
	t := textinput.New()
	t.Focus()
	t.CharLimit = 32
	t.Placeholder = path.Base(repo.Repopath)

	return reporemove{
		parent:  parent,
		repo:    repo,
		confirm: t,
	}
}
