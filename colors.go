package tuiplayground

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#0080ff")).Bold(true)
	actionStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8000ff"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff"))
)
