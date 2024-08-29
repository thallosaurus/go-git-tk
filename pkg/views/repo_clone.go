package views

import (
	"errors"
	"fmt"
	"go-git-tk/pkg/config"
	"go-git-tk/pkg/gitlib"
	"go-git-tk/pkg/layouts"
	"path"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type repoclone struct {
	parent   Richmodel
	repo     gitlib.Repo
	input    textinput.Model
	viewport viewport.Model
}

var (
	repoclone_cancel key.Binding = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)
	repoclone_accept key.Binding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "clone repo"),
	)
)

func OpenRepoClone(parent Richmodel, repo gitlib.Repo) repoclone {
	input := textinput.New()
	input.Focus()

	vp := viewport.New(layouts.GetContentInnerWidth(), layouts.GetContentInnerHeight()-3)
	vp.SetContent("Enter the path of the new Repository and press enter.")
	//vp.Style = layouts.ContentStyle

	basename := path.Base(repo.Repopath)
	input.SetValue(config.Conf.RepoWorkdir + "/" + strings.TrimSuffix(basename, filepath.Ext(basename)))

	return repoclone{
		parent:   parent,
		repo:     repo,
		input:    input,
		viewport: vp,
	}
}

func (rr repoclone) Init() tea.Cmd {
	if config.Conf.RepoWorkdir == "" {
		return ChangeView(errorview{
			Parent: rr.parent,
			Err:    errors.New("RepoWorkdir config option is invalid"),
		})
	} else {
		return nil
	}
}

func (r repoclone) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, repoclone_cancel):
			return r, ChangeView(r.parent)

		case key.Matches(msg, repoclone_accept):
			_, err := r.repo.CloneToWorkdir()

			if err != nil {
				return r, ChangeView(errorview{
					Parent: r,
					Err:    err,
				})
			} else {
				r.viewport.SetYOffset(0)
				return r, ChangeView(r.parent)
			}

		}
	}

	var cmd tea.Cmd
	r.input, cmd = r.input.Update(msg)
	return r, cmd
}

func (rr repoclone) GetHeaderString() string {
	return "Clone Repository to Workdir"
}

func (rr repoclone) View() string {
	return fmt.Sprintf("%s\n%s", rr.viewport.View(), rr.input.View())
}

func (rr repoclone) GetKeymapString() []key.Binding {
	return []key.Binding{
		//	reporename_accept,
		//	reporename_cancel,
		repoclone_accept,
		repoclone_cancel,
	}
}
