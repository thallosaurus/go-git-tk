package views

import (
	"fmt"
	"go-git-tk/pkg/gitlib"
	"go-git-tk/pkg/layouts"
	"log"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type repo_import struct {
	parent Richmodel
	//repo     gitlib.Repo
	input    textinput.Model
	viewport viewport.Model
}

var (
	repoimport_cancel key.Binding = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)
	repoimport_accept key.Binding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "clone repo"),
	)
)

func OpenRepoImport(parent Richmodel) repo_import {
	input := textinput.New()
	input.Focus()

	vp := viewport.New(layouts.GetContentInnerWidth(), layouts.GetContentInnerHeight()-3)
	vp.SetContent("Enter the URL of the Repository and press enter to import it.")
	//vp.Style = layouts.ContentStyle

	//basename := path.Base(repo.Repopath)
	//input.SetValue(config.Conf.RepoWorkdir + "/" + strings.TrimSuffix(basename, filepath.Ext(basename)))

	return repo_import{
		parent: parent,
		//repo:     repo,
		input:    input,
		viewport: vp,
	}
}

func (k repo_import) Init() tea.Cmd {
	return nil
}

func (k repo_import) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, repoclone_cancel):
			return k, ChangeView(k.parent)

		case key.Matches(msg, repoclone_accept):
			//_, err := k.repo.CloneToWorkdir()
			wd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			r, err := gitlib.ImportRemoteRepo(wd+"/repos/"+path.Base(k.input.Value()), k.input.Value())

			if err != nil {
				return k, ChangeView(errorview{
					Parent: k,
					Err:    err,
				})
			} else {
				k.viewport.SetYOffset(0)
				return k, ChangeView(MakeRepoView(k.parent, *r))
			}

		}
	}

	var cmd tea.Cmd
	k.input, cmd = k.input.Update(msg)
	return k, cmd
}

func (k repo_import) GetHeaderString() string {
	return "Import Remote Repo"
}

func (k repo_import) View() string {
	return fmt.Sprintf("%s\n%s", k.viewport.View(), k.input.View())
}

func (k repo_import) GetKeymapString() []key.Binding {
	return []key.Binding{
		repoimport_accept,
		repoimport_cancel,
	}
}
