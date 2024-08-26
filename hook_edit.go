package tuiplayground

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const hookChmod = 0755

type hookedit struct {
	parent richmodel
	cursor int
	repo   repo
	action []Action
}

type (
	openeditor struct {
		path string
	}
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
		switch m := m.String(); m {
		case "esc":
			return h, changeView(h.parent)

		case "up", "down":
			if m == "up" {
				if h.cursor > 0 {
					h.cursor--
				}
			} else {
				if h.cursor < len(h.action)-1 {
					h.cursor++
				}
			}

		case "enter":
			return h, h.action[h.cursor]()
		}

	case openeditor:
		if err := openEditor(m.path); err != nil {
			return h, changeView(errorview{
				Parent: h,
				Err:    err,
			})
		}

		return h, tea.Batch(tea.ShowCursor, changeView(h.parent))
	}
	return h, nil
}

func (h hookedit) View() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Edit Hooks for Repo %s\n\n", h.repo.GetName()))

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

func OpenHookEdit(parent richmodel, r repo) richmodel {
	var a []Action

	a = append(a, func() tea.Cmd {
		return func() tea.Msg {
			return openeditor{
				path: r.repopath + "/hooks/post-receive",
			}
		}
	})
	a = append(a, func() tea.Cmd {
		return func() tea.Msg {
			return openeditor{
				path: r.repopath + "/hooks/pre-receive",
			}
		}
	})
	a = append(a, func() tea.Cmd {
		return func() tea.Msg {
			return openeditor{
				path: r.repopath + "/hooks/pre-commit",
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

func openEditor(hook string) error {

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

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return false
	} else {
		return true
	}
}

func (h hookedit) GetKeymapString() string {
	return "up/down - select, enter - open editor, esc - back"
}
