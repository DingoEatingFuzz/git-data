package gitjson

import (
	"errors"
	"net/url"
	"os"

	"dingoeatingfuzz/git.json/pkg/gitjson"
	scripts "dingoeatingfuzz/git.json/pkg/gitjson/scripts"

	"github.com/spf13/cobra"
)

var SilentErr = errors.New("SilentErr")
var UseDisk bool

var RootCmd = &cobra.Command{
	Use:           "gitjson",
	Short:         "gitjson turns a git repo and a GitHub project into machine readable ndjson files",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var dir string
		var repoUrl string

		if IsUrl(args[0]) {
			if UseDisk {
				tmp, err := os.MkdirTemp("", "gitjson-")
				if err != nil {
					return err
				}
				dir = tmp
			} else {
				dir = ":memory:"
			}
			repoUrl = args[0]
		} else {
			dir = args[0]
		}

		repo := &gitjson.Git{
			RepoUrl: repoUrl,
			Dir:     dir,
		}

		// TODO: stick this in a go routine (all git and vcs operations should be concurrent)
		repo.Clone()

		// Create runner and run
		runner := &gitjson.Runner{
			Scripts: []gitjson.Script{
				&scripts.AllCommits{},
				&scripts.AllCommitsWithFiles{},
				&scripts.GitHubAllIssues{},
			},
			Git: repo,
		}

		runner.Run()

		return nil
	},
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return SilentErr
	})

	RootCmd.Flags().BoolVarP(
		&UseDisk, "use-disk", "u", false,
		"Clone repo to disk instead of in memory",
	)

	// TODO: add destination dir flag (assume ./ by default)
}
