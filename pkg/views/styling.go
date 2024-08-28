package views

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("170")).PaddingLeft(1).PaddingRight(1).PaddingBottom(1).PaddingTop(1)

	itemStyle     = lipgloss.NewStyle().PaddingLeft(2)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))

	dangerousStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Bold(true)
	emptyStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	keymapStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	mainStyle      = lipgloss.NewStyle().Padding(1, 1, 1, 1)
)

func getViewportHeight() int {
	return term_height - (4 + keymapStyle.GetHorizontalFrameSize())
}

func getInnerViewportWidth() int {
	return term_width - 2
}
