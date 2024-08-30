package layouts

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

func colorizeForegroundBorder(s lipgloss.Style, c string) lipgloss.Style {
	return s.BorderForeground(lipgloss.Color(c))
}

func colorizeForeground(s lipgloss.Style, c string) lipgloss.Style {
	return s.Foreground(lipgloss.Color(c))
}
func colorizeBackground(s lipgloss.Style, c string) lipgloss.Style {
	return s.Background(lipgloss.Color(c))
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

func SetAppColor(c string) {
	//ColorTint = c
	MainStyle = colorizeForegroundBorder(MainStyle, c)
	//HeaderStyle = colorizeForeground(HeaderStyle, c)
	HeaderStyle = colorizeBackground(HeaderStyle, c)
	HeaderStyle = colorizeForegroundBorder(HeaderStyle, c)
	ContentStyle = colorizeForegroundBorder(ContentStyle, c)
	FooterStyle = colorizeForegroundBorder(FooterStyle, c)
	FooterStyle = colorizeForeground(FooterStyle, c)
	SelectedStyle = colorizeForeground(SelectedStyle, c)
	ListStyle = colorizeForeground(ListStyle, c)
}

func GetFooterWidth() int {
	return GetContentInnerWidth() + MainStyle.GetHorizontalFrameSize()
}

var (
	term_width   = 0
	term_height  = 0
	MainStyle    = lipgloss.NewStyle()
	HeaderStyle  = lipgloss.NewStyle().Bold(true).Padding(1)
	ContentStyle = lipgloss.NewStyle().Padding(1)
	FooterStyle  = lipgloss.NewStyle()

	ItemStyle     = lipgloss.NewStyle().PaddingLeft(1)
	SelectedStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).PaddingLeft(1)
	SubItemStyle  = lipgloss.NewStyle().PaddingLeft(2)

	ListStyle = lipgloss.NewStyle()

	DangerousStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Bold(true)
	EmptyStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	//ColorTint      = "#ffffff"
)
