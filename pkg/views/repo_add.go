package views

import (
	"go-git-tk/pkg/gitlib"

	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	newrepo_ok_key key.Binding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter/⏎", "create"),
	)
	newrepo_cancel_key key.Binding = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)
	newrepo_select_textview key.Binding = key.NewBinding(
		key.WithKeys("up", "prev input"),
		key.WithKeys("shift+tab"),
		key.WithKeys("down", "next input"),
		key.WithKeys("tab"),
	)
)

type newrepo struct {
	parent     Richmodel
	focusIndex int
	inputs     []textinput.Model
}

func NewRepoView(parent Richmodel) newrepo {
	r := newrepo{
		focusIndex: 0,
		inputs:     make([]textinput.Model, 1),
		parent:     parent,
	}

	var t textinput.Model
	for i := range r.inputs {
		t = textinput.New()
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Repository Name"
			t.Focus()
		}

		r.inputs[i] = t
	}

	return r
}

func (n newrepo) Init() tea.Cmd {
	return textinput.Blink
}

func (n newrepo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(m, newrepo_cancel_key):
			return n, ChangeView(MakeHomeList())

		case key.Matches(m, newrepo_select_textview):
			s := m.String()

			if s == "up" || s == "shift+tab" {
				if n.focusIndex < len(n.inputs) {
					n.focusIndex++
				}
			} else {
				if n.focusIndex > 0 {
					n.focusIndex--
				}
			}

			cmds := make([]tea.Cmd, len(n.inputs))
			for i := 0; i <= len(n.inputs)-1; i++ {
				if i == n.focusIndex {
					cmds[i] = n.inputs[i].Focus()
					continue
				}

				n.inputs[i].Blur()
			}
			return n, tea.Batch(cmds...)

		case key.Matches(m, newrepo_ok_key):
			// create new repository
			repoName := n.inputs[0].Value()

			if len(strings.TrimSpace(repoName)) == 0 {
				return n, nil
			}

			wd, err := os.Getwd()
			if err != nil {
				return n, ChangeView(errorview{
					Parent: n,
					Err:    err,
				})
			}

			path := wd + "/repos/" + sanitize_name(repoName) + ".git"

			repo, err := gitlib.MakeNewRepo(path)

			if err != nil {
				return n, ChangeView(errorview{
					Parent: n,
					Err:    err,
				})
			}

			return n, ChangeView(MakeRepoView(n.parent, *repo))

		default:
			return n, n.updateInputs(msg)
		}
	}

	return n, nil
}

func (n newrepo) View() string {
	var s string
	s += "New Repository\n\n"

	for i := range n.inputs {
		s += n.inputs[i].View() + "\n"
	}

	return s
}

func (m newrepo) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func sanitize_name(s string) string {
	return strings.ReplaceAll(s, " ", "-")
}

func (h newrepo) GetKeymapString() []key.Binding {
	return []key.Binding{
		newrepo_ok_key,
		newrepo_cancel_key,
		newrepo_select_textview,
	}
}
