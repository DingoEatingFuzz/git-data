package gitjson

import (
	"bufio"
	"context"
	"dingoeatingfuzz/git.json/pkg/gitjson"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GitHubAllIssues struct{}

type issue struct {
	Node struct {
		Title        githubv4.String
		Author       user
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

type GitHubIssue struct {
	Title          string    `json:"title"`
	Author         string    `json:"author"`
	CreatedAt      time.Time `json:"createdAt"`
	ClosedAt       time.Time `json:"closedAt"`
	Closed         bool      `json:"closed"`
	Locked         bool      `json:"locked"`
	Participants   []string  `json:"participants"`
	CommentsCount  int       `json:"commentsCount"`
	ReactionsCount int       `json:"reactionsCount"`
}

type user struct {
	Login githubv4.String
}

func (ai *GitHubAllIssues) Source() gitjson.Source {
	return gitjson.GitHubSource
}

func (ai *GitHubAllIssues) Name() string {
	return "All GitHub Issues"
}

func (ai *GitHubAllIssues) Run(git *gitjson.Git, progress func(string, float64, bool)) {
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
			Issues struct {
				Edges    []issue
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"issues(first: $num, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"repo":   githubv4.String(repo),
		"num":    githubv4.Int(50),
		"cursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	var allIssues []issue

	for {
		err := client.Query(context.Background(), &q, variables)
		if err != nil {
			progress(fmt.Sprintf("Woah error: %v", err), 0, false)
			return
		}

		allIssues = append(allIssues, q.Repository.Issues.Edges...)

		if !q.Repository.Issues.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(q.Repository.Issues.PageInfo.EndCursor)
	}

	f, err := os.Create("github-all-issues.ndjson")
	if err != nil {
		progress("Cannot create a file, aborting", 0, false)
		return
	}
	defer f.Close()

	// TODO: Should scripts be responsible for writing files? Or should they send bytes to a channel?
	w := bufio.NewWriter(f)
	length := len(allIssues)

	for i, issue := range allIssues {
		var participants []string
		for _, p := range issue.Node.Participants.Nodes {
			participants = append(participants, string(p.Login))
		}

		row := &GitHubIssue{
			Title:          string(issue.Node.Title),
			Author:         string(issue.Node.Author.Login),
			CreatedAt:      issue.Node.CreatedAt.Time,
			ClosedAt:       issue.Node.ClosedAt.Time,
			Closed:         bool(issue.Node.Closed),
			Locked:         bool(issue.Node.Locked),
			Participants:   participants,
			CommentsCount:  int(issue.Node.Comments.TotalCount),
			ReactionsCount: int(issue.Node.Reactions.TotalCount),
		}

		str, err := json.Marshal(row)
		if err != nil {
			continue
		}

		w.Write(str)

		progress(fmt.Sprintf("%d of %d issues", i, length), float64(i)/float64(length), i == length)
	}
}
