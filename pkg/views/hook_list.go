package views

import (
	"go-git-tk/pkg/gitlib"
	"go-git-tk/pkg/layouts"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type hook_list struct {
	list   list.Model
	parent Richmodel
	repo   gitlib.Repo
}

type openeditor struct {
	path string
}

const hookChmod = 0755

var (
	edit_hook key.Binding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "edit hook"),
	)
	cancel_key key.Binding = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)
)

func MakeHookList(parent Richmodel, r gitlib.Repo) hook_list {
	items := []list.Item{
		hookListItem{
			label: "Post-Receive",
			file:  "post-receive",
		},
		hookListItem{
			label: "Pre-Receive",
			file:  "pre-receive",
		},
		hookListItem{
			label: "Pre-Commit",
			file:  "pre-commit",
		},
	}

	l := list.New(items, hookListDelegate{}, layouts.GetContentInnerWidth(), layouts.GetContentInnerHeight())

	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	//list.SetStatusBarItemName("repo", "repos")
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()

	return hook_list{
		list:   l,
		parent: parent,
		repo:   r,
	}
}

func (k hook_list) GetHeaderString() string {
	return "Select Hook to edit"
}
func (k hook_list) Init() tea.Cmd {
	return nil
}

func (k hook_list) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		/*k.list.SetWidth(msg.Width)
		k.list.SetHeight(getViewportHeight())*/

		k.list.SetWidth(layouts.GetContentInnerWidth())
		k.list.SetHeight(layouts.GetContentInnerHeight())
		return k, nil

	case tea.KeyMsg:
		if k.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, edit_hook):
			item, _ := k.list.SelectedItem().(hookListItem)
			e := makeEditorEvent(k.repo.Repopath + "/hooks/" + item.file)
			return k, e

		case key.Matches(msg, cancel_key):
			return k, ChangeView(k.parent)
		}

	case openeditor:
		if err := OpenEditor(msg.path); err != nil {
			return k, ChangeView(errorview{
				Parent: k,
				Err:    err,
			})
		}

		return k, tea.Batch(tea.ShowCursor, ChangeView(k))
	}

	var cmd tea.Cmd
	k.list, cmd = k.list.Update(msg)
	return k, cmd
}

func (k hook_list) View() string {
	return k.list.View()
}

func (k hook_list) GetKeymapString() []key.Binding {
	return []key.Binding{
		edit_hook,
		cancel_key,
	}
}

func makeEditorEvent(path string) tea.Cmd {
	return func() tea.Msg {
		return openeditor{
			path,
		}
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
