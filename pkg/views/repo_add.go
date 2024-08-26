package views

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-git/go-git/v5"
)

type newrepo struct {
	parent     richmodel
	focusIndex int
	inputs     []textinput.Model
}

func NewRepoView(parent richmodel) newrepo {
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
		switch m.String() {
		case "esc":
			return n, changeView(NewHomeView())

		case "up", "shift+tab", "down", "tab":
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

		case "enter":
			// create new repository
			repoName := n.inputs[0].Value()

			wd, err := os.Getwd()
			if err != nil {
				return n, changeView(errorview{
					Parent: n,
					Err:    err,
				})
			}

			path := wd + "/repos/" + sanitize_name(repoName) + ".git"

			r, err := createGitRepo(path)

			if err != nil {
				//log.Panic(err)
				return n, changeView(errorview{
					Err:    err,
					Parent: n,
				})
			}

			repo := repo{
				git:      r,
				repopath: path,
			}

			return n, changeView(MakeRepoView(n.parent, repo))

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

func createGitRepo(repoPath string) (*git.Repository, error) {
	//path := "./repos/" + sanitize_name(repoName)

	return git.PlainInitWithOptions(repoPath, &git.PlainInitOptions{
		Bare: true,
	})
}

func (h newrepo) GetKeymapString() string {
	return "enter - create, esc - back"
}
