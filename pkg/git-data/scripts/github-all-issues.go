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

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GitHubAllIssues struct{}

type issue struct {
	Node struct {
		Title        githubv4.String
		Author       user
		Url          githubv4.String
		CreatedAt    githubv4.DateTime
		ClosedAt     githubv4.DateTime
		Closed       githubv4.Boolean
		Locked       githubv4.Boolean
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

type rateLimit struct {
	Remaining githubv4.Int
	Used      githubv4.Int
	ResetAt   githubv4.DateTime
}

type user struct {
	Login githubv4.String
}

type GitHubIssue struct {
	Owner          string    `json:"owner"`
	Repo           string    `json:"repo"`
	Title          string    `json:"title"`
	Author         string    `json:"author"`
	Url            string    `json:"url"`
	CreatedAt      time.Time `json:"createdAt"`
	ClosedAt       time.Time `json:"closedAt"`
	Closed         bool      `json:"closed"`
	Locked         bool      `json:"locked"`
	Participants   []string  `json:"participants"`
	CommentsCount  int       `json:"commentsCount"`
	ReactionsCount int       `json:"reactionsCount"`
}

func (ai *GitHubAllIssues) Source() gitdata.Source {
	return gitdata.GitHubSource
}

func (ai *GitHubAllIssues) Name() string {
	return "All GitHub Issues"
}

func (ai *GitHubAllIssues) Run(git *gitdata.Git, config *gitdata.RunnerConfig, progress func(string, float64, bool)) {
	progress("Started", 0, false)
	r, _ := regexp.Compile("github.com/(.+?)/(.+?)(/|\\.git)?$")
	matches := r.FindStringSubmatch(git.RepoUrl)
	owner := matches[1]
	repo := matches[2]

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)

	rateLimiter, rErr := github_ratelimit.NewRateLimitWaiterClient(nil)
	if rErr != nil {
		progress(fmt.Sprintf("Error making rate limiter: %v", rErr), 0, false)
	}

	tripperCtx := context.WithValue(context.Background(), oauth2.HTTPClient, rateLimiter)
	httpClient := oauth2.NewClient(tripperCtx, src)

	client := githubv4.NewClient(httpClient)

	var q struct {
		Repository struct {
			Issues struct {
				Edges    []issue
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"issues(first: $num, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
		RateLimit rateLimit
	}

	var c struct {
		Repository struct {
			Issues struct {
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

	length := int(c.Repository.Issues.TotalCount)

	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"repo":   githubv4.String(repo),
		"num":    githubv4.Int(50),
		"cursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	f, err := os.Create(path.Join(config.DataDir, config.FilePrefix+"github-all-issues"+config.FileSuffix+".ndjson"))
	if err != nil {
		progress("Cannot create a file, aborting", 0, false)
		return
	}
	defer f.Close()

	// TODO: Should scripts be responsible for writing files? Or should they send bytes to a channel?
	w := bufio.NewWriter(f)
	enc := json.NewEncoder(w)
	curr := 0
	retries := 0

	for {
		err := client.Query(context.Background(), &q, variables)
		if err != nil {
			progress(fmt.Sprintf("Paging error (waiting 30s to retry): %v", err), 0, false)
			time.Sleep(30 * time.Second)
			retries++
			if retries >= 10 {
				progress(fmt.Sprintf("Giving up after 10 retries: %v", variables), 0, false)
				break
			} else {
				continue
			}
		}

		for _, issue := range q.Repository.Issues.Edges {
			// Count position across batch requests
			curr++
			var participants []string
			for _, p := range issue.Node.Participants.Nodes {
				participants = append(participants, string(p.Login))
			}

			row := &GitHubIssue{
				Owner:          owner,
				Repo:           repo,
				Title:          string(issue.Node.Title),
				Author:         string(issue.Node.Author.Login),
				Url:            string(issue.Node.Url),
				CreatedAt:      issue.Node.CreatedAt.Time,
				ClosedAt:       issue.Node.ClosedAt.Time,
				Closed:         bool(issue.Node.Closed),
				Locked:         bool(issue.Node.Locked),
				Participants:   participants,
				CommentsCount:  int(issue.Node.Comments.TotalCount),
				ReactionsCount: int(issue.Node.Reactions.TotalCount),
				Labels:         labels,
			}

			err := enc.Encode(row)

			if err != nil {
				continue
			}

			progress(fmt.Sprintf("%d of %d issues (rate limit: %d)", curr, length, q.RateLimit.Remaining), float64(curr)/float64(length), false)
		}

		if !q.Repository.Issues.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(q.Repository.Issues.PageInfo.EndCursor)

		if q.RateLimit.Remaining < 10 {
			diff := q.RateLimit.ResetAt.Sub(time.Now())
			progress(fmt.Sprintf("%d of %d issues Hit the rate limit! Waiting %f seconds", curr, length, diff.Seconds()), float64(curr)/float64(length), false)
			time.Sleep(diff)
		}
	}

	w.Flush()
	progress("No more pages", 1, true)
}
