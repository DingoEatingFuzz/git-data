package gitjson

import (
	"errors"
	"fmt"

	"dingoeatingfuzz/git.json/pkg/gitjson"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

var SilentErr = errors.New("SilentErr")

var RootCmd = &cobra.Command{
	Use:           "gitjson",
	Short:         "gitjson turns a git repo and a GitHub project into machine readable ndjson files",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO if the arg is a path, we should assume that the intent is to use an existing checkout
		// dir, err := os.MkdirTemp("", "gitjson-")
		// if err != nil {
		// 	return err
		// }

		repo := &gitjson.Git{
			RepoUrl: args[0],
			Dir:     "dev-cache", //dir,
		}

		// TODO: stick this in a go routine (all git and vcs operations should be concurrent)
		r2 := repo.Clone()
		fmt.Println(fmt.Sprintf("repo %v", repo))
		fmt.Println(fmt.Sprintf("r2 %v", r2))
		count := 0
		iter, _ := repo.Repo.Log(&gogit.LogOptions{})
		_ = iter.ForEach(func(c *object.Commit) error {
			// fmt.Println(c)
			count = count + 1
			return nil
		})

		fmt.Println(fmt.Sprintf("%d commits", count))
		return nil
	},
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
}
