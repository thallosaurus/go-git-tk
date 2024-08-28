package views

import (
	"fmt"
	"go-git-tk/pkg/gitlib"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
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
	vp := viewport.New(term_width, getViewportHeight())
	vp.SetContent(getViewportContent(repo))
	vp.Style = mainStyle

	return repoview{
		repo:     repo,
		parent:   p,
		viewport: vp,
	}
}

func (r repoview) Init() tea.Cmd {
	r.viewport.Width = getInnerViewportWidth()
	r.viewport.Height = getViewportHeight()
	return nil
}

func (r repoview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		r.viewport.Width = getInnerViewportWidth()
		r.viewport.Height = getViewportHeight()
	case tea.KeyMsg:
		switch {
		case key.Matches(m, repoview_cancel):
			return r, ChangeView(r.parent)

		case key.Matches(m, repoview_hook_edit):
			return r, ChangeView(MakeHookList(r, r.repo))

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

	return fmt.Sprintf("%s\n%s", titleStyle.Render("Manage Repository"), r.viewport.View())
}

func (r repoview) GetKeymapString() []key.Binding {
	return []key.Binding{
		repoview_cancel,
		repoview_rename,
		repoview_hook_edit,
		repoview_delete,
	}
}

func listToView(branches []string) string {
	var s string
	if len(branches) != 0 {

		for _, b := range branches {
			s += fmt.Sprintf("- %s\n", b)
		}
	} else {
		s += emptyStyle.Render(" <empty>") + "\n"
	}

	return s
}

func getViewportContent(repo gitlib.Repo) string {
	var s string
	s += fmt.Sprintf("%s %s\n", selectedStyle.Render("Name:"), repo.GetName())
	s += fmt.Sprintf("%s %s@%s:%s\n", selectedStyle.Render("Repo URL:"), ssh_user, ssh_base_domain, repo.GetName())
	s += "\n"

	branches, err := repo.GetBranches()
	if err != nil {
		log.Panic(err)
	}

	s += selectedStyle.Render("Branches:") + "\n"
	s += listToView(branches)
	s += "\n"

	s += selectedStyle.Render("Tags:") + "\n"
	tags, err := repo.GetTags()
	if err != nil {
		log.Panic(err)
	}
	s += listToView(tags)
	s += "\n"

	c, err := repo.GetCommitters()
	if err != nil {
		log.Panic(err)
	}
	s += "\n"

	//s += "nl"
	s += selectedStyle.Render("Committers:") + "\n"
	s += listToView(c)
	s += "\n"

	s += selectedStyle.Render("Readme:") + "\n"

	rr, err := repo.GetReadme()
	if err != nil {
		//log.Panic(err)
		s += emptyStyle.Render("No readme published")
	} else {
		s += rr
	}

	return s
}

func renderReadme(md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	ast := markdown.Parse([]byte(md), p)

	c := ast.AsContainer()

	mdc := c.GetChildren()

	var sb strings.Builder
	for _, v := range mdc {
		sb.Write(v.AsContainer().Content)
	}

	return sb.String()
}
