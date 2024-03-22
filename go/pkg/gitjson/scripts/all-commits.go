package gitjson

import (
	gitjson "dingoeatingfuzz/git.json/pkg/gitjson"
	"encoding/json"
	"fmt"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type AllCommits struct{}

type GitCommit struct {
	Author         string    `json:"author"`
	AuthorEmail    string    `json:"authorEmail"`
	Committer      string    `json:"committer"`
	CommitterEmail string    `json:"committerEmail"`
	Date           time.Time `json:"date"`
	Message        string    `json:"message"`
}

func (ac *AllCommits) Source() gitjson.Source {
	return gitjson.GitSource
}

func (ac *AllCommits) Run(git *gitjson.Git, progress func(string, float64)) {
	count := 0
	curr := 0
	skipped := 0

	// I wish there was a better way to do this, thought go-git would be more feature complete
	countIter, _ := git.Repo.Log(&gogit.LogOptions{})
	_ = countIter.ForEach(func(c *object.Commit) error {
		count = count + 1
		return nil
	})

	progress(fmt.Sprintf("Logging %d commits in main branch", count), 0)

	iter, _ := git.Repo.Log(&gogit.LogOptions{})
	_ = iter.ForEach(func(c *object.Commit) error {
		curr = curr + 1

		commit := &GitCommit{
			Author:         c.Author.Name,
			AuthorEmail:    c.Author.Email,
			Committer:      c.Committer.Name,
			CommitterEmail: c.Committer.Email,
			Date:           c.Committer.When,
			Message:        c.Message,
		}

		str, err := json.Marshal(commit)
		if err != nil {
			skipped = skipped + 1
		}

		if curr%1000 == 0 {
			progress(string(str), float64(curr)/float64(count))
		}

		// Log the commit to a file
		return nil
	})
}
