package views

import (
	"fmt"
	"go-git-tk/pkg/gitlib"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type repoview struct {
	repo     gitlib.Repo
	parent   Richmodel
	viewport viewport.Model
}

type removeRepo struct {
	path string
}

var (
	repoview_cancel key.Binding = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)
	repoview_hook_edit key.Binding = key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "edit hooks"),
	)
	repoview_delete key.Binding = key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete repo"),
	)
	repoview_rename key.Binding = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename repo"),
	)
)

func MakeRepoView(p Richmodel, repo gitlib.Repo) repoview {
	//targets := make(map[string]func() tea.Cmd, 0)
	vp := viewport.New(term_width, term_height-2)
	vp.SetContent(getViewportContent(repo))

	return repoview{
		repo:     repo,
		parent:   p,
		viewport: vp,
	}
}

func (r repoview) Init() tea.Cmd {
	return nil
}

func (r repoview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		r.viewport.Width = m.Width
		r.viewport.Height = m.Height - 2
	case tea.KeyMsg:
		switch {
		case key.Matches(m, repoview_cancel):
			return r, ChangeView(r.parent)

		case key.Matches(m, repoview_hook_edit):
			return r, ChangeView(OpenHookEdit(r, r.repo))

		case key.Matches(m, repoview_delete):
			return r, ChangeView(ConfirmRepoRemove(r, r.repo))

		case key.Matches(m, repoview_rename):
			return r, ChangeView(OpenRepoRename(r, r.repo))
		}

	}

	var cmd tea.Cmd
	r.viewport, cmd = r.viewport.Update(msg)
	return r, cmd
}

func (r repoview) View() string {
	//return ""
	return fmt.Sprintf("%s\n%s", "Manage Repository", r.viewport.View())
}

func (r repoview) GetKeymapString() []key.Binding {
	return []key.Binding{
		repoview_cancel,
		repoview_rename,
		repoview_hook_edit,
		repoview_delete,
	}
}

func getViewportContent(repo gitlib.Repo) string {
	var s string
	s += fmt.Sprintf("Name: %s\n", repo.GetName())
	s += fmt.Sprintf("Repo URL: %s@%s:%s\n", ssh_user, ssh_base_domain, repo.GetName())
	s += "\n"

	branches, err := repo.GetBranches()
	if err != nil {
		log.Panic(err)
	}

	s += "Branches:"
	if len(branches) != 0 {

		for _, b := range branches {
			s += fmt.Sprintf("\n- %s\n", b)
		}
	} else {
		s += emptyStyle.Render(" <empty>") + "\n"
	}
	//	s += viewIter(branches)

	s += "Tags:"
	tags, err := repo.GetTags()
	if err != nil {
		log.Panic(err)
	}
	if len(tags) != 0 {

		for _, t := range tags {
			s += fmt.Sprintf("\n- %s\n", t)
		}
	} else {
		s += emptyStyle.Render(" <empty>") + "\n"
	}

	c, err := repo.GetCommitters()
	if err != nil {
		log.Panic(err)
	}

	//s += "nl"
	s += "Committers:"

	for _, email := range c {
		s += fmt.Sprintf("\n - %s", email)
	}

	if len(c) == 0 {
		s += emptyStyle.Render(" <empty>") + "\n"
	} else {

	}
	s += "\n"

	return s
}
