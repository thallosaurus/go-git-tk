package views

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type repoview struct {
	repo          repo
	cursor        int
	cursorTargets map[string]func() tea.Cmd
	parent        richmodel
}

type removeRepo struct {
	path string
}

func MakeRepoView(p richmodel, repo repo) repoview {
	targets := make(map[string]func() tea.Cmd, 0)

	return repoview{
		repo:          repo,
		cursor:        0,
		cursorTargets: targets,
		parent:        p,
	}
}

func (r repoview) Init() tea.Cmd {
	return nil
}

func (r repoview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "esc":
			return r, changeView(r.parent)

		case "h":
			return r, changeView(OpenHookEdit(r, r.repo))

		case "d":
			return r, changeView(ConfirmRepoRemove(r, r.repo))

		case "r":
			return r, changeView(OpenRepoRename(r, r.repo))
		}

	}
	return r, nil
}

func (r repoview) View() string {
	var s string

	s += "Manage Repository\n\n"

	s += fmt.Sprintf("Name: %s\n", r.repo.GetName())
	branches, err := r.repo.git.Branches()
	if err != nil {
		log.Panic(err)
	}
	s += "\n"

	s += fmt.Sprintf("Repo URL: %s@%s:%s\n", ssh_user, ssh_base_domain, r.repo.GetName())

	s += "Branches:"
	s += viewIter(branches)

	s += "Tags:"
	tags, err := r.repo.git.Tags()
	if err != nil {
		log.Panic(err)
	}
	s += viewIter(tags)

	c, err := r.repo.git.CommitObjects()
	if err != nil {
		log.Panic(err)
	}

	//s += "nl"
	s += "Committers:"
	s += viewCommitters(c)

	s += "\n"

	return s
}

func (r repoview) GetKeymapString() string {
	return "h - edit hooks, r - rename repo, d - delete repo, esc - back"
}

func viewIter(b storer.ReferenceIter) string {
	var sb strings.Builder
	b.ForEach(func(r *plumbing.Reference) error {
		sb.WriteString(fmt.Sprintf("\n- %s\n", string(r.Name())))

		return nil
	})

	if sb.Len() == 0 {
		sb.WriteString(emptyStyle.Render(" <empty>") + "\n")
	}

	return sb.String()
}

func viewCommitters(c object.CommitIter) string {
	var sb strings.Builder
	committers := make(map[string]string)
	c.ForEach(func(c *object.Commit) error {
		committers[c.Author.Name] = c.Author.Email
		return nil
	})

	i := 0
	for author, email := range committers {
		sb.WriteString(fmt.Sprintf("\n - %s <%s>", author, email))
		i++
	}

	if i == 0 {
		sb.WriteString(emptyStyle.Render(" <empty>") + "\n")
	}

	return sb.String()
}
