package views

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type homeview struct {
	cursor int
	rows   []repo
}

type richmodel interface {
	tea.Model
	GetKeymapString() string
}

type Action func() tea.Cmd

func NewHomeView() richmodel {
	return homeview{
		cursor: 0,
		rows:   make([]repo, 0),
	}
}

func newAction() tea.Cmd {
	return changeView(NewRepoView(NewHomeView()))
}

type open_key_mgmt struct{}

func openKeyMgmt() tea.Cmd {
	return func() tea.Msg {
		return open_key_mgmt{}
	}

}

type (
	updateViewEvent struct{}
)

func updateView() tea.Cmd {
	return func() tea.Msg {
		return updateViewEvent{}
	}
}

func (h homeview) Init() tea.Cmd {
	return tea.Batch(tea.ShowCursor, updateView())
}

func openRepo(p richmodel, r repo) tea.Cmd {
	return changeView(MakeRepoView(p, r))
}

func (h homeview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "up":
			if h.cursor > 0 {
				h.cursor--
			}
			return h, nil

		case "down":
			if h.cursor < len(h.rows)-1 {
				h.cursor++
			}
			return h, nil

		case "enter":
			return h, openRepo(h, h.rows[h.cursor])

		case "n":
			return h, newAction()

		case "k":
			return h, openKeyMgmt()
		}

	case open_key_mgmt:
		home := os.Getenv("HOME")

		if err := openEditor(home + "/.ssh/authorized_keys"); err != nil {
			return h, changeView(errorview{
				Parent: h,
				Err:    err,
			})
		}

		return h, tea.Batch(tea.ShowCursor, changeView(h))

	case updateViewEvent:
		wd, err := os.Getwd()
		if err != nil {
			return h, changeView(errorview{
				Parent: h,
				Err:    err,
			})
		}

		if !pathExists(wd + "/repos") {
			if err = os.Mkdir(wd+"/repos", 0755); err != nil {
				return h, changeView(errorview{
					Parent: h,
					Err:    err,
				})
			}
		}

		r, err := GetRepos(wd)

		if err != nil {
			return h, changeView(errorview{
				Parent: h,
				Err:    err,
			})
		}

		h.rows = r
		return h, nil
	}
	return h, nil
}

func (h homeview) View() string {
	var s string

	s += titleStyle.Render("Git Server Toolkit")
	s += "\n\n"

	s += titleStyle.Render("Select the repository you want to manage")
	s += "\n\n"

	for i, r := range h.rows {
		var cur string
		if h.cursor == i {
			cur = "*"
		} else {
			cur = " "
		}

		s += selectedStyle.Render(fmt.Sprintf("[%s] %s", cur, r.GetName()))
		s += "\n"

	}

	s += "\n\n"

	return s
}

func (h homeview) GetKeymapString() string {
	return "enter - select repo, up/down - move, k - access keys, n - create repository"
}
