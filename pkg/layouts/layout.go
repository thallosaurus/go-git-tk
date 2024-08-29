package layouts

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Describes the layout of the outer origin view

func GetOriginInnerWidth() int {
	return term_width - MainStyle.GetHorizontalFrameSize() + 6
}

func GetOriginInnerHeight() int {
	return term_height - MainStyle.GetVerticalFrameSize() - 2
}

func GetContentInnerWidth() int {
	return GetOriginInnerWidth() - ContentStyle.GetHorizontalFrameSize() - HeaderStyle.GetHorizontalFrameSize() - FooterStyle.GetHorizontalFrameSize()
}

func GetContentInnerHeight() int {
	return GetOriginInnerHeight() - ContentStyle.GetVerticalFrameSize() - HeaderStyle.GetVerticalFrameSize() - FooterStyle.GetVerticalFrameSize()
}

func UpdateTermSize(msg tea.WindowSizeMsg) {
	term_width = msg.Width
	term_height = msg.Height
}

func borderize(s lipgloss.Style, v bool) lipgloss.Style {
	return s.Border(lipgloss.NormalBorder(), v)
}
func TurnOnDebugBorders() {
	MainStyle = borderize(MainStyle, true)
	HeaderStyle = borderize(HeaderStyle, true)
	ContentStyle = borderize(ContentStyle, true)
	FooterStyle = borderize(FooterStyle, true)
}
func TurnOffDebugBorders() {
	MainStyle = borderize(MainStyle, false)
	HeaderStyle = borderize(HeaderStyle, false)
	ContentStyle = borderize(ContentStyle, false)
	FooterStyle = borderize(FooterStyle, false)
}

var (
	term_width   = 0
	term_height  = 0
	MainStyle    = lipgloss.NewStyle().Padding(0).BorderForeground(lipgloss.Color("170"))
	HeaderStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("170")).Padding(1).BorderForeground(lipgloss.Color("170"))
	ContentStyle = lipgloss.NewStyle().Padding(1).BorderForeground(lipgloss.Color("170"))
	FooterStyle  = lipgloss.NewStyle().BorderForeground(lipgloss.Color("170"))

	ItemStyle     = lipgloss.NewStyle().PaddingLeft(2)
	SelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))

	DangerousStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Bold(true)
	EmptyStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
)
