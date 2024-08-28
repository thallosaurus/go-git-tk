package gitlib

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage/memory"
)

func MakeNewRepo(path string) (*Repo, error) {
	_, err := createGitRepo(path)
	if err != nil {
		//repo.SetDescription(repo.GetName() + " repository")
		return nil, err
	} else {
		rr, err := importExistingRepo(path)

		return rr, err
	}
}

func importExistingRepo(path string) (*Repo, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	fs := memfs.New()

	_, err = git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: path,
	})

	repo := &Repo{
		git:        r,
		clonedRepo: fs,
		Repopath:   path,
	}

	return repo, nil
}

func createGitRepo(repoPath string) (*git.Repository, error) {
	//path := "./repos/" + sanitize_name(repoName)

	return git.PlainInitWithOptions(repoPath, &git.PlainInitOptions{
		Bare: true,
	})
}

type Repo struct {
	git        *git.Repository
	clonedRepo billy.Filesystem
	Repopath   string
}

func (r Repo) GetReadme() (string, error) {
	if _, err := r.clonedRepo.Stat("README.md"); err != nil {
		// doesnt exist
		return "", err
	} else {

		f, err := r.clonedRepo.Open("README.md")
		if err != nil {
			return "", err
		}

		b, err := io.ReadAll(f)
		if err != nil {
			return "", err
		}

		return string(b), nil
	}
}

func (r Repo) GetBranches() ([]string, error) {
	branches, err := r.git.Branches()
	if err != nil {
		return nil, err
	}

	return iterStringArray(branches), nil
}

func iterStringArray(iter storer.ReferenceIter) []string {
	var s []string

	iter.ForEach(func(r *plumbing.Reference) error {
		s = append(s, r.Name().String())
		return nil
	})

	return s
}

func (r Repo) GetTags() ([]string, error) {
	tags, err := r.git.Tags()
	if err != nil {
		log.Panic(err)
	}

	return iterStringArray(tags), nil
}

func (r Repo) GetCommitters() ([]string, error) {
	c, err := r.git.CommitObjects()
	if err != nil {
		return nil, err
	}

	return iterCommittersStringArray(c), nil
}

func iterCommittersStringArray(c object.CommitIter) []string {
	var s []string
	committers := make(map[string]string)
	c.ForEach(func(c *object.Commit) error {
		committers[c.Author.Name] = c.Author.Email
		return nil
	})

	for author, email := range committers {
		s = append(s, fmt.Sprintf("%s <%s>", author, email))
	}

	return s
}

func (r Repo) Rename(newname string) error {
	base := path.Dir(r.Repopath)
	nn := base + "/" + newname

	if err := os.Rename(r.Repopath, nn); err != nil {
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

func (r Repo) GetName() string {
	return path.Base(r.Repopath)
}

func (r Repo) GetDescription() (string, error) {
	f, err := os.ReadFile(r.Repopath + "/description")

	if err != nil {
		return "", err
	}

	return string(f), nil
}

func (r Repo) SetDescription(desc string) error {
	return os.WriteFile(r.Repopath+"/description", []byte(desc), 0644)
}

func GetRepos(root string) ([]Repo, error) {
	//workdir := "./repos"

	workdir := root + "/repos"

	entries, err := os.ReadDir(workdir)

	if err != nil {
		return nil, err
	}

	result := make([]Repo, 0)
	for _, e := range entries {
		if r := extractRepo(workdir + "/" + e.Name()); r != nil {
			result = append(result, *r)
		}
	}

	return result, nil
}

func extractRepo(p string) *Repo {
	/*return &Repo{
		git:      r,
		Repopath: p,
	}*/

	f, err := importExistingRepo(p)
	if err != nil {
		log.Panic(err)
	}

	return f
}
