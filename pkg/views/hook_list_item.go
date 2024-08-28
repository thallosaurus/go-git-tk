package views

type hookListItem struct {
	label string
	file  string
}

func (i hookListItem) FilterValue() string { return i.label }
