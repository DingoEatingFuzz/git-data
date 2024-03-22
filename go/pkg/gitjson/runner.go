package gitjson

import "fmt"

type source int

const (
	GitSource source = iota
	GitHubSource
)

type Script interface {
	Source() source
	Run(git *Git, progress func(string, float64))
}

type Runner struct {
	Scripts []Script
	Git     *Git
}

func (r *Runner) Run() {
	groups := map[source][]Script{}

	for _, s := range r.Scripts {
		groups[s.Source()] = append(groups[s.Source()], s)
	}

	if len(groups[GitSource]) > 0 {
		// Print header
		// Run all scripts concurrently, attaching a progress element
		fmt.Println("Git Sources")
		for _, s := range groups[GitSource] {
			s.Run(r.Git, func(msg string, progress float64) {
				fmt.Println(fmt.Sprintf("(%f) %v", progress, msg))
			})
		}
	}
}
