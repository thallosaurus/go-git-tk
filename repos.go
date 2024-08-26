package tuiplayground

import (
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
)

type repo struct {
	git      *git.Repository
	repopath string
}

func (r repo) Rename(newname string) error {
	base := path.Dir(r.repopath)
	nn := base + "/" + newname

	if err := os.Rename(r.repopath, nn); err != nil {
		return err
	}

	rr, err := git.PlainOpen(nn)
	if err != nil {
		return err
	} else {
		r.git = rr
		return nil
	}
}

func (r repo) GetName() string {
	return path.Base(r.repopath)
}

func (r repo) GetDescription() (string, error) {
	f, err := os.ReadFile(r.repopath + "/description")

	if err != nil {
		return "", err
	}

	return string(f), nil
}

func GetRepos(root string) ([]repo, error) {
	//workdir := "./repos"

	workdir := root + "/repos"

	entries, err := os.ReadDir(workdir)

	if err != nil {
		return nil, err
	}

	result := make([]repo, 0)
	for _, e := range entries {
		if r := extractRepo(workdir + "/" + e.Name()); r != nil {
			result = append(result, *r)
		}
	}

	return result, nil
}

func extractRepo(p string) *repo {
	r, err := git.PlainOpen(p)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &repo{
		git:      r,
		repopath: p,
	}
}
