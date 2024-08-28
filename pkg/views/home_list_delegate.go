package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

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
			return selectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
