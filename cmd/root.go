package cmd

import (
	"errors"
	"net/url"
	"os"

	"dingoeatingfuzz/git-data/pkg/git-data"
	scripts "dingoeatingfuzz/git-data/pkg/git-data/scripts"

	"github.com/spf13/cobra"
)

var SilentErr = errors.New("SilentErr")
var UseDisk bool
var DataDir string
var FilePrefix string
var FileSuffix string

var RootCmd = &cobra.Command{
	Use:           "git-data",
	Short:         "git-data turns a git repo and a GitHub project into machine readable ndjson files",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var dir string
		var repoUrl string

		if IsUrl(args[0]) {
			if UseDisk {
				tmp, err := os.MkdirTemp("", "gitdata-")
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

		repo := &gitdata.Git{
			RepoUrl: repoUrl,
			Dir:     dir,
		}

		repo.Clone()

		// Create runner and run
		runner := &gitdata.Runner{
			Scripts: []gitdata.Script{
				&scripts.AllCommits{},
				&scripts.AllCommitsWithFiles{},
				&scripts.GitHubAllIssues{},
				&scripts.GitHubAllPulls{},
			},
			Git: repo,
		}

		runner.Run(gitdata.RunnerConfig{
			DataDir:    DataDir,
			FilePrefix: FilePrefix,
			FileSuffix: FileSuffix,
		})

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

	RootCmd.Flags().StringVarP(
		&DataDir, "data-dir", "d", "",
		"Where to write exported files",
	)

	RootCmd.Flags().StringVarP(
		&FilePrefix, "prefix", "p", "",
		"A common prefix for all exported file names",
	)

	RootCmd.Flags().StringVarP(
		&FileSuffix, "suffix", "s", "",
		"A common suffx for all exported file names",
	)
}
