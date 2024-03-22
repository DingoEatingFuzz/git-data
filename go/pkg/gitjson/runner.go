package gitjson

import "fmt"

type Source int

const (
	GitSource Source = iota
	GitHubSource
)

// TODO: thread bubbletea in here somehow
type Script interface {
	Source() Source
	// TODO: If progress is going to channel bytes, maybe we need lifecycle hooks
	Run(git *Git, progress func(string, float64))
}

type Runner struct {
	Scripts []Script
	Git     *Git
}

func (r *Runner) Run() {
	groups := map[Source][]Script{}

	for _, s := range r.Scripts {
		groups[s.Source()] = append(groups[s.Source()], s)
	}

	if len(groups[GitSource]) > 0 {
		// Print header
		// Run all scripts concurrently, attaching a progress element
		fmt.Println("Git Sources")
		for _, s := range groups[GitSource] {
			s.Run(r.Git, func(msg string, progress float64) {
				fmt.Println(fmt.Sprintf("(%d) %v", int(progress*100), msg))
			})
		}
	}
}
