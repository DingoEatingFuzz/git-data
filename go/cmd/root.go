package gitjson

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"dingoeatingfuzz/git.json/pkg/gitjson"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

var SilentErr = errors.New("SilentErr")
var IsDevMode bool

var RootCmd = &cobra.Command{
	Use:           "gitjson",
	Short:         "gitjson turns a git repo and a GitHub project into machine readable ndjson files",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO if the arg is a path, we should assume that the intent is to use an existing checkout
		var dir string
		var repoUrl string

		if IsUrl(args[0]) {
			tmp, err := os.MkdirTemp("", "gitjson-")
			if err != nil {
				return err
			}
			dir = tmp
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
		count := 0
		iter, _ := repo.Repo.Log(&gogit.LogOptions{})
		_ = iter.ForEach(func(c *object.Commit) error {
			count = count + 1
			return nil
		})

		fmt.Println(fmt.Sprintf("%d commits", count))
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

	// RootCmd.Flags().BoolVarP(
	// 	&IsDevMode, "dev-mode", "dd", false,
	// 	"Run in dev mode (don't use tmp directories)",
	// )
}
