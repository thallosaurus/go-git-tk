package views

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thallosaurus/go-git-tk/pkg/gitlib"
)

type repoview struct {
	repo          gitlib.Repo
	cursor        int
	cursorTargets map[string]func() tea.Cmd
	parent        richmodel
}

type removeRepo struct {
	path string
}

func MakeRepoView(p richmodel, repo gitlib.Repo) repoview {
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

		case "l":
			return r, changeView(MakeListTest())
		}

	}
	return r, nil
}

func (r repoview) View() string {
	var s string

	s += "Manage Repository\n\n"

	s += fmt.Sprintf("Name: %s\n", r.repo.GetName())
	s += fmt.Sprintf("Repo URL: %s@%s:%s\n", ssh_user, ssh_base_domain, r.repo.GetName())
	s += "\n"

	branches, err := r.repo.GetBranches()
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
	tags, err := r.repo.GetTags()
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

	c, err := r.repo.GetCommitters()
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

func (r repoview) GetKeymapString() string {
	return "h - edit hooks, r - rename repo, d - delete repo, esc - back"
}
