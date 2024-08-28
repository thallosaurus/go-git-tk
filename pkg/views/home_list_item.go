package views

import "go-git-tk/pkg/gitlib"

type homeListItem struct {
	repo gitlib.Repo
}

func (i homeListItem) FilterValue() string { return i.repo.GetName() }
