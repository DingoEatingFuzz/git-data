package gitjson

import (
	"fmt"
	"os"

	"dingoeatingfuzz/git.json/pkg/gitjson"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitjson",
	Short: "gitjson turns a git repo and a GitHub project into machine readable ndjson files",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO these should come from an args parser
		repo := &gitjson.Git{
			RepoUrl: "https://github.com/hashicorp/vagrant.git",
			Dir:     "/tmp/foo",
		}
		// TODO how to make this async?
		repo.Init()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
