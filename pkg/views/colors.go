package views

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().Bold(true)
	//actionStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff"))
	selectedStyle  = lipgloss.NewStyle()
	helpStyle      = lipgloss.NewStyle()
	dangerousStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Bold(true)
	emptyStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
)
