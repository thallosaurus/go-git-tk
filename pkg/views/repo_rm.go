package views

import (
	"go-git-tk/pkg/gitlib"
	"go-git-tk/pkg/layouts"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	reporm_cancel key.Binding = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)
	reporm_accept key.Binding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "remove repo"),
	)
)

type reporemove struct {
	parent   Richmodel
	repo     gitlib.Repo
	confirm  textinput.Model
	viewport viewport.Model
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
	case tea.WindowSizeMsg:
		r.viewport.Width = layouts.GetContentInnerWidth()
		r.viewport.Height = layouts.GetContentInnerHeight() - 3
	case tea.KeyMsg:
		switch {
		case key.Matches(m, reporm_cancel):
			return r, ChangeView(r.parent)

		case key.Matches(m, reporm_accept):
			if r.confirm.Value() == path.Base(r.repo.Repopath) {
				return r, finalRemoveRepo(r.repo.Repopath)
			}

		}

	case removeRepo:
		err := os.RemoveAll(m.path)

		if err != nil {
			return r, ChangeView(errorview{
				Parent: r,
				Err:    err,
			})
		} else {
			return r, ChangeView(MakeHomeList())
		}
	}

	var tcmd tea.Cmd
	r.confirm, tcmd = r.confirm.Update(msg)

	var vcmd tea.Cmd
	r.viewport, vcmd = r.viewport.Update(msg)

	return r, tea.Batch(tcmd, vcmd)
}

func getRepoRmViewportContent() string {
	var sb strings.Builder

	sb.WriteString("Are you sure to remove the Repo?\n\n\n")
	sb.WriteString(layouts.DangerousStyle.Render("You will lose everything in this repository! No undo!"))
	sb.WriteString("\n\n\n")
	sb.WriteString("Type the name of the repository below and press enter.\n\n")
	return sb.String()
}

func (r reporemove) GetHeaderString() string {
	return "Remove Repository " + r.repo.GetName()
}

func (r reporemove) View() string {
	var sb string

	//sb += fmt.Sprintf("%s\n%s", r.viewport.View(), r.confirm.View())
	sb += r.viewport.View()
	sb += "\n"
	sb += r.confirm.View()
	return sb
}

func (r reporemove) GetKeymapString() []key.Binding {
	return []key.Binding{
		reporm_accept,
		reporm_cancel,
	}
}

func ConfirmRepoRemove(parent Richmodel, repo gitlib.Repo) reporemove {
	vp := viewport.New(layouts.GetContentInnerWidth(), layouts.GetContentInnerHeight()-3)
	vp.SetContent(getRepoRmViewportContent())
	//vp.Style = layouts.ContentStyle

	t := textinput.New()
	t.Focus()
	t.CharLimit = 32
	t.Placeholder = path.Base(repo.Repopath)

	return reporemove{
		parent:   parent,
		repo:     repo,
		confirm:  t,
		viewport: vp,
	}
}
