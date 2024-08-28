package views

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"go-git-tk/pkg/gitlib"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const hookChmod = 0755

type Action func() tea.Cmd

type hookedit struct {
	parent Richmodel
	cursor int
	repo   gitlib.Repo
	action []Action
}

type (
	openeditor struct {
		path string
	}
)

var (
	hookedit_cancel key.Binding = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)
	hookedit_accept key.Binding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "remove repo"),
	)
	hookedit_move key.Binding = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "move up"),
		key.WithKeys("down"),
		key.WithHelp("down", "move down"),
	)
)

func hookLabels() []string {
	return []string{
		"Post-Receive",
		"Pre-Receive",
		"Pre-Commit",
	}
}

func (h hookedit) Init() tea.Cmd {
	return nil
}

func (h hookedit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(m, hookedit_cancel):
			return h, ChangeView(h.parent)

		case key.Matches(m, hookedit_move):
			if m.String() == "up" {
				if h.cursor > 0 {
					h.cursor--
				}
			} else {
				if h.cursor < len(h.action)-1 {
					h.cursor++
				}
			}

		case key.Matches(m, hookedit_accept):
			return h, h.action[h.cursor]()
		}

	case openeditor:
		if err := OpenEditor(m.path); err != nil {
			return h, ChangeView(errorview{
				Parent: h,
				Err:    err,
			})
		}

		return h, tea.Batch(tea.ShowCursor, ChangeView(h.parent))
	}
	return h, nil
}

func (h hookedit) View() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Edit Hooks for Repo %s\n\n", h.repo.GetName()))

	sb.WriteString("Select the hook you wish to edit. If the file doesn't exist it gets created automatically.\n\n")

	i := 0
	for _, label := range hookLabels() {
		var c string

		if i == h.cursor {
			c = ">"
		} else {
			c = " "
		}

		sb.WriteString(fmt.Sprintf("%s %s\n", c, label))

		i++
	}

	return sb.String()
}

func OpenHookEdit(parent Richmodel, r gitlib.Repo) Richmodel {
	var a []Action

	a = append(a, func() tea.Cmd {
		return func() tea.Msg {
			return openeditor{
				path: r.Repopath + "/hooks/post-receive",
			}
		}
	})
	a = append(a, func() tea.Cmd {
		return func() tea.Msg {
			return openeditor{
				path: r.Repopath + "/hooks/pre-receive",
			}
		}
	})
	a = append(a, func() tea.Cmd {
		return func() tea.Msg {
			return openeditor{
				path: r.Repopath + "/hooks/pre-commit",
			}
		}
	})

	return hookedit{
		parent: parent,
		repo:   r,
		cursor: 0,
		action: a,
	}
}

func OpenEditor(hook string) error {

	// setup hooks folder
	if !pathExists(path.Dir(hook)) {
		// path/to/whatever does not exist
		err := os.Mkdir(path.Dir(hook), 0755)
		if err != nil {
			return err
		}
	}

	var editor string
	editor = os.Getenv("EDITOR")

	if editor == "" {
		editor = "nano"
	}

	//hookname := fmt.Sprintf("%s/hooks/%s", repo.repopath, hook)
	cmd := exec.Command(editor, hook)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		return err
	} else {
		if pathExists(hook) {
			return os.Chmod(hook, hookChmod)
		} else {
			return nil
		}
	}
}

func SetupRepos(wd string) {
	if !pathExists(wd + "/repos") {
		if err := os.Mkdir(wd+"/repos", 0755); err != nil {
			log.Fatal(err)
		}
	}
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return false
	} else {
		return true
	}
}

func (h hookedit) GetKeymapString() []key.Binding {
	return []key.Binding{}
}
