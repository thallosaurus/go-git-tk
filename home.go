package tuiplayground

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

type (
	updateViewEvent struct{}
)

func updateView() tea.Cmd {
	return func() tea.Msg {
		return updateViewEvent{}
	}
}

func actions() []Action {
	var fns []Action
	fns = append(fns, newAction)

	return fns
}

func actionsLabels() []string {
	return []string{
		"Create New Repository",
	}
}

func (h homeview) Init() tea.Cmd {
	return tea.Batch(tea.ShowCursor, updateView())
}

func (h homeview) getRowsLength() int {
	return len(actions()) + len(h.rows)
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
			if h.cursor < h.getRowsLength()-1 {
				h.cursor++
			}
			return h, nil

		case "esc":
			return h, closeView()

		case "enter":
			if h.cursor < len(actionsLabels()) {
				a := actions()
				return h, a[h.cursor]()
			} else {
				return h, openRepo(h, h.rows[h.cursor-len(actions())])
			}
		}

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

	s += titleStyle.Render("Repository Manager")

	s += "\n\n"

	actionlabels := actionsLabels()

	for i, _ := range actions() {
		var cur string
		if h.cursor == i {
			cur = "*"
		} else {
			cur = " "
		}

		s += actionStyle.Render(fmt.Sprintf("[%s] %s ", cur, actionlabels[i]))
		s += "\n"
	}
	s += "\n\n"

	for i, r := range h.rows {
		var cur string
		if h.cursor == len(actionsLabels())+i {
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
	return "enter - select repo, up/down - move, esc - close view (debug)"
}
