package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type hookListDelegate struct{}

func (h hookListDelegate) Height() int {
	return 1
}

func (h hookListDelegate) Spacing() int {
	return 0
}

func (d hookListDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d hookListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(hookListItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.label)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
