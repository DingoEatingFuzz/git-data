package gitjson

import (
	"bufio"
	gitjson "dingoeatingfuzz/git.json/pkg/gitjson"
	"encoding/json"
	"fmt"
	"os"
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
	Hash           string    `json:"hash"`
}

func (ac *AllCommits) Source() gitjson.Source {
	return gitjson.GitSource
}

func (ac *AllCommits) Run(git *gitjson.Git, progress func(string, float64)) {
	count := 0
	curr := 0
	skipped := 0
	totalBytes := 0

	// I wish there was a better way to do this, thought go-git would be more feature complete
	countIter, logErr := git.Repo.Log(&gogit.LogOptions{})
	if logErr != nil {
		progress(fmt.Sprintf("Woah error: %v", logErr), 0)
		return
	}

	err := countIter.ForEach(func(c *object.Commit) error {
		count += 1
		return nil
	})

	if err != nil {
		progress(fmt.Sprintf("Woah error: %v", err), 0)
		return
	}

	f, err := os.Create("all-commits.ndjson")
	if err != nil {
		progress("Cannot create a file, aborting", 0)
		return
	}
	defer f.Close()

	// TODO: Should scripts be responsible for writing files? Or should they send bytes to a channel?
	w := bufio.NewWriter(f)

	iter, _ := git.Repo.Log(&gogit.LogOptions{})
	_ = iter.ForEach(func(c *object.Commit) error {
		curr += 1

		commit := &GitCommit{
			Author:         c.Author.Name,
			AuthorEmail:    c.Author.Email,
			Committer:      c.Committer.Name,
			CommitterEmail: c.Committer.Email,
			Date:           c.Committer.When,
			Message:        c.Message,
			Hash:           c.Hash.String(),
		}

		str, err := json.Marshal(commit)
		if err != nil {
			skipped += 1
			return nil
		}

		b, werr := w.Write(str)
		if werr != nil {
			skipped += 1
			return nil
		}

		totalBytes += b
		progress(fmt.Sprintf("%d of %d commits", curr, count), float64(curr)/float64(count))

		return nil
	})

	w.Flush()
}
