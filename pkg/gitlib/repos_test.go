package gitlib_test

import (
	"go-git-tk/pkg/gitlib"
	"os"
	"testing"
)

func TestMakeNewRepo(t *testing.T) {
	repoName := "__testing_single_commit.git"
	wd := os.Getenv("ROOT_CWD")

	t.Log("cwd: " + wd)

	r, err := gitlib.GetRepos(wd)

	if err != nil {
		t.Error(err)
	}

	var found *gitlib.Repo

	for _, repo := range r {
		if repo.GetName() == repoName {
			found = &repo
		}
	}

	if found == nil {
		t.Error("Repo names didn't match")
	}
}
