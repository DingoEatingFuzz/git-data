package gitjson

import (
	"fmt"
	"io"
	"os"

	gogit "github.com/go-git/go-git/v5"
)

type Git struct {
	RepoUrl string
	Dir     string
	Repo    *gogit.Repository
}

func (g *Git) Clone() *Git {
	f, err := os.Open(g.Dir)
	if err != nil {
		panic(fmt.Sprintf("Could not open directory '%v': %v", g.Dir, err))
	}
	defer f.Close()

	// Check if dir == :memory: and clone using memory storage

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

		g.Repo = repo
	}

	return g
}
