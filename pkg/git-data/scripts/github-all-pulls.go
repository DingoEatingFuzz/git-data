package gitdata

import (
	"bufio"
	"context"
	"dingoeatingfuzz/git-data/pkg/git-data"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GitHubAllPulls struct{}

type pullRequest struct {
	Node struct {
		Title        githubv4.String
		Author       user
		Url          githubv4.String
		CreatedAt    githubv4.DateTime
		ClosedAt     githubv4.DateTime
		Closed       githubv4.Boolean
		Locked       githubv4.Boolean
		Merged       githubv4.Boolean
		Additions    githubv4.Int
		Deletions    githubv4.Int
		ChangedFiles githubv4.Int
		Participants struct {
			Nodes []user
		} `graphql:"participants(first: 100)"`
		Comments struct {
			TotalCount githubv4.Int
		}
		Reactions struct {
			TotalCount githubv4.Int
		}
	}
}

type GitHubPullRequest struct {
	Owner          string    `json:"owner"`
	Repo           string    `json:"repo"`
	Title          string    `json:"title"`
	Author         string    `json:"author"`
	Url            string    `json:"url"`
	CreatedAt      time.Time `json:"createdAt"`
	ClosedAt       time.Time `json:"closedAt"`
	Closed         bool      `json:"closed"`
	Locked         bool      `json:"locked"`
	Merged         bool      `json:"merged"`
	Additions      int       `json:"additions"`
	Deletions      int       `json:"deletions"`
	ChangedFiles   int       `json:"changedFiles"`
	Participants   []string  `json:"participants"`
	CommentsCount  int       `json:"commentsCount"`
	ReactionsCount int       `json:"reactionsCount"`
}

func (ai *GitHubAllPulls) Source() gitdata.Source {
	return gitdata.GitHubSource
}

func (ai *GitHubAllPulls) Name() string {
	return "All GitHub Pull Requests"
}

func (ai *GitHubAllPulls) Run(git *gitdata.Git, config *gitdata.RunnerConfig, progress func(string, float64, bool)) {
	progress("Started", 0, false)
	r, _ := regexp.Compile("github.com/(.+?)/(.+?)(/|\\.git)?$")
	matches := r.FindStringSubmatch(git.RepoUrl)
	owner := matches[1]
	repo := matches[2]

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	var q struct {
		Repository struct {
			PullRequests struct {
				Edges    []pullRequest
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"pullRequests(first: $num, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
		RateLimit rateLimit
	}

	var c struct {
		Repository struct {
			PullRequests struct {
				TotalCount githubv4.Int
			}
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	err := client.Query(context.Background(), &c, map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
	})

	if err != nil {
		progress(fmt.Sprintf("Woah error: %v", err), 0, false)
		return
	}

	length := int(c.Repository.PullRequests.TotalCount)

	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"repo":   githubv4.String(repo),
		"num":    githubv4.Int(50),
		"cursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	f, err := os.Create(path.Join(config.DataDir, config.FilePrefix+"github-all-pull-requests"+config.FileSuffix+".ndjson"))
	if err != nil {
		progress("Cannot create a file, aborting", 0, false)
		return
	}
	defer f.Close()

	// TODO: Should scripts be responsible for writing files? Or should they send bytes to a channel?
	w := bufio.NewWriter(f)
	curr := 0

	for {
		err := client.Query(context.Background(), &q, variables)
		if err != nil {
			progress(fmt.Sprintf("Woah error: %v", err), 0, false)
			return
		}

		for _, pull := range q.Repository.PullRequests.Edges {
			// Count position across batch requests
			curr++
			var participants []string
			for _, p := range pull.Node.Participants.Nodes {
				participants = append(participants, string(p.Login))
			}

			row := &GitHubPullRequest{
				Owner:          owner,
				Repo:           repo,
				Title:          string(pull.Node.Title),
				Author:         string(pull.Node.Author.Login),
				Url:            string(pull.Node.Url),
				CreatedAt:      pull.Node.CreatedAt.Time,
				ClosedAt:       pull.Node.ClosedAt.Time,
				Closed:         bool(pull.Node.Closed),
				Locked:         bool(pull.Node.Locked),
				Merged:         bool(pull.Node.Merged),
				Additions:      int(pull.Node.Additions),
				Deletions:      int(pull.Node.Deletions),
				ChangedFiles:   int(pull.Node.ChangedFiles),
				Participants:   participants,
				CommentsCount:  int(pull.Node.Comments.TotalCount),
				ReactionsCount: int(pull.Node.Reactions.TotalCount),
			}

			str, err := json.Marshal(row)
			if err != nil {
				continue
			}

			w.Write(str)
			progress(fmt.Sprintf("%d of %d pulls", curr, length), float64(curr)/float64(length), false)
		}

		if !q.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(q.Repository.PullRequests.PageInfo.EndCursor)

		if q.RateLimit.Remaining < 10 {
			diff := q.RateLimit.ResetAt.Sub(time.Now())
			progress(fmt.Sprintf("%d of %d pulls Hit the rate limit! Waiting %f seconds", curr, length, diff.Seconds()), float64(curr)/float64(length), false)
			time.Sleep(diff)
		}
	}

	w.Flush()
	progress("No more pages", 1, true)
}
