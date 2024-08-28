package views

import (
	"fmt"
	"go-git-tk/pkg/gitlib"

	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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
	input      textinput.Model
	viewport   viewport.Model
}

func NewRepoView(parent Richmodel) newrepo {
	vp := viewport.New(getInnerViewportWidth(), getViewportHeight()-3)
	vp.SetContent("Enter the name of the new Repository and press enter")
	vp.Style = mainStyle

	ti := textinput.New()
	ti.CharLimit = 32

	ti.Placeholder = "Repository Name"
	ti.Focus()

	r := newrepo{
		focusIndex: 0,
		input:      ti,
		parent:     parent,
		viewport:   vp,
	}

	return r
}

func (n newrepo) Init() tea.Cmd {
	return textinput.Blink
}

func (n newrepo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		n.viewport.Width = getInnerViewportWidth()
		n.viewport.Height = getViewportHeight() - 3
	case tea.KeyMsg:
		switch {
		case key.Matches(m, newrepo_cancel_key):
			return n, ChangeView(MakeHomeList())

		case key.Matches(m, newrepo_ok_key):
			// create new repository
			repoName := n.input.Value()

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

		}
	}

	return n, n.updateInputs(msg)
}

func (n newrepo) View() string {
	var sb string

	sb += fmt.Sprintf("%s\n%s\n%s", titleStyle.Render("New Repository "), n.viewport.View(), mainStyle.Render(n.input.View()))

	return sb
}

func (m *newrepo) updateInputs(msg tea.Msg) tea.Cmd {

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	var tcmd tea.Cmd
	m.input, tcmd = m.input.Update(msg)

	var vcmd tea.Cmd
	m.viewport, vcmd = m.viewport.Update(msg)

	return tea.Batch(tcmd, vcmd)
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
