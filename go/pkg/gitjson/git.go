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
}

func (g Git) Init() Git {
	f, err := os.Open(g.Dir)
	if err != nil {
		panic(fmt.Sprintf("Could not open directory '%v': %v", g.Dir, err))
	}
	defer f.Close()

	// Check if dir == :memory: and clone using memory storage

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		// Directory is empty, clone
		_, gitErr := gogit.PlainClone(g.Dir, false, &gogit.CloneOptions{
			URL:      g.RepoUrl,
			Progress: os.Stdout,
		})

		if gitErr != nil {
			panic(fmt.Sprintf("Could not clone repo: %v", gitErr))
		}
	}

	return g
}
