package views

import (
	"go-git-tk/pkg/gitlib"
	"go-git-tk/pkg/layouts"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type repo_sync struct {
	parent   Richmodel
	repo     gitlib.Repo
	input    textinput.Model
	viewport viewport.Model
}

var (
	reposync_cancel key.Binding = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)
	reposync_accept key.Binding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "clone repo"),
	)
)

func OpenRepoSync(parent Richmodel, r gitlib.Repo) repo_sync {
	input := textinput.New()
	input.Focus()

	vp := viewport.New(layouts.GetContentInnerWidth(), layouts.GetContentInnerHeight()-3)
	vp.SetContent("Press Enter to sync this repo with the remote origin")
	//vp.Style = layouts.ContentStyle

	//basename := path.Base(repo.Repopath)
	//input.SetValue(config.Conf.RepoWorkdir + "/" + strings.TrimSuffix(basename, filepath.Ext(basename)))

	return repo_sync{
		parent:   parent,
		repo:     r,
		input:    input,
		viewport: vp,
	}
}

func (k repo_sync) Init() tea.Cmd {
	return nil
}

func (k repo_sync) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, reposync_cancel):
			return k, ChangeView(k.parent)

		case key.Matches(msg, reposync_accept):
			if err := k.repo.PushToOrigin(); err != nil {
				return k, ChangeView(errorview{
					Parent: k,
					Err:    err,
				})
			} else {
				return k, ChangeView(MakeRepoView(k.parent, k.repo))
			}

		}
	}

	var cmd tea.Cmd
	k.input, cmd = k.input.Update(msg)
	return k, cmd
}

func (k repo_sync) GetHeaderString() string {
	return "Sync Repo with Origin"
}

func (k repo_sync) View() string {
	return k.viewport.View()
}

func (k repo_sync) GetKeymapString() []key.Binding {
	return []key.Binding{
		reposync_accept,
		reposync_cancel,
	}
}
