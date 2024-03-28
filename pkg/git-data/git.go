package gitdata

import (
	"fmt"
	"io"
	"os"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Git struct {
	RepoUrl string
	Dir     string
	Repo    *gogit.Repository
}

func (g *Git) Clone() *Git {
	if g.Dir == ":memory:" {
		fmt.Println("Cloning in memory")
		repo, gitErr := gogit.Clone(memory.NewStorage(), nil, &gogit.CloneOptions{
			URL:      g.RepoUrl,
			Progress: os.Stdout,
		})

		if gitErr != nil {
			panic(fmt.Sprintf("Could not clone repo: %v", gitErr))
		}

		g.Repo = repo
	} else {
		f, err := os.Open(g.Dir)
		if err != nil {
			panic(fmt.Sprintf("Could not open directory '%v': %v", g.Dir, err))
		}
		defer f.Close()

		_, err = f.Readdirnames(1)
		if err == io.EOF {
			fmt.Println("Cloning")
			// Directory is empty, clone
			repo, gitErr := gogit.PlainClone(g.Dir, false, &gogit.CloneOptions{
				URL:      g.RepoUrl,
				Progress: os.Stdout,
			})

			if gitErr != nil {
				panic(fmt.Sprintf("Could not clone repo: %v", gitErr))
			}

			g.Repo = repo
		} else {
			fmt.Println("Existing")
			repo, gitErr := gogit.PlainOpen(g.Dir)

			if gitErr != nil {
				panic(fmt.Sprintf("Could not open repo: %v", gitErr))
			}

			remotes, rErr := repo.Remotes()
			if rErr != nil {
				panic(fmt.Sprintf("Could not read remotes for repo: %v", rErr))
			}

			for _, r := range remotes {
				// Conventionally true in most cases
				if r.Config().Name == "origin" {
					g.RepoUrl = r.Config().URLs[0]
					break
				}
			}

			if g.RepoUrl == "" {
				fmt.Println("!! No RepoUrl provided or determined for repo")
			}

			g.Repo = repo
		}
	}

	return g
}
