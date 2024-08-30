package views

import (
	"fmt"
	"go-git-tk/pkg/layouts"
	"io"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type homeListDelegate struct{}

func (h homeListDelegate) Height() int {
	return 2
}

func (h homeListDelegate) Spacing() int {
	return 1
}

func (d homeListDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d homeListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(homeListItem)
	if !ok {
		return
	}

	desc, err := i.repo.GetDescription()
	if err != nil {
		log.Panic(err)
		return
	}

	//draw title

	var title string
	var subtitle string
	var style lipgloss.Style

	title = fmt.Sprintf("%d. %s", index+1, i.repo.GetName())
	subtitle = fmt.Sprintf("   %s", desc)
	if index == m.Index() {
		style = layouts.SelectedStyle
		//	return layouts.SelectedStyle.Render("> " + strings.Join(s, " "))
	} else {
		style = layouts.ItemStyle

	}

	fmt.Fprint(w, style.Render(fmt.Sprintf("%s\n%s", title, subtitle)))
}
