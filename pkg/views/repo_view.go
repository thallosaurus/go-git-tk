package views

import (
	"fmt"
	"go-git-tk/pkg/config"
	"go-git-tk/pkg/gitlib"
	"go-git-tk/pkg/layouts"
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

	full_tags bool
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
	repoview_clone key.Binding = key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "clone to new workdir"),
	)
	repoview_browse key.Binding = key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "browse repo"),
		key.WithDisabled(),
	)
	repoview_push_origin key.Binding = key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "push changes to origin"),
	)
	repoview_toggle_tag_collapse key.Binding = key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("ctrl+t", "toggle tag collapse"),
	)
)

func MakeRepoView(p Richmodel, repo gitlib.Repo) repoview {
	//targets := make(map[string]func() tea.Cmd, 0)
	vp := viewport.New(layouts.GetContentInnerWidth(), layouts.GetContentInnerHeight())
	//vp.Style = layouts.MainStyle

	return repoview{
		repo:      repo,
		parent:    p,
		viewport:  vp,
		full_tags: true,
	}
}

func (r repoview) GetHeaderString() string {
	return "Manage Repository"
}

func (r repoview) Init() tea.Cmd {
	r.viewport.Width = layouts.GetContentInnerWidth()
	r.viewport.Height = layouts.GetContentInnerHeight()

	o, _ := r.repo.HasRemoteOrigin()
	repoview_push_origin.SetEnabled(o)
	/*if r.repo.IsRepoEmpty() {
		repoview_clone.SetEnabled(false)
	}*/
	return nil
}

func (r repoview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		r.viewport.Width = layouts.GetContentInnerWidth()
		r.viewport.Height = layouts.GetContentInnerHeight()
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

		case key.Matches(m, repoview_clone):
			return r, ChangeView(OpenRepoClone(r, r.repo))

		case key.Matches(m, repoview_push_origin):
			return r, ChangeView(OpenRepoSync(r, r.repo))

		case key.Matches(m, repoview_toggle_tag_collapse):
			r.full_tags = !r.full_tags
			return r, nil

		}

	}

	r.viewport.SetContent(getViewportContent(r.repo, r))

	var cmd tea.Cmd
	r.viewport, cmd = r.viewport.Update(msg)
	return r, cmd
}

func (r repoview) View() string {
	//return ""

	return r.viewport.View()
}

func (r repoview) GetKeymapString() []key.Binding {
	return []key.Binding{
		repoview_cancel,
		repoview_rename,
		repoview_hook_edit,
		repoview_delete,
		repoview_clone,
		repoview_push_origin,
	}
}

func listToView(branches []string, full bool) string {
	var s string
	if !full {
		var lIndex int
		if len(branches) > 5 {
			lIndex = 5
		} else {
			lIndex = len(branches)
		}
		branches = branches[0:lIndex]
	}

	if len(branches) != 0 {
		for _, b := range branches {
			s += fmt.Sprintf("- %s\n", b)
		}
	} else {
		s += layouts.EmptyStyle.Render(" <empty>") + "\n"
	}

	return s
}

func getRepoURL(repo gitlib.Repo) string {
	var url string
	if config.Conf.EnableSSH {
		var name string
		if config.Conf.ShowFullRepoPath {
			name = repo.Repopath
		} else {
			name = repo.GetName()
		}
		url = fmt.Sprintf("%s@%s:%s", config.Conf.SSHUser, config.Conf.SSHBaseDomain, name)
	} else {
		url = repo.Repopath
	}

	return url
}

func getViewportContent(repo gitlib.Repo, view repoview) string {
	var s string
	s += fmt.Sprintf("%s %s\n", layouts.ListStyle.Render("Name:"), repo.GetName())

	s += fmt.Sprintf("%s %s\n", layouts.ListStyle.Render("Repo URL:"), getRepoURL(repo))
	s += "\n"

	branches, err := repo.GetBranches()
	if err != nil {
		log.Panic(err)
	}

	s += layouts.ListStyle.Render("Branches:") + "\n"

	s += listToView(branches, false)
	s += "\n"

	s += layouts.ListStyle.Render("Tags:") + "\n"
	tags, err := repo.GetTags()
	if err != nil {
		log.Panic(err)
	}

	s += listToView(tags, view.full_tags)
	//s += "\n"

	c, err := repo.GetCommitters()
	if err != nil {
		log.Panic(err)
	}
	s += "\n"

	//s += "nl"
	s += layouts.ListStyle.Render("Committers:") + "\n"
	s += listToView(c, false)
	s += "\n"

	s += layouts.ListStyle.Render("Readme:") + "\n"

	rr, err := repo.GetReadme()
	if err != nil {
		//log.Panic(err)
		s += layouts.EmptyStyle.Render(" No readme published")
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
