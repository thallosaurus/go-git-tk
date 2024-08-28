package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type listtest struct {
	list list.Model
}

var (
	titleStyleA       = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyleA        = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

func (l listtest) Init() tea.Cmd {
	l.list.SetWidth(term_width)
	l.list.SetHeight(term_height)
	return nil
}
func (l listtest) View() string {
	return l.list.View()
}
func (l listtest) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.list.SetWidth(msg.Width)
		l.list.SetHeight(msg.Height)
		return l, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return l, ChangeView(MakeHomeList())
		}
	}
	var cmd tea.Cmd
	l.list, cmd = l.list.Update(msg)
	return l, cmd
}
func (l listtest) GetKeymapString() string {
	return ""
}

type item string

func (i item) FilterValue() string { return "" }

func MakeListTest() listtest {
	items := []list.Item{
		item("Test"),
		item("Test"),
		item("Test"),
		item("Test"),
		item("Test"),
	}
	list := list.New(items, itemDelegate{}, term_width, term_height-4)
	list.Title = "Testing Lists"
	list.DisableQuitKeybindings()
	list.SetShowHelp(false)
	return listtest{
		list,
	}
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
