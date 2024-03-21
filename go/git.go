package main

import (
	"fmt"
	"io"
	"os"

	gogit "github.com/go-git/go-git/v5"
)

type git struct {
	repoUrl string
	dir     string
}

func (g git) Init() git {
	f, err := os.Open(g.dir)
	if err != nil {
		panic(fmt.Sprintf("Could not open directory '%v': %v", g.dir, err))
	}
	defer f.Close()

	// Check if dir == :memory: and clone using memory storage

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		// Directory is empty, clone
		_, gitErr := gogit.PlainClone(g.dir, false, &gogit.CloneOptions{
			URL:      g.repoUrl,
			Progress: os.Stdout,
		})

		if gitErr != nil {
			panic(fmt.Sprintf("Could not clone repo: %v", gitErr))
		}
	}

	return g
}
