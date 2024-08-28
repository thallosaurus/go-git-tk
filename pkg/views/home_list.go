package views

import (
	"fmt"
	"go-git-tk/pkg/gitlib"

	"io"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type homeListItem struct {
	repo gitlib.Repo
}

func (i homeListItem) FilterValue() string { return i.repo.GetName() }

type home_list struct {
	list list.Model
}

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
)

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
		i.list.SetWidth(msg.Width)
		i.list.SetHeight(msg.Height - 1)
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
	return nil
}

func (i home_list) GetKeymapString() []key.Binding {
	return []key.Binding{
		enter_repo,
		new_repo,
		key_mgmt_key,
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

	list := list.New(items, homeListDelegate{}, term_width, term_height-1)
	list.Title = "Git Server Toolkit"

	list.SetShowStatusBar(false)
	//list.SetStatusBarItemName("repo", "repos")
	list.SetShowHelp(false)
	list.DisableQuitKeybindings()

	return home_list{
		list: list,
	}
}

type homeListDelegate struct{}

func (h homeListDelegate) Height() int {
	return 1
}

func (h homeListDelegate) Spacing() int {
	return 0
}

func (d homeListDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d homeListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(homeListItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.repo.GetName())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
