package views

import (
	"go-git-tk/pkg/gitlib"
	"go-git-tk/pkg/layouts"

	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	enter_repo key.Binding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter/⏎", "enter repo"),
	)
	new_repo key.Binding = key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new repo"),
	)
	key_mgmt_key key.Binding = key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "key management"),
	)
	import_repo key.Binding = key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "import repo"),
		key.WithDisabled(),
	)
)

type home_list struct {
	list list.Model
}

func (i home_list) GetHeaderString() string {
	return "Git Server Toolkit"
}

func (i home_list) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case open_key_mgmt:
		home := os.Getenv("HOME")

		if err := OpenEditor(home + "/.ssh/authorized_keys"); err != nil {
			return i, ChangeView(errorview{
				Parent: i,
				Err:    err,
			})
		}

		return i, tea.Batch(tea.ShowCursor, ChangeView(i))

	case tea.WindowSizeMsg:
		//i.list.SetWidth(msg.Width)
		//i.list.SetHeight(getViewportHeight())
		i.list.SetWidth(layouts.GetContentInnerWidth())
		i.list.SetHeight(layouts.GetContentInnerHeight())
		return i, nil

	case tea.KeyMsg:
		if i.list.FilterState() == list.Filtering {
			break
		}

		switch {

		case key.Matches(msg, enter_repo):
			item, _ := i.list.SelectedItem().(homeListItem)
			//i.choice = string(i)
			//return i, openRepo(i, item.repo)
			return i, ChangeView(MakeRepoView(i, item.repo))

		case key.Matches(msg, new_repo):
			return i, ChangeView(NewRepoView(MakeHomeList()))

		case key.Matches(msg, key_mgmt_key):
			return i, openKeyMgmt()

			/*case key.Matches(msg, import_repo):
			return i, ChangeView(repo_import{})*/
		}
	}

	var cmd tea.Cmd
	i.list, cmd = i.list.Update(msg)
	return i, cmd
}

func openKeyMgmt() tea.Cmd {
	return func() tea.Msg {
		return open_key_mgmt{}
	}

}

func (i home_list) View() string {
	return i.list.View()
}

func (i home_list) Init() tea.Cmd {
	i.list.SetWidth(layouts.GetContentInnerWidth())
	i.list.SetHeight(layouts.GetContentInnerHeight())

	return tea.Batch(tea.ShowCursor, tea.WindowSize())
}

func (i home_list) GetKeymapString() []key.Binding {
	return []key.Binding{
		enter_repo,
		new_repo,
		key_mgmt_key,
		//import_repo,
	}
}

type open_key_mgmt struct{}

func MakeHomeList() home_list {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	SetupRepos(wd)

	i, err := gitlib.GetRepos(wd)
	if err != nil {
		log.Fatal(err)
	}

	var items []list.Item

	for _, repo := range i {
		items = append(items, homeListItem{repo})
	}

	list := list.New(items, homeListDelegate{}, layouts.GetContentInnerWidth(), layouts.GetContentInnerHeight())

	list.SetShowStatusBar(false)
	list.SetShowTitle(false)
	//list.SetStatusBarItemName("repo", "repos")
	list.SetShowHelp(false)
	list.DisableQuitKeybindings()

	return home_list{
		list: list,
	}
}
